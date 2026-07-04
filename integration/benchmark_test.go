package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"integration/valkyrie"
	"strconv"
	"testing"
	"time"
)

func generateCUID() string {
	now := uint64(time.Now().UnixMilli())
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	buf := make([]byte, 1, 33)
	buf[0] = 'c'
	buf = strconv.AppendUint(buf, now, 16)
	
	const hextable = "0123456789abcdef"
	for _, v := range b {
		buf = append(buf, hextable[v>>4], hextable[v&0x0f])
	}
	return string(buf)
}

func TestCreationBenchmark(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	const iterations = 500

	// 1. Raw SQL Insert (Write Only)
	t.Logf("Running %d iterations of Raw SQL Insert (Write-only)...", iterations)
	startRawWrite := time.Now()
	for i := 0; i < iterations; i++ {
		id := "raw-w-" + strconv.Itoa(i)
		email := fmt.Sprintf("raw-w-%d@example.com", i)
		phone := fmt.Sprintf("+12345%d", i)
		role := "student"

		_, err := db.Raw().ExecContext(ctx,
			"INSERT INTO User (id, email, phoneNum, role) VALUES (?, ?, ?, ?)",
			id, email, phone, role,
		)
		if err != nil {
			t.Fatalf("Raw SQL write failed: %v", err)
		}
	}
	rawWriteDuration := time.Since(startRawWrite)
	t.Logf("Raw SQL Write-only: %s (avg %s/op)", rawWriteDuration, rawWriteDuration/iterations)

	// 2. ORM Create (which inserts and scans back the record)
	t.Logf("Running %d iterations of ORM Create...", iterations)
	startORM := time.Now()
	for i := 0; i < iterations; i++ {
		_, err := db.User.Create(valkyrie.UserCreateInput{
			Email:    fmt.Sprintf("orm-%d@example.com", i),
			PhoneNum: fmt.Sprintf("+12345%d", i),
		}).Exec(ctx)
		if err != nil {
			t.Fatalf("ORM create failed: %v", err)
		}
	}
	ormDuration := time.Since(startORM)
	t.Logf("ORM Create (Insert + Scan): %s (avg %s/op)", ormDuration, ormDuration/iterations)

	// 3. Raw SQL Insert + Select (to match ORM behavior of returning the record)
	t.Logf("Running %d iterations of Raw SQL Insert + Scan...", iterations)
	startRawRead := time.Now()
	for i := range iterations {
		id := "raw-r-" + strconv.Itoa(i)
		email := fmt.Sprintf("raw-r-%d@example.com", i)
		phone := fmt.Sprintf("+12345%d", i)
		role := "student"

		_, err := db.Raw().ExecContext(ctx,
			"INSERT INTO User (id, email, phoneNum, role) VALUES (?, ?, ?, ?)",
			id, email, phone, role,
		)
		if err != nil {
			t.Fatalf("Raw SQL write failed: %v", err)
		}

		var res valkyrie.User
		err = db.Raw().QueryRowContext(ctx,
			"SELECT id, email, phoneNum, role, referredById FROM User WHERE id = ?",
			id,
		).Scan(&res.Id, &res.Email, &res.PhoneNum, &res.Role, &res.ReferredById)
		if err != nil {
			t.Fatalf("Raw SQL select failed: %v", err)
		}
	}
	rawReadDuration := time.Since(startRawRead)
	t.Logf("Raw SQL Insert + Scan: %s (avg %s/op)", rawReadDuration, rawReadDuration/iterations)
}

func BenchmarkORMCreate(b *testing.B) {
	db, cleanup := setupTestDB(nil) // setupTestDB supports t being nil/non-nil for cleanup
	defer cleanup()
	ctx := context.Background()

	for i := 0; b.Loop(); i++ {
		_, err := db.User.Create(valkyrie.UserCreateInput{
			Email:    fmt.Sprintf("bench-orm-%d@example.com", i),
			PhoneNum: fmt.Sprintf("+98765%d", i),
		}).Exec(ctx)
		if err != nil {
			b.Fatalf("ORM create failed: %v", err)
		}
	}
}

func BenchmarkRawSQLCreate(b *testing.B) {
	db, cleanup := setupTestDB(nil)
	defer cleanup()
	ctx := context.Background()

	for i := 0; b.Loop(); i++ {
		id := "bench-raw-" + strconv.Itoa(i)
		email := fmt.Sprintf("bench-raw-%d@example.com", i)
		phone := fmt.Sprintf("+98765%d", i)
		role := "student"

		_, err := db.Raw().ExecContext(ctx,
			"INSERT INTO User (id, email, phoneNum, role) VALUES (?, ?, ?, ?)",
			id, email, phone, role,
		)
		if err != nil {
			b.Fatalf("Raw SQL write failed: %v", err)
		}
	}
}

func BenchmarkRawSQLCreateWithScan(b *testing.B) {
	db, cleanup := setupTestDB(nil)
	defer cleanup()
	ctx := context.Background()

	for i := 0; b.Loop(); i++ {
		id := generateCUID()
		email := fmt.Sprintf("bench-raw-scan-%d@example.com", i)
		phone := fmt.Sprintf("+98765%d", i)
		role := "student"

		_, err := db.Raw().ExecContext(ctx,
			"INSERT INTO User (id, email, phoneNum, role) VALUES (?, ?, ?, ?)",
			id, email, phone, role,
		)
		if err != nil {
			b.Fatalf("Raw SQL write failed: %v", err)
		}

		var res valkyrie.User
		err = db.Raw().QueryRowContext(ctx,
			"SELECT id, email, phoneNum, role, referredById FROM User WHERE id = ?",
			id,
		).Scan(&res.Id, &res.Email, &res.PhoneNum, &res.Role, &res.ReferredById)
		if err != nil {
			b.Fatalf("Raw SQL select failed: %v", err)
		}
	}
}
