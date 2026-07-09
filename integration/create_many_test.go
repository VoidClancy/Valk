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
