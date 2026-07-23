package main

import (
	"context"
	"encoding/json"
	"fmt"
	"integration/valk"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRelationLoadChildHoldsFK(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u, err := db.User.Create().
		SetEmail("parent@example.com").
		SetPhoneNum("+111111111").
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	_, err = db.Post.Create().
		SetTitle("Post 1").
		SetContent("Content 1").
		SetAuthorId(u.Id).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create post 1: %v", err)
	}
	_, err = db.Post.Create().
		SetTitle("Post 2").
		SetContent("Content 2").
		SetAuthorId(u.Id).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create post 2: %v", err)
	}

	_, err = db.Profile.Create().
		SetBio("My bio").
		SetUserId(u.Id).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create profile: %v", err)
	}

	u2, err := db.User.Create().
		SetEmail("parent2@example.com").
		SetPhoneNum("+222222222").
		Select(valk.UserSelect{
			Id:    true,
			Email: true,
			Posts: &valk.PostSelect{
				Id:    true,
				Title: true,
			},
			Profile: &valk.ProfileSelect{
				Id:  true,
				Bio: true,
			},
		}).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user2: %v", err)
	}

	b, _ := json.MarshalIndent(u2, "", "  ")
	fmt.Println(string(b))

	if len(u2.Posts) != 0 {
		t.Errorf("expected 0 posts for new user, got %d", len(u2.Posts))
	}
	if u2.Profile != nil {
		t.Errorf("expected nil profile for new user, got %+v", u2.Profile)
	}
}

func TestRelationLoadParentHoldsFK(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u, err := db.User.Create().
		SetEmail("author@example.com").
		SetPhoneNum("+222222222").
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	p, err := db.Post.Create().
		SetTitle("My Post").
		SetAuthorId(u.Id).
		Select(valk.PostSelect{
			Id:       true,
			Title:    true,
			AuthorId: true,
			Author: &valk.UserSelect{
				Id:    true,
				Email: true,
			},
		}).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	b, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(b))

	if p.Title != "My Post" {
		t.Errorf("expected title 'My Post', got '%s'", p.Title)
	}
	if p.Author == nil {
		t.Fatalf("expected author to be loaded, got nil")
	}
	if p.Author.Id != u.Id {
		t.Errorf("expected author id '%s', got '%s'", u.Id, p.Author.Id)
	}
	if p.Author.Email != "author@example.com" {
		t.Errorf("expected author email 'author@example.com', got '%s'", p.Author.Email)
	}
}

func TestRelationLoadSelfRelation(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	referrer, err := db.User.Create().
		SetEmail("referrer@example.com").
		SetPhoneNum("+333333333").
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create referrer: %v", err)
	}

	referred, err := db.User.Create().
		SetEmail("referred@example.com").
		SetPhoneNum("+444444444").
		SetReferredById(referrer.Id).
		Select(valk.UserSelect{
			Id:           true,
			Email:        true,
			ReferredById: true,
			ReferredBy: &valk.UserSelect{
				Id:    true,
				Email: true,
			},
		}).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create referred user: %v", err)
	}

	b, _ := json.MarshalIndent(referred, "", "  ")
	fmt.Println(string(b))

	if referred.ReferredBy == nil {
		t.Fatalf("expected referredBy to be loaded, got nil")
	}
	if referred.ReferredBy.Id != referrer.Id {
		t.Errorf("expected referredBy id '%s', got '%s'", referrer.Id, referred.ReferredBy.Id)
	}
	if referred.ReferredBy.Email != "referrer@example.com" {
		t.Errorf("expected referredBy email 'referrer@example.com', got '%s'", referred.ReferredBy.Email)
	}
}

func TestRelationLoadDeepNesting(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u, err := db.User.Create().
		SetEmail("deep@example.com").
		SetPhoneNum("+555555555").
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	p, err := db.Post.Create().
		SetTitle("Deep Post").
		SetAuthorId(u.Id).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	_, err = db.Comment.Create().
		SetTextify(42).
		SetDummy3("d3").
		SetDummy1(1).
		SetDummy2("d2").
		SetPostId(p.Id).
		SetAuthorId(u.Id).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}

	p2, err := db.Post.Create().
		SetTitle("Another Post").
		SetAuthorId(u.Id).
		Select(valk.PostSelect{
			Id:    true,
			Title: true,
			Author: &valk.UserSelect{
				Id:    true,
				Email: true,
				Posts: &valk.PostSelect{
					Id:    true,
					Title: true,
					Comments: &valk.CommentSelect{
						Id:      true,
						Textify: true,
					},
				},
			},
		}).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create post2: %v", err)
	}

	b, _ := json.MarshalIndent(p2, "", "  ")
	fmt.Println(string(b))

	if p2.Author == nil {
		t.Fatalf("expected author to be loaded")
	}
	if p2.Author.Email != "deep@example.com" {
		t.Errorf("expected author email 'deep@example.com', got '%s'", p2.Author.Email)
	}
	if len(p2.Author.Posts) != 2 {
		t.Fatalf("expected 2 posts on author, got %d", len(p2.Author.Posts))
	}
	var deepPost *valk.Post
	for _, post := range p2.Author.Posts {
		if post.Title == "Deep Post" {
			deepPost = post
			break
		}
	}
	if deepPost == nil {
		t.Fatalf("expected to find 'Deep Post' in author's posts")
	}
	if len(deepPost.Comments) != 1 {
		t.Fatalf("expected 1 comment on Deep Post, got %d", len(deepPost.Comments))
	}
	if deepPost.Comments[0].Textify != 42 {
		t.Errorf("expected comment textify 42, got %d", deepPost.Comments[0].Textify)
	}
}
