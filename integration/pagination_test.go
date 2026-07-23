package main

import (
	"context"
	"testing"

	"integration/valk"
	"integration/valk/category"
	"integration/valk/categoryToPost"
	"integration/valk/comment"
	"integration/valk/post"
	"integration/valk/user"
)

func TestPagination_OrderBy_AscAndDesc(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	u1, err := db.User.Create().SetEmail("pag_c@example.com").SetPhoneNum("+111").SetLoginCount(30).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create u1: %v", err)
	}
	t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(u1.Id)).Exec(ctx) })

	u2, err := db.User.Create().SetEmail("pag_a@example.com").SetPhoneNum("+222").SetLoginCount(10).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create u2: %v", err)
	}
	t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(u2.Id)).Exec(ctx) })

	u3, err := db.User.Create().SetEmail("pag_b@example.com").SetPhoneNum("+333").SetLoginCount(20).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create u3: %v", err)
	}
	t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(u3.Id)).Exec(ctx) })

	// Test 1: OrderBy Email ASC
	ascUsers, err := db.User.FindMany(
		valk.Or(
			user.Email.EQ("pag_a@example.com"),
			user.Email.EQ("pag_b@example.com"),
			user.Email.EQ("pag_c@example.com"),
		),
	).OrderBy(user.Email.Asc()).Exec(ctx)
	if err != nil {
		t.Fatalf("FindMany OrderBy Asc failed: %v", err)
	}
	if len(ascUsers) != 3 {
		t.Fatalf("expected 3 users, got %d", len(ascUsers))
	}
	if ascUsers[0].Email != "pag_a@example.com" || ascUsers[1].Email != "pag_b@example.com" || ascUsers[2].Email != "pag_c@example.com" {
		t.Errorf("unexpected ASC order: got [%s, %s, %s]", ascUsers[0].Email, ascUsers[1].Email, ascUsers[2].Email)
	}

	// Test 2: OrderBy LoginCount DESC
	descUsers, err := db.User.FindMany(
		valk.Or(
			user.Email.EQ("pag_a@example.com"),
			user.Email.EQ("pag_b@example.com"),
			user.Email.EQ("pag_c@example.com"),
		),
	).OrderBy(user.LoginCount.Desc()).Exec(ctx)
	if err != nil {
		t.Fatalf("FindMany OrderBy Desc failed: %v", err)
	}
	if len(descUsers) != 3 {
		t.Fatalf("expected 3 users, got %d", len(descUsers))
	}
	if descUsers[0].LoginCount != 30 || descUsers[1].LoginCount != 20 || descUsers[2].LoginCount != 10 {
		t.Errorf("unexpected DESC order: got [%d, %d, %d]", descUsers[0].LoginCount, descUsers[1].LoginCount, descUsers[2].LoginCount)
	}
}

func TestPagination_Cursor(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	u1, err := db.User.Create().SetEmail("cur_1@example.com").SetPhoneNum("+101").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create u1: %v", err)
	}
	t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(u1.Id)).Exec(ctx) })

	u2, err := db.User.Create().SetEmail("cur_2@example.com").SetPhoneNum("+102").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create u2: %v", err)
	}
	t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(u2.Id)).Exec(ctx) })

	u3, err := db.User.Create().SetEmail("cur_3@example.com").SetPhoneNum("+103").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create u3: %v", err)
	}
	t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(u3.Id)).Exec(ctx) })

	// Fetch page using Cursor on u1
	nextPage, err := db.User.FindMany(
		valk.Or(
			user.Email.EQ("cur_1@example.com"),
			user.Email.EQ("cur_2@example.com"),
			user.Email.EQ("cur_3@example.com"),
		),
	).OrderBy(user.Email.Asc()).Cursor(user.Email.EQ("cur_1@example.com")).Take(2).Exec(ctx)
	if err != nil {
		t.Fatalf("FindMany Cursor failed: %+v", err)
	}

	if len(nextPage) != 2 {
		t.Fatalf("expected 2 users on next page, got %d", len(nextPage))
	}
	if nextPage[0].Email != "cur_2@example.com" || nextPage[1].Email != "cur_3@example.com" {
		t.Errorf("unexpected cursor results: got [%s, %s]", nextPage[0].Email, nextPage[1].Email)
	}
}

func TestPagination_NegativeTake(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	u1, err := db.User.Create().SetEmail("neg_1@example.com").SetPhoneNum("+201").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create u1: %v", err)
	}
	t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(u1.Id)).Exec(ctx) })

	u2, err := db.User.Create().SetEmail("neg_2@example.com").SetPhoneNum("+202").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create u2: %v", err)
	}
	t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(u2.Id)).Exec(ctx) })

	u3, err := db.User.Create().SetEmail("neg_3@example.com").SetPhoneNum("+203").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create u3: %v", err)
	}
	t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(u3.Id)).Exec(ctx) })

	// Fetch 2 items BEFORE neg_3 using Take(-2) with Cursor(neg_3)
	prevPage, err := db.User.FindMany(
		valk.Or(
			user.Email.EQ("neg_1@example.com"),
			user.Email.EQ("neg_2@example.com"),
			user.Email.EQ("neg_3@example.com"),
		),
	).OrderBy(user.Email.Asc()).Cursor(user.Email.EQ("neg_3@example.com")).Take(-2).Exec(ctx)
	if err != nil {
		t.Fatalf("FindMany Negative Take failed: %v", err)
	}

	if len(prevPage) != 2 {
		t.Fatalf("expected 2 users on previous page, got %d", len(prevPage))
	}
	if prevPage[0].Email != "neg_1@example.com" || prevPage[1].Email != "neg_2@example.com" {
		t.Errorf("unexpected negative take results: got [%s, %s]", prevPage[0].Email, prevPage[1].Email)
	}
}

func TestPagination_AllFourCombined(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Seed 5 users: comb_1 .. comb_5
	emails := []string{"comb_1@ex.com", "comb_2@ex.com", "comb_3@ex.com", "comb_4@ex.com", "comb_5@ex.com"}
	for i, e := range emails {
		usr, err := db.User.Create().SetEmail(e).SetPhoneNum("+30" + string(rune('1'+i))).Exec(ctx)
		if err != nil {
			t.Fatalf("failed to seed %s: %v", e, err)
		}
		id := usr.Id
		t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(id)).Exec(ctx) })
	}

	// All 4 Combined: Cursor(comb_1), Skip(1), Take(2), OrderBy(Email ASC)
	// After comb_1: list is [comb_2, comb_3, comb_4, comb_5].
	// Skip(1): skips comb_2 -> remaining is [comb_3, comb_4, comb_5].
	// Take(2): returns [comb_3, comb_4].
	results, err := db.User.FindMany(
		user.Email.Contains("comb_"),
	).
		OrderBy(user.Email.Asc()).
		Cursor(user.Email.EQ("comb_1@ex.com")).
		Skip(1).
		Take(2).
		Exec(ctx)
	if err != nil {
		t.Fatalf("FindMany All Four Combined failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 users, got %d", len(results))
	}
	if results[0].Email != "comb_3@ex.com" || results[1].Email != "comb_4@ex.com" {
		t.Errorf("unexpected results for All Four Combined: got [%s, %s]", results[0].Email, results[1].Email)
	}
}

func TestPagination_FindFirst_TakeAndCursor(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	u1, _ := db.User.Create().SetEmail("ff_1@ex.com").SetPhoneNum("+401").Exec(ctx)
	t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(u1.Id)).Exec(ctx) })

	u2, _ := db.User.Create().SetEmail("ff_2@ex.com").SetPhoneNum("+402").Exec(ctx)
	t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(u2.Id)).Exec(ctx) })

	u3, _ := db.User.Create().SetEmail("ff_3@ex.com").SetPhoneNum("+403").Exec(ctx)
	t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(u3.Id)).Exec(ctx) })

	// FindFirst with Take(-1) and Cursor(ff_3) -> should return ff_2 (the single item immediately before ff_3)
	lastBefore, err := db.User.FindFirst(user.Email.Contains("ff_")).
		OrderBy(user.Email.Asc()).
		Cursor(user.Email.EQ("ff_3@ex.com")).
		Take(-1).
		Exec(ctx)
	if err != nil {
		t.Fatalf("FindFirst Take(-1) failed: %v", err)
	}
	if lastBefore == nil || lastBefore.Email != "ff_2@ex.com" {
		t.Fatalf("expected ff_2@ex.com, got %v", lastBefore)
	}
}

func TestPagination_NestedRelation_CursorAndOrderBy(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	usr, err := db.User.Create().SetEmail("rel_pag@ex.com").SetPhoneNum("+500").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(usr.Id)).Exec(ctx) })

	// Seed 4 posts for this user
	p1, _ := db.Post.Create().SetTitle("Post 1").SetAuthorId(usr.Id).Exec(ctx)
	t.Cleanup(func() { _, _ = db.Post.Delete(post.Id.EQ(p1.Id)).Exec(ctx) })

	p2, _ := db.Post.Create().SetTitle("Post 2").SetAuthorId(usr.Id).Exec(ctx)
	t.Cleanup(func() { _, _ = db.Post.Delete(post.Id.EQ(p2.Id)).Exec(ctx) })

	p3, _ := db.Post.Create().SetTitle("Post 3").SetAuthorId(usr.Id).Exec(ctx)
	t.Cleanup(func() { _, _ = db.Post.Delete(post.Id.EQ(p3.Id)).Exec(ctx) })

	p4, _ := db.Post.Create().SetTitle("Post 4").SetAuthorId(usr.Id).Exec(ctx)
	t.Cleanup(func() { _, _ = db.Post.Delete(post.Id.EQ(p4.Id)).Exec(ctx) })

	// Query user with Posts query builder using OrderBy, Cursor, and Take
	res, err := db.User.FindUnique(user.Id.EQ(usr.Id)).Select(valk.UserSelect{
		Email: true,
		Posts: post.Query().
			OrderBy(post.Title.Asc()).
			Cursor(post.Id.EQ(p1.Id)).
			Take(2).
			Select(post.Select{
				Id:    true,
				Title: true,
			}),
	}).Exec(ctx)
	if err != nil {
		t.Fatalf("FindUnique with nested relation cursor failed: %v", err)
	}

	if len(res.Posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(res.Posts))
	}
	if res.Posts[0].Title != "Post 2" || res.Posts[1].Title != "Post 3" {
		t.Errorf("unexpected nested relation posts: got [%s, %s]", res.Posts[0].Title, res.Posts[1].Title)
	}
}

func TestPagination_DeepNestedRelation(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	usr, _ := db.User.Create().SetEmail("deep_pag@ex.com").SetPhoneNum("+600").Exec(ctx)
	t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(usr.Id)).Exec(ctx) })

	pst, _ := db.Post.Create().SetTitle("Deep Post").SetAuthorId(usr.Id).Exec(ctx)
	t.Cleanup(func() { _, _ = db.Post.Delete(post.Id.EQ(pst.Id)).Exec(ctx) })

	c1, err := db.Comment.Create().SetTextify(100).SetDummy1(1).SetDummy2("d2_1").SetDummy3("d3_1").SetAuthorId(usr.Id).SetPostId(pst.Id).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create comment 1: %v", err)
	}
	t.Cleanup(func() { _, _ = db.Comment.Delete(comment.Id.EQ(c1.Id)).Exec(ctx) })

	c2, err := db.Comment.Create().SetTextify(200).SetDummy1(2).SetDummy2("d2_2").SetDummy3("d3_2").SetAuthorId(usr.Id).SetPostId(pst.Id).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create comment 2: %v", err)
	}
	t.Cleanup(func() { _, _ = db.Comment.Delete(comment.Id.EQ(c2.Id)).Exec(ctx) })

	// Query User -> Posts -> Comments with OrderBy and Take on both levels
	u, err := db.User.FindUnique(user.Id.EQ(usr.Id)).Select(valk.UserSelect{
		Email: true,
		Posts: post.Query().Take(1).Select(post.Select{
			Title: true,
			Comments: comment.Query().
				OrderBy(comment.Textify.Desc()).
				Take(1).
				Select(comment.Select{
					Textify: true,
				}),
		}),
	}).Exec(ctx)
	if err != nil {
		t.Fatalf("Deep nested relation query failed: %v", err)
	}

	if len(u.Posts) != 1 || len(u.Posts[0].Comments) != 1 {
		t.Fatalf("expected 1 post and 1 comment, got %d posts and %v comments", len(u.Posts), len(u.Posts[0].Comments))
	}
	if u.Posts[0].Comments[0].Textify != 200 {
		t.Errorf("expected highest comment textify 200, got %d", u.Posts[0].Comments[0].Textify)
	}
}

func TestPagination_OutOfBoundsAndEmpty(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	usr, _ := db.User.Create().SetEmail("oob@ex.com").SetPhoneNum("+700").Exec(ctx)
	t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(usr.Id)).Exec(ctx) })

	// Query cursor past last item
	res, err := db.User.FindMany(user.Email.EQ("oob@ex.com")).
		OrderBy(user.Email.Asc()).
		Cursor(user.Email.EQ("oob@ex.com")).
		Take(5).
		Exec(ctx)
	if err != nil {
		t.Fatalf("out of bounds cursor should not error: %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("expected 0 items for out of bounds cursor, got %d", len(res))
	}
}

func TestPagination_CompoundPK(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// 1. Seed User & Category
	usr, err := db.User.Create().SetEmail("cp_author@ex.com").SetPhoneNum("+800").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	t.Cleanup(func() { _, _ = db.User.Delete(user.Id.EQ(usr.Id)).Exec(ctx) })

	cat, err := db.Category.Create().SetName("Tech").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create category: %v", err)
	}
	t.Cleanup(func() { _, _ = db.Category.Delete(category.Id.EQ(cat.Id)).Exec(ctx) })

	// 2. Seed 3 Posts with explicit sorted IDs
	p1, _ := db.Post.Create().SetId("post-cp-1").SetTitle("Post 1").SetAuthorId(usr.Id).Exec(ctx)
	t.Cleanup(func() { _, _ = db.Post.Delete(post.Id.EQ(p1.Id)).Exec(ctx) })

	p2, _ := db.Post.Create().SetId("post-cp-2").SetTitle("Post 2").SetAuthorId(usr.Id).Exec(ctx)
	t.Cleanup(func() { _, _ = db.Post.Delete(post.Id.EQ(p2.Id)).Exec(ctx) })

	p3, _ := db.Post.Create().SetId("post-cp-3").SetTitle("Post 3").SetAuthorId(usr.Id).Exec(ctx)
	t.Cleanup(func() { _, _ = db.Post.Delete(post.Id.EQ(p3.Id)).Exec(ctx) })

	// 3. Seed 3 CategoryToPost records with compound PK (postId, categoryId)
	cp1, err := db.CategoryToPost.Create().SetPostId(p1.Id).SetCategoryId(cat.Id).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create cp1: %v", err)
	}
	t.Cleanup(func() {
		_, _ = db.CategoryToPost.Delete(categoryToPost.PostIdCategoryIdUnique(p1.Id, cat.Id)).Exec(ctx)
	})

	cp2, err := db.CategoryToPost.Create().SetPostId(p2.Id).SetCategoryId(cat.Id).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create cp2: %v", err)
	}
	t.Cleanup(func() {
		_, _ = db.CategoryToPost.Delete(categoryToPost.PostIdCategoryIdUnique(p2.Id, cat.Id)).Exec(ctx)
	})

	cp3, err := db.CategoryToPost.Create().SetPostId(p3.Id).SetCategoryId(cat.Id).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create cp3: %v", err)
	}
	t.Cleanup(func() {
		_, _ = db.CategoryToPost.Delete(categoryToPost.PostIdCategoryIdUnique(p3.Id, cat.Id)).Exec(ctx)
	})

	_ = cp1
	_ = cp2
	_ = cp3

	// Test cursor pagination on compound primary key model
	nextCPs, err := db.CategoryToPost.FindMany(
		categoryToPost.CategoryId.EQ(cat.Id),
	).
		OrderBy(categoryToPost.PostId.Asc()).
		Cursor(categoryToPost.PostIdCategoryIdUnique(p1.Id, cat.Id)).
		Take(2).
		Exec(ctx)
	if err != nil {
		t.Fatalf("Compound PK Cursor pagination failed: %v", err)
	}

	if len(nextCPs) != 2 {
		t.Fatalf("expected 2 records, got %d", len(nextCPs))
	}
	if nextCPs[0].PostId != p2.Id || nextCPs[1].PostId != p3.Id {
		t.Errorf("unexpected compound PK results: got [%s, %s]", nextCPs[0].PostId, nextCPs[1].PostId)
	}
}
