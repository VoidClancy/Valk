package main

import (
	"context"
	"encoding/json"
	"fmt"
	"integration/valkyrie"
	"testing"
)

func TestCreateMany(t *testing.T) {
	ctx := context.Background()
	client, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("CreateMany returns correct count", func(t *testing.T) {
		count, err := client.User.CreateMany([]valkyrie.UserCreateInput{
			{
				Email:    "bulk1@example.com",
				PhoneNum: "+111",
			},
			{
				Email:    "bulk2@example.com",
				PhoneNum: "+222",
			},
			{
				Email:    "bulk3@example.com",
				PhoneNum: "+333",
			},
		}).Exec(ctx)

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
		author, err := client.User.Create(valkyrie.UserCreateInput{
			Email:    "author@example.com",
			PhoneNum: "+444",
		}).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to create author: %v", err)
		}

		posts, err := client.Post.CreateManyAndReturn([]valkyrie.PostCreateInput{
			{
				Title:    "Post One",
				AuthorId: author.Id,
			},
			{
				Title:    "Post Two",
				AuthorId: author.Id,
			},
		}).Select(valkyrie.PostSelect{
			Id:    true,
			Title: true,
			Author: &valkyrie.UserSelect{
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
