package main

import (
	"context"
	"encoding/json"
	"fmt"
	"integration/valk"
	"integration/valk/post"
	"integration/valk/user"
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

		var phone string
		err = client.Raw().QueryRowContext(ctx,
			`SELECT phoneNum FROM "User" WHERE email = 'hooked@example.com'`,
		).Scan(&phone)
		if err != nil {
			t.Fatalf("query failed: %v", err)
		}
		if phone != "+188888888" {
			t.Errorf("expected PhoneNum to be mutated to '+188888888', got %q", phone)
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
