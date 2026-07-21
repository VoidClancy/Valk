package main

import (
	"context"
	"integration/valk"
	"integration/valk/user"
	"testing"
)

func TestCountBasic(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// Initial count should be 0
	cnt, err := db.User.Count().Exec(ctx)
	if err != nil {
		t.Fatalf("failed to count users: %v", err)
	}
	if cnt != 0 {
		t.Errorf("expected 0 users, got %d", cnt)
	}

	// Seed some users
	_, err = db.User.CreateMany(
		db.User.Create().SetEmail("student1@example.com").SetPhoneNum("+111").SetRole(valk.UserRole.Student),
		db.User.Create().SetEmail("student2@example.com").SetPhoneNum("+222").SetRole(valk.UserRole.Student),
		db.User.Create().SetEmail("teacher1@example.com").SetPhoneNum("+333").SetRole(valk.UserRole.Teacher),
	).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed users: %v", err)
	}

	// Total count
	cnt, err = db.User.Count().Exec(ctx)
	if err != nil {
		t.Fatalf("failed to count users: %v", err)
	}
	if cnt != 3 {
		t.Errorf("expected 3 users, got %d", cnt)
	}

	// Count with filters
	studentsCount, err := db.User.Count(user.Role.EQ(valk.UserRole.Student)).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to count students: %v", err)
	}
	if studentsCount != 2 {
		t.Errorf("expected 2 student users, got %d", studentsCount)
	}

	teachersCount, err := db.User.Count(user.Role.EQ(valk.UserRole.Teacher)).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to count teachers: %v", err)
	}
	if teachersCount != 1 {
		t.Errorf("expected 1 teacher user, got %d", teachersCount)
	}
}

func TestCountLimitOffset(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// Seed 5 users
	var builders []*user.CreateBuilder
	for i := 1; i <= 5; i++ {
		builders = append(builders, db.User.Create().
			SetEmail(string(rune('a'+i))+"@example.com").
			SetPhoneNum("+"+string(rune('0'+i))),
		)
	}
	_, err := db.User.CreateMany(builders...).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed users: %v", err)
	}

	// Count with Limit only (Take)
	cnt, err := db.User.Count().Take(3).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to count with limit: %v", err)
	}
	if cnt != 3 {
		t.Errorf("expected 3 users, got %d", cnt)
	}

	// Count with Offset only (Skip)
	cnt, err = db.User.Count().Skip(2).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to count with offset: %v", err)
	}
	if cnt != 3 {
		t.Errorf("expected 3 users (5 total - 2 skip), got %d", cnt)
	}

	// Count with Limit and Offset (Take and Skip)
	cnt, err = db.User.Count().Take(2).Skip(2).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to count with limit/offset: %v", err)
	}
	if cnt != 2 {
		t.Errorf("expected 2 users (5 total - 2 skip, limited to 2), got %d", cnt)
	}
}

func TestCountHooks(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// Register count extension/hook that multiplies count by 2
	db.User.Use(user.Extension{
		Count: func(ctx context.Context, params valk.QueryParams[valk.User], next valk.UserCountQuery) (int64, error) {
			res, err := next(ctx, params)
			if err != nil {
				return 0, err
			}
			return res * 2, nil
		},
	})

	// Seed 2 users
	_, err := db.User.CreateMany(
		db.User.Create().SetEmail("u1@example.com").SetPhoneNum("+111"),
		db.User.Create().SetEmail("u2@example.com").SetPhoneNum("+222"),
	).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed users: %v", err)
	}

	cnt, err := db.User.Count().Exec(ctx)
	if err != nil {
		t.Fatalf("failed to count: %v", err)
	}

	// Hook should double 2 to 4
	if cnt != 4 {
		t.Errorf("expected count to be 4 (doubled by hook), got %d", cnt)
	}
}
