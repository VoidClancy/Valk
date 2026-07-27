package main

import (
	"benchmark/valk"
	"benchmark/valk/user"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
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
		_, err := db.User.Create().
			SetId(fmt.Sprintf("valk-cwds-new-id-%d", i)).
			SetEmail(fmt.Sprintf("valk-cwds-new-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("valk-cwds-new-phone-%d", i)).
			SetRole(valk.UserRoleTypeStudent).
			SetReferredById(parentID).
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
			inputs[j] = db.User.Create().
				SetId(fmt.Sprintf("valk-cmwds-new-id-%d-%d", i, j)).
				SetEmail(fmt.Sprintf("valk-cmwds-new-%d-%d@example.com", i, j)).
				SetPhoneNum(fmt.Sprintf("valk-cmwds-new-phone-%d-%d", i, j)).
				SetRole(valk.UserRoleTypeStudent).
				SetReferredById(parentID)
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
		_, err := db.User.Create().
			SetId(fmt.Sprintf("valk-uwds-id-%d", i)).
			SetEmail(fmt.Sprintf("valk-uwds-%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("valk-uwds-phone-%d", i)).
			SetRole(valk.UserRoleTypeStudent).
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("valk-uwds-parent-id-%d", i%500)
		email := fmt.Sprintf("valk-uwds-%d@example.com", i%seedCount)
		_, err := db.User.Create().
			SetId(fmt.Sprintf("valk-uwds-id-new-%d", i)).
			SetEmail(email).
			SetPhoneNum(fmt.Sprintf("valk-uwds-phone-new-%d", i)).
			SetRole(valk.UserRoleTypeStudent).
			SetLoginCount(int32(i)).
			SetReferredById(parentID).
			OnConflict(user.Email).
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

func benchValkyrieHooksCreate(b *testing.B) {
	ctx := context.Background()
	db := initValkDB(b, ctx)
	defer db.Close()

	db.User.Use(valk.UserExtension{
		Create: func(ctx context.Context, args *valk.UserCreateArgs, next valk.UserCreateQuery) (*valk.User, error) {
			if args.Data != nil {
				args.Data.Email = strings.ToLower(args.Data.Email)
				if args.Data.Password != nil {
					hash := sha256.Sum256([]byte(*args.Data.Password))
					hexStr := hex.EncodeToString(hash[:])
					args.Data.Password = &hexStr
				}
			}
			return next(ctx, args)
		},
	})

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		pwd := "mypassword"
		_, err := db.User.Create().
			SetId(fmt.Sprintf("valk-c-hook-%d", i)).
			SetEmail(fmt.Sprintf("Valk-C-Hook-%d@Example.com", i)).
			SetPhoneNum(fmt.Sprintf("valk-c-hook-phone-%d", i)).
			SetPassword(pwd).
			SetRole(valk.UserRoleTypeStudent).
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchValkyrieHooksUpdate(b *testing.B) {
	ctx := context.Background()
	db := initValkDB(b, ctx)
	defer db.Close()

	seedData(db.Raw(), "valk-u-hook")

	db.User.Use(valk.UserExtension{
		Update: func(ctx context.Context, where valk.UniquePredicate[valk.User], additional []valk.PredicateOf[valk.User], assignments []valk.FieldAssignment, selects *valk.UserSelect, omits *valk.UserOmit, next valk.UserUpdateQuery) (*valk.User, error) {
			for i, a := range assignments {
				if a.Col == "email" {
					if s, ok := a.Val.(string); ok {
						assignments[i].Val = strings.ToLower(s)
					}
				}
			}
			return next(ctx, where, additional, assignments, selects, omits)
		},
	})

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, err := db.User.Update(user.Id.EQ(fmt.Sprintf("valk-u-hook-id-%d", i%seedCount))).
			SetEmail(fmt.Sprintf("NEW-VALK-U-HOOK-%d@EXAMPLE.COM", i)).
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchValkyrieHooksFindUnique(b *testing.B) {
	ctx := context.Background()
	db := initValkDB(b, ctx)
	defer db.Close()

	seedData(db.Raw(), "valk-f-hook")

	db.User.Use(valk.UserExtension{
		FindUnique: func(ctx context.Context, args *valk.UserFindUniqueArgs, next valk.UserFindUniqueQuery) (*valk.User, error) {
			u, err := next(ctx, args)
			if err == nil && u != nil {
				u.Email += "-read"
			}
			return u, err
		},
	})

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, err := db.User.FindUnique(user.Id.EQ(fmt.Sprintf("valk-f-hook-id-%d", i%seedCount))).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchValkyrieHooksDelete(b *testing.B) {
	ctx := context.Background()
	db := initValkDB(b, ctx)
	defer db.Close()

	seedDeleteData(db.Raw(), "valk", 50000)

	db.User.Use(valk.UserExtension{
		Delete: func(ctx context.Context, where valk.UniquePredicate[valk.User], selects *valk.UserSelect, omits *valk.UserOmit, next valk.UserDeleteQuery) (*valk.User, error) {
			u, err := next(ctx, where, selects, omits)
			if err == nil && u != nil {
				u.Id = strings.ToUpper(u.Id)
			}
			return u, err
		},
	})

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, err := db.User.Delete(user.Id.EQ(fmt.Sprintf("valk-d-hook-%d", i%50000))).Exec(ctx)
		if err != nil && err != sql.ErrNoRows {
			b.Fatal(err)
		}
	}
}
