package main

import (
	"context"
	"integration/valk"
	"integration/valk/categoryToPost"
	"integration/valk/user"
	"testing"
)

func TestUpsert_OnConflict(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateMany with OnConflict(user.Email).Ignore()", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := client.User.Create().SetEmail("ignore@example.com").SetPhoneNum("+111").SetLoginCount(10).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to insert initial: %v", err)
		}

		affected, err := client.User.CreateMany(
			client.User.Create().SetEmail("ignore@example.com").SetPhoneNum("+222").SetLoginCount(20),
			client.User.Create().SetEmail("new@example.com").SetPhoneNum("+333").SetLoginCount(30),
		).OnConflict(user.Email).Ignore().Exec(ctx)

		if err != nil {
			t.Fatalf("CreateMany OnConflict.Ignore failed: %v", err)
		}
		if affected != 1 {
			t.Errorf("expected 1 row affected (1 inserted, 1 ignored), got %d", affected)
		}

		dbUser, err := client.User.FindUnique(user.Email.EQ("ignore@example.com")).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to query ignore@example.com: %v", err)
		}
		if dbUser.LoginCount != 10 {
			t.Errorf("expected loginCount to remain 10, got %d", dbUser.LoginCount)
		}
		if dbUser.PhoneNum != "+111" {
			t.Errorf("expected phoneNum to remain +111, got %q", dbUser.PhoneNum)
		}

		newDbUser, err := client.User.FindUnique(user.Email.EQ("new@example.com")).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to query new@example.com: %v", err)
		}
		if newDbUser.LoginCount != 30 {
			t.Errorf("expected new user loginCount 30, got %d", newDbUser.LoginCount)
		}
	})

	t.Run("CreateMany with OnConflict(user.Email).UpdateNewValues()", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := client.User.Create().SetEmail("update_new@example.com").SetPhoneNum("+111").SetLoginCount(10).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to insert initial: %v", err)
		}

		affected, err := client.User.CreateMany(
			client.User.Create().SetEmail("update_new@example.com").SetPhoneNum("+222").SetLoginCount(20),
		).OnConflict(user.Email).UpdateNewValues().Exec(ctx)

		if err != nil {
			t.Fatalf("CreateMany OnConflict.UpdateNewValues failed: %v", err)
		}
		if affected != 1 {
			t.Errorf("expected 1 row affected/updated, got %d", affected)
		}

		dbUser, err := client.User.FindUnique(user.Email.EQ("update_new@example.com")).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to query: %v", err)
		}
		if dbUser.LoginCount != 20 {
			t.Errorf("expected loginCount to be updated to 20, got %d", dbUser.LoginCount)
		}
		if dbUser.PhoneNum != "+222" {
			t.Errorf("expected phoneNum to be updated to +222, got %q", dbUser.PhoneNum)
		}
	})

	t.Run("CreateMany with OnConflict(user.Email).Update(custom)", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := client.User.Create().SetEmail("custom@example.com").SetPhoneNum("+111").SetLoginCount(10).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to insert initial: %v", err)
		}

		affected, err := client.User.CreateMany(
			client.User.Create().SetEmail("custom@example.com").SetPhoneNum("+222").SetLoginCount(20),
		).OnConflict(user.Email).Update(func(u *valk.UserUpsert) {
			u.PhoneNum.Update()
			u.LoginCount.Add(5)
		}).Exec(ctx)

		if err != nil {
			t.Fatalf("CreateMany OnConflict.Update failed: %v", err)
		}
		if affected != 1 {
			t.Errorf("expected 1 row affected, got %d", affected)
		}

		dbUser, err := client.User.FindUnique(user.Email.EQ("custom@example.com")).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to query: %v", err)
		}
		if dbUser.PhoneNum != "+222" {
			t.Errorf("expected phoneNum to be updated to +222, got %q", dbUser.PhoneNum)
		}
		if dbUser.LoginCount != 15 {
			t.Errorf("expected loginCount to be 10 + 5 = 15, got %d", dbUser.LoginCount)
		}
	})

	t.Run("CreateMany with OnConflict(user.EmailPhone).UpdateNewValues()", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := client.User.Create().SetEmail("composite@example.com").SetPhoneNum("+888").SetLoginCount(10).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to insert initial: %v", err)
		}

		affected, err := client.User.CreateMany(
			client.User.Create().SetEmail("composite@example.com").SetPhoneNum("+888").SetLoginCount(25),
		).OnConflict(user.EmailPhone).UpdateNewValues().Exec(ctx)

		if err != nil {
			t.Fatalf("CreateMany OnConflict(user.EmailPhone) failed: %v", err)
		}
		if affected != 1 {
			t.Errorf("expected 1 row affected/updated, got %d", affected)
		}

		dbUser, err := client.User.FindUnique(user.Email.EQ("composite@example.com")).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to query: %v", err)
		}
		if dbUser.LoginCount != 25 {
			t.Errorf("expected loginCount to be updated to 25, got %d", dbUser.LoginCount)
		}
	})

	t.Run("Single Create with OnConflict.Update()", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := client.User.Create().SetEmail("single@example.com").SetPhoneNum("+123").SetLoginCount(5).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to insert initial: %v", err)
		}

		res, err := client.User.Create().
			SetEmail("single@example.com").
			SetPhoneNum("+456").
			SetLoginCount(10).
			OnConflict(user.Email).
			Update(func(u *valk.UserUpsert) {
				u.PhoneNum.Update()
				u.LoginCount.Add(20)
			}).
			Exec(ctx)

		if err != nil {
			t.Fatalf("single upsert failed: %v", err)
		}
		if res.PhoneNum != "+456" {
			t.Errorf("expected updated phoneNum +456, got %q", res.PhoneNum)
		}
		if res.LoginCount != 25 {
			t.Errorf("expected updated loginCount 25 (5 + 20), got %d", res.LoginCount)
		}
	})

	t.Run("Single Create with OnConflict.Ignore()", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := client.User.Create().SetEmail("ignore@example.com").SetPhoneNum("+111").Exec(ctx)
		if err != nil {
			t.Fatalf("failed to insert initial: %v", err)
		}

		res, err := client.User.Create().
			SetEmail("ignore@example.com").
			SetPhoneNum("+222").
			OnConflict(user.Email).
			Ignore().
			Exec(ctx)

		if err != nil {
			t.Fatalf("single ignore failed with error: %v", err)
		}
		if res != nil {
			t.Errorf("expected returned record to be nil on ignore conflict, got %+v", res)
		}
	})

	t.Run("CreateMany with empty slice", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		affected, err := client.User.CreateMany().
			OnConflict(user.Email).
			Update(func(u *valk.UserUpsert) {
				u.LoginCount.Add(1)
			}).
			Exec(ctx)

		if err != nil {
			t.Fatalf("expected no error for empty CreateMany, got: %v", err)
		}
		if affected != 0 {
			t.Errorf("expected 0 affected rows, got %d", affected)
		}
	})

	t.Run("Single Create with multi-assignment custom update", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		p := "initialPass"
		_, err := client.User.Create().
			SetEmail("multi@example.com").
			SetPhoneNum("+100").
			SetPassword(p).
			SetLoginCount(10).
			Exec(ctx)
		if err != nil {
			t.Fatalf("failed to insert initial: %v", err)
		}

		newRole := valk.UserRoleTypeTeacher
		strRole := string(newRole)
		res, err := client.User.Create().
			SetEmail("multi@example.com").
			SetPhoneNum("+200").
			SetPassword("ignoredNewPass").
			SetRoleOptional(newRole).
			OnConflict(user.Email).
			Update(func(u *valk.UserUpsert) {
				u.PhoneNum.Update()
				u.LoginCount.Add(5)
				u.Password.Set(nil)
				u.RoleOptional.Set(&strRole)
			}).
			Exec(ctx)

		if err != nil {
			t.Fatalf("multi-assignment upsert failed: %v", err)
		}

		if res.PhoneNum != "+200" {
			t.Errorf("expected phoneNum to be updated to +200, got %q", res.PhoneNum)
		}
		if res.LoginCount != 15 {
			t.Errorf("expected loginCount to be 10 + 5 = 15, got %d", res.LoginCount)
		}
		if res.Password != nil {
			t.Errorf("expected password to be NULL (nil), got %v", res.Password)
		}
		if res.RoleOptional == nil || *res.RoleOptional != valk.UserRoleTypeTeacher {
			t.Errorf("expected roleOptional to be TEACHER, got %v", res.RoleOptional)
		}
	})

	t.Run("CreateMany with duplicate records in same batch", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		affected, err := client.User.CreateMany(
			client.User.Create().SetEmail("batch-dupe@example.com").SetPhoneNum("+111"),
			client.User.Create().SetEmail("batch-dupe@example.com").SetPhoneNum("+222"),
		).OnConflict(user.Email).Ignore().Exec(ctx)

		if err != nil {
			t.Fatalf("expected duplicate ignore in same batch to succeed, got error: %v", err)
		}
		if affected != 1 {
			t.Errorf("expected 1 row affected, got %d", affected)
		}

		_, err2 := client.User.CreateMany(
			client.User.Create().SetEmail("batch-dupe2@example.com").SetPhoneNum("+333"),
			client.User.Create().SetEmail("batch-dupe2@example.com").SetPhoneNum("+444"),
		).OnConflict(user.Email).UpdateNewValues().Exec(ctx)

		if err2 != nil {
			t.Logf("CreateMany duplicate update failed as expected on this database provider: %v", err2)
		} else {
			t.Logf("CreateMany duplicate update succeeded on this database provider")
		}
	})

	t.Run("Upsert custom update violating constraints", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := client.User.Create().SetEmail("fk-violate@example.com").SetPhoneNum("+111").Exec(ctx)
		if err != nil {
			t.Fatalf("failed to insert initial: %v", err)
		}

		badFk := "does-not-exist"
		_, err = client.User.Create().
			SetEmail("fk-violate@example.com").
			SetPhoneNum("+222").
			OnConflict(user.Email).
			Update(func(u *valk.UserUpsert) {
				u.ReferredById.Set(&badFk)
			}).
			Exec(ctx)

		if err == nil {
			t.Errorf("expected foreign key violation error, got nil")
		} else {
			t.Logf("got expected foreign key violation error: %v", err)
		}

		dbUser, err := client.User.FindUnique(user.Email.EQ("fk-violate@example.com")).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to query: %v", err)
		}
		if dbUser.PhoneNum != "+111" {
			t.Errorf("expected phoneNum to remain +111 after failed transaction, got %q", dbUser.PhoneNum)
		}
	})

	t.Run("Composite primary key writes and upserts on CategoryToPost", func(t *testing.T) {
		client, cleanup := setupTestDB(t)
		defer cleanup()

		usr, err := client.User.Create().
			SetEmail("catpost@example.com").
			SetPhoneNum("+9999").
			Exec(ctx)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		pst, err := client.Post.Create().
			SetTitle("GraphQL vs REST").
			SetAuthorId(usr.Id).
			Exec(ctx)
		if err != nil {
			t.Fatalf("failed to create post: %v", err)
		}

		cat, err := client.Category.Create().
			SetName("Tech").
			Exec(ctx)
		if err != nil {
			t.Fatalf("failed to create category: %v", err)
		}

		link, err := client.CategoryToPost.Create().
			SetPostId(pst.Id).
			SetCategoryId(cat.Id).
			Exec(ctx)
		if err != nil {
			t.Fatalf("failed to create CategoryToPost link: %v", err)
		}
		if link.PostId != pst.Id || link.CategoryId != cat.Id {
			t.Errorf("incorrect link values: %+v", link)
		}

		ignoredLink, err := client.CategoryToPost.Create().
			SetPostId(pst.Id).
			SetCategoryId(cat.Id).
			OnConflict(categoryToPost.PostIdCategoryId).
			Ignore().
			Exec(ctx)
		if err != nil {
			t.Fatalf("OnConflict.Ignore on composite key failed: %v", err)
		}
		if ignoredLink != nil {
			t.Errorf("expected ignored link to return nil, got: %+v", ignoredLink)
		}

		affected, err := client.CategoryToPost.CreateMany(
			client.CategoryToPost.Create().SetPostId(pst.Id).SetCategoryId(cat.Id),
		).OnConflict(categoryToPost.PostIdCategoryId).Ignore().Exec(ctx)
		if err != nil {
			t.Fatalf("CreateMany OnConflict.Ignore on composite key failed: %v", err)
		}
		if affected != 0 {
			t.Errorf("expected 0 rows affected since it already exists, got %d", affected)
		}
	})
}
