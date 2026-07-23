package main

import (
	"context"
	"encoding/json"
	"fmt"
	"integration/valk"
	"integration/valk/comment"
	"integration/valk/post"
	"integration/valk/user"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestFindUniqueWithNoFieldsSet(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().SetEmail("onlyuser@example.com").SetPhoneNum("000").Exec(ctx)
	if err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	res, err := db.User.FindUnique(valk.UniquePredicate[valk.User]{}).Exec(ctx)
	if err == nil && res != nil {
		t.Errorf("FindUnique with a zero-value where matched a row unexpectedly (%+v); it should require at least one unique field or return an error", res)
	}
}

func TestFindUniqueConflictingCompoundFields(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().SetEmail("a@example.com").SetPhoneNum("111").Exec(ctx)
	if err != nil {
		t.Fatalf("seed a failed: %v", err)
	}
	_, err = db.User.Create().SetEmail("b@example.com").SetPhoneNum("222").Exec(ctx)
	if err != nil {
		t.Fatalf("seed b failed: %v", err)
	}

	res, err := db.User.FindUnique(user.EmailPhoneUnique("a@example.com", "222")).Exec(ctx)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if res != nil {
		t.Errorf("expected nil since no single row matches both email a@example.com and phone 222, got: %+v", res)
	}
}

func TestSelectWithNoFieldsSet(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().SetEmail("empty_select@example.com").SetPhoneNum("333").Exec(ctx)
	if err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	res, err := db.User.FindUnique(user.Email.EQ("empty_select@example.com")).Select(valk.UserSelect{}).Exec(ctx)
	if err != nil {
		t.Fatalf("empty select produced an error instead of degrading gracefully: %v", err)
	}
	if res == nil {
		t.Fatal("expected a non-nil result even with no fields selected")
	}
}

func TestOmitAllFields(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().SetEmail("omit_all@example.com").SetPhoneNum("334").Exec(ctx)
	if err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	res, err := db.User.FindUnique(user.Email.EQ("omit_all@example.com")).Omit(valk.UserOmit{
		Id:       true,
		Email:    true,
		PhoneNum: true,
		Password: true,
		Role:     true,
	}).Exec(ctx)
	if err != nil {
		t.Fatalf("omitting every field produced an error instead of degrading gracefully: %v", err)
	}
	if res == nil {
		t.Fatal("expected a non-nil result even with everything omitted")
	}
}

func TestOmitIdFieldStillAllowsFilterById(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u, err := db.User.Create().SetEmail("omit_id@example.com").SetPhoneNum("335").Exec(ctx)
	if err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	res, err := db.User.FindUnique(user.Id.EQ(u.Id)).Omit(valk.UserOmit{Id: true}).Exec(ctx)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if res == nil {
		t.Fatal("expected to find the user by id even though id is omitted from the returned columns")
	}
	if res.Id != "" {
		t.Errorf("expected omitted Id field to be zero-value, got %q", res.Id)
	}
}

func TestRelationLoadWithNoRelatedRows(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().SetEmail("noposts@example.com").SetPhoneNum("444").Exec(ctx)
	if err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	res, err := db.User.FindUnique(user.Email.EQ("noposts@example.com")).Select(valk.UserSelect{
		Email: true,
		Posts: &valk.PostSelect{Title: true},
	}).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to find user with empty relation: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil user")
	}
	if len(res.Posts) != 0 {
		t.Errorf("expected 0 related posts, got %d", len(res.Posts))
	}
}

func TestFindUniqueRelationLoad(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	author, err := db.User.Create().SetEmail("unique_rel@example.com").SetPhoneNum("445").Exec(ctx)
	if err != nil {
		t.Fatalf("seed author failed: %v", err)
	}
	_, err = db.Post.Create().SetTitle("Unique Rel Post").SetAuthorId(author.Id).Exec(ctx)
	if err != nil {
		t.Fatalf("seed post failed: %v", err)
	}

	res, err := db.User.FindUnique(user.Email.EQ("unique_rel@example.com")).Select(valk.UserSelect{
		Email: true,
		Posts: &valk.PostSelect{Title: true},
	}).Exec(ctx)
	if err != nil {
		t.Fatalf("FindUnique with relation load failed: %v", err)
	}
	if res == nil || len(res.Posts) != 1 || res.Posts[0].Title != "Unique Rel Post" {
		t.Errorf("expected FindUnique to load the relation the same way FindMany does, got: %+v", res)
	}
}

func TestFindFirstRelationLoad(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	author, err := db.User.Create().SetEmail("first_rel@example.com").SetPhoneNum("446").Exec(ctx)
	if err != nil {
		t.Fatalf("seed author failed: %v", err)
	}
	_, err = db.Post.Create().SetTitle("First Rel Post").SetAuthorId(author.Id).Exec(ctx)
	if err != nil {
		t.Fatalf("seed post failed: %v", err)
	}

	res, err := db.User.FindFirst(user.Email.EQ("first_rel@example.com")).Select(valk.UserSelect{
		Email: true,
		Posts: &valk.PostSelect{Title: true},
	}).Exec(ctx)
	if err != nil {
		t.Fatalf("FindFirst with relation load failed: %v", err)
	}
	if res == nil || len(res.Posts) != 1 || res.Posts[0].Title != "First Rel Post" {
		t.Errorf("expected FindFirst to load the relation the same way FindMany does, got: %+v", res)
	}
}

func TestContextAlreadyCancelled(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := db.User.FindMany().Exec(ctx)
	if err == nil {
		t.Error("expected error when querying with an already-cancelled context, got nil")
	}
}

func TestContextDeadlineExceeded(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	time.Sleep(time.Millisecond)

	_, err := db.User.FindMany().Exec(ctx)
	if err == nil {
		t.Error("expected error for an expired context deadline, got nil")
	}
}

func TestDuplicateUniqueCreateFails(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().SetEmail("dup@example.com").SetPhoneNum("555").Exec(ctx)
	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	_, err = db.User.Create().SetEmail("dup@example.com").SetPhoneNum("556").Exec(ctx)
	if err == nil {
		t.Error("expected a unique constraint violation on duplicate email, got nil error")
	}
}

func TestConcurrentDuplicateCreateOnlyOneSucceeds(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	const n = 10
	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0
	errCount := 0

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := db.User.Create().SetEmail("race@example.com").SetPhoneNum("race-phone").Exec(ctx)
			mu.Lock()
			defer mu.Unlock()
			if err == nil {
				successCount++
			} else {
				errCount++
			}
		}()
	}
	wg.Wait()

	if successCount != 1 {
		t.Errorf("expected exactly 1 concurrent create to succeed under the unique constraint, got %d successes and %d errors", successCount, errCount)
	}
}

func TestWhitespacePaddedEmailNotTreatedAsDuplicate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().SetEmail("dup2@example.com").SetPhoneNum("601").Exec(ctx)
	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}
	_, err = db.User.Create().SetEmail(" dup2@example.com").SetPhoneNum("602").Exec(ctx)
	if err != nil {
		t.Fatalf("expected leading-whitespace email to be treated as a distinct value, create failed: %v", err)
	}

	res, err := db.User.FindUnique(user.Email.EQ(" dup2@example.com")).Exec(ctx)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if res == nil || res.PhoneNum != "602" {
		t.Errorf("expected exact match including leading whitespace, since no normalization should be silently applied, got: %+v", res)
	}
}

func TestEmailCaseSensitivity(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().SetEmail("CaseTest@Example.com").SetPhoneNum("603").Exec(ctx)
	if err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	res, err := db.User.FindUnique(user.Email.EQ("casetest@example.com")).Exec(ctx)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if res != nil {
		t.Errorf("query matched a differently-cased email (%+v); confirm this is an intentional case-insensitive collation and not an accidental DB default", res)
	}
}

func TestControlCharacterInFilter(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	res, err := db.User.FindFirst(user.Email.EQ("test\x00null@example.com")).Exec(ctx)
	if err == nil {
		t.Fatalf("expected validation error for query with an embedded null byte, got nil")
	}
	if res != nil {
		t.Errorf("expected no match for a value containing a null byte, got: %+v", res)
	}
}

func TestVeryLongEmailValue(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	longEmail := strings.Repeat("a", 5000) + "@example.com"

	_, err := db.User.Create().SetEmail(longEmail).SetPhoneNum("999").Exec(ctx)
	if err != nil {
		t.Fatalf("create with a very long email failed: %v", err)
	}

	res, err := db.User.FindUnique(user.Email.EQ(longEmail)).Exec(ctx)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if res == nil || res.Email != longEmail {
		gotLen := 0
		if res != nil {
			gotLen = len(res.Email)
		}
		t.Errorf("long email value was not stored/retrieved exactly, expected len=%d got len=%d", len(longEmail), gotLen)
	}
}

func TestCreateWithEmptyStringEmail(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().SetEmail("").SetPhoneNum("800").Exec(ctx)
	if err != nil {
		return
	}

	res, err := db.User.FindUnique(user.Email.EQ("")).Exec(ctx)
	if err != nil {
		t.Fatalf("query for empty-string email failed: %v", err)
	}
	if res == nil {
		t.Error("empty string email was accepted on create but cannot be queried back via FindUnique")
	}
}

func TestOptionalEnumNullVsValueFilter(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	adminRole := valk.UserRole.Admin
	_, err := db.User.Create().SetEmail("role_set@example.com").SetPhoneNum("700").SetRoleOptional(adminRole).Exec(ctx)
	if err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	_, err = db.User.Create().SetEmail("role_unset@example.com").SetPhoneNum("701").Exec(ctx)
	if err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	nullRoleUsers, err := db.User.FindMany(user.RoleOptional.IsNull()).Exec(ctx)
	if err != nil {
		t.Fatalf("query for null role failed: %v", err)
	}
	if len(nullRoleUsers) != 1 || nullRoleUsers[0].Email != "role_unset@example.com" {
		t.Errorf("expected exactly role_unset@example.com for a null-role filter, got: %+v", nullRoleUsers)
	}

	setRoleUsers, err := db.User.FindMany(user.RoleOptional.EQ(adminRole)).Exec(ctx)
	if err != nil {
		t.Fatalf("query for admin role failed: %v", err)
	}
	if len(setRoleUsers) != 1 || setRoleUsers[0].Email != "role_set@example.com" {
		t.Errorf("expected exactly role_set@example.com for an admin-role filter, got: %+v", setRoleUsers)
	}
}

func TestSQLInjectionVariants(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().SetEmail("injection_target@example.com").SetPhoneNum("900").Exec(ctx)
	if err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	payloads := []string{
		"' OR '1'='1",
		"'; DROP TABLE \"User\"; --",
		"' UNION SELECT * FROM \"User\" --",
		"\\'; --",
		"%' OR '1'='1",
		"injection_target@example.com'--",
		"' OR 1=1#",
	}

	for _, payload := range payloads {
		t.Run(payload, func(t *testing.T) {
			res, err := db.User.FindFirst(user.Email.EQ(payload)).Exec(ctx)
			if err != nil {
				t.Fatalf("query crashed on payload %q: %v", payload, err)
			}
			if res != nil {
				t.Errorf("payload %q unexpectedly matched a row: %+v", payload, res)
			}
		})
	}

	sanity, err := db.User.FindUnique(user.Email.EQ("injection_target@example.com")).Exec(ctx)
	if err != nil || sanity == nil {
		t.Fatalf("sanity check failed after injection attempts: the seed row should still exist, err=%v res=%+v", err, sanity)
	}
}

func TestFindManyReturnsEmptySliceNotNilWhenNoMatches(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	res, err := db.User.FindMany(user.Email.EQ("definitely_not_present@example.com")).Exec(ctx)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("expected 0 results, got %d", len(res))
	}
}

func TestCompoundUniqueWithOneFieldMatchingWrongRow(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().SetEmail("compound_a@example.com").SetPhoneNum("701a").Exec(ctx)
	if err != nil {
		t.Fatalf("seed a failed: %v", err)
	}
	_, err = db.User.Create().SetEmail("compound_b@example.com").SetPhoneNum("701b").Exec(ctx)
	if err != nil {
		t.Fatalf("seed b failed: %v", err)
	}

	res, err := db.User.FindUnique(user.EmailPhoneUnique("compound_a@example.com", "701b")).Exec(ctx)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if res != nil {
		t.Errorf("expected nil since email and phone belong to different rows, got: %+v", res)
	}
}

func TestCompoundUniqueConstraintEdgeCases(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().SetEmail("compound_edge@example.com").SetPhoneNum("800a").Exec(ctx)
	if err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	t.Run("Happy path", func(t *testing.T) {
		res, err := db.User.FindUnique(user.EmailPhoneUnique("compound_edge@example.com", "800a")).Exec(ctx)
		if err != nil {
			t.Fatalf("happy path failed: %v", err)
		}
		if res == nil || res.Email != "compound_edge@example.com" {
			t.Errorf("expected to retrieve seeded row, got: %+v", res)
		}
	})

	t.Run("SQL Injection in one compound field", func(t *testing.T) {
		res, err := db.User.FindUnique(user.EmailPhoneUnique("compound_edge@example.com", "800a' OR '1'='1")).Exec(ctx)
		if err != nil {
			t.Fatalf("query failed: %v", err)
		}
		if res != nil {
			t.Errorf("expected no match due to SQL injection payload, but got row: %+v", res)
		}
	})

	t.Run("Partial mismatch", func(t *testing.T) {
		res, err := db.User.FindUnique(user.EmailPhoneUnique("compound_edge@example.com", "wrong_phone")).Exec(ctx)
		if err != nil {
			t.Fatalf("query failed: %v", err)
		}
		if res != nil {
			t.Errorf("expected nil on partial mismatch, got: %+v", res)
		}
	})

	t.Run("Empty strings in both compound fields", func(t *testing.T) {
		res, err := db.User.FindUnique(user.EmailPhoneUnique("", "")).Exec(ctx)
		if err != nil {
			t.Fatalf("query failed: %v", err)
		}
		if res != nil {
			t.Errorf("expected nil for empty strings, got: %+v", res)
		}
	})

	t.Run("Control characters", func(t *testing.T) {
		res, err := db.User.FindUnique(user.EmailPhoneUnique("compound_edge@example.com\x00", "800a\r\n")).Exec(ctx)
		if err == nil {
			t.Fatalf("expected validation error for control-character mutated fields, got nil")
		}
		if res != nil {
			t.Errorf("expected nil for control-character mutated fields, got: %+v", res)
		}
	})
}

func TestJsonField(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u, err := db.User.Create().SetEmail("json_test_user@example.com").SetPhoneNum("555-json").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	p, err := db.Post.Create().SetTitle("JSON Post").SetAuthorId(u.Id).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	metaVal := json.RawMessage(`{"tags":["valkyrie","orm"],"version":1}`)
	c, err := db.Comment.Create().SetTextify(1).SetDummy3("dummy3").SetDummy1(10).SetDummy2("dummy2").SetPostId(p.Id).SetAuthorId(u.Id).SetMeta(metaVal).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create comment with JSON: %v", err)
	}

	// 1. FindFirst by exact JSON matching
	found, err := db.Comment.FindFirst(comment.Meta.EQ(metaVal)).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to find comment by json EQ: %v", err)
	}
	if found == nil || found.Id != c.Id {
		t.Errorf("expected to find comment %s, got %v", c.Id, found)
	}

	// 2. FindMany using IN operator with JSON slices
	foundMany, err := db.Comment.FindMany(comment.Meta.In([]json.RawMessage{metaVal})).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to find comments by json IN: %v", err)
	}
	if len(foundMany) != 1 || foundMany[0].Id != c.Id {
		t.Errorf("expected 1 comment, got %d comments", len(foundMany))
	}
}

func TestFindManyTakeAndSkip(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	var usersToCreate []*user.CreateBuilder

	for i := range 5 {
		usersToCreate = append(usersToCreate, db.User.Create().
			SetEmail(fmt.Sprintf("user%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("99%d", i)),
		)
	}

	_, err := db.User.CreateMany(usersToCreate...).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create users: %v", err)

	}

	resTake, err := db.User.FindMany().Take(2).Exec(ctx)
	if err != nil {
		t.Fatalf("FindMany Take(2) failed: %v", err)
	}
	if len(resTake) != 2 {
		t.Errorf("expected 2 users, got %d", len(resTake))
	}

	resSkip, err := db.User.FindMany().Skip(3).Exec(ctx)
	if err != nil {
		t.Fatalf("FindMany Skip(3) failed: %v", err)
	}
	if len(resSkip) != 2 {
		t.Errorf("expected 2 users on skip 3, got %d", len(resSkip))
	}

	resTakeSkip, err := db.User.FindMany().Take(2).Skip(1).Exec(ctx)
	if err != nil {
		t.Fatalf("FindMany Take(2) Skip(1) failed: %v", err)
	}
	if len(resTakeSkip) != 2 {
		t.Errorf("expected 2 users, got %d", len(resTakeSkip))
	}
}

func TestFindFirstSkip(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	var usersToCreate []*user.CreateBuilder
	for i := range 5 {
		usersToCreate = append(usersToCreate, db.User.Create().
			SetEmail(fmt.Sprintf("ff_user%d@example.com", i)).
			SetPhoneNum(fmt.Sprintf("88%d", i)),
		)
	}
	_, err := db.User.CreateMany(usersToCreate...).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create users: %v", err)
	}

	res, err := db.User.FindFirst().Skip(2).Exec(ctx)
	if err != nil {
		t.Fatalf("FindFirst with Skip failed: %v", err)
	}
	if res == nil {
		t.Fatal("expected to find a user")
	}
	if !strings.HasPrefix(res.Email, "ff_user") {
		t.Errorf("expected found user to be one of the seeded users, got %q", res.Email)
	}
}

func TestRelationSubqueries(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	uA, err := db.User.Create().SetEmail("user_a@example.com").SetPhoneNum("111").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed User A: %v", err)
	}
	uB, err := db.User.Create().SetEmail("user_b@example.com").SetPhoneNum("222").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed User B: %v", err)
	}
	uC, err := db.User.Create().SetEmail("user_c@example.com").SetPhoneNum("333").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed User C: %v", err)
	}

	p1, err := db.Post.Create().SetTitle("Alpha").SetPublished(true).SetAuthorId(uA.Id).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed Post 1: %v", err)
	}
	_, err = db.Post.Create().SetTitle("Beta").SetPublished(false).SetAuthorId(uA.Id).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed Post 2: %v", err)
	}
	p3, err := db.Post.Create().SetTitle("Gamma").SetPublished(true).SetAuthorId(uA.Id).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed Post 3: %v", err)
	}
	_, err = db.Post.Create().SetTitle("Delta").SetPublished(true).SetAuthorId(uA.Id).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed Post 4: %v", err)
	}

	_, err = db.Comment.Create().SetTextify(100).SetDummy3("c1").SetDummy1(1).SetDummy2("d").SetPostId(p1.Id).SetAuthorId(uB.Id).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed Comment 1: %v", err)
	}
	_, err = db.Comment.Create().SetTextify(200).SetDummy3("c2").SetDummy1(2).SetDummy2("d").SetPostId(p1.Id).SetAuthorId(uA.Id).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed Comment 2: %v", err)
	}
	_, err = db.Comment.Create().SetTextify(300).SetDummy3("c3").SetDummy1(3).SetDummy2("d").SetPostId(p1.Id).SetAuthorId(uC.Id).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed Comment 3: %v", err)
	}

	_, err = db.Comment.Create().SetTextify(400).SetDummy3("c4").SetDummy1(4).SetDummy2("d").SetPostId(p3.Id).SetAuthorId(uB.Id).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed Comment 4: %v", err)
	}

	t.Run("Deep Nesting with Multi-Level Filters, Sorting, and Pagination", func(t *testing.T) {
		res, err := db.User.FindFirst(user.Email.EQ("user_a@example.com")).Select(user.Select{
			Email: true,
			Posts: post.Query().
				Where(post.Published.EQ(true)).
				OrderBy(post.Title.Asc()).
				Take(2).
				Skip(1).
				Select(post.Select{
					Title: true,
					Comments: comment.Query().
						OrderBy(comment.Textify.Desc()).
						Take(1).
						Select(comment.Select{
							Textify: true,
							Dummy3:  true,
							Author: &user.Select{
								Email: true,
							},
						}),
				}),
		}).Exec(ctx)

		if err != nil {
			t.Fatalf("deep query failed: %v", err)
		}
		if res == nil {
			t.Fatal("expected User A to be found")
		}

		if len(res.Posts) != 2 {
			t.Fatalf("expected exactly 2 posts, got %d", len(res.Posts))
		}

		pFirst := res.Posts[0]
		pSecond := res.Posts[1]
		if pFirst.Title != "Delta" || pSecond.Title != "Gamma" {
			t.Errorf("expected posts 'Delta' and 'Gamma', got %q and %q", pFirst.Title, pSecond.Title)
		}

		if len(pFirst.Comments) != 0 {
			t.Errorf("expected 'Delta' to have 0 comments, got %d", len(pFirst.Comments))
		}

		if len(pSecond.Comments) != 1 {
			t.Fatalf("expected 'Gamma' to have 1 comment, got %d", len(pSecond.Comments))
		}
		c := pSecond.Comments[0]
		if c.Dummy3 != "c4" || c.Textify != 400 {
			t.Errorf("expected comment 'c4' with textify 400, got %+v", c)
		}
		if c.Author == nil || c.Author.Email != "user_b@example.com" {
			t.Errorf("expected comment author to be user_b@example.com, got %+v", c.Author)
		}
	})

	t.Run("Extreme Subquery Pagination Limits (Take=0, Skip=100)", func(t *testing.T) {
		res, err := db.User.FindFirst(user.Email.EQ("user_a@example.com")).Select(user.Select{
			Email: true,
			Posts: post.Query().Take(0).Select(post.Select{Title: true}),
		}).Exec(ctx)
		if err != nil {
			t.Fatalf("query failed: %v", err)
		}
		if res == nil || len(res.Posts) != 0 {
			t.Errorf("expected 0 posts with Take(0), got: %+v", res)
		}

		resSkip, err := db.User.FindFirst(user.Email.EQ("user_a@example.com")).Select(user.Select{
			Email: true,
			Posts: post.Query().Skip(100).Select(post.Select{Title: true}),
		}).Exec(ctx)
		if err != nil {
			t.Fatalf("query failed: %v", err)
		}
		if resSkip == nil || len(resSkip.Posts) != 0 {
			t.Errorf("expected 0 posts with Skip(100), got: %+v", resSkip)
		}
	})

	t.Run("Complex Logical Operator Filtering inside Subquery", func(t *testing.T) {
		res, err := db.User.FindFirst(user.Email.EQ("user_a@example.com")).Select(user.Select{
			Email: true,
			Posts: post.Query().
				Where(post.And(
					post.Published.EQ(true),
					post.Or(
						post.Title.EQ("Alpha"),
						post.Title.EQ("Gamma"),
					),
				)).
				OrderBy(post.Title.Asc()).
				Select(post.Select{Title: true}),
		}).Exec(ctx)

		if err != nil {
			t.Fatalf("query failed: %v", err)
		}
		if res == nil || len(res.Posts) != 2 || res.Posts[0].Title != "Alpha" || res.Posts[1].Title != "Gamma" {
			t.Errorf("expected posts 'Alpha' and 'Gamma', got: %+v", res)
		}
	})
}

func TestFindUniqueExtended(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	adminRole := valk.UserRole.Admin
	studentRole := valk.UserRole.Student

	_, err := db.User.Create().SetEmail("ext_admin@example.com").SetPhoneNum("1111").SetRoleOptional(adminRole).Exec(ctx)
	if err != nil {
		t.Fatalf("seed admin failed: %v", err)
	}

	_, err = db.User.Create().SetEmail("ext_student@example.com").SetPhoneNum("2222").SetRoleOptional(studentRole).Exec(ctx)
	if err != nil {
		t.Fatalf("seed student failed: %v", err)
	}

	t.Run("Match unique and non-unique", func(t *testing.T) {
		res, err := db.User.FindUnique(
			user.Email.EQ("ext_admin@example.com"),
			user.RoleOptional.EQ(adminRole),
			user.PhoneNum.EQ("1111"),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("query failed: %v", err)
		}
		if res == nil || res.Email != "ext_admin@example.com" {
			t.Errorf("expected to find admin user, got: %+v", res)
		}
	})

	t.Run("Mismatch on non-unique field", func(t *testing.T) {
		res, err := db.User.FindUnique(
			user.Email.EQ("ext_admin@example.com"),
			user.RoleOptional.EQ(studentRole), // Mismatching role
		).Exec(ctx)
		if err != nil {
			t.Fatalf("query failed: %v", err)
		}
		if res != nil {
			t.Errorf("expected nil result on role mismatch, got: %+v", res)
		}
	})
}
