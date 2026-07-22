package main

import (
	"context"
	"integration/valk"
	"integration/valk/user"
	"testing"
)

func TestDeleteMany(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.CreateMany(
		db.User.Create().SetEmail("delete1@example.com").SetPhoneNum("+111").SetRole(valk.UserRole.Student),
		db.User.Create().SetEmail("delete2@example.com").SetPhoneNum("+222").SetRole(valk.UserRole.Student),
		db.User.Create().SetEmail("keep1@example.com").SetPhoneNum("+333").SetRole(valk.UserRole.Admin),
	).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed users: %v", err)
	}

	count, err := db.User.DeleteMany(user.Role.EQ(valk.UserRole.Student)).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to delete students: %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 deleted users, got %d", count)
	}

	var remainingCount int
	err = db.Raw().QueryRowContext(ctx, query(
		`SELECT COUNT(*) FROM "User"`,
		`SELECT COUNT(*) FROM "User"`,
	)).Scan(&remainingCount)
	if err != nil {
		t.Fatalf("failed to count remaining users: %v", err)
	}

	if remainingCount != 1 {
		t.Errorf("expected 1 remaining user in DB, got %d", remainingCount)
	}

	var remainingEmail string
	err = db.Raw().QueryRowContext(ctx, query(
		`SELECT "email" FROM "User"`,
		`SELECT "email" FROM "User"`,
	)).Scan(&remainingEmail)
	if err != nil {
		t.Fatalf("failed to get remaining user: %v", err)
	}

	if remainingEmail != "keep1@example.com" {
		t.Errorf("expected remaining user to be keep1@example.com, got %s", remainingEmail)
	}
}

func TestDeleteMany_NoMatches(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.CreateMany(
		db.User.Create().SetEmail("keep1@example.com").SetPhoneNum("+333").SetRole(valk.UserRole.Admin),
	).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed users: %v", err)
	}

	count, err := db.User.DeleteMany(user.Role.EQ(valk.UserRole.Student)).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to delete students: %v", err)
	}

	if count != 0 {
		t.Errorf("expected 0 deleted users, got %d", count)
	}
}

func TestDeleteBasic(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u, err := db.User.Create().
		SetEmail("todelete@example.com").
		SetPhoneNum("+123456").
		SetRole(valk.UserRole.Student).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	deleted, err := db.User.Delete(user.Email.EQ("todelete@example.com")).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	if deleted.Email != "todelete@example.com" {
		t.Errorf("expected deleted user email to be todelete@example.com, got %s", deleted.Email)
	}
	if deleted.Id != u.Id {
		t.Errorf("expected deleted user id to be %s, got %s", u.Id, deleted.Id)
	}

	var count int
	err = db.Raw().QueryRowContext(ctx, `SELECT COUNT(*) FROM "User" WHERE email = 'todelete@example.com'`).Scan(&count)
	if err != nil {
		t.Fatalf("failed to query count: %v", err)
	}
	if count != 0 {
		t.Errorf("expected user to be completely deleted from DB, got count %d", count)
	}
}

func TestDeleteSelectOmit(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().
		SetEmail("selectdelete@example.com").
		SetPhoneNum("+55555").
		SetRole(valk.UserRole.Student).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	deleted, err := db.User.Delete(user.Email.EQ("selectdelete@example.com")).
		Select(valk.UserSelect{Email: true}).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	if deleted.Email != "selectdelete@example.com" {
		t.Errorf("expected email to be populated, got %s", deleted.Email)
	}
	if deleted.PhoneNum != "" {
		t.Errorf("expected phoneNum to be omitted (empty), got %s", deleted.PhoneNum)
	}
}

func TestDeleteWithRelations(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u, err := db.User.Create().
		SetEmail("reldelete@example.com").
		SetPhoneNum("+777").
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	_, err = db.Profile.Create().
		SetBio("My bio").
		SetUserId(u.Id).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create profile: %v", err)
	}

	deleted, err := db.User.Delete(user.Id.EQ(u.Id)).
		Select(valk.UserSelect{
			Email: true,
			Profile: &valk.ProfileSelect{
				Bio: true,
			},
		}).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	if deleted.Profile == nil {
		t.Errorf("expected related profile to be loaded, got nil")
	} else if deleted.Profile.Bio == nil || *deleted.Profile.Bio != "My bio" {
		t.Errorf("expected loaded profile bio to be 'My bio', got %v", deleted.Profile.Bio)
	}
}

func TestDeleteNotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Delete(user.Email.EQ("notfound@example.com")).Exec(ctx)
	if err == nil {
		t.Fatal("expected error when deleting non-existent row, got nil")
	}
}

func TestDeleteHooks(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	hookCalled := false
	db.User.Use(user.Extension{
		Delete: func(ctx context.Context, where valk.UniquePredicate[valk.User], selects *valk.UserSelect, omits *valk.UserOmit, next valk.UserDeleteQuery) (*valk.User, error) {
			hookCalled = true
			return next(ctx, where, selects, omits)
		},
	})

	_, err := db.User.Create().
		SetEmail("hookdelete@example.com").
		SetPhoneNum("+999").
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	_, err = db.User.Delete(user.Email.EQ("hookdelete@example.com")).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	if !hookCalled {
		t.Errorf("expected delete hook to be called, but it wasn't")
	}
}
