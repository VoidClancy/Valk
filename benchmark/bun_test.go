package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/schema"
)

type UserBun struct {
	bun.BaseModel `bun:"table:User"`
	Id            string   `bun:"id,pk"`
	Email         string   `bun:"email,unique,notnull"`
	PhoneNum      string   `bun:"phoneNum,unique,notnull"`
	Password      *string  `bun:"password"`
	Role          string   `bun:"role,default:'STUDENT'"`
	RoleOptional  *string  `bun:"roleOptional"`
	LoginCount    int32    `bun:"loginCount,default:0"`
	ReferredById  *string  `bun:"referredById"`
	ReferredBy    *UserBun `bun:"rel:belongs-to,join:referredById=id"`
}

func openBun(b *testing.B) *bun.DB {
	b.Helper()
	sqldb, err := sql.Open(activeDialect.Driver, activeDialect.DSN)
	if err != nil {
		b.Fatal(err)
	}
	sqldb.SetMaxOpenConns(80)
	sqldb.SetMaxIdleConns(80)

	var dialect schema.Dialect
	if activeDialect.Name == "postgres" {
		dialect = pgdialect.New()
	} else {
		dialect = sqlitedialect.New()
	}
	return bun.NewDB(sqldb, dialect)
}

func seedBun(db *bun.DB, prefix string) {
	ctx := context.Background()
	for i := range seedCount {
		_, err := db.NewInsert().Model(&UserBun{
			Id:       fmt.Sprintf("%s-id-%d", prefix, i),
			Email:    fmt.Sprintf("%s-user-%d@example.com", prefix, i),
			PhoneNum: fmt.Sprintf("%s-phone-%d", prefix, i),
			Role:     "STUDENT",
		}).Exec(ctx)
		if err != nil {
			panic(fmt.Sprintf("seed %s: %v", prefix, err))
		}
	}
}

func benchBunCreate(b *testing.B) {
	ctx := context.Background()
	db := openBun(b)
	defer db.Close()
	createSchema(db.DB)

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, err := db.NewInsert().Model(&UserBun{
			Id:       fmt.Sprintf("bun-create-%d", i),
			Email:    fmt.Sprintf("bun-create-%d@example.com", i),
			PhoneNum: fmt.Sprintf("bun-create-phone-%d", i),
			Role:     "STUDENT",
		}).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchBunCreateMany(b *testing.B) {
	ctx := context.Background()
	db := openBun(b)
	defer db.Close()
	createSchema(db.DB)

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		users := make([]*UserBun, 10)
		for j := 0; j < 10; j++ {
			n := i*10 + j
			users[j] = &UserBun{
				Id:         fmt.Sprintf("bun-cmany-%d", n),
				Email:      fmt.Sprintf("bun-cmany-%d@example.com", n),
				PhoneNum:   fmt.Sprintf("bun-cmany-phone-%d", n),
				Role:       "STUDENT",
				LoginCount: 0,
			}
		}
		_, err := db.NewInsert().Model(&users).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchBunCreateManyAndReturn(b *testing.B) {
	ctx := context.Background()
	db := openBun(b)
	defer db.Close()
	createSchema(db.DB)

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		users := make([]*UserBun, 10)
		for j := 0; j < 10; j++ {
			n := i*10 + j
			users[j] = &UserBun{
				Id:         fmt.Sprintf("bun-cmar-%d", n),
				Email:      fmt.Sprintf("bun-cmar-%d@example.com", n),
				PhoneNum:   fmt.Sprintf("bun-cmar-phone-%d", n),
				Role:       "STUDENT",
				LoginCount: 0,
			}
		}
		_, err := db.NewInsert().Model(&users).Returning("*").Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchBunFindUnique(b *testing.B) {
	ctx := context.Background()
	db := openBun(b)
	defer db.Close()
	createSchema(db.DB)
	seedBun(db, "bun-fu")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		var u UserBun
		err := db.NewSelect().Model(&u).
			Where("email = ?", fmt.Sprintf("bun-fu-user-%d@example.com", i%seedCount)).
			Scan(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchBunFindFirst(b *testing.B) {
	ctx := context.Background()
	db := openBun(b)
	defer db.Close()
	createSchema(db.DB)
	seedBun(db, "bun-ff")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		var u UserBun
		err := db.NewSelect().Model(&u).
			Where("email = ?", fmt.Sprintf("bun-ff-user-%d@example.com", i%seedCount)).
			Limit(1).
			Scan(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchBunFindMany(b *testing.B) {
	ctx := context.Background()
	db := openBun(b)
	defer db.Close()
	createSchema(db.DB)
	seedBun(db, "bun-fm")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		var users []UserBun
		err := db.NewSelect().Model(&users).
			OrderExpr("id ASC").
			Limit(10).
			Offset(i % seedCount).
			Scan(ctx)
		if err != nil {
			b.Fatal(err)
		}
		if len(users) == 0 {
			b.Fatal("expected at least one user")
		}
	}
}

func benchBunUpsert(b *testing.B) {
	ctx := context.Background()
	db := openBun(b)
	defer db.Close()
	createSchema(db.DB)
	seedBun(db, "bun-ups")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		email := fmt.Sprintf("bun-ups-user-%d@example.com", i%seedCount)
		_, err := db.NewInsert().Model(&UserBun{
			Id:         fmt.Sprintf("bun-ups-id-%d", i),
			Email:      email,
			PhoneNum:   fmt.Sprintf("bun-ups-phone-new-%d", i),
			Role:       "STUDENT",
			LoginCount: int32(i),
		}).On("CONFLICT (email) DO UPDATE").
			Set(`"phoneNum" = EXCLUDED."phoneNum"`).
			Set(`"role" = EXCLUDED."role"`).
			Set(`"loginCount" = EXCLUDED."loginCount"`).
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchBunReadDeepRelation(b *testing.B) {
	ctx := context.Background()
	db := openBun(b)
	defer db.Close()
	createSchema(db.DB)
	seedRelations(db.DB, "bun-rdr")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		email := fmt.Sprintf("bun-rdr-grand-%d@example.com", i%500)
		var user UserBun
		err := db.NewSelect().Model(&user).
			Where(`"user_bun"."email" = ?`, email).
			Relation("ReferredBy.ReferredBy").
			Scan(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchBunCreateWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	db := openBun(b)
	defer db.Close()
	createSchema(db.DB)
	seedRelations(db.DB, "bun-cwds")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("bun-cwds-parent-id-%d", i%500)
		_, err := db.NewInsert().Model(&UserBun{
			Id:           fmt.Sprintf("bun-cwds-new-id-%d", i),
			Email:        fmt.Sprintf("bun-cwds-new-%d@example.com", i),
			PhoneNum:     fmt.Sprintf("bun-cwds-new-phone-%d", i),
			Role:         "STUDENT",
			ReferredById: &parentID,
		}).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
		var referredBy UserBun
		err = db.NewSelect().Model(&referredBy).
			Where(`"user_bun"."id" = ?`, parentID).
			Relation("ReferredBy").
			Scan(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchBunCreateManyAndReturnWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	db := openBun(b)
	defer db.Close()
	createSchema(db.DB)
	seedRelations(db.DB, "bun-cmwds")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("bun-cmwds-parent-id-%d", i%500)
		users := make([]UserBun, 10)
		ids := make([]string, 10)
		for j := range 10 {
			id := fmt.Sprintf("bun-cmwds-new-id-%d-%d", i, j)
			ids[j] = id
			users[j] = UserBun{
				Id:           id,
				Email:        fmt.Sprintf("bun-cmwds-new-%d-%d@example.com", i, j),
				PhoneNum:     fmt.Sprintf("bun-cmwds-new-phone-%d-%d", i, j),
				Role:         "STUDENT",
				ReferredById: &parentID,
			}
		}
		_, err := db.NewInsert().Model(&users).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
		var fetched []UserBun
		err = db.NewSelect().Model(&fetched).
			Where(`"user_bun"."id" IN (?)`, bun.In(ids)).
			Relation("ReferredBy.ReferredBy").
			Scan(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchBunUpsertWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	db := openBun(b)
	defer db.Close()
	createSchema(db.DB)
	seedRelations(db.DB, "bun-uwds")

	for i := range seedCount {
		_, err := db.NewInsert().Model(&UserBun{
			Id:       fmt.Sprintf("bun-uwds-id-%d", i),
			Email:    fmt.Sprintf("bun-uwds-%d@example.com", i),
			PhoneNum: fmt.Sprintf("bun-uwds-phone-%d", i),
			Role:     "STUDENT",
		}).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("bun-uwds-parent-id-%d", i%500)
		email := fmt.Sprintf("bun-uwds-%d@example.com", i%seedCount)
		_, err := db.NewInsert().Model(&UserBun{
			Id:           fmt.Sprintf("bun-uwds-id-new-%d", i),
			Email:        email,
			PhoneNum:     fmt.Sprintf("bun-uwds-phone-new-%d", i),
			Role:         "STUDENT",
			LoginCount:   int32(i),
			ReferredById: &parentID,
		}).On("CONFLICT (email) DO UPDATE").
			Set(`"phoneNum" = EXCLUDED."phoneNum"`).
			Set(`"role" = EXCLUDED."role"`).
			Set(`"loginCount" = EXCLUDED."loginCount"`).
			Set(`"referredById" = EXCLUDED."referredById"`).
			Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
		var fetched UserBun
		err = db.NewSelect().Model(&fetched).
			Where(`"user_bun"."email" = ?`, email).
			Relation("ReferredBy.ReferredBy").
			Scan(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// UserBunHooks is used in hook-only benchmarks to avoid polluting UserBun
type UserBunHooks struct {
	bun.BaseModel `bun:"table:User"`
	Id            string   `bun:"id,pk"`
	Email         string   `bun:"email,unique,notnull"`
	PhoneNum      string   `bun:"phoneNum,unique,notnull"`
	Password      *string  `bun:"password"`
	Role          string   `bun:"role,default:'STUDENT'"`
	RoleOptional  *string  `bun:"roleOptional"`
	LoginCount    int32    `bun:"loginCount,default:0"`
	ReferredById  *string  `bun:"referredById"`
}

func (u *UserBunHooks) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		u.Email = strings.ToLower(u.Email)
		if u.Password != nil {
			hash := sha256.Sum256([]byte(*u.Password))
			hexStr := hex.EncodeToString(hash[:])
			u.Password = &hexStr
		}
	case *bun.UpdateQuery:
		if u.Email != "" {
			u.Email = strings.ToLower(u.Email)
		}
	case *bun.DeleteQuery:
		u.Id = strings.ToUpper(u.Id)
	}
	return nil
}

func (u *UserBunHooks) AfterScanRow(ctx context.Context) error {
	u.Email += "-read"
	return nil
}

func benchBunHooksCreate(b *testing.B) {
	db := openBun(b)
	defer db.Close()
	createSchema(db.DB)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		pwd := "mypassword"
		u := UserBunHooks{
			Id:       fmt.Sprintf("bun-c-hook-%d", i),
			Email:    fmt.Sprintf("Bun-C-Hook-%d@Example.com", i),
			PhoneNum: fmt.Sprintf("bun-c-hook-phone-%d", i),
			Password: &pwd,
			Role:     "STUDENT",
		}
		_, err := db.NewInsert().Model(&u).Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchBunHooksUpdate(b *testing.B) {
	db := openBun(b)
	defer db.Close()
	createSchema(db.DB)
	ctx := context.Background()

	seedData(db.DB, "bun-u-hook")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		u := UserBunHooks{
			Id:    fmt.Sprintf("bun-u-hook-id-%d", i%seedCount),
			Email: fmt.Sprintf("NEW-BUN-U-HOOK-%d@EXAMPLE.COM", i),
		}
		_, err := db.NewUpdate().Model(&u).Column("email").WherePK().Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchBunHooksFindUnique(b *testing.B) {
	db := openBun(b)
	defer db.Close()
	createSchema(db.DB)
	ctx := context.Background()

	seedData(db.DB, "bun-f-hook")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		var fetched UserBunHooks
		err := db.NewSelect().Model(&fetched).Where("id = ?", fmt.Sprintf("bun-f-hook-id-%d", i%seedCount)).Scan(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchBunHooksDelete(b *testing.B) {
	db := openBun(b)
	defer db.Close()
	createSchema(db.DB)
	ctx := context.Background()

	seedDeleteData(db.DB, "bun", 50000)

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		u := UserBunHooks{Id: fmt.Sprintf("bun-d-hook-%d", i)}
		_, err := db.NewDelete().Model(&u).WherePK().Exec(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}
