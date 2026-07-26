package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserGORM struct {
	Id           string        `gorm:"column:id;primaryKey"`
	Email        string        `gorm:"column:email;uniqueIndex;not null"`
	PhoneNum     string        `gorm:"column:phoneNum;uniqueIndex;not null"`
	Password     *string       `gorm:"column:password"`
	Role         string        `gorm:"column:role;default:student"`
	RoleOptional *string       `gorm:"column:roleOptional"`
	Profile      *ProfileGORM  `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE"`
	Posts        []PostGORM    `gorm:"foreignKey:AuthorId;references:Id"`
	Comments     []CommentGORM `gorm:"foreignKey:AuthorId;references:Id"`
	LoginCount   int32         `gorm:"column:loginCount;default:0"`
	ReferredById *string       `gorm:"column:referredById"`
	ReferredBy   *UserGORM     `gorm:"foreignKey:ReferredById;references:Id"`
	Referrals    []UserGORM    `gorm:"foreignKey:ReferredById;references:Id"`
}

func (UserGORM) TableName() string {
	return "User"
}

type ProfileGORM struct {
	Id        string    `gorm:"column:id;primaryKey"`
	Bio       *string   `gorm:"column:bio"`
	UserId    string    `gorm:"column:userId;uniqueIndex;not null"`
	User      *UserGORM `gorm:"foreignKey:UserId;references:Id"`
	CreatedAt time.Time `gorm:"column:createdAt;default:CURRENT_TIMESTAMP"`
}

func (ProfileGORM) TableName() string {
	return "Profile"
}

type PostGORM struct {
	Id         string               `gorm:"column:id;primaryKey"`
	Title      string               `gorm:"column:title;not null"`
	Content    *string              `gorm:"column:content"`
	Published  bool                 `gorm:"column:published;default:false"`
	AuthorId   string               `gorm:"column:authorId;not null"`
	Author     *UserGORM            `gorm:"foreignKey:AuthorId;references:Id"`
	Comments   []CommentGORM        `gorm:"foreignKey:PostId;references:Id"`
	Categories []CategoryToPostGORM `gorm:"foreignKey:PostId;references:Id"`
}

func (PostGORM) TableName() string {
	return "Post"
}

type CommentGORM struct {
	Id       string    `gorm:"column:id;primaryKey"`
	Textify  int       `gorm:"column:textify;not null"`
	Dummy3   string    `gorm:"column:dummy3;not null"`
	Dummy1   int       `gorm:"column:dummy1;not null"`
	Dummy2   string    `gorm:"column:dummy2;not null"`
	PostId   string    `gorm:"column:postId;not null"`
	Post     *PostGORM `gorm:"foreignKey:PostId;references:Id"`
	AuthorId string    `gorm:"column:authorId;not null"`
	Author   *UserGORM `gorm:"foreignKey:AuthorId;references:Id"`
	Meta     *string   `gorm:"column:meta"`
}

func (CommentGORM) TableName() string {
	return "Comment"
}

type CategoryGORM struct {
	Id    int                  `gorm:"column:id;primaryKey;autoIncrement"`
	Name  string               `gorm:"column:name;uniqueIndex;not null"`
	Posts []CategoryToPostGORM `gorm:"foreignKey:CategoryId;references:Id"`
}

func (CategoryGORM) TableName() string {
	return "Category"
}

type CategoryToPostGORM struct {
	Id         int           `gorm:"column:id;primaryKey;autoIncrement"`
	PostId     string        `gorm:"column:postId;uniqueIndex:idx_post_category,priority:1;not null"`
	CategoryId int           `gorm:"column:categoryId;uniqueIndex:idx_post_category,priority:2;not null"`
	Post       *PostGORM     `gorm:"foreignKey:PostId;references:Id"`
	Category   *CategoryGORM `gorm:"foreignKey:CategoryId;references:Id"`
}

func (CategoryToPostGORM) TableName() string {
	return "CategoryToPost"
}

func openGORM(b *testing.B) *gorm.DB {
	b.Helper()
	var dialector gorm.Dialector
	if activeDialect.Name == "postgres" {
		dialector = postgres.Open(activeDialect.DSN)
	} else {
		dialector = sqlite.Open(activeDialect.DSN)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		b.Fatal(err)
	}

	if activeDialect.Name == "postgres" {
		sqlDB, err := db.DB()
		if err == nil {
			resetPostgres(sqlDB)
		}
	}

	if err := db.AutoMigrate(&UserGORM{}, &ProfileGORM{}, &PostGORM{}, &CommentGORM{}, &CategoryGORM{}, &CategoryToPostGORM{}); err != nil {
		b.Fatal(err)
	}
	return db
}

func seedGORM(db *gorm.DB, prefix string) {
	for i := range seedCount {
		db.Create(&UserGORM{
			Id:       fmt.Sprintf("%s-id-%d", prefix, i),
			Email:    fmt.Sprintf("%s-user-%d@example.com", prefix, i),
			PhoneNum: fmt.Sprintf("%s-phone-%d", prefix, i),
			Role:     "STUDENT",
		})
	}
}

func benchGORMCreate(b *testing.B) {
	ctx := context.Background()
	db := openGORM(b)
	sqlDB, err := db.DB()
	if err != nil {
		b.Fatal(err)
	}
	defer sqlDB.Close()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		user := UserGORM{
			Id:       fmt.Sprintf("gorm-create-%d", i),
			Email:    fmt.Sprintf("gorm-create-%d@example.com", i),
			PhoneNum: fmt.Sprintf("gorm-create-phone-%d", i),
			Role:     "STUDENT",
		}
		if err := db.WithContext(ctx).Create(&user).Error; err != nil {
			b.Fatal(err)
		}
	}
}

func benchGORMCreateMany(b *testing.B) {
	ctx := context.Background()
	db := openGORM(b)
	sqlDB, err := db.DB()
	if err != nil {
		b.Fatal(err)
	}
	defer sqlDB.Close()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		users := make([]UserGORM, 10)
		for j := 0; j < 10; j++ {
			n := i*10 + j
			users[j] = UserGORM{
				Id:         fmt.Sprintf("gorm-cmany-%d", n),
				Email:      fmt.Sprintf("gorm-cmany-%d@example.com", n),
				PhoneNum:   fmt.Sprintf("gorm-cmany-phone-%d", n),
				Role:       "STUDENT",
				LoginCount: 0,
			}
		}
		if err := db.WithContext(ctx).Create(&users).Error; err != nil {
			b.Fatal(err)
		}
	}
}

func benchGORMCreateManyAndReturn(b *testing.B) {
	ctx := context.Background()
	db := openGORM(b)
	sqlDB, err := db.DB()
	if err != nil {
		b.Fatal(err)
	}
	defer sqlDB.Close()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		users := make([]UserGORM, 10)
		for j := 0; j < 10; j++ {
			n := i*10 + j
			users[j] = UserGORM{
				Id:         fmt.Sprintf("gorm-cmar-%d", n),
				Email:      fmt.Sprintf("gorm-cmar-%d@example.com", n),
				PhoneNum:   fmt.Sprintf("gorm-cmar-phone-%d", n),
				Role:       "STUDENT",
				LoginCount: 0,
			}
		}
		if err := db.WithContext(ctx).Create(&users).Error; err != nil {
			b.Fatal(err)
		}
		for _, u := range users {
			if u.Id == "" {
				b.Fatal("expected returned IDs")
			}
		}
	}
}

func benchGORMFindUnique(b *testing.B) {
	ctx := context.Background()
	db := openGORM(b)
	sqlDB, err := db.DB()
	if err != nil {
		b.Fatal(err)
	}
	defer sqlDB.Close()
	seedGORM(db, "gorm-fu")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		var user UserGORM
		email := fmt.Sprintf("gorm-fu-user-%d@example.com", i%seedCount)
		if err := db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
			b.Fatal(err)
		}
	}
}

func benchGORMFindFirst(b *testing.B) {
	ctx := context.Background()
	db := openGORM(b)
	sqlDB, err := db.DB()
	if err != nil {
		b.Fatal(err)
	}
	defer sqlDB.Close()
	seedGORM(db, "gorm-ff")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		var user UserGORM
		email := fmt.Sprintf("gorm-ff-user-%d@example.com", i%seedCount)
		if err := db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
			b.Fatal(err)
		}
	}
}

func benchGORMFindMany(b *testing.B) {
	ctx := context.Background()
	db := openGORM(b)
	sqlDB, err := db.DB()
	if err != nil {
		b.Fatal(err)
	}
	defer sqlDB.Close()
	seedGORM(db, "gorm-fm")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		var users []UserGORM
		if err := db.WithContext(ctx).Order("id").Limit(10).Offset(i % seedCount).Find(&users).Error; err != nil {
			b.Fatal(err)
		}
		if len(users) == 0 {
			b.Fatal("expected at least one user")
		}
	}
}

func benchGORMUpsert(b *testing.B) {
	ctx := context.Background()
	db := openGORM(b)
	sqlDB, err := db.DB()
	if err != nil {
		b.Fatal(err)
	}
	defer sqlDB.Close()
	seedGORM(db, "gorm-ups")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		email := fmt.Sprintf("gorm-ups-user-%d@example.com", i%seedCount)
		user := UserGORM{
			Id:         fmt.Sprintf("gorm-ups-id-%d", i),
			Email:      email,
			PhoneNum:   fmt.Sprintf("gorm-ups-phone-new-%d", i),
			Role:       "STUDENT",
			LoginCount: int32(i),
		}
		if err := db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "email"}},
			DoUpdates: clause.AssignmentColumns([]string{"phoneNum", "role", "loginCount"}),
		}).Create(&user).Error; err != nil {
			b.Fatal(err)
		}
	}
}

func benchGORMReadDeepRelation(b *testing.B) {
	ctx := context.Background()
	db := openGORM(b)
	sqlDB, err := db.DB()
	if err != nil {
		b.Fatal(err)
	}
	defer sqlDB.Close()
	seedRelations(sqlDB, "gorm-rdr")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		var user UserGORM
		email := fmt.Sprintf("gorm-rdr-grand-%d@example.com", i%500)
		if err := db.WithContext(ctx).Preload("ReferredBy.ReferredBy").Where("email = ?", email).First(&user).Error; err != nil {
			b.Fatal(err)
		}
	}
}

func benchGORMCreateWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	db := openGORM(b)
	sqlDB, err := db.DB()
	if err != nil {
		b.Fatal(err)
	}
	defer sqlDB.Close()
	seedRelations(sqlDB, "gorm-cwds")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("gorm-cwds-parent-id-%d", i%500)
		user := UserGORM{
			Id:           fmt.Sprintf("gorm-cwds-new-id-%d", i),
			Email:        fmt.Sprintf("gorm-cwds-new-%d@example.com", i),
			PhoneNum:     fmt.Sprintf("gorm-cwds-new-phone-%d", i),
			Role:         "STUDENT",
			ReferredById: &parentID,
		}
		if err := db.WithContext(ctx).Create(&user).Error; err != nil {
			b.Fatal(err)
		}
		var referredBy UserGORM
		if err := db.WithContext(ctx).Preload("ReferredBy").Where("id = ?", parentID).First(&referredBy).Error; err != nil {
			b.Fatal(err)
		}
		user.ReferredBy = &referredBy
	}
}

func benchGORMCreateManyAndReturnWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	db := openGORM(b)
	sqlDB, err := db.DB()
	if err != nil {
		b.Fatal(err)
	}
	defer sqlDB.Close()
	seedRelations(sqlDB, "gorm-cmwds")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("gorm-cmwds-parent-id-%d", i%500)
		users := make([]UserGORM, 10)
		ids := make([]string, 10)
		for j := 0; j < 10; j++ {
			id := fmt.Sprintf("gorm-cmwds-new-id-%d-%d", i, j)
			ids[j] = id
			users[j] = UserGORM{
				Id:           id,
				Email:        fmt.Sprintf("gorm-cmwds-new-%d-%d@example.com", i, j),
				PhoneNum:     fmt.Sprintf("gorm-cmwds-new-phone-%d-%d", i, j),
				Role:         "STUDENT",
				ReferredById: &parentID,
			}
		}
		if err := db.WithContext(ctx).Create(&users).Error; err != nil {
			b.Fatal(err)
		}
		var fetched []UserGORM
		if err := db.WithContext(ctx).Preload("ReferredBy.ReferredBy").Where("id IN ?", ids).Find(&fetched).Error; err != nil {
			b.Fatal(err)
		}
	}
}

func benchGORMUpsertWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	db := openGORM(b)
	sqlDB, err := db.DB()
	if err != nil {
		b.Fatal(err)
	}
	defer sqlDB.Close()
	seedRelations(sqlDB, "gorm-uwds")

	// Preseed
	for i := 0; i < seedCount; i++ {
		user := UserGORM{
			Id:       fmt.Sprintf("gorm-uwds-id-%d", i),
			Email:    fmt.Sprintf("gorm-uwds-%d@example.com", i),
			PhoneNum: fmt.Sprintf("gorm-uwds-phone-%d", i),
			Role:     "STUDENT",
		}
		if err := db.WithContext(ctx).Create(&user).Error; err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("gorm-uwds-parent-id-%d", i%500)
		email := fmt.Sprintf("gorm-uwds-%d@example.com", i%seedCount)
		user := UserGORM{
			Id:           fmt.Sprintf("gorm-uwds-id-new-%d", i),
			Email:        email,
			PhoneNum:     fmt.Sprintf("gorm-uwds-phone-new-%d", i),
			Role:         "STUDENT",
			LoginCount:   int32(i),
			ReferredById: &parentID,
		}
		if err := db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "email"}},
			DoUpdates: clause.AssignmentColumns([]string{"phoneNum", "role", "loginCount", "referredById"}),
		}).Create(&user).Error; err != nil {
			b.Fatal(err)
		}
		var fetched UserGORM
		if err := db.WithContext(ctx).Preload("ReferredBy.ReferredBy").Where("email = ?", email).First(&fetched).Error; err != nil {
			b.Fatal(err)
		}
	}
}

func benchGORMHooksCreate(b *testing.B) {
	db := openGORM(b)
	db.Callback().Create().Before("gorm:create").Register("bench_hook_create", func(tx *gorm.DB) {
		if u, ok := tx.Statement.Dest.(*UserGORM); ok {
			u.Email = strings.ToLower(u.Email)
			if u.Password != nil {
				hash := sha256.Sum256([]byte(*u.Password))
				hexStr := hex.EncodeToString(hash[:])
				u.Password = &hexStr
			}
		}
	})
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		pwd := "mypassword"
		u := UserGORM{
			Id:       fmt.Sprintf("gorm-c-hook-%d", i),
			Email:    fmt.Sprintf("Gorm-C-Hook-%d@Example.com", i),
			PhoneNum: fmt.Sprintf("gorm-c-hook-phone-%d", i),
			Password: &pwd,
			Role:     "STUDENT",
		}
		if err := db.WithContext(ctx).Create(&u).Error; err != nil {
			b.Fatal(err)
		}
	}
}

func benchGORMHooksUpdate(b *testing.B) {
	db := openGORM(b)
	db.Callback().Update().Before("gorm:update").Register("bench_hook_update", func(tx *gorm.DB) {
		if u, ok := tx.Statement.Dest.(*UserGORM); ok && u.Email != "" {
			u.Email = strings.ToLower(u.Email)
		}
	})
	ctx := context.Background()

	sqlDB, _ := db.DB()
	seedData(sqlDB, "gorm-u-hook")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		u := UserGORM{
			Email: fmt.Sprintf("NEW-GORM-U-HOOK-%d@EXAMPLE.COM", i),
		}
		if err := db.WithContext(ctx).Model(&UserGORM{}).Where("id = ?", fmt.Sprintf("gorm-u-hook-id-%d", i%seedCount)).Updates(&u).Error; err != nil {
			b.Fatal(err)
		}
	}
}

func benchGORMHooksFindUnique(b *testing.B) {
	db := openGORM(b)
	db.Callback().Query().After("gorm:query").Register("bench_hook_find", func(tx *gorm.DB) {
		if u, ok := tx.Statement.Dest.(*UserGORM); ok {
			u.Email += "-read"
		}
	})
	ctx := context.Background()

	sqlDB, _ := db.DB()
	seedData(sqlDB, "gorm-f-hook")

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		var fetched UserGORM
		if err := db.WithContext(ctx).Where("id = ?", fmt.Sprintf("gorm-f-hook-id-%d", i%seedCount)).First(&fetched).Error; err != nil {
			b.Fatal(err)
		}
	}
}

func benchGORMHooksDelete(b *testing.B) {
	db := openGORM(b)
	db.Callback().Delete().Before("gorm:delete").Register("bench_hook_delete", func(tx *gorm.DB) {
		if u, ok := tx.Statement.Dest.(*UserGORM); ok {
			u.Id = strings.ToUpper(u.Id)
		}
	})
	ctx := context.Background()

	sqlDB, _ := db.DB()
	seedDeleteData(sqlDB, "gorm", 50000)

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		if err := db.WithContext(ctx).Where("id = ?", fmt.Sprintf("gorm-d-hook-%d", i)).Delete(&UserGORM{Id: fmt.Sprintf("gorm-d-hook-%d", i)}).Error; err != nil {
			b.Fatal(err)
		}
	}
}
