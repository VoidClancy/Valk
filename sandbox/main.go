package main

import (
	"context"
	"fmt"
	"log"
	"sandbox/valkyrie"

	_ "modernc.org/sqlite"
)

func main() {

	db, err := valkyrie.Open("sqlite", "file::memory:?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	// access to the underlying *sql.DB
	rawDB := db.Raw()
	rawDB.SetMaxOpenConns(10)

	ctx := context.Background()

	if err := db.RunMigrations(ctx); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	fmt.Println("Valkyrie database client successfully opened!")
	fmt.Println("Delegates registered:")
	fmt.Printf("  UserDelegate:           %+v\n", *db.User)
	fmt.Printf("  ProfileDelegate:        %+v\n", *db.Profile)
	fmt.Printf("  PostDelegate:           %+v\n", *db.Post)
	fmt.Printf("  CommentDelegate:        %+v\n", *db.Comment)
	fmt.Printf("  CategoryDelegate:       %+v\n", *db.Category)
	fmt.Printf("  CategoryToPostDelegate: %+v\n", *db.CategoryToPost)

	fmt.Printf("  UserRole.Admin:         %v\n", db.UserRole.Admin)

	err = db.Transaction(ctx, func(tx *valkyrie.Tx) error {
		fmt.Println("Block-based Transaction: started successfully!")

		// tx.User.Create(...)

		//rawTx := tx.Raw()

		return nil
	})
	if err != nil {
		log.Fatalf("Block-based Transaction failed: %v", err)
	}
	fmt.Println("Block-based Transaction: committed successfully!")

	// Manual Transaction
	manualTx(db, ctx)

}

func manualTx(db *valkyrie.DB, ctx context.Context) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("Manual Transaction: failed to begin: %v", err)
	}
	defer tx.Rollback()

	fmt.Println("Manual Transaction: started successfully!")
	// tx.User.Create(...)

	if err := tx.Commit(); err != nil {
		log.Fatalf("Manual Transaction: commit failed: %v", err)
	}
	fmt.Println("Manual Transaction: committed successfully!")
}
