package main

import (
	"context"

	"testing"

	"integration/valk"
	"integration/valk/post"
	"integration/valk/user"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func seedUsers(t *testing.T, ctx context.Context, db *valk.DB) ([]*valk.User, error) {
	t.Helper()

	u1, err := db.User.Create().
		SetEmail("alpha@example.com").
		SetPhoneNum("1111111111").
		SetLoginCount(10).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	u2, err := db.User.Create().
		SetEmail("beta@domain.org").
		SetPhoneNum("2222222222").
		SetLoginCount(25).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	u3, err := db.User.Create().
		SetEmail("gamma@example.com").
		SetPhoneNum("3333333333").
		SetLoginCount(50).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return []*valk.User{u1, u2, u3}, nil
}
func TestPredicateOperators(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	users, err := seedUsers(t, ctx, db)
	if err != nil {
		t.Fatalf("failed to seed users: %v", err)
	}

	u1, u2, u3 := users[0], users[1], users[2]
	t.Run("NotIn", func(t *testing.T) {
		res, err := db.User.FindMany(user.Email.NotIn([]string{u1.Email, u2.Email})).Exec(ctx)
		if err != nil {
			t.Fatalf("NotIn failed: %v", err)
		}
		if len(res) != 1 || res[0].Id != u3.Id {
			t.Errorf("expected u3, got %v", res)
		}
	})

	t.Run("Between", func(t *testing.T) {
		res, err := db.User.FindMany(user.LoginCount.Between(20, 60)).Exec(ctx)
		if err != nil {
			t.Fatalf("Between failed: %v", err)
		}
		if len(res) != 2 {
			t.Errorf("expected 2 records, got %d", len(res))
		}
	})

	t.Run("HasPrefix", func(t *testing.T) {
		res, err := db.User.FindMany(user.Email.HasPrefix("alp")).Exec(ctx)
		if err != nil {
			t.Fatalf("HasPrefix failed: %v", err)
		}
		if len(res) != 1 || res[0].Id != u1.Id {
			t.Errorf("expected u1, got %v", res)
		}
	})

	t.Run("HasSuffix", func(t *testing.T) {
		res, err := db.User.FindMany(user.Email.HasSuffix("@example.com")).Exec(ctx)
		if err != nil {
			t.Fatalf("HasSuffix failed: %v", err)
		}
		if len(res) != 2 {
			t.Errorf("expected 2 records, got %d", len(res))
		}
	})

	t.Run("ILike", func(t *testing.T) {
		res, err := db.User.FindMany(user.Email.ILike("ALPHA@%")).Exec(ctx)
		if err != nil {
			t.Fatalf("ILike failed: %v", err)
		}
		if len(res) != 1 || res[0].Id != u1.Id {
			t.Errorf("expected u1, got %v", res)
		}
	})

	t.Run("Logical Composition And/Or/Not", func(t *testing.T) {
		// (Email ends with @example.com AND loginCount = 50) OR (phoneNum = 2222222222)
		comp := valk.Or(
			valk.And(
				user.Email.HasSuffix("@example.com"),
				user.LoginCount.EQ(50),
			),
			user.PhoneNum.EQ("2222222222"),
		)
		res, err := db.User.FindMany(comp).Exec(ctx)
		if err != nil {
			t.Fatalf("Logical And/Or failed: %v", err)
		}
		if len(res) != 2 {
			t.Errorf("expected 2 records, got %d", len(res))
		}

		// NOT (Email ends with @example.com)
		notRes, err := db.User.FindMany(valk.Not(user.Email.HasSuffix("@example.com"))).Exec(ctx)
		if err != nil {
			t.Fatalf("Logical Not failed: %v", err)
		}
		if len(notRes) != 1 || notRes[0].Id != u2.Id {
			t.Errorf("expected u2, got %v", notRes)
		}
	})

	t.Run("IsNull and IsNotNull on optional field", func(t *testing.T) {
		nullRes, err := db.User.FindMany(user.ReferredById.IsNull()).Exec(ctx)
		if err != nil {
			t.Fatalf("IsNull failed: %v", err)
		}
		if len(nullRes) != 3 {
			t.Errorf("expected 3 users with null ReferredById, got %d", len(nullRes))
		}

		notNullRes, err := db.User.FindMany(user.ReferredById.IsNotNull()).Exec(ctx)
		if err != nil {
			t.Fatalf("IsNotNull failed: %v", err)
		}
		if len(notNullRes) != 0 {
			t.Errorf("expected 0 users with non-null ReferredById, got %d", len(notNullRes))
		}
	})
}

func TestArrayOperators(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	users, err := seedUsers(t, ctx, db)
	if err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	p1, err := db.Post.Create().
		SetTitle("Go ORM").
		SetAuthorId(users[0].Id).
		SetTags([]string{"golang", "orm", "database"}).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed creating p1: %v", err)
	}

	p2, err := db.Post.Create().
		SetTitle("Python Web").
		SetAuthorId(users[1].Id).
		SetTags([]string{"python", "web"}).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed creating p2: %v", err)
	}
	_ = p2

	p3, err := db.Post.Create().
		SetTitle("Empty Post").
		SetAuthorId(users[2].Id).
		SetTags([]string{}).
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed creating p3: %v", err)
	}

	t.Run("Has", func(t *testing.T) {
		res, err := db.Post.FindMany(post.Tags.Has("golang")).Exec(ctx)
		if err != nil {
			t.Fatalf("Has failed: %v", err)
		}
		if len(res) != 1 || res[0].Id != p1.Id {
			t.Errorf("expected p1, got %v", res)
		}
	})

	t.Run("HasEvery", func(t *testing.T) {
		res, err := db.Post.FindMany(post.Tags.HasEvery([]string{"golang", "orm"})).Exec(ctx)
		if err != nil {
			t.Fatalf("HasEvery failed: %v", err)
		}
		if len(res) != 1 || res[0].Id != p1.Id {
			t.Errorf("expected p1, got %v", res)
		}
	})

	t.Run("HasSome", func(t *testing.T) {
		res, err := db.Post.FindMany(post.Tags.HasSome([]string{"python", "golang"})).Exec(ctx)
		if err != nil {
			t.Fatalf("HasSome failed: %v", err)
		}
		if len(res) != 2 {
			t.Errorf("expected 2 records, got %d", len(res))
		}
	})

	t.Run("IsEmpty", func(t *testing.T) {
		res, err := db.Post.FindMany(post.Tags.IsEmpty()).Exec(ctx)
		if err != nil {
			t.Fatalf("IsEmpty failed: %v", err)
		}
		if len(res) != 1 || res[0].Id != p3.Id {
			t.Errorf("expected p3, got %v", res)
		}
	})
}
