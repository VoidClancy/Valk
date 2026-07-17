package main

import (
	"benchmark/ent"
	"benchmark/ent/user"
	"context"
	"fmt"
	"testing"
)

func openEnt(b *testing.B) *ent.Client {
	b.Helper()
	client, err := ent.Open("sqlite3", "file::memory:?cache=shared&_fk=1")
	if err != nil {
		b.Fatal(err)
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
