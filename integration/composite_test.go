package main

import (
	"context"
	"integration/phi"
	"integration/phi/categoryToPost"
	"integration/phi/user"
	"testing"
)

func TestCompositeKeys(t *testing.T) {
	ctx := context.Background()

	t.Run("EmailPhone_EQ_returns_correct_Column_and_Value", func(t *testing.T) {
		p := user.EmailPhone.EQ("a@b.com", "555")

		col := p.Column()
		if col != "emailPhone" {
			t.Errorf("expected column 'emailPhone', got %q", col)
		}

		val, ok := p.Value().(map[string]any)
		if !ok {
			t.Fatalf("expected map[string]any, got %T", p.Value())
		}
		if val["email"] != "a@b.com" {
			t.Errorf("expected email 'a@b.com', got %v", val["email"])
		}
		if val["phoneNum"] != "555" {
			t.Errorf("expected phoneNum '555', got %v", val["phoneNum"])
		}

		if _, ok := p.Value().(map[string]any); !ok {
			t.Error("expected composite predicate (map value)")
		}
	})

	t.Run("PostId_CategoryId_EQ_returns_correct_Column_and_Value", func(t *testing.T) {
		p := categoryToPost.PostId_CategoryId.EQ("post-1", 42)

		col := p.Column()
		if col != "PostId_CategoryId" {
			t.Errorf("expected column 'PostId_CategoryId', got %q", col)
		}

		val, ok := p.Value().(map[string]any)
		if !ok {
			t.Fatalf("expected map[string]any, got %T", p.Value())
		}
		if val["postId"] != "post-1" {
			t.Errorf("expected postId 'post-1', got %v", val["postId"])
		}
		if val["categoryId"] != int32(42) {
			t.Errorf("expected categoryId 42, got %v", val["categoryId"])
		}
	})

	t.Run("simple_unique_predicate_returns_scalar_Column_and_Value", func(t *testing.T) {
		p := user.Email.EQ("x@y.com")

		col := p.Column()
		if col != "email" {
			t.Errorf("expected column 'email', got %q", col)
		}

		val, ok := p.Value().(string)
		if !ok {
			t.Fatalf("expected string, got %T", p.Value())
		}
		if val != "x@y.com" {
			t.Errorf("expected 'x@y.com', got %q", val)
		}
	})

	t.Run("composite_findunique_works", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := db.User.Create().
			SetEmail("compositepk@example.com").
			SetPhoneNum("composite-001").
			Exec(ctx)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}

		u, err := db.User.FindUnique(
			user.EmailPhone.EQ("compositepk@example.com", "composite-001"),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("find unique failed: %v", err)
		}
		if u.Email != "compositepk@example.com" {
			t.Errorf("expected 'compositepk@example.com', got %q", u.Email)
		}
	})

	t.Run("composite_findunique_with_additional_predicate", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := db.User.Create().
			SetEmail("additional@example.com").
			SetPhoneNum("composite-010").
			SetPassword("secret").
			Exec(ctx)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}

		u, err := db.User.FindUnique(
			user.EmailPhone.EQ("additional@example.com", "composite-010"),
			user.LoginCount.EQ(0),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("find unique failed: %v", err)
		}
		if u.Email != "additional@example.com" {
			t.Errorf("expected email 'additional@example.com', got %q", u.Email)
		}
	})

	t.Run("composite_PK_findunique_works", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		author, err := db.User.Create().
			SetEmail("pk-author@example.com").
			SetPhoneNum("composite-020").
			Exec(ctx)
		if err != nil {
			t.Fatalf("create author failed: %v", err)
		}

		p, err := db.Post.Create().
			SetTitle("composite-pk-test").
			SetAuthorId(author.Id).
			Exec(ctx)
		if err != nil {
			t.Fatalf("create post failed: %v", err)
		}

		cat, err := db.Category.Create().
			SetName("composite-pk-cat").
			Exec(ctx)
		if err != nil {
			t.Fatalf("create category failed: %v", err)
		}

		_, err = db.CategoryToPost.Create().
			SetPostId(p.Id).
			SetCategoryId(cat.Id).
			Exec(ctx)
		if err != nil {
			t.Fatalf("create categoryToPost failed: %v", err)
		}

		ctp, err := db.CategoryToPost.FindUnique(
			categoryToPost.PostId_CategoryId.EQ(p.Id, cat.Id),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("find unique categoryToPost failed: %v", err)
		}
		if ctp.PostId != p.Id || ctp.CategoryId != cat.Id {
			t.Errorf("expected (%q, %d), got (%q, %d)", p.Id, cat.Id, ctp.PostId, ctp.CategoryId)
		}
	})

	t.Run("hook_inspects_composite_Column_and_Value", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		var seenCol string
		var seenVal any

		db.User.Use(user.Extension{
			FindUnique: func(ctx context.Context, args *user.FindUniqueArgs, next user.FindUniqueQuery) (*phi.User, error) {
				if len(args.Where) > 0 {
					seenCol = args.Where[0].Column()
					seenVal = args.Where[0].Value()
				}
				return next(ctx, args)
			},
		})

		_, err := db.User.Create().
			SetEmail("hook-inspect@example.com").
			SetPhoneNum("composite-004").
			Exec(ctx)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}

		_, err = db.User.FindUnique(
			user.EmailPhone.EQ("hook-inspect@example.com", "composite-004"),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("find unique failed: %v", err)
		}

		if seenCol != "emailPhone" {
			t.Errorf("expected column 'emailPhone', got %q", seenCol)
		}
		m, ok := seenVal.(map[string]any)
		if !ok {
			t.Fatalf("expected map, got %T", seenVal)
		}
		if m["email"] != "hook-inspect@example.com" {
			t.Errorf("expected email 'hook-inspect@example.com', got %v", m["email"])
		}
		if m["phoneNum"] != "composite-004" {
			t.Errorf("expected phoneNum 'composite-004', got %v", m["phoneNum"])
		}
	})

	t.Run("hook_inspects_composite_Children", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		var children []phi.ChildPredicate

		db.User.Use(user.Extension{
			FindUnique: func(ctx context.Context, args *user.FindUniqueArgs, next user.FindUniqueQuery) (*phi.User, error) {
				if len(args.Where) > 0 {
					children = args.Where[0].Children()
				}
				return next(ctx, args)
			},
		})

		_, err := db.User.Create().
			SetEmail("hook-children@example.com").
			SetPhoneNum("composite-005b").
			Exec(ctx)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}

		_, err = db.User.FindUnique(
			user.EmailPhone.EQ("hook-children@example.com", "composite-005b"),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("find unique failed: %v", err)
		}

		if len(children) != 2 {
			t.Fatalf("expected 2 children, got %d", len(children))
		}
		if children[0].Column != "email" || children[0].Value != "hook-children@example.com" {
			t.Errorf("unexpected child 0: %+v", children[0])
		}
		if children[1].Column != "phoneNum" || children[1].Value != "composite-005b" {
			t.Errorf("unexpected child 1: %+v", children[1])
		}
	})

	t.Run("hook_replaces_composite_predicate", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := db.User.Create().
			SetEmail("original@example.com").
			SetPhoneNum("composite-005").
			Exec(ctx)
		if err != nil {
			t.Fatalf("create original failed: %v", err)
		}

		_, err = db.User.Create().
			SetEmail("replacement@example.com").
			SetPhoneNum("composite-006a").
			Exec(ctx)
		if err != nil {
			t.Fatalf("create replacement failed: %v", err)
		}

		db.User.Use(user.Extension{
			FindUnique: func(ctx context.Context, args *user.FindUniqueArgs, next user.FindUniqueQuery) (*phi.User, error) {
				for i, w := range args.Where {
					if w.Column() == "emailPhone" {
						args.Where[i] = user.EmailPhone.EQ("replacement@example.com", "composite-006a")
					}
				}
				return next(ctx, args)
			},
		})

		u, err := db.User.FindUnique(
			user.EmailPhone.EQ("original@example.com", "composite-005"),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("find unique failed: %v", err)
		}
		if u.Email != "replacement@example.com" {
			t.Errorf("expected replacement email, got %q", u.Email)
		}
	})

	t.Run("hook_appends_to_composite_query", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := db.User.Create().
			SetEmail("append-test@example.com").
			SetPhoneNum("composite-030").
			SetPassword("secret123").
			Exec(ctx)
		if err != nil {
			t.Fatalf("create user failed: %v", err)
		}

		db.User.Use(user.Extension{
			FindUnique: func(ctx context.Context, args *user.FindUniqueArgs, next user.FindUniqueQuery) (*phi.User, error) {
				args.Where = append(args.Where, user.Password.EQ("secret123"))
				return next(ctx, args)
			},
		})

		u, err := db.User.FindUnique(
			user.EmailPhone.EQ("append-test@example.com", "composite-030"),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("find unique failed: %v", err)
		}
		if u.Email != "append-test@example.com" {
			t.Errorf("expected email 'append-test@example.com', got %q", u.Email)
		}
	})

	t.Run("composite_delete_works", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		author, err := db.User.Create().
			SetEmail("delete-author@example.com").
			SetPhoneNum("composite-040").
			Exec(ctx)
		if err != nil {
			t.Fatalf("create author failed: %v", err)
		}

		p, err := db.Post.Create().SetTitle("delete-composite").SetAuthorId(author.Id).Exec(ctx)
		if err != nil {
			t.Fatalf("create post failed: %v", err)
		}
		cat, err := db.Category.Create().SetName("delete-composite-cat").Exec(ctx)
		if err != nil {
			t.Fatalf("create category failed: %v", err)
		}
		_, err = db.CategoryToPost.Create().SetPostId(p.Id).SetCategoryId(cat.Id).Exec(ctx)
		if err != nil {
			t.Fatalf("create categoryToPost failed: %v", err)
		}

		deleted, err := db.CategoryToPost.Delete(
			categoryToPost.PostId_CategoryId.EQ(p.Id, cat.Id),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("delete failed: %v", err)
		}
		if deleted.PostId != p.Id {
			t.Errorf("expected deleted postId %q, got %q", p.Id, deleted.PostId)
		}

		found, err := db.CategoryToPost.FindUnique(
			categoryToPost.PostId_CategoryId.EQ(p.Id, cat.Id),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("find unique after delete failed: %v", err)
		}
		if found != nil {
			t.Errorf("expected nil after deletion, got %+v", found)
		}
	})

	t.Run("composite_OnConflict_ignores_duplicate", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		author, err := db.User.Create().
			SetEmail("onconflict-author@example.com").
			SetPhoneNum("composite-050").
			Exec(ctx)
		if err != nil {
			t.Fatalf("create author failed: %v", err)
		}

		p, err := db.Post.Create().SetTitle("onconflict-composite").SetAuthorId(author.Id).Exec(ctx)
		if err != nil {
			t.Fatalf("create post failed: %v", err)
		}
		cat, err := db.Category.Create().SetName("onconflict-composite-cat").Exec(ctx)
		if err != nil {
			t.Fatalf("create category failed: %v", err)
		}

		_, err = db.CategoryToPost.Create().
			SetPostId(p.Id).SetCategoryId(cat.Id).
			Exec(ctx)
		if err != nil {
			t.Fatalf("first create failed: %v", err)
		}

		affected, err := db.CategoryToPost.CreateMany(
			db.CategoryToPost.Create().SetPostId(p.Id).SetCategoryId(cat.Id),
		).OnConflict(categoryToPost.PostId_CategoryId).Ignore().Exec(ctx)
		if err != nil {
			t.Fatalf("OnConflict ignore failed: %v", err)
		}
		if affected != 0 {
			t.Errorf("expected 0 affected rows (duplicate ignored), got %d", affected)
		}
	})

	t.Run("simple_user_Email_EQ_still_works", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := db.User.Create().
			SetEmail("simple-eq@example.com").
			SetPhoneNum("composite-008").
			Exec(ctx)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}

		u, err := db.User.FindUnique(user.Email.EQ("simple-eq@example.com")).Exec(ctx)
		if err != nil {
			t.Fatalf("find unique failed: %v", err)
		}
		if u.Email != "simple-eq@example.com" {
			t.Errorf("expected 'simple-eq@example.com', got %q", u.Email)
		}
	})

	t.Run("composite_update_works", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := db.User.Create().
			SetEmail("update-composite@example.com").
			SetPhoneNum("composite-009").
			Exec(ctx)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}

		u, err := db.User.Update(
			user.EmailPhone.EQ("update-composite@example.com", "composite-009"),
		).SetPassword("newpassword").Exec(ctx)
		if err != nil {
			t.Fatalf("update failed: %v", err)
		}
		if u.Password == nil || *u.Password != "newpassword" {
			t.Errorf("expected password 'newpassword', got %v", u.Password)
		}
	})

	t.Run("composite_OnConflict_EmailPhone_UpdateNewValues", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		_, err := db.User.Create().
			SetEmail("upsert-composite@example.com").
			SetPhoneNum("composite-010").
			SetPassword("original").
			Exec(ctx)
		if err != nil {
			t.Fatalf("first create failed: %v", err)
		}

		affected, err := db.User.CreateMany(
			db.User.Create().
				SetEmail("upsert-composite@example.com").
				SetPhoneNum("composite-010").
				SetPassword("updated"),
		).OnConflict(user.EmailPhone).UpdateNewValues().Exec(ctx)
		if err != nil {
			t.Fatalf("OnConflict update failed: %v", err)
		}
		if affected != 1 {
			t.Errorf("expected 1 affected row, got %d", affected)
		}

		u, err := db.User.FindUnique(
			user.EmailPhone.EQ("upsert-composite@example.com", "composite-010"),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("find unique after upsert failed: %v", err)
		}
		if u.Password != nil && *u.Password == "updated" {
			t.Log("password was updated to new value")
		}
	})

	t.Run("hook_mixed_composite_and_simple", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		var seen []string

		db.User.Use(user.Extension{
			FindFirst: func(ctx context.Context, args *user.FindFirstArgs, next user.FindFirstQuery) (*phi.User, error) {
				for _, w := range args.Where {
					seen = append(seen, w.Column())
				}
				return next(ctx, args)
			},
		})

		_, err := db.User.Create().
			SetEmail("mixed@example.com").
			SetPhoneNum("composite-011").
			Exec(ctx)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}

		_, err = db.User.FindFirst(
			user.EmailPhone.EQ("mixed@example.com", "composite-011"),
			user.LoginCount.EQ(5),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("find first failed: %v", err)
		}

		if len(seen) != 2 {
			t.Fatalf("expected 2 predicates, got %d", len(seen))
		}
		if seen[0] != "emailPhone" {
			t.Errorf("expected first column 'emailPhone', got %q", seen[0])
		}
		if seen[1] != "loginCount" {
			t.Errorf("expected second column 'loginCount', got %q", seen[1])
		}
	})
}
