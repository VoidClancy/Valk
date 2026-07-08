package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"integration/valk"
	"integration/valk/category"
	"integration/valk/categoryToPost"
	"integration/valk/comment"
	"integration/valk/post"
	"integration/valk/profile"
	"integration/valk/user"

	"log"

	_ "modernc.org/sqlite"
)

type SeedData struct {
	ReferrerId string
	ReferredId string
	PostId     string
	Meta1      json.RawMessage
	Meta2      json.RawMessage
}

func seed(db *valk.DB, ctx context.Context) *SeedData {
	db.User.BeforeCreate(func(ctx context.Context, uc *user.Create) error {
		return errors.New("AAAAAAAAH")
	})

	referrer, err := db.User.Create(user.Create{
		Email:    "referrer@example.com",
		PhoneNum: "555-0001",
		Password: new("pass123"),
		Role:     &valk.UserRole.Admin,
	}).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create referrer: %v", err)
	}

	referred, err := db.User.Create(user.Create{
		Email:        "referred@example.com",
		PhoneNum:     "555-0002",
		Password:     new("pass456"),
		Role:         &valk.UserRole.Student,
		ReferredById: &referrer.Id,
	}).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create referred: %v", err)
	}

	prof, err := db.Profile.Create(profile.Create{
		Bio:    new("BLEH"),
		UserId: referred.Id,
	}).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create profile: %v", err)
	}
	_ = prof

	p, err := db.Post.Create(post.Create{
		Title:    "Valkyrie ORM Deep Dive",
		Content:  new("skrrrt"),
		AuthorId: referred.Id,
	}).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create post: %v", err)
	}

	cat, err := db.Category.Create(category.Create{
		Name: "Programming",
	}).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create category: %v", err)
	}

	_, err = db.CategoryToPost.Create(categoryToPost.Create{
		PostId:     p.Id,
		CategoryId: cat.Id,
	}).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create CategoryToPost: %v", err)
	}

	meta1 := json.RawMessage(`{"rating":5,"verified":true}`)
	_, err = db.Comment.Create(comment.Create{
		Textify:  100,
		Dummy3:   "dummy_val_1",
		Dummy1:   42,
		Dummy2:   "dummy_val_2",
		PostId:   p.Id,
		AuthorId: referrer.Id,
		Meta:     &meta1,
	}).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create comment 1: %v", err)
	}

	meta2 := json.RawMessage(`{"rating":4,"verified":false}`)
	_, err = db.Comment.Create(comment.Create{
		Textify:  200,
		Dummy3:   "dummy_val_3",
		Dummy1:   84,
		Dummy2:   "dummy_val_4",
		PostId:   p.Id,
		AuthorId: referred.Id,
		Meta:     &meta2,
	}).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create comment 2: %v", err)
	}

	return &SeedData{
		ReferrerId: referrer.Id,
		ReferredId: referred.Id,
		PostId:     p.Id,
		Meta1:      meta1,
		Meta2:      meta2,
	}
}

func main() {
	db := openConn()
	defer db.Close()

	rawDB := db.Raw()
	rawDB.SetMaxOpenConns(10)
	ctx := context.Background()

	runMigrations(db, ctx)

	fmt.Println("=== Seeding Data ===")
	data := seed(db, ctx)
	fmt.Println("Seeding complete.")
	fmt.Println()

	all, err := db.User.FindMany().Exec(ctx)
	if err != nil {
		log.Fatalf("failed to get all users: %v", err)
	}
	fmt.Println("=== ALL ===")

	printJSON(all)

	bleh, err := db.User.FindMany(user.And(
		user.Email.Contains("@example"),
		user.Or(
			user.Email.EQ("referred@example.com"),
			user.PhoneNum.EQ("+1111"),
		),
	)).Select(user.Select{
		Id: true,
	}).Exec(ctx)
	printJSON(bleh)

	fmt.Println("=== QUERY 1: Deep Nested Select ===")
	resUser, err := db.User.FindFirst(
		user.Email.EQ("referred@example.com"),
	).Select(user.Select{
		Email: true,
		Profile: &profile.Select{
			Bio: true,
		},
		ReferredBy: &user.Select{
			Email:    true,
			PhoneNum: true,
		},
		Posts: &post.Select{
			Title: true,
			Comments: &comment.Select{
				Textify: true,
				Meta:    true,
				Author: &user.Select{
					Email: true,
				},
			},
		},
	}).Exec(ctx)

	if err != nil {
		log.Fatalf("Query 1 failed: %v", err)
	}
	printJSON(resUser)
	fmt.Println()

	fmt.Println("=== QUERY 2: Omit Nested Fields ===")
	resPost, err := db.Post.FindFirst(
		post.Title.Like("%Valkyrie%"),
	).
		Select(post.Select{
			Title:     true,
			Published: true,
			Comments: &comment.Select{
				Textify: true,
				Meta:    true,
			},
		}).
		Exec(ctx)

	if err != nil {
		log.Fatalf("Query 2 failed: %v", err)
	}
	printJSON(resPost)
	fmt.Println()

	fmt.Println("=== QUERY 3: Filtering with Relations ===")
	resComments, err := db.Comment.FindMany(
		comment.Meta.EQ(data.Meta1),
	).Select(comment.Select{
		Textify: true,
		Meta:    true,
		Post: &post.Select{
			Title: true,
			Author: &user.Select{
				Email: true,
			},
		},
	}).Exec(ctx)

	if err != nil {
		log.Fatalf("Query 3 failed: %v", err)
	}
	printJSON(resComments)
	fmt.Println()
}

func openConn() *valk.DB {
	db, err := valk.Open("sqlite", "file::memory:?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	return db
}
func runMigrations(db *valk.DB, ctx context.Context) {
	if err := db.RunMigrations(ctx); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}

func runManualTransaction(db *valk.DB, ctx context.Context) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("Manual Transaction: failed to begin: %v", err)
		return
	}
	defer tx.Rollback()

	fmt.Println("Manual Transaction: started successfully")
	author, err := tx.User.Create(user.Create{
		Email:    "clancySizer@gmail.com",
		PhoneNum: "+1234567890",
	}).Exec(ctx)
	if err != nil {
		fmt.Printf("failed to create user: %+v", err)
		return
	}

	postWithAuthor, err := tx.Post.Create(post.Create{
		Title:    "A Post",
		AuthorId: author.Id,
	}).Select(post.Select{
		Id:    true,
		Title: true,
		Author: &user.Select{
			Email: true,
		},
	}).Exec(ctx)
	if err != nil {
		fmt.Printf("failed to create Post: %+v", err)
		return
	}

	b, _ := json.MarshalIndent(postWithAuthor, "", "  ")
	fmt.Println(string(b))

	if err := tx.Commit(); err != nil {
		log.Printf("Manual Transaction: commit failed: %v", err)
		return
	}
	fmt.Println("Manual Transaction: committed successfully")
}

func runBlockBasedTransaction(db *valk.DB, ctx context.Context) {
	err := db.Transaction(ctx, func(tx *valk.Tx) error {
		fmt.Println("Block-based Transaction: started successfully")

		author, err := tx.User.Create(user.Create{
			Email:    "clancySizer@gmail.com",
			PhoneNum: "+1234567890",
		}).Exec(ctx)
		if err != nil {
			return err
		}

		postWithAuthor, err := tx.Post.Create(post.Create{
			Title:    "A Post",
			AuthorId: author.Id,
		}).Select(post.Select{
			Id:    true,
			Title: true,
			Author: &user.Select{
				Email: true,
			},
		}).Exec(ctx)
		if err != nil {
			return err
		}

		b, _ := json.MarshalIndent(postWithAuthor, "", "  ")
		fmt.Println(string(b))
		return nil
	})
	if err != nil {
		fmt.Printf("Block-based Transaction failed: %v", err)
	}
	fmt.Println("Block-based Transaction: committed successfully")
}

func printJSON(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
