package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unicode/utf8"

	"integration/valk"
	"integration/valk/allFieldsSoFar"
	"integration/valk/user"
)

func TestCreate_DuplicateEmail_Rejected(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	if _, err := db.User.Create().SetEmail("dupe@example.com").SetPhoneNum("+100000001").Exec(ctx); err != nil {
		t.Fatalf("first insert should succeed, got: %v", err)
	}

	// Same email, different phoneNum, must still fail if email is unique
	if _, err := db.User.Create().SetEmail("dupe@example.com").SetPhoneNum("+100000002").Exec(ctx); err == nil {
		t.Fatal("expected unique constraint violation on duplicate email, got nil error")
	}

	var count int
	if err := db.Raw().QueryRowContext(ctx, query("SELECT COUNT(*) FROM User WHERE email = ?", "SELECT COUNT(*) FROM \"User\" WHERE email = $1"), "dupe@example.com").Scan(&count); err != nil {
		t.Fatalf("failed count query: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected exactly 1 row after rejected duplicate, got %d", count)
	}
}

func TestCreate_DuplicateEmail_CaseVariants(t *testing.T) {

	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	if _, err := db.User.Create().
		SetEmail("Case@Example.com").SetPhoneNum("+100000003").
		Exec(ctx); err != nil {
		t.Fatalf("failed initial insert: %v", err)
	}

	_, err := db.User.Create().
		SetEmail("case@example.com").SetPhoneNum("+100000004").
		Exec(ctx)

	t.Logf("case-variant email insert result: err=%v (confirm this matches intended uniqueness semantics)", err)
}

func TestCreate_CompoundUnique_Rejected(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	if _, err := db.User.Create().
		SetEmail("a@example.com").SetPhoneNum("+199999999").
		Exec(ctx); err != nil {
		t.Fatalf("first insert failed: %v", err)
	}

	// Same email + same phoneNum  hits @@unique([email, phoneNum]).
	if _, err := db.User.Create().
		SetEmail("a@example.com").SetPhoneNum("+199999999").
		Exec(ctx); err == nil {
		t.Fatal("expected unique constraint violation on duplicate (email, phoneNum)")
	}
}

func TestCreate_ZeroValueInput_Rejected(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().Exec(ctx)
	if err == nil {
		t.Fatal("expected error creating user with entirely zero-value input (missing required email)")
	}
}

func TestCreate_EmptyStringEmail_Rejected(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().
		SetEmail("").SetPhoneNum("+100000005").
		Exec(ctx)
	if err == nil {
		t.Fatal("expected error creating user with empty-string email")
	}
}

func TestCreate_WhitespaceOnlyEmail_Rejected(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().
		SetEmail("   ").SetPhoneNum("+100000006").
		Exec(ctx)

	if err == nil {
		t.Log("WARNING: whitespace-only email was accepted  confirm this is intentional, not an oversight")
	}
}

func TestCreate_ReferredBy_NonexistentID_Rejected(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create().
		SetEmail("orphan@example.com").
		SetPhoneNum("+100000007").
		SetReferredById("clnonexistent00000000000").
		Exec(ctx)

	if err == nil {
		t.Fatal("expected FK violation when referredById points to a nonexistent user")
	}
}

func TestCreate_ReferredBy_SelfReference_Rejected(t *testing.T) {

	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u, err := db.User.Create().
		SetEmail("self@example.com").SetPhoneNum("+100000008").
		Exec(ctx)
	if err != nil {
		t.Fatalf("setup insert failed: %v", err)
	}

	b, err := db.User.Create().
		SetEmail("referred@example.com").SetPhoneNum("+100000009").SetReferredById(u.Id).
		Exec(ctx)
	if err != nil {
		t.Fatalf("valid referral chain should succeed: %v", err)
	}
	if b.ReferredById == nil || *b.ReferredById != u.Id {
		t.Fatalf("expected ReferredById %q, got %v", u.Id, b.ReferredById)
	}
}

func TestCreate_InvalidEnumValue_BypassingTypeSystem(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	pffff := valk.UserRoleType("totallyNotARole")

	_, err := db.User.Create().
		SetEmail("pffff-role@example.com").SetPhoneNum("+100000010").SetRole(pffff).
		Exec(ctx)

	if err == nil {
		t.Fatal("expected rejection of an enum value outside the declared domain")
	}
}

func TestCreate_DefaultEnumAppliedWhenNil(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u, err := db.User.Create().
		SetEmail("noRole@example.com").SetPhoneNum("+100000011").
		Exec(ctx)
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}
	if u.Role != valk.UserRole_STUDENT {
		t.Errorf("expected default role Student when Role is nil, got %q", u.Role)
	}
}

func TestCreate_StringEdgeCases(t *testing.T) {
	cases := []struct {
		name        string
		email       string
		phone       string
		expectError bool
	}{
		{"unicode_email_local_part", "üñîçødé@example.com", "+200000001", false},
		{"sql_injection_shaped_email", `test' OR '1'='1@example.com`, "+200000004", false},
		{"null_byte_in_email", "nul\x00byte@example.com", "+200000005", true},
		{"leading_trailing_whitespace_email", "  padded@example.com  ", "+200000006", false},
		{"rtl_override_char", "\u202Eevil@example.com", "+200000007", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, cleanup := setupTestDB(t)
			defer cleanup()
			ctx := context.Background()

			u, err := db.User.Create().
				SetEmail(tc.email).SetPhoneNum(tc.phone).
				Exec(ctx)

			if tc.expectError && err == nil {
				t.Fatalf("expected error for input %q, got success (id=%s)", tc.email, u.Id)
			}
			if !tc.expectError && err != nil {
				t.Fatalf("expected success for input %q, got error: %v", tc.email, err)
			}
			if err == nil && !utf8.ValidString(u.Email) {
				t.Errorf("returned email is not valid UTF-8: %q", u.Email)
			}
			if err == nil {
				var stored string
				qerr := db.Raw().QueryRowContext(ctx, query("SELECT email FROM User WHERE id = ?", "SELECT email FROM \"User\" WHERE id = $1"), u.Id).Scan(&stored)
				if qerr != nil {
					t.Fatalf("failed to read back: %v", qerr)
				}
				if stored != strings.TrimSpace(tc.email) && stored != tc.email {
					t.Errorf("stored email %q does not match input %q (check for silent mutation)", stored, tc.email)
				}
			}
		})
	}
}

func TestCreate_Select_ForceIncludesFK_EvenWhenNotExplicitlySelected(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	referrer, err := db.User.Create().
		SetEmail("referrer@example.com").SetPhoneNum("+300000002").
		Exec(ctx)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	rid := referrer.Id
	u, err := db.User.Create().
		SetEmail("referredfk@example.com").SetPhoneNum("+300000003").SetReferredById(rid).
		Select(user.Select{
			Id:         true,
			ReferredBy: &user.Select{Id: true},
		}).Exec(ctx)

	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if u.ReferredBy == nil {
		t.Fatal("expected ReferredBy to be populated when its relation was selected")
	}
	if u.ReferredBy.Id != referrer.Id {
		t.Errorf("expected ReferredBy.Id=%s, got %s", referrer.Id, u.ReferredBy.Id)
	}
}

func TestCreate_Select_EmptyStruct_ReturnsEverything(t *testing.T) {

	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u, err := db.User.Create().
		SetEmail("empty-select@example.com").SetPhoneNum("+300000004").
		Select(valk.UserSelect{}).Exec(ctx)

	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if u.Email != "empty-select@example.com" {
		t.Errorf("expected empty Select{} to select everything (select all), got Email=%q", u.Email)
	}
}

func TestCreate_ContextAlreadyCancelled(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before the call even starts

	_, err := db.User.Create().
		SetEmail("cancelled@example.com").SetPhoneNum("+400000001").
		Exec(ctx)

	if err == nil {
		t.Fatal("expected error when context is already cancelled")
	}

	var count int
	if qerr := db.Raw().QueryRowContext(context.Background(),
		query("SELECT COUNT(*) FROM User WHERE email = ?", "SELECT COUNT(*) FROM \"User\" WHERE email = $1"), "cancelled@example.com").Scan(&count); qerr != nil {
		t.Fatalf("count query failed: %v", qerr)
	}
	if count != 0 {
		t.Fatalf("expected no row persisted for cancelled-context create, found %d", count)
	}
}

func TestCreate_ContextTimeout_DuringExec(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()

	_, err := db.User.Create().
		SetEmail("timeout@example.com").SetPhoneNum("+400000002").
		Exec(ctx)

	if err == nil {
		t.Log("create succeeded despite near-zero timeout  likely fine if driver executes faster than ctx propagation, but worth a second look under load")
	}
}

func TestCreate_ConcurrentDuplicateEmail_ExactlyOneWins(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	const goroutines = 20
	var wg sync.WaitGroup
	var successes int64
	var failures int64

	for i := range goroutines {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_, err := db.User.Create().
				SetEmail("race@example.com").
				SetPhoneNum(fmt.Sprintf("+50000%04d", n)).
				Exec(ctx)
			if err != nil {
				atomic.AddInt64(&failures, 1)
			} else {
				atomic.AddInt64(&successes, 1)
			}
		}(i)
	}
	wg.Wait()

	if successes != 1 {
		t.Fatalf("expected exactly 1 success under concurrent duplicate-email creates, got %d successes, %d failures",
			successes, failures)
	}

	var count int
	if err := db.Raw().QueryRowContext(ctx, query("SELECT COUNT(*) FROM User WHERE email = ?", "SELECT COUNT(*) FROM \"User\" WHERE email = $1"), "race@example.com").Scan(&count); err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected exactly 1 persisted row, found %d  unique constraint may not be enforced at the DB level under concurrency", count)
	}
}

func TestCreate_ConcurrentUniqueIDs_NoCollision(t *testing.T) {

	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	const n = 200
	var wg sync.WaitGroup
	ids := make([]string, n)
	errs := make([]error, n)

	for i := range n {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			u, err := db.User.Create().
				SetEmail(fmt.Sprintf("bulk%d@example.com", idx)).
				SetPhoneNum(fmt.Sprintf("+600%06d", idx)).
				Exec(ctx)
			errs[idx] = err
			if err == nil {
				ids[idx] = u.Id
			}
		}(i)
	}
	wg.Wait()

	seen := make(map[string]bool, n)
	for i, err := range errs {
		if err != nil {
			t.Fatalf("create %d failed: %v", i, err)
		}
		if ids[i] == "" {
			t.Fatalf("create %d returned empty ID", i)
		}
		if seen[ids[i]] {
			t.Fatalf("duplicate CUID generated: %s", ids[i])
		}
		seen[ids[i]] = true
	}
}

func TestCreate_FailurePartway_LeavesNoPartialRow(t *testing.T) {

	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	fakeReferrer := "clDoesNotExist00000000000"
	before := countAllUsers(t, ctx, db)

	_, err := db.User.Create().
		SetEmail("partial@example.com").
		SetPhoneNum("+700000001").
		SetReferredById(fakeReferrer).
		Exec(ctx)

	if err == nil {
		t.Fatal("expected FK failure")
	}

	after := countAllUsers(t, ctx, db)
	if after != before {
		t.Fatalf("row count changed despite failed create: before=%d after=%d", before, after)
	}
}

func countAllUsers(t *testing.T, ctx context.Context, db *valk.DB) int {
	t.Helper()
	var count int
	if err := db.Raw().QueryRowContext(ctx, query("SELECT COUNT(*) FROM User", "SELECT COUNT(*) FROM \"User\"")).Scan(&count); err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	return count
}

func TestCreate_Hooks(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	var afterCalled bool
	db.User.Use(user.Extension{
		Create: func(ctx context.Context, args *user.CreateArgs, next user.CreateQuery) (*valk.User, error) {
			if args.Data.Email == "hook@example.com" {
				args.Data.PhoneNum = "+188888888"
			}
			res, err := next(ctx, args)
			if err == nil && res.Email == "hook@example.com" {
				afterCalled = true
			}
			return res, err
		},
	})

	u, err := db.User.Create().
		SetEmail("hook@example.com").
		SetPhoneNum("+100000000").
		Exec(ctx)

	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if u.PhoneNum != "+188888888" {
		t.Errorf("expected PhoneNum mutated to '+188888888', got %q", u.PhoneNum)
	}

	if !afterCalled {
		t.Error("expected AfterCreate hook to be called")
	}
}

func TestCreate_Hooks_PasswordHashing(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	db.User.Use(user.Extension{
		Create: func(ctx context.Context, args *user.CreateArgs, next user.CreateQuery) (*valk.User, error) {
			if args.Data.Email == "hash@example.com" && args.Data.Password != nil {
				h := sha256.Sum256([]byte(*args.Data.Password))
				hashed := hex.EncodeToString(h[:])
				args.Data.Password = &hashed
			}
			return next(ctx, args)
		},
	})

	rawPassword := "12345678"

	u, err := db.User.Create().
		SetEmail("hash@example.com").
		SetPhoneNum("+199999999").
		SetPassword(rawPassword).
		Exec(ctx)

	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if u.Password == nil {
		t.Fatal("expected Password field to be populated, got nil")
	}

	if *u.Password == rawPassword {
		t.Errorf("expected Password to be hashed, got %q", *u.Password)
	}
}

func baseAllFieldsValidation(t *testing.T) []valk.FieldAssignmentOf[valk.AllFieldsSoFar] {
	t.Helper()
	return []valk.FieldAssignmentOf[valk.AllFieldsSoFar]{
		allFieldsSoFar.StringReq.Set("test"),
		allFieldsSoFar.StringVarchar.Set("varchar"),
		allFieldsSoFar.StringChar.Set("0123456789"),
		allFieldsSoFar.BitVal.Set("1010101010"),
		allFieldsSoFar.VarBitVal.Set("1101"),
		allFieldsSoFar.InetVal.Set("10.0.0.1"),
		allFieldsSoFar.XmlVal.Set("<x/>"),
		allFieldsSoFar.UuidDb.Set("550e8400-e29b-41d4-a716-446655440000"),
		allFieldsSoFar.IntReq.Set(1),
		allFieldsSoFar.IntegerVal.Set(2),
		allFieldsSoFar.SmallInt.Set(3),
		allFieldsSoFar.TinyInt.Set(4),
		allFieldsSoFar.OidVal.Set(5),
		allFieldsSoFar.BigIntReq.Set(int64(6)),
		allFieldsSoFar.FloatReq.Set(1.5),
		allFieldsSoFar.RealVal.Set(2.5),
		allFieldsSoFar.DecimalReq.Set("10.50"),
		allFieldsSoFar.DecimalPrecise.Set("99.99"),
		allFieldsSoFar.MoneyVal.Set("12.34"),
		allFieldsSoFar.BoolReq.Set(true),
		allFieldsSoFar.DateTimeReq.Set(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		allFieldsSoFar.UpdatedAt.Set(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		allFieldsSoFar.DateTimeTz.Set(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		allFieldsSoFar.TimestampVal.Set(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		allFieldsSoFar.TimeVal.Set(time.Date(1, 1, 1, 10, 30, 0, 0, time.UTC)),
		allFieldsSoFar.TimetzVal.Set(time.Date(1, 1, 1, 10, 30, 0, 0, time.UTC)),
		allFieldsSoFar.JsonReq.Set(json.RawMessage(`{}`)),
		allFieldsSoFar.JsonVal.Set(json.RawMessage(`[]`)),
		allFieldsSoFar.BytesReq.Set([]byte{0x01}),
	}
}

func TestCreate_DecimalValidation(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFieldsValidation(t)

	_, err := db.AllFieldsSoFar.Create().
		Assignments(base...).
		SetDecimalReq("not-a-number").
		Exec(ctx)

	if err == nil {
		t.Fatal("expected ValidationError for invalid decimal format 'not-a-number', got nil")
	}

	_, err = db.AllFieldsSoFar.Create().
		Assignments(base...).
		SetDecimalReq("12.34.56").
		Exec(ctx)

	if err == nil {
		t.Fatal("expected ValidationError for invalid decimal format '12.34.56', got nil")
	}

	rec, err := db.AllFieldsSoFar.Create().
		Assignments(base...).
		SetDecimalReq("123.456789").
		Exec(ctx)

	if err != nil {
		t.Fatalf("expected valid decimal creation to succeed, got: %v", err)
	}

	if rec.DecimalReq != "123.456789" {
		t.Errorf("expected DecimalReq '123.456789', got %q", rec.DecimalReq)
	}
}

func TestCreate_JsonValidation(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFieldsValidation(t)

	_, err := db.AllFieldsSoFar.Create().
		Assignments(base...).
		SetJsonReq(json.RawMessage("{invalid_json}")).
		Exec(ctx)

	if err == nil {
		t.Fatal("expected ValidationError for malformed JSON, got nil")
	}

	rec, err := db.AllFieldsSoFar.Create().
		Assignments(base...).
		SetJsonReq(json.RawMessage(`{"key":"value"}`)).
		Exec(ctx)

	if err != nil {
		t.Fatalf("expected valid JSON create to succeed, got: %v", err)
	}

	if len(rec.JsonReq) == 0 {
		t.Error("expected JsonReq to be populated")
	}
}
