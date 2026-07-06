package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"integration/valk"

	"log"

	_ "modernc.org/sqlite"
)

func hashPassword(pass string) string {
	h := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(h[:])

}
func main() {

	db := openConn()
	defer db.Close()

	rawDB := db.Raw()
	rawDB.SetMaxOpenConns(10)

	ctx := context.Background()

	runMigrations(db, ctx)

	db.User.BeforeCreate(func(ctx context.Context, i *valk.UserCreate) error {
		if i.Password != nil {
			hash := hashPassword(*i.Password)
			i.Password = &hash
			return nil
		}
		return nil
	})

	db.User.AfterCreate(func(ctx context.Context, u *valk.User) error {
		fmt.Printf("HASHED PASS: %s \n", *u.Password)
		return nil
	})
	author, err := db.User.Create(valk.UserCreate{
		Email:    "clancySizer@gmail.com",
		PhoneNum: "+1234567890",
		Password: new("veryStrongPassword"),
		Role:     &valk.UserRole.Admin,
	}).Exec(ctx)

	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	post, err := db.Post.Create(valk.PostCreate{

		Content:  new("eheheh"),
		Title:    "some post",
		AuthorId: author.Id,
	}).Exec(ctx)

	if err != nil {
		log.Fatalf("failed to create user: %v", err)

	}

	comment, err := db.Comment.Create(valk.CommentCreate{
		Textify:  42,
		Dummy3:   "d3",
		Dummy1:   1,
		Dummy2:   "d2",
		PostId:   post.Id,
		AuthorId: author.Id,
	}).Select(valk.CommentSelect{
		Id:      true,
		Textify: true,
		Dummy3:  true,
		Dummy1:  true,
		Dummy2:  true,
		PostId:  true,
		Author:  &valk.UserSelect{},
		Post:    &valk.PostSelect{},
	}).Exec(ctx)

	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	usersCount, err := db.User.CreateMany([]valk.UserCreate{
		{Email: "cl@gm.com"}, {Email: "cc@gg.com"},
	}).Exec(ctx)
	fmt.Printf("\nCREATED %d USERS\n", usersCount)
	fmt.Println("COMMENT:")
	printJSON(comment)

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
	author, err := tx.User.Create(valk.UserCreate{
		Email:    "clancySizer@gmail.com",
		PhoneNum: "+1234567890",
	}).Exec(ctx)
	if err != nil {
		fmt.Printf("failed to create user: %+v", err)
		return

	}

	postWithAuthor, err := tx.Post.Create(valk.PostCreate{
		Title:    "A Post",
		AuthorId: author.Id,
	}).Select(valk.PostSelect{
		Id:    true,
		Title: true,
		Author: &valk.UserSelect{

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

		author, err := tx.User.Create(valk.UserCreate{
			Email:    "clancySizer@gmail.com",
			PhoneNum: "+1234567890",
		}).Exec(ctx)
		if err != nil {
			return err
		}

		postWithAuthor, err := tx.Post.Create(valk.PostCreate{
			Title:    "A Post",
			AuthorId: author.Id,
		}).Select(valk.PostSelect{
			Id:    true,
			Title: true,
			Author: &valk.UserSelect{

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
