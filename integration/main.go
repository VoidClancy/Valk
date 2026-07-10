package main

import (
	"context"
	"encoding/json"
	"fmt"
	"integration/valk"
	"integration/valk/category"
	"integration/valk/categoryToPost"
	"integration/valk/comment"
	"integration/valk/post"
	"integration/valk/profile"
	"integration/valk/user"
	"time"

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

	db.User.BeforeCreate(func(ctx context.Context, user *valk.UserCreate) error {
		user.Role = new(valk.UserRole.Admin)
		return nil
	})

	db.User.AfterCreateMany(func(ctx context.Context, usersInput []valk.UserCreate, count int64) error {
		var emails []string
		for _, u := range usersInput {
			emails = append(emails, u.Email)
		}
		fmt.Printf("AfterCreateMany: %d users created (emails: %v)\n", count, emails)
		return nil
	})
	db.User.AfterCreate(func(ctx context.Context, users []*valk.User) error {
		for _, u := range users {
			fmt.Printf("AfterCreate: user %s (role=%s)\n", u.Email, u.Role)
		}
		return nil
	})

	var usersToCreate []valk.RecordInput

	for i := range 20 {
		usersToCreate = append(usersToCreate, user.Record(
			user.Email.Set(fmt.Sprintf("email-%d", i)),
			user.PhoneNum.Set(fmt.Sprintf("555-%d", i)),
			user.Password.Set(fmt.Sprintf("password-%d", i)),
		))
	}

	users, err := db.User.CreateManyAndReturn(usersToCreate...).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create users: %v", err)
	}
	fmt.Printf("CreateManyAndReturn: %d users returned with auto-generated IDs\n", len(users))

	if _, err := db.User.CreateMany(
		user.Record(
			user.Email.Set("test"),
			user.PhoneNum.Set("555-test"),
			user.Password.Set("passwd"),
		),
		user.Record(
			user.Email.Set("again"),
			user.PhoneNum.Set("555-again"),
			user.Password.Set("123456"),
		),
	).Exec(ctx); err != nil {
		log.Fatalf("failed to CreateMany: %v", err)
	}
	referrer, err := db.User.Create(
		user.Email.Set("referrer@example.com"),
		user.PhoneNum.Set("555-0001"),
		user.Password.Set("pass123"),
		user.Role.Set(valk.UserRole.Student),
	).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create referrer: %v", err)
	}

	referred, err := db.User.Create(
		user.Email.Set("referred@example.com"),
		user.PhoneNum.Set("555-0002"),
		user.Password.Set("pass456"),
		user.Role.Set(valk.UserRole.Student),
		user.ReferredById.Set(referrer.Id),
	).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create referred: %v", err)
	}

	prof, err := db.Profile.Create(
		profile.Bio.Set("BLEH"),
		profile.UserId.Set(referred.Id),
		profile.CreatedAt.Set(time.Now()),
	).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create profile: %v", err)
	}
	fmt.Println("PROFILE:")
	printJSON(prof)

	categoryTest, err := db.Category.Create(
		category.Name.Set("TEST"),
	).Exec(ctx)

	if err != nil {
		log.Fatalf("failed to create category: %v", err)
	}
	fmt.Println("CATEGORY:")
	printJSON(categoryTest)

	AnothercategoryTest, err := db.Category.Create(
		category.Id.Set(24),
		category.Name.Set("TEST2"),
	).Exec(ctx)

	if err != nil {
		log.Fatalf("failed to create category: %v", err)
	}
	fmt.Println("CATEGORY2:")
	printJSON(AnothercategoryTest)

	p, err := db.Post.Create(
		post.Title.Set("Valkyrie ORM Deep Dive"),
		post.Content.Set("skrrrt"),
		post.AuthorId.Set(referred.Id),
	).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create post: %v", err)
	}

	cat, err := db.Category.Create(
		category.Name.Set("Programming"),
	).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create category: %v", err)
	}

	_, err = db.CategoryToPost.Create(
		categoryToPost.PostId.Set(p.Id),
		categoryToPost.CategoryId.Set(cat.Id),
	).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create CategoryToPost: %v", err)
	}

	meta1 := json.RawMessage(`{"rating":5,"verified":true}`)
	_, err = db.Comment.Create(
		comment.Textify.Set(100),
		comment.Dummy3.Set("dummy_val_1"),
		comment.Dummy1.Set(42),
		comment.Dummy2.Set("dummy_val_2"),
		comment.PostId.Set(p.Id),
		comment.AuthorId.Set(referrer.Id),
		comment.Meta.Set(meta1),
	).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create comment 1: %v", err)
	}

	meta2 := json.RawMessage(`{"rating":4,"verified":false}`)
	_, err = db.Comment.Create(
		comment.Textify.Set(200),
		comment.Dummy3.Set("dummy_val_3"),
		comment.Dummy1.Set(84),
		comment.Dummy2.Set("dummy_val_4"),
		comment.PostId.Set(p.Id),
		comment.AuthorId.Set(referred.Id),
		comment.Meta.Set(meta2),
	).Exec(ctx)
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
	db, err := valk.Open("sqlite", "file::memory:?_pragma=foreign_keys(1)&_time_format=sqlite")
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
	author, err := tx.User.Create(
		user.Email.Set("clancySizer@gmail.com"),
		user.PhoneNum.Set("+1234567890"),
	).Exec(ctx)
	if err != nil {
		fmt.Printf("failed to create user: %+v", err)
		return
	}

	postWithAuthor, err := tx.Post.Create(
		post.Title.Set("A Post"),
		post.AuthorId.Set(author.Id),
	).Select(post.Select{
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

		author, err := tx.User.Create(
			user.Email.Set("clancySizer@gmail.com"),
			user.PhoneNum.Set("+1234567890"),
		).Exec(ctx)
		if err != nil {
			return err
		}

		postWithAuthor, err := tx.Post.Create(
			post.Title.Set("A Post"),
			post.AuthorId.Set(author.Id),
		).Select(post.Select{
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
