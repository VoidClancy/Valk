package main

import (
	"context"
	"encoding/json"
	"integration/valk"
	"integration/valk/comment"
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

	_, err := db.User.Create(valk.UserCreate{
		Email:    "onlyuser@example.com",
		PhoneNum: "000",
	}).Exec(ctx)
	if err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	res, err := db.User.FindUnique(nil).Exec(ctx)
	if err == nil && res != nil {
		t.Errorf("FindUnique with a zero-value where matched a row unexpectedly (%+v); it should require at least one unique field or return an error", res)
	}
}

func TestFindUniqueConflictingCompoundFields(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create(valk.UserCreate{Email: "a@example.com", PhoneNum: "111"}).Exec(ctx)
	if err != nil {
		t.Fatalf("seed a failed: %v", err)
	}
	_, err = db.User.Create(valk.UserCreate{Email: "b@example.com", PhoneNum: "222"}).Exec(ctx)
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

	_, err := db.User.Create(valk.UserCreate{Email: "empty_select@example.com", PhoneNum: "333"}).Exec(ctx)
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

	_, err := db.User.Create(valk.UserCreate{Email: "omit_all@example.com", PhoneNum: "334"}).Exec(ctx)
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

	u, err := db.User.Create(valk.UserCreate{Email: "omit_id@example.com", PhoneNum: "335"}).Exec(ctx)
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

	_, err := db.User.Create(valk.UserCreate{Email: "noposts@example.com", PhoneNum: "444"}).Exec(ctx)
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

	author, err := db.User.Create(valk.UserCreate{Email: "unique_rel@example.com", PhoneNum: "445"}).Exec(ctx)
	if err != nil {
		t.Fatalf("seed author failed: %v", err)
	}
	_, err = db.Post.Create(valk.PostCreate{Title: "Unique Rel Post", AuthorId: author.Id}).Exec(ctx)
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

	author, err := db.User.Create(valk.UserCreate{Email: "first_rel@example.com", PhoneNum: "446"}).Exec(ctx)
	if err != nil {
		t.Fatalf("seed author failed: %v", err)
	}
	_, err = db.Post.Create(valk.PostCreate{Title: "First Rel Post", AuthorId: author.Id}).Exec(ctx)
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

	_, err := db.User.Create(valk.UserCreate{Email: "dup@example.com", PhoneNum: "555"}).Exec(ctx)
	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	_, err = db.User.Create(valk.UserCreate{Email: "dup@example.com", PhoneNum: "556"}).Exec(ctx)
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
			_, err := db.User.Create(valk.UserCreate{
				Email:    "race@example.com",
				PhoneNum: "race-phone",
			}).Exec(ctx)
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

	_, err := db.User.Create(valk.UserCreate{Email: "dup2@example.com", PhoneNum: "601"}).Exec(ctx)
	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}
	_, err = db.User.Create(valk.UserCreate{Email: " dup2@example.com", PhoneNum: "602"}).Exec(ctx)
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

	_, err := db.User.Create(valk.UserCreate{Email: "CaseTest@Example.com", PhoneNum: "603"}).Exec(ctx)
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

	_, err := db.User.Create(valk.UserCreate{Email: longEmail, PhoneNum: "999"}).Exec(ctx)
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

	_, err := db.User.Create(valk.UserCreate{
		Email:    "",
		PhoneNum: "800",
	}).Exec(ctx)
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
	_, err := db.User.Create(valk.UserCreate{
		Email:        "role_set@example.com",
		PhoneNum:     "700",
		RoleOptional: &adminRole,
	}).Exec(ctx)
	if err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	_, err = db.User.Create(valk.UserCreate{
		Email:    "role_unset@example.com",
		PhoneNum: "701",
	}).Exec(ctx)
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

	_, err := db.User.Create(valk.UserCreate{Email: "injection_target@example.com", PhoneNum: "900"}).Exec(ctx)
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

	_, err := db.User.Create(valk.UserCreate{Email: "compound_a@example.com", PhoneNum: "701a"}).Exec(ctx)
	if err != nil {
		t.Fatalf("seed a failed: %v", err)
	}
	_, err = db.User.Create(valk.UserCreate{Email: "compound_b@example.com", PhoneNum: "701b"}).Exec(ctx)
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

	_, err := db.User.Create(valk.UserCreate{
		Email:    "compound_edge@example.com",
		PhoneNum: "800a",
	}).Exec(ctx)
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

	u, err := db.User.Create(user.Create{
		Email:    "json_test_user@example.com",
		PhoneNum: "555-json",
	}).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	p, err := db.Post.Create(valk.PostCreate{
		Title:    "JSON Post",
		AuthorId: u.Id,
	}).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	metaVal := json.RawMessage(`{"tags":["valkyrie","orm"],"version":1}`)
	c, err := db.Comment.Create(valk.CommentCreate{
		Textify:  1,
		Dummy3:   "dummy3",
		Dummy1:   10,
		Dummy2:   "dummy2",
		PostId:   p.Id,
		AuthorId: u.Id,
		Meta:     &metaVal,
	}).Exec(ctx)
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
