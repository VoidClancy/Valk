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
}
