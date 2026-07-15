package main

import (
	"context"
	"encoding/json"
	"fmt"
	"integration/valk"
	"integration/valk/user"
	"strings"
	"testing"
)

func TestCreateMany_Hooks(t *testing.T) {
	ctx := context.Background()

	t.Run("BeforeCreate hook mutates input during CreateMany", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		client.User.BeforeCreate(func(ctx context.Context, input *valk.UserCreate) error {
			if input.Email == "hooked@example.com" {
				input.PhoneNum = "+188888888"
			}
			return nil
		})

		count, err := client.User.CreateMany(
			client.User.Create().SetEmail("hooked@example.com").SetPhoneNum("+100000000"),
			client.User.Create().SetEmail("normal@example.com").SetPhoneNum("+200000000"),
		).Exec(ctx)

		if err != nil {
			t.Fatalf("CreateMany failed: %v", err)
		}
		if count != 2 {
			t.Fatalf("expected count 2, got %d", count)
		}
		usr, err := client.User.FindUnique(user.Email.EQ("hooked@example.com")).Exec(ctx)

		if err != nil {
			t.Fatalf("query failed: %v", err)
		}
		if usr.PhoneNum != "+188888888" {
			t.Errorf("expected PhoneNum to be mutated to '+188888888', got %q", usr.PhoneNum)
		}
	})

	t.Run("BeforeCreate hook mutates input during CreateManyAndReturn", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		client.User.BeforeCreate(func(ctx context.Context, input *valk.UserCreate) error {
			if input.Email == "hooked@example.com" {
				input.PhoneNum = "+188888888"
			}
			return nil
		})

		users, err := client.User.CreateManyAndReturn(
			client.User.Create().SetEmail("hooked@example.com").SetPhoneNum("+100000001"),
			client.User.Create().SetEmail("normal@example.com").SetPhoneNum("+200000001"),
		).Exec(ctx)

		if err != nil {
			t.Fatalf("CreateManyAndReturn failed: %v", err)
		}
		if len(users) != 2 {
			t.Fatalf("expected 2 users, got %d", len(users))
		}

		for _, u := range users {
			if u.Email == "hooked@example.com" && u.PhoneNum != "+188888888" {
				t.Errorf("expected hooked user PhoneNum to be '+188888888', got %q", u.PhoneNum)
			}
		}
	})

	t.Run("BeforeCreate hook error aborts CreateMany", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		client.User.BeforeCreate(func(ctx context.Context, input *valk.UserCreate) error {
			if input.Email == "reject@example.com" {
				return fmt.Errorf("hook rejected: %s", input.Email)
			}
			return nil
		})

		_, err := client.User.CreateMany(
			client.User.Create().SetEmail("good@example.com").SetPhoneNum("+300000000"),
			client.User.Create().SetEmail("reject@example.com").SetPhoneNum("+300000001"),
		).Exec(ctx)

		if err == nil {
			t.Fatal("expected error from hook rejection, got nil")
		}

		var count int
		if err := client.Raw().QueryRowContext(ctx, `SELECT count(*) FROM "User"`).Scan(&count); err != nil {
			t.Fatalf("failed to scan user count: %v", err)
		}
		if count != 0 {
			t.Fatalf("expected 0 rows after aborted CreateMany, got %d", count)
		}
	})

	t.Run("AfterCreate hook fires during CreateManyAndReturn", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		var afterCalled bool
		var gotRole valk.UserRoleType
		client.User.AfterCreate(func(ctx context.Context, users []*valk.User) error {
			afterCalled = true
			if len(users) > 0 {
				gotRole = users[0].Role
			}
			return nil
		})

		users, err := client.User.CreateManyAndReturn(
			client.User.Create().SetEmail("after@example.com").SetPhoneNum("+500000000"),
		).Exec(ctx)

		if err != nil {
			t.Fatalf("CreateManyAndReturn failed: %v", err)
		}
		if len(users) != 1 {
			t.Fatalf("expected 1 user, got %d", len(users))
		}
		if !afterCalled {
			t.Fatal("expected AfterCreate hook to be called")
		}
		if gotRole != users[0].Role {
			t.Errorf("AfterCreate received role %v, expected %v", gotRole, users[0].Role)
		}
	})

	t.Run("AfterCreate hook error aborts CreateManyAndReturn", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		client.User.AfterCreate(func(ctx context.Context, users []*valk.User) error {
			return fmt.Errorf("after hook rejected")
		})

		var count int
		if err := client.Raw().QueryRowContext(ctx, `SELECT count(*) FROM "User"`).Scan(&count); err != nil {
			t.Fatalf("failed to scan user count: %v", err)
		}
		prevCount := count

		_, err := client.User.CreateManyAndReturn(
			client.User.Create().SetEmail("aftershoot@example.com").SetPhoneNum("+600000000"),
		).Exec(ctx)

		if err == nil {
			t.Fatal("expected error from AfterCreate rejection, got nil")
		}
		if !strings.Contains(err.Error(), "after hook rejected") {
			t.Errorf("expected 'after hook rejected' in error, got %v", err)
		}

		if err := client.Raw().QueryRowContext(ctx, `SELECT count(*) FROM "User"`).Scan(&count); err != nil {
			t.Fatalf("failed to scan user count: %v", err)
		}
		if count != prevCount+1 {
			t.Fatalf("expected %d rows (insert still committed), got %d", prevCount+1, count)
		}
	})

	t.Run("AfterCreateMany hook receives structs and count", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		var gotUsers []valk.UserCreate
		var gotCount int64
		client.User.AfterCreateMany(func(ctx context.Context, users []valk.UserCreate, count int64) error {
			gotUsers = users
			gotCount = count
			return nil
		})

		count, err := client.User.CreateMany(
			client.User.Create().SetEmail("bulk1@example.com").SetPhoneNum("+700000001"),
			client.User.Create().SetEmail("bulk2@example.com").SetPhoneNum("+700000002"),
		).Exec(ctx)

		if err != nil {
			t.Fatalf("CreateMany failed: %v", err)
		}
		if count != 2 {
			t.Fatalf("expected count 2, got %d", count)
		}
		if gotCount != 2 {
			t.Fatalf("expected hook count 2, got %d", gotCount)
		}
		if len(gotUsers) != 2 {
			t.Fatalf("expected 2 users in hook, got %d", len(gotUsers))
		}
		if gotUsers[0].Email != "bulk1@example.com" {
			t.Errorf("expected email 'bulk1@example.com', got %q", gotUsers[0].Email)
		}
		if gotUsers[1].Email != "bulk2@example.com" {
			t.Errorf("expected email 'bulk2@example.com', got %q", gotUsers[1].Email)
		}
	})

	t.Run("AfterCreateMany hook error does not rollback inserts", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		client.User.AfterCreateMany(func(ctx context.Context, users []valk.UserCreate, count int64) error {
			return fmt.Errorf("after create many rejected")
		})

		_, err := client.User.CreateMany(
			client.User.Create().SetEmail("ghost@example.com").SetPhoneNum("+800000000"),
		).Exec(ctx)

		if err == nil {
			t.Fatal("expected error from AfterCreateMany rejection, got nil")
		}
		if !strings.Contains(err.Error(), "after create many rejected") {
			t.Errorf("expected 'after create many rejected' in error, got %v", err)
		}

		var count int
		if err := client.Raw().QueryRowContext(ctx, query(
			`SELECT count(*) FROM "User" WHERE email = 'ghost@example.com'`,
			`SELECT count(*) FROM "User" WHERE email = 'ghost@example.com'`,
		)).Scan(&count); err != nil {
			t.Fatalf("failed to scan user count: %v", err)
		}
		if count != 1 {
			t.Fatalf("expected row to still exist (insert committed before hook), got %d", count)
		}
	})
}

func TestCreateMany(t *testing.T) {
	ctx := context.Background()
	client, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("CreateMany returns correct count", func(t *testing.T) {
		count, err := client.User.CreateMany(
			client.User.Create().SetEmail("bulk1@example.com").SetPhoneNum("+111"),
			client.User.Create().SetEmail("bulk2@example.com").SetPhoneNum("+222"),
			client.User.Create().SetEmail("bulk3@example.com").SetPhoneNum("+333"),
		).Exec(ctx)

		if err != nil {
			t.Fatalf("CreateMany failed: %v", err)
		}

		if count != 3 {
			t.Errorf("expected count 3, got %d", count)
		}

		var dbCount int
		err = client.Raw().QueryRowContext(ctx, `SELECT count(*) FROM "User"`).Scan(&dbCount)
		if err != nil {
			t.Fatalf("Raw SQL query failed: %v", err)
		}
		if dbCount != 3 {
			t.Errorf("expected 3 users in db, got %d", dbCount)
		}
	})

	t.Run("CreateManyAndReturn works and supports Select", func(t *testing.T) {
		author, err := client.User.Create(
			user.Email.Set("author@example.com"),
			user.PhoneNum.Set("+444"),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to create author: %v", err)
		}

		posts, err := client.Post.CreateManyAndReturn(
			client.Post.Create().SetTitle("Post One").SetAuthorId(author.Id),
			client.Post.Create().SetTitle("Post Two").SetAuthorId(author.Id),
		).Select(valk.PostSelect{
			Id:    true,
			Title: true,
			Author: &valk.UserSelect{
				Email: true,
			},
		}).Exec(ctx)

		if err != nil {
			t.Fatalf("CreateManyAndReturn failed: %v", err)
		}

		if len(posts) != 2 {
			t.Fatalf("expected 2 posts, got %d", len(posts))
		}

		if posts[0].Title != "Post One" || posts[1].Title != "Post Two" {
			t.Errorf("unexpected post titles: %v, %v", posts[0].Title, posts[1].Title)
		}

		for _, post := range posts {
			if post.Author == nil {
				t.Fatalf("expected Author to be loaded, got nil")
			}
			if post.Author.Email != "author@example.com" {
				t.Errorf("expected author email 'author@example.com', got %s", post.Author.Email)
			}
		}

		bytes, _ := json.MarshalIndent(posts, "", "  ")
		fmt.Println(string(bytes))
	})
}

func TestCreateMany_MixedDefaults(t *testing.T) {
	ctx := context.Background()
	client, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("CreateMany mixed defaults grouping fallback works", func(t *testing.T) {
		count, err := client.User.CreateMany(

			client.User.Create().SetEmail("admin@example.com").SetPhoneNum("+100").SetRole(valk.UserRoleTypeAdmin),
			client.User.Create().SetEmail("student@example.com").SetPhoneNum("+200"),
			client.User.Create().SetEmail("teacher@example.com").SetPhoneNum("+300").SetRole(valk.UserRoleTypeTeacher),
			client.User.Create().SetEmail("student2@example.com").SetPhoneNum("+400"),
		).Exec(ctx)

		if err != nil {
			t.Fatalf("CreateMany failed: %v", err)
		}
		if count != 4 {
			t.Fatalf("expected count 4, got %d", count)
		}

		admin, err := client.User.FindUnique(user.Email.EQ("admin@example.com")).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to find admin: %v", err)
		}
		if admin.Role != valk.UserRoleTypeAdmin {
			t.Errorf("expected admin role ADMIN, got %v", admin.Role)
		}

		student, err := client.User.FindUnique(user.Email.EQ("student@example.com")).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to find student: %v", err)
		}
		if student.Role != valk.UserRoleTypeStudent {
			t.Errorf("expected student role STUDENT, got %v", student.Role)
		}

		teacher, err := client.User.FindUnique(user.Email.EQ("teacher@example.com")).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to find teacher: %v", err)
		}
		if teacher.Role != valk.UserRoleTypeTeacher {
			t.Errorf("expected teacher role TEACHER, got %v", teacher.Role)
		}

		student2, err := client.User.FindUnique(user.Email.EQ("student2@example.com")).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to find student2: %v", err)
		}
		if student2.Role != valk.UserRoleTypeStudent {
			t.Errorf("expected student2 role STUDENT, got %v", student2.Role)
		}
	})

	t.Run("CreateManyAndReturn mixed defaults grouping fallback works", func(t *testing.T) {
		users, err := client.User.CreateManyAndReturn(
			client.User.Create().SetEmail("admin-ret@example.com").SetPhoneNum("+100-ret").SetRole(valk.UserRoleTypeAdmin),
			client.User.Create().SetEmail("student-ret@example.com").SetPhoneNum("+200-ret"),
			client.User.Create().SetEmail("teacher-ret@example.com").SetPhoneNum("+300-ret").SetRole(valk.UserRoleTypeTeacher),
			client.User.Create().SetEmail("student2-ret@example.com").SetPhoneNum("+400-ret"),
		).Exec(ctx)

		if err != nil {
			t.Fatalf("CreateManyAndReturn failed: %v", err)
		}
		if len(users) != 4 {
			t.Fatalf("expected 4 users, got %d", len(users))
		}

		for _, u := range users {
			switch u.Email {
			case "admin-ret@example.com":
				if u.Role != valk.UserRoleTypeAdmin {
					t.Errorf("returned admin expected ADMIN, got %v", u.Role)
				}
			case "student-ret@example.com":
				if u.Role != valk.UserRoleTypeStudent {
					t.Errorf("returned student expected STUDENT, got %v", u.Role)
				}
			case "teacher-ret@example.com":
				if u.Role != valk.UserRoleTypeTeacher {
					t.Errorf("returned teacher expected TEACHER, got %v", u.Role)
				}
			case "student2-ret@example.com":
				if u.Role != valk.UserRoleTypeStudent {
					t.Errorf("returned student2 expected STUDENT, got %v", u.Role)
				}
			}
		}
	})
}

func TestCreateMany_SkipDuplicates(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateMany with SkipDuplicates skips existing records", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := client.User.Create().SetEmail("dup@example.com").SetPhoneNum("+100").Exec(ctx)
		if err != nil {
			t.Fatalf("failed to create initial user: %v", err)
		}

		count, err := client.User.CreateMany(
			client.User.Create().SetEmail("dup@example.com").SetPhoneNum("+101"), // duplicate email
			client.User.Create().SetEmail("unique@example.com").SetPhoneNum("+102"),
		).SkipDuplicates().Exec(ctx)

		if err != nil {
			t.Fatalf("CreateMany with SkipDuplicates failed: %v", err)
		}

		if count != 1 {
			t.Errorf("expected count 1, got %d", count)
		}

		var dbCount int
		err = client.Raw().QueryRowContext(ctx, `SELECT count(*) FROM "User"`).Scan(&dbCount)
		if err != nil {
			t.Fatalf("Raw SQL query failed: %v", err)
		}
		if dbCount != 2 {
			t.Errorf("expected 2 users in db, got %d", dbCount)
		}

		dupUser, err := client.User.FindUnique(user.Email.EQ("dup@example.com")).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to query dupUser: %v", err)
		}
		if dupUser.PhoneNum != "+100" {
			t.Errorf("expected original phone number '+100', got %q", dupUser.PhoneNum)
		}
	})

	t.Run("CreateManyAndReturn with SkipDuplicates only returns newly inserted records", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := client.User.Create().SetEmail("dup-ret@example.com").SetPhoneNum("+200").Exec(ctx)
		if err != nil {
			t.Fatalf("failed to create initial user: %v", err)
		}

		users, err := client.User.CreateManyAndReturn(
			client.User.Create().SetEmail("dup-ret@example.com").SetPhoneNum("+201"),
			client.User.Create().SetEmail("unique-ret@example.com").SetPhoneNum("+202"),
		).SkipDuplicates().Exec(ctx)

		if err != nil {
			t.Fatalf("CreateManyAndReturn with SkipDuplicates failed: %v", err)
		}

		if len(users) != 1 {
			t.Fatalf("expected 1 returned user, got %d", len(users))
		}
		if users[0].Email != "unique-ret@example.com" {
			t.Errorf("expected returned user email 'unique-ret@example.com', got %q", users[0].Email)
		}

		var dbCount int
		err = client.Raw().QueryRowContext(ctx, `SELECT count(*) FROM "User"`).Scan(&dbCount)
		if err != nil {
			t.Fatalf("Raw SQL query failed: %v", err)
		}
		if dbCount != 2 {
			t.Errorf("expected 2 users in db, got %d", dbCount)
		}
	})
}
