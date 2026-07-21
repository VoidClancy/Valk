package main

import (
	"benchmark/ent"
	"benchmark/ent/user"
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func openEnt(b *testing.B) *ent.Client {
	b.Helper()
	var driverName string
	var dsn string
	if activeDialect.Name == "postgres" {
		driverName = "postgres"
		dsn = activeDialect.DSN
	} else {
		driverName = "sqlite3"
		dsn = activeDialect.DSN + "&_fk=1"
	}

	client, err := ent.Open(driverName, dsn)
	if err != nil {
		b.Fatal(err)
	}

	// Reset PG schema
	if activeDialect.Name == "postgres" {
		db, err := sql.Open("postgres", activeDialect.DSN)
		if err == nil {
			resetPostgres(db)
			db.Close()
		}
	}

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		b.Fatal(err)
	}
	return client
}

func seedEnt(b *testing.B, client *ent.Client, prefix string) {
	b.Helper()
	ctx := context.Background()
	for i := range seedCount {
		_, err := client.User.Create().
			SetID(fmt.Sprintf("%s-id-%d", prefix, i)).
			SetEmail(fmt.Sprintf("%s-user-%d@example.com", prefix, i)).
			SetPhoneNum(fmt.Sprintf("%s-phone-%d", prefix, i)).
			SetRole(user.DefaultRole).
			Save(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchEntCreate(b *testing.B) {
	ctx := context.Background()
	client := openEnt(b)
	defer client.Close()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, err := client.User.Create().
			SetID(fmt.Sprintf("ent-%d", i)).
			SetEmail(fmt.Sprintf("ent-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("ent-phone-%d", i)).
			SetRole(user.DefaultRole).
			Save(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchEntCreateMany(b *testing.B) {
	ctx := context.Background()
	client := openEnt(b)
	defer client.Close()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		builders := make([]*ent.UserCreate, 10)
		for j := 0; j < 10; j++ {
			n := i*10 + j
			builders[j] = client.User.Create().
				SetID(fmt.Sprintf("ent-cmany-%d", n)).
				SetEmail(fmt.Sprintf("ent-cmany-%d@example.com", n)).
				SetPhoneNum(fmt.Sprintf("ent-cmany-phone-%d", n)).
				SetRole(user.DefaultRole).
				SetLoginCount(0)
		}
		_, err := client.User.CreateBulk(builders...).Save(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchEntCreateManyAndReturn(b *testing.B) {
	ctx := context.Background()
	client := openEnt(b)
	defer client.Close()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		builders := make([]*ent.UserCreate, 10)
		for j := 0; j < 10; j++ {
			n := i*10 + j
			builders[j] = client.User.Create().
				SetID(fmt.Sprintf("ent-cmar-%d", n)).
				SetEmail(fmt.Sprintf("ent-cmar-%d@example.com", n)).
				SetPhoneNum(fmt.Sprintf("ent-cmar-phone-%d", n)).
				SetRole(user.DefaultRole).
				SetLoginCount(0)
		}
		users, err := client.User.CreateBulk(builders...).Save(ctx)
		if err != nil {
			b.Fatal(err)
		}
		if len(users) != 10 {
			b.Fatalf("expected 10 users, got %d", len(users))
		}
	}
}

func benchEntFindUnique(b *testing.B) {
	ctx := context.Background()
	client := openEnt(b)
	defer client.Close()
	seedEnt(b, client, "ent-fu")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, err := client.User.Query().
			Where(user.Email(fmt.Sprintf("ent-fu-user-%d@example.com", i%seedCount))).
			Only(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchEntFindFirst(b *testing.B) {
	ctx := context.Background()
	client := openEnt(b)
	defer client.Close()
	seedEnt(b, client, "ent-ff")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, err := client.User.Query().
			Where(user.Email(fmt.Sprintf("ent-ff-user-%d@example.com", i%seedCount))).
			First(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchEntFindMany(b *testing.B) {
	ctx := context.Background()
	client := openEnt(b)
	defer client.Close()
	seedEnt(b, client, "ent-fm")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		users, err := client.User.Query().
			Limit(10).
			Offset(i % seedCount).
			All(ctx)
		if err != nil {
			b.Fatal(err)
		}
		if len(users) == 0 {
			b.Fatal("expected at least one user")
		}
	}
}

func benchEntUpsert(b *testing.B) {
	ctx := context.Background()
	client := openEnt(b)
	defer client.Close()
	seedEnt(b, client, "ent-ups")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		email := fmt.Sprintf("ent-ups-user-%d@example.com", i%seedCount)
		if err := client.User.Create().
			SetID(fmt.Sprintf("ent-ups-id-%d", i)).
			SetEmail(email).
			SetPhoneNum(fmt.Sprintf("ent-ups-phone-new-%d", i)).
			SetRole(user.DefaultRole).
			SetLoginCount(int32(i)).
			OnConflictColumns("email").
			UpdateNewValues().
			Exec(ctx); err != nil {
			b.Fatal(err)
		}
	}
}

func benchEntReadDeepRelation(b *testing.B) {
	ctx := context.Background()
	client := openEnt(b)
	defer client.Close()

	db := openDB(b)
	seedRelations(db, "ent-rdr")
	db.Close()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		email := fmt.Sprintf("ent-rdr-grand-%d@example.com", i%500)
		_, err := client.User.Query().
			Where(user.Email(email)).
			WithReferredBy(func(q *ent.UserQuery) {
				q.WithReferredBy()
			}).
			Only(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchEntCreateWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	client := openEnt(b)
	defer client.Close()

	db := openDB(b)
	seedRelations(db, "ent-cwds")
	db.Close()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("ent-cwds-parent-id-%d", i%500)
		u, err := client.User.Create().
			SetID(fmt.Sprintf("ent-cwds-new-id-%d", i)).
			SetEmail(fmt.Sprintf("ent-cwds-new-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("ent-cwds-new-phone-%d", i)).
			SetRole("STUDENT").
			SetReferredByID(parentID).
			Save(ctx)
		if err != nil {
			b.Fatal(err)
		}
		referredBy, err := u.QueryReferredBy().WithReferredBy().Only(ctx)
		if err != nil {
			b.Fatal(err)
		}
		u.Edges.ReferredBy = referredBy
	}
}

func benchEntCreateManyAndReturnWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	client := openEnt(b)
	defer client.Close()

	db := openDB(b)
	seedRelations(db, "ent-cmwds")
	db.Close()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("ent-cmwds-parent-id-%d", i%500)
		builders := make([]*ent.UserCreate, 10)
		ids := make([]string, 10)
		for j := 0; j < 10; j++ {
			id := fmt.Sprintf("ent-cmwds-new-id-%d-%d", i, j)
			ids[j] = id
			builders[j] = client.User.Create().
				SetID(id).
				SetEmail(fmt.Sprintf("ent-cmwds-new-%d-%d@example.com", i, j)).
				SetPhoneNum(fmt.Sprintf("ent-cmwds-new-phone-%d-%d", i, j)).
				SetRole("STUDENT").
				SetReferredByID(parentID)
		}
		_, err := client.User.CreateBulk(builders...).Save(ctx)
		if err != nil {
			b.Fatal(err)
		}
		_, err = client.User.Query().
			Where(user.IDIn(ids...)).
			WithReferredBy(func(q *ent.UserQuery) {
				q.WithReferredBy()
			}).
			All(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchEntUpsertWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	client := openEnt(b)
	defer client.Close()

	db := openDB(b)
	seedRelations(db, "ent-uwds")
	db.Close()

	// Preseed
	for i := range seedCount {
		_, err := client.User.Create().
			SetID(fmt.Sprintf("ent-uwds-id-%d", i)).
			SetEmail(fmt.Sprintf("ent-uwds-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("ent-uwds-phone-%d", i)).
			SetRole("STUDENT").
			Save(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("ent-uwds-parent-id-%d", i%500)
		email := fmt.Sprintf("ent-uwds-%d@example.com", i%seedCount)
		err := client.User.Create().
			SetID(fmt.Sprintf("ent-uwds-id-new-%d", i)).
			SetEmail(email).
			SetPhoneNum(fmt.Sprintf("ent-uwds-phone-new-%d", i)).
			SetRole("STUDENT").
			SetLoginCount(int32(i)).
			SetReferredByID(parentID).
			OnConflictColumns("email").
			UpdateNewValues().
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
		_, err = client.User.Query().
			Where(user.Email(email)).
			WithReferredBy(func(q *ent.UserQuery) {
				q.WithReferredBy()
			}).
			Only(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}
