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

	db, err := valkyrie.Open("sqlite", "file::memory:?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	rawDB := db.Raw()
	rawDB.SetMaxOpenConns(10)

	ctx := context.Background()

	if err := db.RunMigrations(ctx); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	fmt.Println("Valkyrie database client successfully opened")

	author, err := db.User.Create(valkyrie.UserCreateInput{
		Email:    "clancySizer@gmail.com",
		PhoneNum: "+1234567890",
		Role:     &valkyrie.UserRole.Admin,
	}).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	fmt.Printf("CREATED AUTHOR ID: %q\n", author.Id)

	postWithAuthor, err := db.Post.Create(valkyrie.PostCreateInput{
		AuthorId: author.Id,
	}).Select(valkyrie.PostSelect{
		Id:    true,
		Title: true,
		Author: &valkyrie.UserSelect{
			Id:    false,
			Email: true,
		},
	}).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create post: %v", err)
	}

	b, _ := json.MarshalIndent(postWithAuthor, "", "  ")
	fmt.Println(string(b))

	err = db.Transaction(ctx, func(tx *valkyrie.Tx) error {
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

	// Manual Transaction
	err = manualTx(db, ctx)
	if err != nil {
		fmt.Printf("Manual Transaction failed: %v", err)
	}
}

func manualTx(db *valkyrie.DB, ctx context.Context) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("Manual Transaction: failed to begin: %v", err)
	}
	defer tx.Rollback()

	fmt.Println("Manual Transaction: started successfully")
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

	if err := tx.Commit(); err != nil {
		log.Fatalf("Manual Transaction: commit failed: %v", err)
	}
	fmt.Println("Manual Transaction: committed successfully")

	return nil
}
