package main

import (
	"context"
	"integration/valk"
	"integration/valk/user"
	"testing"
)

func TestUpdate_SingleRow(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u1, err := db.User.Create().
		SetEmail("update1@example.com").
		SetPhoneNum("+10001").
		SetRole(valk.UserRole.Student).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	updated, err := db.User.Update(user.Id.EQ(u1.Id)).
		SetEmail("updated1@example.com").
		SetRole(valk.UserRole.Admin).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	if updated.Email != "updated1@example.com" {
		t.Errorf("expected email updated1@example.com, got %s", updated.Email)
	}
	if updated.Role != valk.UserRole.Admin {
		t.Errorf("expected role Admin, got %s", updated.Role)
	}

	// Verify in DB
	fetched, err := db.User.FindUnique(user.Id.EQ(u1.Id)).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to fetch user from DB: %v", err)
	}
	if fetched.Email != "updated1@example.com" {
		t.Errorf("expected DB email updated1@example.com, got %s", fetched.Email)
	}
}

func TestUpdate_WithAdditionalPredicates(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u1, err := db.User.Create().
		SetEmail("add_pred@example.com").
		SetPhoneNum("+10009").
		SetRole(valk.UserRole.Student).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Update with UniquePredicate + Additional Predicate matching
	updated, err := db.User.Update(user.Id.EQ(u1.Id), user.Role.EQ(valk.UserRole.Student)).
		SetEmail("updated_add_pred@example.com").
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to update user with additional predicate: %v", err)
	}
	if updated.Email != "updated_add_pred@example.com" {
		t.Errorf("expected email updated_add_pred@example.com, got %s", updated.Email)
	}

	// Update with UniquePredicate + Additional Predicate NOT matching (should fail with sql.ErrNoRows or return nil)
	_, err = db.User.Update(user.Id.EQ(u1.Id), user.Role.EQ(valk.UserRole.Admin)).
		SetEmail("should_not_update@example.com").
		Exec(ctx)
	if err == nil {
		t.Fatalf("expected error updating when additional predicate does not match")
	}
}

func TestUpdate_WithNestedSelect(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u1, err := db.User.Create().
		SetEmail("update_select@example.com").
		SetPhoneNum("+10002").
		SetRole(valk.UserRole.Student).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	_, err = db.Post.Create().
		SetTitle("Post Title 1").
		SetContent("Content 1").
		SetAuthorId(u1.Id).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	updated, err := db.User.Update(user.Id.EQ(u1.Id)).
		SetEmail("new_email_select@example.com").
		Select(valk.UserSelect{
			Email: true,
			Posts: &valk.PostSelect{},
		}).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to update user with select: %v", err)
	}

	if updated.Email != "new_email_select@example.com" {
		t.Errorf("expected updated email new_email_select@example.com, got %s", updated.Email)
	}
	if len(updated.Posts) != 1 {
		t.Fatalf("expected 1 post loaded on updated user, got %d", len(updated.Posts))
	}
	if updated.Posts[0].Title != "Post Title 1" {
		t.Errorf("expected post title 'Post Title 1', got %s", updated.Posts[0].Title)
	}
}

func TestUpdateMany_Count(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.CreateMany(
		db.User.Create().SetEmail("many1@example.com").SetPhoneNum("+20001").SetRole(valk.UserRole.Student).SetLoginCount(1),
		db.User.Create().SetEmail("many2@example.com").SetPhoneNum("+20002").SetRole(valk.UserRole.Student).SetLoginCount(1),
		db.User.Create().SetEmail("many3@example.com").SetPhoneNum("+20003").SetRole(valk.UserRole.Admin).SetLoginCount(1),
	).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create users: %v", err)
	}

	count, err := db.User.UpdateMany(user.Role.EQ(valk.UserRole.Student)).
		SetLoginCount(50).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to updateMany students: %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 affected rows, got %d", count)
	}

	students, err := db.User.FindMany(user.Role.EQ(valk.UserRole.Student)).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to fetch students: %v", err)
	}
	for _, s := range students {
		if s.LoginCount != 50 {
			t.Errorf("expected loginCount 50 for student %s, got %d", s.Email, s.LoginCount)
		}
	}
}

func TestUpdateManyAndReturn_Basic(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.CreateMany(
		db.User.Create().SetEmail("ret1@example.com").SetPhoneNum("+30001").SetRole(valk.UserRole.Student).SetLoginCount(5),
		db.User.Create().SetEmail("ret2@example.com").SetPhoneNum("+30002").SetRole(valk.UserRole.Student).SetLoginCount(5),
	).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed users: %v", err)
	}

	updatedUsers, err := db.User.UpdateManyAndReturn(user.Role.EQ(valk.UserRole.Student)).
		SetLoginCount(99).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to updateManyAndReturn: %v", err)
	}

	if len(updatedUsers) != 2 {
		t.Fatalf("expected 2 returned updated users, got %d", len(updatedUsers))
	}
	for _, u := range updatedUsers {
		if u.LoginCount != 99 {
			t.Errorf("expected loginCount 99 for %s, got %d", u.Email, u.LoginCount)
		}
	}
}

func TestUpdateManyAndReturn_WithRelations(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u1, err := db.User.Create().
		SetEmail("ret_rel1@example.com").
		SetPhoneNum("+40001").
		SetRole(valk.UserRole.Student).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	_, err = db.Post.Create().
		SetTitle("Student Post").
		SetContent("Hello").
		SetAuthorId(u1.Id).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	updatedUsers, err := db.User.UpdateManyAndReturn(user.Id.EQ(u1.Id)).
		SetLoginCount(777).
		Select(valk.UserSelect{
			LoginCount: true,
			Posts:      &valk.PostSelect{},
		}).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to updateManyAndReturn with relations: %v", err)
	}

	if len(updatedUsers) != 1 {
		t.Fatalf("expected 1 returned updated user, got %d", len(updatedUsers))
	}
	if updatedUsers[0].LoginCount != 777 {
		t.Errorf("expected loginCount 777, got %d", updatedUsers[0].LoginCount)
	}
	if len(updatedUsers[0].Posts) != 1 {
		t.Fatalf("expected 1 post loaded on returned user, got %d", len(updatedUsers[0].Posts))
	}
	if updatedUsers[0].Posts[0].Title != "Student Post" {
		t.Errorf("expected post title 'Student Post', got %s", updatedUsers[0].Posts[0].Title)
	}
}

func TestUpdate_NoAssignments(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u, err := db.User.Create().
		SetEmail("no_assign@example.com").
		SetPhoneNum("+50001").
		SetRole(valk.UserRole.Student).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	_, err = db.Post.Create().
		SetTitle("Existing Post").
		SetContent("Hello").
		SetAuthorId(u.Id).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Update with 0 assignments should just return current record + relations
	updated, err := db.User.Update(user.Id.EQ(u.Id)).
		Select(valk.UserSelect{
			Email: true,
			Posts: &valk.PostSelect{},
		}).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to update with 0 assignments: %v", err)
	}

	if updated.Email != "no_assign@example.com" {
		t.Errorf("expected email no_assign@example.com, got %s", updated.Email)
	}
	if len(updated.Posts) != 1 {
		t.Errorf("expected 1 post, got %d", len(updated.Posts))
	}
}

func TestUpdate_NonExistentRow(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Update(user.Id.EQ("non-existent-id-999")).
		SetEmail("nobody@example.com").
		Exec(ctx)
	if err == nil {
		t.Fatalf("expected error updating non-existent row, got nil")
	}
}

func TestUpdate_InsideUserTransaction(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u, err := db.User.Create().
		SetEmail("in_tx_base@example.com").
		SetPhoneNum("+60001").
		SetRole(valk.UserRole.Student).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	_, err = db.Post.Create().
		SetTitle("Tx Post").
		SetContent("In Tx").
		SetAuthorId(u.Id).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	err = db.Transaction(ctx, func(tx *valk.Tx) error {
		updated, err := tx.User.Update(user.Id.EQ(u.Id)).
			SetEmail("in_tx_updated@example.com").
			Select(valk.UserSelect{
				Email: true,
				Posts: &valk.PostSelect{},
			}).
			Exec(ctx)
		if err != nil {
			return err
		}
		if updated.Email != "in_tx_updated@example.com" {
			t.Errorf("expected updated email in tx, got %s", updated.Email)
		}
		if len(updated.Posts) != 1 {
			t.Errorf("expected 1 post in tx, got %d", len(updated.Posts))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("transaction failed: %v", err)
	}
}

func TestUpdateMany_EdgeCases(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// 0 assignments -> returns 0, nil
	count, err := db.User.UpdateMany(user.Role.EQ(valk.UserRole.Student)).Exec(ctx)
	if err != nil {
		t.Fatalf("updateMany with 0 assignments failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 count for 0 assignments, got %d", count)
	}

	// 0 matching rows -> returns 0, nil
	count, err = db.User.UpdateMany(user.Email.EQ("nonexistent_many@example.com")).
		SetLoginCount(100).
		Exec(ctx)
	if err != nil {
		t.Fatalf("updateMany with 0 matching rows failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 count for 0 matching rows, got %d", count)
	}
}

func TestUpdateManyAndReturn_EdgeCases(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// 0 matching rows -> returns empty slice
	results, err := db.User.UpdateManyAndReturn(user.Email.EQ("nonexistent_ret@example.com")).
		SetLoginCount(100).
		Exec(ctx)
	if err != nil {
		t.Fatalf("updateManyAndReturn 0 matches failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty slice for 0 matches, got %d", len(results))
	}
}
