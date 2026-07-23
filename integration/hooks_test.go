package main

import (
	"context"
	"errors"
	"fmt"
	"integration/valk"
	"integration/valk/user"
	"strings"
	"testing"
)

var errCustomUnique = errors.New("custom unique error")
var errShortCircuit = errors.New("short circuit error")

type ctxKey string

func TestHooks(t *testing.T) {
	ctx := context.Background()

	t.Run("chaining_order_and_context_propagation", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		var executionOrder []string

		db.User.Use(user.Extension{
			Create: func(ctx context.Context, input *user.CreateInput, next user.CreateQuery) (*valk.User, error) {
				executionOrder = append(executionOrder, "first_pre")
				input.Email = input.Email + "-first"
				ctx = context.WithValue(ctx, ctxKey("key1"), "val1")

				res, err := next(ctx, input)

				executionOrder = append(executionOrder, "first_post")
				return res, err
			},
		})

		db.User.Use(user.Extension{
			Create: func(ctx context.Context, input *user.CreateInput, next user.CreateQuery) (*valk.User, error) {
				executionOrder = append(executionOrder, "second_pre")
				input.Email = input.Email + "-second"

				val1 := ctx.Value(ctxKey("key1"))
				if val1 != "val1" {
					return nil, fmt.Errorf("missing key1")
				}
				ctx = context.WithValue(ctx, ctxKey("key2"), "val2")

				res, err := next(ctx, input)

				executionOrder = append(executionOrder, "second_post")
				return res, err
			},
		})

		db.User.Use(user.Extension{
			Create: func(ctx context.Context, input *user.CreateInput, next user.CreateQuery) (*valk.User, error) {
				executionOrder = append(executionOrder, "third_pre")
				input.Email = input.Email + "-third"

				val2 := ctx.Value(ctxKey("key2"))
				if val2 != "val2" {
					return nil, fmt.Errorf("missing key2")
				}

				res, err := next(ctx, input)

				executionOrder = append(executionOrder, "third_post")
				return res, err
			},
		})

		u, err := db.User.Create().
			SetEmail("chain").
			SetPhoneNum("111").
			Exec(ctx)

		if err != nil {
			t.Fatalf("create failed: %v", err)
		}

		expectedOrder := []string{"first_pre", "second_pre", "third_pre", "third_post", "second_post", "first_post"}
		if len(executionOrder) != len(expectedOrder) {
			t.Fatalf("expected order len %d, got %d", len(expectedOrder), len(executionOrder))
		}
		for i, v := range expectedOrder {
			if executionOrder[i] != v {
				t.Errorf("at index %d expected %q, got %q", i, v, executionOrder[i])
			}
		}

		if !strings.HasSuffix(u.Email, "-first-second-third") {
			t.Errorf("expected suffix, got %q", u.Email)
		}
	})

	t.Run("short_circuit_create", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		db.User.Use(user.Extension{
			Create: func(ctx context.Context, input *user.CreateInput, next user.CreateQuery) (*valk.User, error) {
				if input.Email == "short-circuit@example.com" {
					return &valk.User{
						Id:    "mocked-id",
						Email: "short-circuit@example.com",
					}, nil
				}
				return next(ctx, input)
			},
		})

		u, err := db.User.Create().
			SetEmail("short-circuit@example.com").
			SetPhoneNum("222").
			Exec(ctx)

		if err != nil {
			t.Fatalf("expected success, got err: %v", err)
		}

		if u.Id != "mocked-id" {
			t.Errorf("expected mocked-id, got %q", u.Id)
		}

		var count int
		if err := db.Raw().QueryRowContext(ctx, query(`SELECT count(*) FROM "User"`, `SELECT count(*) FROM "User"`)).Scan(&count); err != nil {
			t.Fatalf("scan failed: %v", err)
		}
		if count != 0 {
			t.Errorf("database should be empty, got %d rows", count)
		}
	})

	t.Run("short_circuit_with_error", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		db.User.Use(user.Extension{
			Create: func(ctx context.Context, input *user.CreateInput, next user.CreateQuery) (*valk.User, error) {
				return nil, errShortCircuit
			},
		})

		_, err := db.User.Create().
			SetEmail("error@example.com").
			SetPhoneNum("333").
			Exec(ctx)

		if !errors.Is(err, errShortCircuit) {
			t.Errorf("expected errShortCircuit, got %v", err)
		}
	})

	t.Run("database_error_interception", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		db.User.Use(user.Extension{
			Create: func(ctx context.Context, input *user.CreateInput, next user.CreateQuery) (*valk.User, error) {
				res, err := next(ctx, input)
				if err != nil {
					return nil, errCustomUnique
				}
				return res, nil
			},
		})

		_, err := db.User.Create().SetEmail("uniq@example.com").SetPhoneNum("444").Exec(ctx)
		if err != nil {
			t.Fatalf("first insert failed: %v", err)
		}

		_, err = db.User.Create().SetEmail("uniq@example.com").SetPhoneNum("445").Exec(ctx)
		if !errors.Is(err, errCustomUnique) {
			t.Errorf("expected errCustomUnique, got %v", err)
		}
	})

	t.Run("createmany_mutation_and_abort", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		db.User.Use(user.Extension{
			CreateMany: func(ctx context.Context, args []*user.CreateInput, next user.CreateManyQuery) (int64, error) {
				for _, input := range args {
					if input.Email == "invalid@example.com" {
						return 0, fmt.Errorf("rejected invalid email")
					}
					input.Email = strings.ToUpper(input.Email)
				}
				return next(ctx, args)
			},
		})

		_, err := db.User.CreateMany(
			db.User.Create().SetEmail("user1@example.com").SetPhoneNum("501"),
			db.User.Create().SetEmail("invalid@example.com").SetPhoneNum("502"),
		).Exec(ctx)

		if err == nil || !strings.Contains(err.Error(), "rejected invalid email") {
			t.Errorf("expected invalid email rejection, got %v", err)
		}

		var count int
		if err := db.Raw().QueryRowContext(ctx, query(`SELECT count(*) FROM "User"`, `SELECT count(*) FROM "User"`)).Scan(&count); err != nil {
			t.Fatalf("scan failed: %v", err)
		}
		if count != 0 {
			t.Errorf("database should be empty, got %d rows", count)
		}

		countVal, err := db.User.CreateMany(
			db.User.Create().SetEmail("user1@example.com").SetPhoneNum("501"),
			db.User.Create().SetEmail("user2@example.com").SetPhoneNum("502"),
		).Exec(ctx)

		if err != nil {
			t.Fatalf("create many failed: %v", err)
		}
		if countVal != 2 {
			t.Errorf("expected count 2, got %d", countVal)
		}

		var emails []string
		rows, err := db.Raw().QueryContext(ctx, func() string {
			if getActiveProvider() == "postgres" {
				return `SELECT email FROM "User" ORDER BY email`
			}
			return `SELECT email FROM "User" ORDER BY email`
		}())
		if err != nil {
			t.Fatalf("query failed: %v", err)
		}
		defer rows.Close()
		for rows.Next() {
			var email string
			if err := rows.Scan(&email); err != nil {
				t.Fatalf("scan failed: %v", err)
			}
			emails = append(emails, email)
		}
		if err := rows.Err(); err != nil {
			t.Fatalf("rows error: %v", err)
		}

		if len(emails) != 2 || emails[0] != "USER1@EXAMPLE.COM" || emails[1] != "USER2@EXAMPLE.COM" {
			t.Errorf("emails not mutated correctly: %v", emails)
		}
	})

	t.Run("createmanyandreturn_mutation", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		db.User.Use(user.Extension{
			CreateManyAndReturn: func(ctx context.Context, args []*user.CreateInput, next user.CreateManyAndReturnQuery) ([]*valk.User, error) {
				for _, input := range args {
					input.Email = strings.ToLower(input.Email)
				}
				res, err := next(ctx, args)
				if err == nil {
					for _, u := range res {
						u.Email = u.Email + "-returned"
					}
				}
				return res, err
			},
		})

		users, err := db.User.CreateManyAndReturn(
			db.User.Create().SetEmail("USER1@EXAMPLE.COM").SetPhoneNum("601"),
			db.User.Create().SetEmail("USER2@EXAMPLE.COM").SetPhoneNum("602"),
		).Exec(ctx)

		if err != nil {
			t.Fatalf("create many and return failed: %v", err)
		}

		if len(users) != 2 {
			t.Fatalf("expected 2 users, got %d", len(users))
		}

		if users[0].Email != "user1@example.com-returned" || users[1].Email != "user2@example.com-returned" {
			t.Errorf("unexpected returned emails: %q, %q", users[0].Email, users[1].Email)
		}

		var emails []string
		rows, err := db.Raw().QueryContext(ctx, query(
			`SELECT email FROM "User" ORDER BY email`,
			`SELECT email FROM "User" ORDER BY email`,
		))
		if err != nil {
			t.Fatalf("query failed: %v", err)
		}
		defer rows.Close()
		for rows.Next() {
			var email string
			if err := rows.Scan(&email); err != nil {
				t.Fatalf("scan failed: %v", err)
			}
			emails = append(emails, email)
		}
		if err := rows.Err(); err != nil {
			t.Fatalf("rows error: %v", err)
		}

		if len(emails) != 2 || emails[0] != "user1@example.com" || emails[1] != "user2@example.com" {
			t.Errorf("stored emails in db are unexpected: %v", emails)
		}
	})

	t.Run("transaction_isolation_and_rollback", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		var hookCalled bool
		db.User.Use(user.Extension{
			Create: func(ctx context.Context, input *user.CreateInput, next user.CreateQuery) (*valk.User, error) {
				hookCalled = true
				if input.Email == "rollback@example.com" {
					res, err := next(ctx, input)
					if err == nil {
						return nil, fmt.Errorf("force rollback")
					}
					return res, err
				}
				return next(ctx, input)
			},
		})

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin transaction: %v", err)
		}
		defer func() { _ = tx.Rollback() }()

		_, err = tx.User.Create().
			SetEmail("rollback@example.com").
			SetPhoneNum("701").
			Exec(ctx)

		if err == nil || !strings.Contains(err.Error(), "force rollback") {
			t.Fatalf("expected rollback error, got: %v", err)
		}
		if !hookCalled {
			t.Error("expected hook to be called inside transaction")
		}

		if err := tx.Rollback(); err != nil {
			t.Fatalf("rollback failed: %v", err)
		}

		var count int
		if err := db.Raw().QueryRowContext(ctx, query(`SELECT count(*) FROM "User"`, `SELECT count(*) FROM "User"`)).Scan(&count); err != nil {
			t.Fatalf("scan failed: %v", err)
		}
		if count != 0 {
			t.Errorf("transaction should have rolled back, but found %d users", count)
		}
	})

	t.Run("find_unique_hook_interception", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		db.User.Use(user.Extension{
			FindUnique: func(ctx context.Context, where valk.UniquePredicate[valk.User], additional []valk.PredicateOf[valk.User], selects *user.Select, omits *user.Omit, next user.FindUniqueQuery) (*valk.User, error) {
				// Short circuit if matching specific email
				if cond, ok := where.Data.Value.(string); ok && cond == "intercept@example.com" {
					return &valk.User{
						Id:    "intercepted-unique-id",
						Email: "intercept@example.com",
					}, nil
				}
				return next(ctx, where, additional, selects, omits)
			},
		})

		u, err := db.User.FindUnique(user.Email.EQ("intercept@example.com")).Exec(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if u.Id != "intercepted-unique-id" {
			t.Errorf("expected intercepted-unique-id, got %q", u.Id)
		}
	})

	t.Run("find_first_hook_interception", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		db.User.Use(user.Extension{
			FindFirst: func(ctx context.Context, params valk.QueryParams[valk.User], selects *user.Select, omits *user.Omit, next user.FindFirstQuery) (*valk.User, error) {
				return &valk.User{
					Id:    "intercepted-first-id",
					Email: "first@example.com",
				}, nil
			},
		})

		u, err := db.User.FindFirst().Exec(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if u.Id != "intercepted-first-id" {
			t.Errorf("expected intercepted-first-id, got %q", u.Id)
		}
	})

	t.Run("find_many_hook_filtering", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := db.User.Create().SetEmail("allow-1").SetPhoneNum("801").Exec(ctx)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
		_, err = db.User.Create().SetEmail("allow-2").SetPhoneNum("802").Exec(ctx)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
		_, err = db.User.Create().SetEmail("block-1").SetPhoneNum("803").Exec(ctx)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		db.User.Use(user.Extension{
			FindMany: func(ctx context.Context, params valk.QueryParams[valk.User], selects *user.Select, omits *user.Omit, next user.FindManyQuery) ([]*valk.User, error) {
				// Inject filter so we only allow emails starting with "allow-"
				params.Where = append(params.Where, user.Email.Like("allow-%"))
				return next(ctx, params, selects, omits)
			},
		})

		users, err := db.User.FindMany().Exec(ctx)
		if err != nil {
			t.Fatalf("find many failed: %v", err)
		}

		if len(users) != 2 {
			t.Fatalf("expected 2 users, got %d", len(users))
		}
		for _, u := range users {
			if !strings.HasPrefix(u.Email, "allow-") {
				t.Errorf("unexpected user returned: %q", u.Email)
			}
		}
	})
}
