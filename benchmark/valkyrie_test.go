package main

import (
	"context"
	"fmt"
	"integration/valk"
	"integration/valk/user"
	"testing"
)

func initValkDB(b *testing.B, ctx context.Context) *valk.DB {
	b.Helper()
	db, err := valk.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		b.Fatal(err)
	}
	_, err = db.Raw().ExecContext(ctx, userSchema)
	if err != nil {
		b.Fatal(err)
	}
	return db
}

func benchValkyrieCreate(b *testing.B) {
	ctx := context.Background()
	db := initValkDB(b, ctx)
	defer db.Close()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, err := db.User.Create().
			SetId(fmt.Sprintf("valk-create-%d", i)).
			SetEmail(fmt.Sprintf("valk-create-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("valk-create-phone-%d", i)).
			SetRole(valk.UserRoleTypeStudent).
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchValkyrieCreateMany(b *testing.B) {
	ctx := context.Background()
	db := initValkDB(b, ctx)
	defer db.Close()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		builders := make([]*user.CreateBuilder, 10)
		for j := 0; j < 10; j++ {
			n := i*10 + j
			builders[j] = db.User.Create().
				SetId(fmt.Sprintf("valk-cmany-%d", n)).
				SetEmail(fmt.Sprintf("valk-cmany-%d@example.com", n)).
				SetPhoneNum(fmt.Sprintf("valk-cmany-phone-%d", n)).
				SetRole(valk.UserRoleTypeStudent).
				SetLoginCount(0)
		}
		_, err := db.User.CreateMany(builders...).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchValkyrieCreateManyAndReturn(b *testing.B) {
	ctx := context.Background()
	db := initValkDB(b, ctx)
	defer db.Close()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		builders := make([]*user.CreateBuilder, 10)
		for j := 0; j < 10; j++ {
			n := i*10 + j
			builders[j] = db.User.Create().
				SetId(fmt.Sprintf("valk-cmar-%d", n)).
				SetEmail(fmt.Sprintf("valk-cmar-%d@example.com", n)).
				SetPhoneNum(fmt.Sprintf("valk-cmar-phone-%d", n)).
				SetRole(valk.UserRoleTypeStudent).
				SetLoginCount(0)
		}
		users, err := db.User.CreateManyAndReturn(builders...).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
		if len(users) != 10 {
			b.Fatalf("expected 10 users, got %d", len(users))
		}
	}
}

func benchValkyrieFindUnique(b *testing.B) {
	ctx := context.Background()
	db := initValkDB(b, ctx)
	defer db.Close()
	for i := 0; i < seedCount; i++ {
		_, err := db.User.Create().
			SetId(fmt.Sprintf("valk-fu-id-%d", i)).
			SetEmail(fmt.Sprintf("valk-fu-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("valk-fu-phone-%d", i)).
			SetRole(valk.UserRoleTypeStudent).
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, err := db.User.FindUnique(
			user.Email.EQ(fmt.Sprintf("valk-fu-%d@example.com", i%seedCount)),
		).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchValkyrieFindFirst(b *testing.B) {
	ctx := context.Background()
	db := initValkDB(b, ctx)
	defer db.Close()
	for i := 0; i < seedCount; i++ {
		_, err := db.User.Create().
			SetId(fmt.Sprintf("valk-ff-id-%d", i)).
			SetEmail(fmt.Sprintf("valk-ff-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("valk-ff-phone-%d", i)).
			SetRole(valk.UserRoleTypeStudent).
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, err := db.User.FindFirst(
			user.Email.EQ(fmt.Sprintf("valk-ff-%d@example.com", i%seedCount)),
		).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchValkyrieFindMany(b *testing.B) {
	ctx := context.Background()
	db := initValkDB(b, ctx)
	defer db.Close()
	for i := 0; i < seedCount; i++ {
		_, err := db.User.Create().
			SetId(fmt.Sprintf("valk-fm-id-%d", i)).
			SetEmail(fmt.Sprintf("valk-fm-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("valk-fm-phone-%d", i)).
			SetRole(valk.UserRoleTypeStudent).
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		users, err := db.User.FindMany().Take(10).Skip(i % seedCount).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
		if len(users) == 0 {
			b.Fatal("expected at least one user")
		}
	}
}

func benchValkyrieUpsert(b *testing.B) {
	ctx := context.Background()
	db := initValkDB(b, ctx)
	defer db.Close()
	for i := 0; i < seedCount; i++ {
		_, err := db.User.Create().
			SetId(fmt.Sprintf("valk-ups-id-%d", i)).
			SetEmail(fmt.Sprintf("valk-ups-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("valk-ups-phone-%d", i)).
			SetRole(valk.UserRoleTypeStudent).
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		email := fmt.Sprintf("valk-ups-%d@example.com", i%seedCount)
		_, err := db.User.Create().
			SetId(fmt.Sprintf("valk-ups-id-new-%d", i)).
			SetEmail(email).
			SetPhoneNum(fmt.Sprintf("valk-ups-phone-new-%d", i)).
			SetRole(valk.UserRoleTypeStudent).
			SetLoginCount(int32(i)).
			OnConflict(user.Email).
			UpdateNewValues().
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}
