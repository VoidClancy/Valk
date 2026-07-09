package main

import (
	"context"
	"encoding/json"
	"fmt"
	"integration/valk"
	"integration/valk/post"
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
			user.Record(user.Email.Set("hooked@example.com"), user.PhoneNum.Set("+100000000")),
			user.Record(user.Email.Set("normal@example.com"), user.PhoneNum.Set("+200000000")),
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
			user.Record(user.Email.Set("hooked@example.com"), user.PhoneNum.Set("+100000001")),
			user.Record(user.Email.Set("normal@example.com"), user.PhoneNum.Set("+200000001")),
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
			user.Record(user.Email.Set("good@example.com"), user.PhoneNum.Set("+300000000")),
			user.Record(user.Email.Set("reject@example.com"), user.PhoneNum.Set("+300000001")),
		).Exec(ctx)

		if err == nil {
			t.Fatal("expected error from hook rejection, got nil")
		}

		var count int
		client.Raw().QueryRowContext(ctx, `SELECT count(*) FROM "User"`).Scan(&count)
		if count != 0 {
			t.Fatalf("expected 0 rows after aborted CreateMany, got %d", count)
		}
	})

	t.Run("AfterCreate hook fires during CreateManyAndReturn", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		var afterCalled bool
		var gotRole valk.UserRoleType
		client.User.AfterCreate(func(ctx context.Context, user *valk.User) error {
			afterCalled = true
			gotRole = user.Role
			return nil
		})

		users, err := client.User.CreateManyAndReturn(
			user.Record(user.Email.Set("after@example.com"), user.PhoneNum.Set("+500000000")),
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

		client.User.AfterCreate(func(ctx context.Context, user *valk.User) error {
			return fmt.Errorf("after hook rejected")
		})

		var count int
		client.Raw().QueryRowContext(ctx, `SELECT count(*) FROM "User"`).Scan(&count)
		prevCount := count

		_, err := client.User.CreateManyAndReturn(
			user.Record(user.Email.Set("aftershoot@example.com"), user.PhoneNum.Set("+600000000")),
		).Exec(ctx)

		if err == nil {
			t.Fatal("expected error from AfterCreate rejection, got nil")
		}
		if !strings.Contains(err.Error(), "after hook rejected") {
			t.Errorf("expected 'after hook rejected' in error, got %v", err)
		}

		client.Raw().QueryRowContext(ctx, `SELECT count(*) FROM "User"`).Scan(&count)
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
			user.Record(
				user.Email.Set("bulk1@example.com"),
				user.PhoneNum.Set("+700000001"),
			),
			user.Record(
				user.Email.Set("bulk2@example.com"),
				user.PhoneNum.Set("+700000002"),
			),
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
			user.Record(
				user.Email.Set("ghost@example.com"),
				user.PhoneNum.Set("+800000000"),
			),
		).Exec(ctx)

		if err == nil {
			t.Fatal("expected error from AfterCreateMany rejection, got nil")
		}
		if !strings.Contains(err.Error(), "after create many rejected") {
			t.Errorf("expected 'after create many rejected' in error, got %v", err)
		}

		var count int
		client.Raw().QueryRowContext(ctx, query(
			`SELECT count(*) FROM "User" WHERE email = 'ghost@example.com'`,
			`SELECT count(*) FROM "User" WHERE email = 'ghost@example.com'`,
		)).Scan(&count)
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
			user.Record(user.Email.Set("bulk1@example.com"), user.PhoneNum.Set("+111")),
			user.Record(user.Email.Set("bulk2@example.com"), user.PhoneNum.Set("+222")),
			user.Record(user.Email.Set("bulk3@example.com"), user.PhoneNum.Set("+333")),
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
			post.Record(post.Title.Set("Post One"), post.AuthorId.Set(author.Id)),
			post.Record(post.Title.Set("Post Two"), post.AuthorId.Set(author.Id)),
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
