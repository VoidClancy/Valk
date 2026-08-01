package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"integration/valk"
	"integration/valk/post"
	"integration/valk/user"
)

func TestErrorHandling(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	t.Run("1. CREATE - Error Domains & Constraints", func(t *testing.T) {
		// A. Client-Side Validation Error (Missing required field: phoneNum)
		t.Run("Validation Error - Required Field Missing", func(t *testing.T) {
			_, err := db.User.Create().SetEmail("valid@example.com").Exec(ctx)
			if err == nil {
				t.Fatalf("expected error when required field PhoneNum is missing, got nil")
			}
			if !valk.IsValidationError(err) {
				t.Errorf("expected valk.IsValidationError(err) to be true, got false for: %v", err)
			}
			var vErr valk.ValidationError
			if !errors.As(err, &vErr) {
				t.Errorf("expected errors.As to match valk.ValidationError, got: %T", err)
			} else if !vErr.HasErrors() {
				t.Errorf("expected ValidationError to contain field errors")
			}
		})

		// B. Client-Side Validation Error (Null byte in string)
		t.Run("Validation Error - Safety Check (Null Byte)", func(t *testing.T) {
			_, err := db.User.Create().
				SetEmail("bad\x00email@example.com").
				SetPhoneNum("+111").
				Exec(ctx)
			if err == nil {
				t.Fatalf("expected validation error for string containing null byte, got nil")
			}
			if !valk.IsValidationError(err) {
				t.Errorf("expected valk.IsValidationError(err) to be true, got false for: %v", err)
			}
		})

		// C. Unique Constraint Violation (Duplicate Email)
		t.Run("Unique Constraint - Single Field", func(t *testing.T) {
			u1, err := db.User.Create().
				SetEmail("unique1@example.com").
				SetPhoneNum("+1001").
				Exec(ctx)
			if err != nil {
				t.Fatalf("failed to create initial user: %v", err)
			}
			if u1 == nil {
				t.Fatalf("expected created user, got nil")
			}

			// Duplicate insert
			_, err = db.User.Create().
				SetEmail("unique1@example.com").
				SetPhoneNum("+1002").
				Exec(ctx)
			if err == nil {
				t.Fatalf("expected unique constraint error on duplicate email insert, got nil")
			}
			if !valk.IsUniqueConstraint(err) {
				t.Errorf("expected valk.IsUniqueConstraint(err) to be true, got false for: %v", err)
			}
			if !valk.IsConstraintError(err) {
				t.Errorf("expected valk.IsConstraintError(err) to be true, got false for: %v", err)
			}

			var cErr *valk.ConstraintError
			if !errors.As(err, &cErr) {
				t.Errorf("expected errors.As to match *valk.ConstraintError, got: %T", err)
			}
		})

		// D. Foreign Key Constraint Violation (Insert Post with non-existent authorId)
		t.Run("Foreign Key Constraint - Non-Existent Parent", func(t *testing.T) {
			_, err := db.Post.Create().
				SetTitle("Orphan Post").
				SetAuthorId("non-existent-user-id-99999").
				Exec(ctx)
			if err == nil {
				t.Fatalf("expected foreign key constraint error when authorId does not exist, got nil")
			}
			if !valk.IsFKConstraint(err) {
				t.Errorf("expected valk.IsFKConstraint(err) to be true, got false for: %v", err)
			}
			if !valk.IsConstraintError(err) {
				t.Errorf("expected valk.IsConstraintError(err) to be true, got false for: %v", err)
			}
		})
	})

	t.Run("2. READ / QUERY - Error Domains", func(t *testing.T) {
		// A. FindUnique missing record -> returns (nil, nil)
		t.Run("FindUnique - Missing Record Returns Nil Without Error", func(t *testing.T) {
			res, err := db.User.FindUnique(user.Id.EQ("missing-id-00000")).Exec(ctx)
			if err != nil {
				t.Fatalf("expected nil error on FindUnique for missing record, got: %v", err)
			}
			if res != nil {
				t.Fatalf("expected nil result for missing record, got: %v", res)
			}
		})

		// B. FindMany missing records -> returns empty slice
		t.Run("FindMany - Missing Records Returns Empty Slice", func(t *testing.T) {
			res, err := db.User.FindMany(user.Email.EQ("nobody@example.com")).Exec(ctx)
			if err != nil {
				t.Fatalf("expected nil error on FindMany for missing records, got: %v", err)
			}
			if len(res) != 0 {
				t.Fatalf("expected 0 records, got: %d", len(res))
			}
		})

		// C. Context Canceled during query
		t.Run("Query - Context Canceled", func(t *testing.T) {
			cCtx, cancel := context.WithCancel(ctx)
			cancel() // cancel context before query execution

			_, err := db.User.FindMany(user.Email.EQ("test@example.com")).Exec(cCtx)
			if err == nil {
				t.Fatalf("expected error on cancelled context, got nil")
			}
			if !errors.Is(err, context.Canceled) {
				t.Errorf("expected errors.Is(err, context.Canceled) to be true, got: %v", err)
			}
		})
	})

	t.Run("3. UPDATE - Error Domains & Constraints", func(t *testing.T) {
		// A. Update non-existent record -> returns NotFoundError
		t.Run("Update - Non-Existent Record Returns NotFoundError", func(t *testing.T) {
			_, err := db.User.Update(user.Id.EQ("non-existent-user-id")).
				SetEmail("new@example.com").
				Exec(ctx)
			if err == nil {
				t.Fatalf("expected NotFoundError on updating missing record, got nil")
			}
			if !valk.IsNotFound(err) {
				t.Errorf("expected valk.IsNotFound(err) to be true, got false for: %v", err)
			}
			var nfErr *valk.NotFoundError
			if !errors.As(err, &nfErr) {
				t.Errorf("expected errors.As to match *valk.NotFoundError, got: %T", err)
			} else if nfErr.Model != "User" {
				t.Errorf("expected NotFoundError.Model to be 'User', got: %s", nfErr.Model)
			}
		})

		// B. Update violating Unique Constraint
		t.Run("Update - Unique Constraint Violation", func(t *testing.T) {
			uA, err := db.User.Create().SetEmail("userA@example.com").SetPhoneNum("+2001").Exec(ctx)
			if err != nil {
				t.Fatalf("failed to create userA: %v", err)
			}
			uB, err := db.User.Create().SetEmail("userB@example.com").SetPhoneNum("+2002").Exec(ctx)
			if err != nil {
				t.Fatalf("failed to create userB: %v", err)
			}

			// Try to update userB's email to userA's email
			_, err = db.User.Update(user.Id.EQ(uB.Id)).
				SetEmail(uA.Email).
				Exec(ctx)
			if err == nil {
				t.Fatalf("expected unique constraint error on update to existing email, got nil")
			}
			if !valk.IsUniqueConstraint(err) {
				t.Errorf("expected valk.IsUniqueConstraint(err) to be true, got false for: %v", err)
			}
		})

		// C. Update violating Foreign Key Constraint
		t.Run("Update - Foreign Key Violation", func(t *testing.T) {
			u, err := db.User.Create().SetEmail("author@example.com").SetPhoneNum("+2003").Exec(ctx)
			if err != nil {
				t.Fatalf("failed to create author: %v", err)
			}
			p, err := db.Post.Create().SetTitle("Post Title").SetAuthorId(u.Id).Exec(ctx)
			if err != nil {
				t.Fatalf("failed to create post: %v", err)
			}

			// Update post's authorId to a non-existent user
			_, err = db.Post.Update(post.Id.EQ(p.Id)).
				SetAuthorId("non-existent-user-id-77777").
				Exec(ctx)
			if err == nil {
				t.Fatalf("expected foreign key constraint error on updating authorId, got nil")
			}
			if !valk.IsFKConstraint(err) {
				t.Errorf("expected valk.IsFKConstraint(err) to be true, got false for: %v", err)
			}
		})
	})

	t.Run("4. DELETE - Error Domains & Constraints", func(t *testing.T) {
		// A. Delete non-existent record -> returns NotFoundError
		t.Run("Delete - Non-Existent Record Returns NotFoundError", func(t *testing.T) {
			_, err := db.User.Delete(user.Id.EQ("missing-user-id-555")).Exec(ctx)
			if err == nil {
				t.Fatalf("expected NotFoundError on deleting missing record, got nil")
			}
			if !valk.IsNotFound(err) {
				t.Errorf("expected valk.IsNotFound(err) to be true, got false for: %v", err)
			}
		})

		// B. Delete parent row with child references -> Foreign Key Constraint Violation
		t.Run("Delete - Foreign Key Violation on Parent Row", func(t *testing.T) {
			u, err := db.User.Create().SetEmail("parent@example.com").SetPhoneNum("+3001").Exec(ctx)
			if err != nil {
				t.Fatalf("failed to create parent user: %v", err)
			}
			_, err = db.Post.Create().SetTitle("Child Post").SetAuthorId(u.Id).Exec(ctx)
			if err != nil {
				t.Fatalf("failed to create child post: %v", err)
			}

			// Attempting to delete user without deleting posts first
			_, err = db.User.Delete(user.Id.EQ(u.Id)).Exec(ctx)
			if err == nil {
				t.Fatalf("expected foreign key constraint error when deleting parent user with posts, got nil")
			}
			if !valk.IsFKConstraint(err) {
				t.Errorf("expected valk.IsFKConstraint(err) to be true, got false for: %v", err)
			}
			if !valk.IsConstraintError(err) {
				t.Errorf("expected valk.IsConstraintError(err) to be true, got false for: %v", err)
			}
		})
	})

	t.Run("5. TRANSACTION - Error Propagation & Rollback", func(t *testing.T) {
		t.Run("Transaction Rollback Preserves Domain Error", func(t *testing.T) {
			err := db.Transaction(ctx, func(tx *valk.Tx) error {
				_, err := tx.User.Create().SetEmail("tx1@example.com").SetPhoneNum("+4001").Exec(ctx)
				if err != nil {
					return err
				}
				// Force a constraint failure inside transaction
				_, err = tx.User.Create().SetEmail("tx1@example.com").SetPhoneNum("+4002").Exec(ctx)
				return err
			})

			if err == nil {
				t.Fatalf("expected error from transaction, got nil")
			}
			if !valk.IsUniqueConstraint(err) {
				t.Errorf("expected valk.IsUniqueConstraint(err) to be true across transaction rollback, got false for: %v", err)
			}

			// Verify first user was rolled back
			u, err := db.User.FindUnique(user.Email.EQ("tx1@example.com")).Exec(ctx)
			if err != nil {
				t.Fatalf("unexpected error finding user after transaction rollback: %v", err)
			}
			if u != nil {
				t.Fatalf("expected user to be rolled back, but found record in database")
			}
		})

		t.Run("Transaction Timeout / Context Deadline Exceeded", func(t *testing.T) {
			tCtx, cancel := context.WithTimeout(ctx, 1*time.Microsecond)
			time.Sleep(2 * time.Microsecond)
			defer cancel()

			err := db.Transaction(tCtx, func(tx *valk.Tx) error {
				_, err := tx.User.Create().SetEmail("timeout@example.com").SetPhoneNum("+4003").Exec(tCtx)
				return err
			})

			if err == nil {
				t.Fatalf("expected timeout error on expired transaction context, got nil")
			}
			if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
				t.Errorf("expected context timeout or cancellation error, got: %v", err)
			}
		})
	})
}
