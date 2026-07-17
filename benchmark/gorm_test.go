package main

import (
	"context"
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserGORM struct {
	Id         string `gorm:"column:id;primaryKey"`
	Email      string `gorm:"column:email;uniqueIndex;not null"`
	PhoneNum   string `gorm:"column:phoneNum;uniqueIndex;not null"`
	Password   *string
	Role       string `gorm:"column:role;default:student"`
	LoginCount int32  `gorm:"column:loginCount;default:0"`
}

func (UserGORM) TableName() string {
	return "User"
}

func openGORM(b *testing.B) *gorm.DB {
	b.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		b.Fatal(err)
	}
	if err := db.AutoMigrate(&UserGORM{}); err != nil {
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
			Role:     "student",
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
			Role:     "student",
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
				Role:       "student",
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
				Role:       "student",
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
			Role:       "student",
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
