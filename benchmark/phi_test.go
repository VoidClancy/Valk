package main

import (
	"benchmark/phi"
	"benchmark/phi/user"
	"context"
	"database/sql"
	"fmt"
	"testing"
)

func initPhiDB(b *testing.B, ctx context.Context) *phi.DB {
	b.Helper()
	db, err := phi.Open(activeDialect.Driver, activeDialect.DSN)
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

func benchPhiCreate(b *testing.B) {
	ctx := context.Background()
	db := initPhiDB(b, ctx)
	defer db.Close()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, err := db.User.Create().
			SetId(fmt.Sprintf("phi-create-%d", i)).
			SetEmail(fmt.Sprintf("phi-create-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("phi-create-phone-%d", i)).
			SetRole(phi.UserRole_STUDENT).
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchPhiCreateMany(b *testing.B) {
	ctx := context.Background()
	db := initPhiDB(b, ctx)
	defer db.Close()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		builders := make([]*user.CreateBuilder, 10)
		for j := 0; j < 10; j++ {
			n := i*10 + j
			builders[j] = db.User.Create().
				SetId(fmt.Sprintf("phi-cmany-%d", n)).
				SetEmail(fmt.Sprintf("phi-cmany-%d@example.com", n)).
				SetPhoneNum(fmt.Sprintf("phi-cmany-phone-%d", n)).
				SetRole(phi.UserRole_STUDENT).
				SetLoginCount(0)
		}
		_, err := db.User.CreateMany(builders...).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchPhiCreateManyAndReturn(b *testing.B) {
	ctx := context.Background()
	db := initPhiDB(b, ctx)
	defer db.Close()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		builders := make([]*user.CreateBuilder, 10)
		for j := range 10 {
			n := i*10 + j
			builders[j] = db.User.Create().
				SetId(fmt.Sprintf("phi-cmar-%d", n)).
				SetEmail(fmt.Sprintf("phi-cmar-%d@example.com", n)).
				SetPhoneNum(fmt.Sprintf("phi-cmar-phone-%d", n)).
				SetRole(phi.UserRole_STUDENT).
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

func benchPhiFindUnique(b *testing.B) {
	ctx := context.Background()
	db := initPhiDB(b, ctx)
	defer db.Close()
	for i := 0; i < seedCount; i++ {
		_, err := db.User.Create().
			SetId(fmt.Sprintf("phi-fu-id-%d", i)).
			SetEmail(fmt.Sprintf("phi-fu-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("phi-fu-phone-%d", i)).
			SetRole(phi.UserRole_STUDENT).
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, err := db.User.FindUnique(
			user.Email.EQ(fmt.Sprintf("phi-fu-%d@example.com", i%seedCount)),
		).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchPhiFindFirst(b *testing.B) {
	ctx := context.Background()
	db := initPhiDB(b, ctx)
	defer db.Close()
	for i := 0; i < seedCount; i++ {
		_, err := db.User.Create().
			SetId(fmt.Sprintf("phi-ff-id-%d", i)).
			SetEmail(fmt.Sprintf("phi-ff-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("phi-ff-phone-%d", i)).
			SetRole(phi.UserRole_STUDENT).
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, err := db.User.FindFirst(
			user.Email.EQ(fmt.Sprintf("phi-ff-%d@example.com", i%seedCount)),
		).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchPhiFindMany(b *testing.B) {
	ctx := context.Background()
	db := initPhiDB(b, ctx)
	defer db.Close()
	for i := 0; i < seedCount; i++ {
		_, err := db.User.Create().
			SetId(fmt.Sprintf("phi-fm-id-%d", i)).
			SetEmail(fmt.Sprintf("phi-fm-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("phi-fm-phone-%d", i)).
			SetRole(phi.UserRole_STUDENT).
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

func benchPhiUpsert(b *testing.B) {
	ctx := context.Background()
	db := initPhiDB(b, ctx)
	defer db.Close()
	for i := 0; i < seedCount; i++ {
		_, err := db.User.Create().
			SetId(fmt.Sprintf("phi-ups-id-%d", i)).
			SetEmail(fmt.Sprintf("phi-ups-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("phi-ups-phone-%d", i)).
			SetRole(phi.UserRole_STUDENT).
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		email := fmt.Sprintf("phi-ups-%d@example.com", i%seedCount)
		_, err := db.User.Create().
			SetId(fmt.Sprintf("phi-ups-id-new-%d", i)).
			SetEmail(email).
			SetPhoneNum(fmt.Sprintf("phi-ups-phone-new-%d", i)).
			SetRole(phi.UserRole_STUDENT).
			SetLoginCount(int32(i)).
			OnConflict(user.Email).
			UpdateNewValues().
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchPhiReadDeepRelation(b *testing.B) {
	ctx := context.Background()
	db := initPhiDB(b, ctx)
	defer db.Close()
	seedRelations(db.Raw(), "phi-rdr")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		email := fmt.Sprintf("phi-rdr-grand-%d@example.com", i%500)
		_, err := db.User.FindUnique(
			user.Email.EQ(email),
		).Select(phi.UserSelect{
			ReferredBy: &phi.UserSelect{
				ReferredBy: &phi.UserSelect{},
			},
		}).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchPhiCreateWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	db := initPhiDB(b, ctx)
	defer db.Close()
	seedRelations(db.Raw(), "phi-cwds")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("phi-cwds-parent-id-%d", i%500)
		_, err := db.User.Create().
			SetId(fmt.Sprintf("phi-cwds-new-id-%d", i)).
			SetEmail(fmt.Sprintf("phi-cwds-new-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("phi-cwds-new-phone-%d", i)).
			SetRole(phi.UserRole_STUDENT).
			SetReferredById(parentID).
			Select(phi.UserSelect{
				ReferredBy: &phi.UserSelect{
					ReferredBy: &phi.UserSelect{},
				},
			}).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchPhiCreateManyAndReturnWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	db := initPhiDB(b, ctx)
	defer db.Close()
	seedRelations(db.Raw(), "phi-cmwds")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("phi-cmwds-parent-id-%d", i%500)
		inputs := make([]*phi.UserCreateBuilder, 10)
		for j := 0; j < 10; j++ {
			inputs[j] = db.User.Create().
				SetId(fmt.Sprintf("phi-cmwds-new-id-%d-%d", i, j)).
				SetEmail(fmt.Sprintf("phi-cmwds-new-%d-%d@example.com", i, j)).
				SetPhoneNum(fmt.Sprintf("phi-cmwds-new-phone-%d-%d", i, j)).
				SetRole(phi.UserRole_STUDENT).
				SetReferredById(parentID)
		}
		_, err := db.User.CreateManyAndReturn(inputs...).Select(phi.UserSelect{
			ReferredBy: &phi.UserSelect{
				ReferredBy: &phi.UserSelect{},
			},
		}).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchPhiUpsertWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	db := initPhiDB(b, ctx)
	defer db.Close()
	seedRelations(db.Raw(), "phi-uwds")

	for i := 0; i < seedCount; i++ {
		_, err := db.User.Create().
			SetId(fmt.Sprintf("phi-uwds-id-%d", i)).
			SetEmail(fmt.Sprintf("phi-uwds-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("phi-uwds-phone-%d", i)).
			SetRole(phi.UserRole_STUDENT).
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("phi-uwds-parent-id-%d", i%500)
		email := fmt.Sprintf("phi-uwds-%d@example.com", i%seedCount)
		_, err := db.User.Create().
			SetId(fmt.Sprintf("phi-uwds-id-new-%d", i)).
			SetEmail(email).
			SetPhoneNum(fmt.Sprintf("phi-uwds-phone-new-%d", i)).
			SetRole(phi.UserRole_STUDENT).
			SetLoginCount(int32(i)).
			SetReferredById(parentID).
			OnConflict(user.Email).
			UpdateNewValues().
			Select(phi.UserSelect{
				ReferredBy: &phi.UserSelect{
					ReferredBy: &phi.UserSelect{},
				},
			}).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchPhiDelete(b *testing.B) {
	ctx := context.Background()
	db := initPhiDB(b, ctx)
	defer db.Close()

	seedDeleteData(db.Raw(), "phi", 50000)

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, err := db.User.Delete(user.Id.EQ(fmt.Sprintf("phi-d-hook-%d", i))).Exec(ctx)
		if err != nil && err != sql.ErrNoRows {
			b.Fatal(err)
		}
	}
}
