package main

import (
	"benchmark/valk"
	"benchmark/valk/user"
	"context"
	"fmt"
	"testing"
)

func initValkDB(b *testing.B, ctx context.Context) *valk.DB {
	b.Helper()
	db, err := valk.Open(activeDialect.Driver, activeDialect.DSN)
	if err != nil {
		b.Fatal(err)
	}

	if activeDialect.Name == "postgres" {
		resetPostgres(db.Raw())

		err = db.RunMigrations(ctx)
		if err != nil {
			b.Fatal(err)
		}
	} else {
		createSQLiteSchema(db.Raw())
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
		for j := range 10 {
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

func benchValkyrieReadDeepRelation(b *testing.B) {
	ctx := context.Background()
	db := initValkDB(b, ctx)
	defer db.Close()
	seedRelations(db.Raw(), "valk-rdr")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		email := fmt.Sprintf("valk-rdr-grand-%d@example.com", i%500)
		_, err := db.User.FindUnique(
			user.Email.EQ(email),
		).Select(valk.UserSelect{
			ReferredBy: &valk.UserSelect{
				ReferredBy: &valk.UserSelect{},
			},
		}).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchValkyrieCreateWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	db := initValkDB(b, ctx)
	defer db.Close()
	seedRelations(db.Raw(), "valk-cwds")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("valk-cwds-parent-id-%d", i%500)
		_, err := db.User.Create(
			user.Id.Set(fmt.Sprintf("valk-cwds-new-id-%d", i)),
			user.Email.Set(fmt.Sprintf("valk-cwds-new-%d@example.com", i)),
			user.PhoneNum.Set(fmt.Sprintf("valk-cwds-new-phone-%d", i)),
			user.Role.Set(valk.UserRoleTypeStudent),
			user.ReferredById.Set(parentID),
		).Select(valk.UserSelect{
			ReferredBy: &valk.UserSelect{
				ReferredBy: &valk.UserSelect{},
			},
		}).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchValkyrieCreateManyAndReturnWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	db := initValkDB(b, ctx)
	defer db.Close()
	seedRelations(db.Raw(), "valk-cmwds")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("valk-cmwds-parent-id-%d", i%500)
		inputs := make([]*valk.UserCreateBuilder, 10)
		for j := 0; j < 10; j++ {
			inputs[j] = db.User.Create(
				user.Id.Set(fmt.Sprintf("valk-cmwds-new-id-%d-%d", i, j)),
				user.Email.Set(fmt.Sprintf("valk-cmwds-new-%d-%d@example.com", i, j)),
				user.PhoneNum.Set(fmt.Sprintf("valk-cmwds-new-phone-%d-%d", i, j)),
				user.Role.Set(valk.UserRoleTypeStudent),
				user.ReferredById.Set(parentID),
			)
		}
		_, err := db.User.CreateManyAndReturn(inputs...).Select(valk.UserSelect{
			ReferredBy: &valk.UserSelect{
				ReferredBy: &valk.UserSelect{},
			},
		}).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchValkyrieUpsertWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	db := initValkDB(b, ctx)
	defer db.Close()
	seedRelations(db.Raw(), "valk-uwds")

	for i := 0; i < seedCount; i++ {
		_, err := db.User.Create(
			user.Id.Set(fmt.Sprintf("valk-uwds-id-%d", i)),
			user.Email.Set(fmt.Sprintf("valk-uwds-%d@example.com", i)),
			user.PhoneNum.Set(fmt.Sprintf("valk-uwds-phone-%d", i)),
			user.Role.Set(valk.UserRoleTypeStudent),
		).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("valk-uwds-parent-id-%d", i%500)
		email := fmt.Sprintf("valk-uwds-%d@example.com", i%seedCount)
		_, err := db.User.Create(
			user.Id.Set(fmt.Sprintf("valk-uwds-id-new-%d", i)),
			user.Email.Set(email),
			user.PhoneNum.Set(fmt.Sprintf("valk-uwds-phone-new-%d", i)),
			user.Role.Set(valk.UserRoleTypeStudent),
			user.LoginCount.Set(int32(i)),
			user.ReferredById.Set(parentID),
		).OnConflict(user.Email).
			UpdateNewValues().
			Select(valk.UserSelect{
				ReferredBy: &valk.UserSelect{
					ReferredBy: &valk.UserSelect{},
				},
			}).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}
