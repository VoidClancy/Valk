package main

import (
	"context"
	"encoding/json"
	"fmt"
	"integration/valkyrie"

	"log"

	_ "modernc.org/sqlite"
)

func main() {

	db := openConn()
	defer db.Close()

	rawDB := db.Raw()
	rawDB.SetMaxOpenConns(10)

	ctx := context.Background()

	runMigrations(db, ctx)

	author, err := db.User.Create(valkyrie.UserCreateInput{
		Email:    "clancySizer@gmail.com",
		PhoneNum: "+1234567890",
		Role:     &valkyrie.UserRole.Admin,
	}).Exec(ctx)

	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	post, err := db.Post.Create(valkyrie.PostCreateInput{

		Content:  new("eheheh"),
		Title:    "some post",
		AuthorId: author.Id,
	}).Exec(ctx)

	if err != nil {
		log.Fatalf("failed to create user: %v", err)

	}

	comment, err := db.Comment.Create(valkyrie.CommentCreateInput{
		Textify:  42,
		Dummy3:   "d3",
		Dummy1:   1,
		Dummy2:   "d2",
		PostId:   post.Id,
		AuthorId: author.Id,
	}).Select(valkyrie.CommentSelect{
		Id:      true,
		Textify: true,
		Dummy3:  true,
		Dummy1:  true,
		Dummy2:  true,
		PostId:  true,
		Author:  &valkyrie.UserSelect{},
		Post:    &valkyrie.PostSelect{},
	}).Exec(ctx)

	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	printJSON(comment)

}
func openConn() *valkyrie.DB {
	db, err := valkyrie.Open("sqlite", "file::memory:?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	return db
}

func runMigrations(db *valkyrie.DB, ctx context.Context) {
	if err := db.RunMigrations(ctx); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}
func runManualTransaction(db *valkyrie.DB, ctx context.Context) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("Manual Transaction: failed to begin: %v", err)
		return
	}
	defer tx.Rollback()

	fmt.Println("Manual Transaction: started successfully")
	author, err := tx.User.Create(valkyrie.UserCreateInput{
		Email:    "clancySizer@gmail.com",
		PhoneNum: "+1234567890",
	}).Exec(ctx)
	if err != nil {
		fmt.Printf("failed to create user: %+v", err)
		return

	}

	postWithAuthor, err := tx.Post.Create(valkyrie.PostCreateInput{
		Title:    "A Post",
		AuthorId: author.Id,
	}).Select(valkyrie.PostSelect{
		Id:    true,
		Title: true,
		Author: &valkyrie.UserSelect{

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

func runBlockBasedTransaction(db *valkyrie.DB, ctx context.Context) {
	err := db.Transaction(ctx, func(tx *valkyrie.Tx) error {
		fmt.Println("Block-based Transaction: started successfully")

		author, err := tx.User.Create(valkyrie.UserCreateInput{
			Email:    "clancySizer@gmail.com",
			PhoneNum: "+1234567890",
		}).Exec(ctx)
		if err != nil {
			return err
		}

		postWithAuthor, err := tx.Post.Create(valkyrie.PostCreateInput{
			Title:    "A Post",
			AuthorId: author.Id,
		}).Select(valkyrie.PostSelect{
			Id:    true,
			Title: true,
			Author: &valkyrie.UserSelect{

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
