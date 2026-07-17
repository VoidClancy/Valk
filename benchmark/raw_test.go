package main

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
)

func benchRawCreate(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, err := db.ExecContext(ctx,
			`INSERT INTO "User" ("id", "email", "phoneNum", "role") VALUES (?, ?, ?, ?)`,
			fmt.Sprintf("raw-create-%d", i),
			fmt.Sprintf("raw-create-%d@example.com", i),
			fmt.Sprintf("raw-phone-%d", i),
			"student",
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchRawCreateMany(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		sqlStr := `INSERT INTO "User" ("id", "email", "phoneNum", "role", "loginCount") VALUES `
		args := make([]any, 0, 50)
		for j := 0; j < 10; j++ {
			if j > 0 {
				sqlStr += ", "
			}
			n := i*10 + j
			sqlStr += "(?, ?, ?, ?, ?)"
			args = append(args,
				fmt.Sprintf("raw-cmany-%d", n),
				fmt.Sprintf("raw-cmany-%d@example.com", n),
				fmt.Sprintf("raw-cmany-phone-%d", n),
				"student",
				0,
			)
		}
		_, err := db.ExecContext(ctx, sqlStr, args...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchRawCreateManyAndReturn(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		sqlStr := `INSERT INTO "User" ("id", "email", "phoneNum", "role", "loginCount") VALUES `
		args := make([]any, 0, 50)
		for j := 0; j < 10; j++ {
			if j > 0 {
				sqlStr += ", "
			}
			n := i*10 + j
			sqlStr += "(?, ?, ?, ?, ?)"
			args = append(args,
				fmt.Sprintf("raw-cmar-%d", n),
				fmt.Sprintf("raw-cmar-%d@example.com", n),
				fmt.Sprintf("raw-cmar-phone-%d", n),
				"student",
				0,
			)
		}
		sqlStr += ` RETURNING "id", "email", "phoneNum", "role", "loginCount"`
		rows, err := db.QueryContext(ctx, sqlStr, args...)
		if err != nil {
			b.Fatal(err)
		}
		if rows.Err() != nil {
			b.Fatal(rows.Err())
		}
		for rows.Next() {
			var id, email, phoneNum, role string
			var loginCount int32
			if err := rows.Scan(&id, &email, &phoneNum, &role, &loginCount); err != nil {
				b.Fatal(err)
			}
		}
		rows.Close()
	}
}

func benchRawFindUnique(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)
	seedData(db, "raw-fu")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		email := fmt.Sprintf("raw-fu-user-%d@example.com", i%seedCount)
		var id, phoneNum, role string
		err := db.QueryRowContext(ctx,
			`SELECT "id", "phoneNum", "role" FROM "User" WHERE "email" = ?`, email,
		).Scan(&id, &phoneNum, &role)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchRawFindFirst(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)
	seedData(db, "raw-ff")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		email := fmt.Sprintf("raw-ff-user-%d@example.com", i%seedCount)
		var id, phoneNum, role string
		err := db.QueryRowContext(ctx,
			`SELECT "id", "phoneNum", "role" FROM "User" WHERE "email" = ? LIMIT 1`, email,
		).Scan(&id, &phoneNum, &role)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			b.Fatal(err)
		}
	}
}

func benchRawFindMany(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)
	seedData(db, "raw-fm")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		offset := i % seedCount
		rows, err := db.QueryContext(ctx,
			`SELECT "id", "email", "phoneNum", "role" FROM "User" ORDER BY "id" LIMIT 10 OFFSET ?`, offset,
		)
		if err != nil {
			b.Fatal(err)
		}
		if rows.Err() != nil {
			b.Fatal(rows.Err())
		}
		for rows.Next() {
			var id, email, phoneNum, role string
			if err := rows.Scan(&id, &email, &phoneNum, &role); err != nil {
				b.Fatal(err)
			}
		}
		rows.Close()
	}
}

func benchRawUpsert(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)
	seedData(db, "raw-ups")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		email := fmt.Sprintf("raw-ups-user-%d@example.com", i%seedCount)
		phone := fmt.Sprintf("raw-ups-phone-%d", i)
		_, err := db.ExecContext(ctx,
			`INSERT INTO "User" ("id", "email", "phoneNum", "role", "loginCount") VALUES (?, ?, ?, ?, ?)
			 ON CONFLICT("email") DO UPDATE SET "phoneNum" = EXCLUDED."phoneNum", "role" = EXCLUDED."role", "loginCount" = EXCLUDED."loginCount"`,
			fmt.Sprintf("raw-ups-id-%d", i),
			email, phone, "student", int32(i),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}
