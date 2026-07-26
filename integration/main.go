package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"integration/valk"
	"integration/valk/comment"
	"integration/valk/post"
	"integration/valk/profile"
	"integration/valk/user"
	"os"
	"strings"
	"time"

	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type SeedData struct {
	ReferrerId string
	ReferredId string
	PostId     string
	Meta1      json.RawMessage
	Meta2      json.RawMessage
}

func dbReset(db *valk.DB) error {
	tx, err := db.Raw().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DROP SCHEMA public CASCADE`); err != nil {
		return err
	}

	if _, err := tx.Exec(`CREATE SCHEMA public`); err != nil {
		return err
	}

	return tx.Commit()
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db := openConn()
	// db := openPGConn()
	defer dbReset(db)
	defer db.Close()
	var test user.CreateBuilder
	_ = test
	rawDB := db.Raw()
	rawDB.SetMaxOpenConns(10)
	ctx := context.Background()

	runMigrations(db, ctx)

	db.User.Use(user.Extension{
		Create: func(ctx context.Context, args *valk.UserCreateArgs, next valk.UserCreateQuery) (*valk.User, error) {
			fmt.Println("[Create Hook] Conflict Target: ", args.ConflictTarget)
			fmt.Println("[Create Hook] Conflict Action: ", args.ConflictAction)
			return next(ctx, args)
		},
		CreateMany: func(ctx context.Context,
			args *valk.UserCreateManyArgs,
			next valk.UserCreateManyQuery) (int64, error) {

			if args.ConflictTarget == user.Email {
				args.ConflictTarget = user.PhoneNum

				args.ConflictAction = user.ConflictUpdate(func(u *valk.UserUpsert) {
					u.Role.Set(string(valk.UserRole.Admin))
				})
			}

			fmt.Println("[CreateMany Hook] Mutated Target: ", args.ConflictTarget)
			fmt.Println("[CreateMany Hook] Mutated Action: ", args.ConflictAction)
			fmt.Println("[CreateMany Hook] IsUpdateCustom? ", args.ConflictAction.IsUpdateCustom())
			return next(ctx, args)
		},
		CreateManyAndReturn: func(ctx context.Context,
			args *valk.UserCreateManyAndReturnArgs,
			next valk.UserCreateManyAndReturnQuery) ([]*valk.User, error) {

			fmt.Println("[CreateManyAndReturn Hook] Conflict Target: ", args.ConflictTarget)
			fmt.Println("[CreateManyAndReturn Hook] Conflict Action: ", args.ConflictAction)
			return next(ctx, args)
		},
	})

	fmt.Println("\n--- Testing Single Create OnConflict ---")
	_, _ = db.User.Create().SetEmail("test1@example.com").SetPhoneNum("111").OnConflict(user.Email).Ignore().Exec(ctx)

	fmt.Println("\n--- Testing CreateMany OnConflict with user.ConflictUpdate ---")
	_, _ = db.User.CreateMany(
		db.User.Create().SetEmail("test2@example.com").SetPhoneNum("222"),
		db.User.Create().SetEmail("test3@example.com").SetPhoneNum("333"),
	).OnConflict(user.Email).Ignore().Exec(ctx)

	fmt.Println("\n--- Testing CreateManyAndReturn OnConflict ---")
	_, _ = db.User.CreateManyAndReturn(
		db.User.Create().SetEmail("test4@example.com").SetPhoneNum("444"),
		db.User.Create().SetEmail("test5@example.com").SetPhoneNum("555"),
	).OnConflict(user.Email).UpdateNewValues().Exec(ctx)

	// runPaginationExamples(db, ctx)
}

func seed(db *valk.DB, ctx context.Context) *SeedData {

	db.User.Use(user.Extension{
		Create: func(ctx context.Context, args *valk.UserCreateArgs, next valk.UserCreateQuery) (*valk.User, error) {
			if args.Data != nil {
				args.Data.Email = strings.ToLower(args.Data.Email)
				if args.Data.Password != nil {
					hash := sha256.Sum256([]byte(*args.Data.Password))
					hexStr := hex.EncodeToString(hash[:])
					args.Data.Password = &hexStr
				}
			}
			return next(ctx, args)
		},
	})

	var usersToCreate []*user.CreateBuilder

	for i := range 20 {
		usersToCreate = append(usersToCreate, db.User.Create().
			SetEmail(fmt.Sprintf("email-%d", i)).
			SetPhoneNum(fmt.Sprintf("555-%d", i)).
			SetPassword(fmt.Sprintf("password-%d", i)),
		)
	}

	_, err := db.User.FindUnique(
		user.EmailPhoneUnique("x@y.com", "+1111"),
	).Select(user.Select{
		Id:       true,
		Email:    true,
		PhoneNum: true,
		Profile:  &profile.Select{},

		Posts: post.Query().Where(post.And(
			post.Title.Contains("super-cool-post"),
			post.Published.EQ(true),
		)).
			Select(post.Select{
				Id:    true,
				Title: true,
				Comments: comment.Query().Where(comment.Or(
					comment.AuthorId.Contains("xyz"),
					comment.AuthorId.Contains("abc"),
				)),
			}),
	}).
		Exec(ctx)

	users, err := db.User.CreateManyAndReturn(usersToCreate...).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create users: %v", err)
	}
	fmt.Printf("CreateManyAndReturn: %d users returned with auto-generated IDs\n", len(users))

	if _, err := db.User.CreateMany(
		db.User.Create().
			SetEmail("test").
			SetPhoneNum("555-test").
			SetPassword("passwd"),
		db.User.Create().
			SetEmail("again").
			SetPhoneNum("555-again").
			SetPassword("123456"),
	).Exec(ctx); err != nil {
		log.Fatalf("failed to CreateMany: %v", err)
	}
	referrer, err := db.User.Create().SetEmail("referrer@example.com").SetPhoneNum("555-0001").SetPassword("pass123").SetRole(valk.UserRole.Student).Select(user.Select{
		Id:    true,
		Email: true,
	}).
		Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create referrer: %v", err)
	}

	referred, err := db.User.Create().SetEmail("referred@example.com").SetPhoneNum("555-0002").SetPassword("pass456").SetRole(valk.UserRole.Student).SetReferredById(referrer.Id).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create referred: %v", err)
	}

	prof, err := db.Profile.Create().SetBio("BLEH").SetUserId(referred.Id).SetCreatedAt(time.Now()).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create profile: %v", err)
	}
	fmt.Println("PROFILE:")
	printJSON(prof)

	categoryTest, err := db.Category.Create().SetName("TEST").Exec(ctx)

	if err != nil {
		log.Fatalf("failed to create category: %v", err)
	}
	fmt.Println("CATEGORY:")
	printJSON(categoryTest)

	p, err := db.Post.Create().SetTitle("Valkyrie ORM Deep Dive").SetContent("skrrrt").SetAuthorId(referred.Id).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create post: %v", err)
	}

	cat, err := db.Category.Create().SetName("Programming").Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create category: %v", err)
	}

	_, err = db.CategoryToPost.Create().SetPostId(p.Id).SetCategoryId(cat.Id).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create CategoryToPost: %v", err)
	}

	meta1 := json.RawMessage(`{"rating":5,"verified":true}`)
	_, err = db.Comment.Create().
		SetTextify(100).
		SetDummy3("dummy_val_1").
		SetDummy1(42).SetDummy2("dummy_val_2").
		SetPostId(p.Id).SetAuthorId(referrer.Id).
		SetMeta(meta1).Select(comment.Select{
		Post: &post.Select{
			Id:    true,
			Title: true,
			Author: &user.Select{
				Id:    true,
				Email: true,
			},
		},
	}).
		Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create comment 1: %v", err)
	}

	meta2 := json.RawMessage(`{"rating":4,"verified":false}`)
	_, err = db.Comment.Create().SetTextify(200).SetDummy3("dummy_val_3").SetDummy1(84).SetDummy2("dummy_val_4").SetPostId(p.Id).SetAuthorId(referred.Id).SetMeta(meta2).Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create comment 2: %v", err)
	}

	return &SeedData{
		ReferrerId: referrer.Id,
		ReferredId: referred.Id,
		PostId:     p.Id,
		Meta1:      meta1,
		Meta2:      meta2,
	}
}

func openConn() *valk.DB {
	db, err := valk.Open("sqlite3", "file::memory:?_pragma=foreign_keys(1)&_time_format=sqlite")

	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	return db
}

func openPGConn() *valk.DB {
	pgUrl := os.Getenv("DATABASE_DIRECT_URL")
	db, err := valk.Open("postgres", pgUrl)
	_, err = db.Raw().Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")

	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	return db
}
func runMigrations(db *valk.DB, ctx context.Context) {
	if err := db.RunMigrations(ctx); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}

func runManualTransaction(db *valk.DB, ctx context.Context) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("Manual Transaction: failed to begin: %v", err)
		return
	}
	defer tx.Rollback()

	fmt.Println("Manual Transaction: started successfully")
	author, err := tx.User.Create().
		SetEmail("clancySizer@gmail.com").
		SetPhoneNum("+1234567890").
		Exec(ctx)
	if err != nil {
		fmt.Printf("failed to create user: %+v", err)
		return
	}

	postWithAuthor, err := tx.Post.Create().
		SetTitle("A Post").
		SetAuthorId(author.Id).
		Select(post.Select{
			Id:    true,
			Title: true,
			Author: &user.Select{
				Email: true,
			},
		}).Exec(ctx)
	if err != nil {
		fmt.Printf("failed to create Post: %+v", err)
		return
	}

	b, _ := json.MarshalIndent(postWithAuthor, "", "  ")
	fmt.Println(string(b))

	if err := tx.Commit(); err != nil {
		log.Printf("Manual Transaction: commit failed: %v", err)
		return
	}
	fmt.Println("Manual Transaction: committed successfully")
}

func runBlockBasedTransaction(db *valk.DB, ctx context.Context) {
	err := db.Transaction(ctx, func(tx *valk.Tx) error {
		fmt.Println("Block-based Transaction: started successfully")

		author, err := tx.User.Create().
			SetEmail("clancySizer@gmail.com").
			SetPhoneNum("+1234567890").
			Exec(ctx)
		if err != nil {
			return err
		}

		postWithAuthor, err := tx.Post.Create().
			SetTitle("A Post").
			SetAuthorId(author.Id).
			Select(post.Select{
				Id:    true,
				Title: true,
				Author: &user.Select{
					Email: true,
				},
			}).Exec(ctx)
		if err != nil {
			return err
		}

		b, _ := json.MarshalIndent(postWithAuthor, "", "  ")
		fmt.Println(string(b))
		return nil
	})
	if err != nil {
		fmt.Printf("Block-based Transaction failed: %v", err)
	}
	fmt.Println("Block-based Transaction: committed successfully")
}

func printJSON(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

func runPaginationExamples(db *valk.DB, ctx context.Context) {
	fmt.Println("=================== PAGINATION & ORDERBY EXAMPLES ===================")

	// Seed 4 test users for pagination
	u1, err := db.User.Create().SetEmail("pag_alpha@example.com").SetPhoneNum("+101").SetLoginCount(50).Exec(ctx)
	if err != nil {
		log.Printf("Pagination Seed u1 failed: %v", err)
		return
	}
	defer db.User.Delete(user.Id.EQ(u1.Id)).Exec(ctx)

	u2, err := db.User.Create().SetEmail("pag_bravo@example.com").SetPhoneNum("+102").SetLoginCount(20).Exec(ctx)
	if err != nil {
		log.Printf("Pagination Seed u2 failed: %v", err)
		return
	}
	defer db.User.Delete(user.Id.EQ(u2.Id)).Exec(ctx)

	u3, err := db.User.Create().SetEmail("pag_charlie@example.com").SetPhoneNum("+103").SetLoginCount(50).Exec(ctx)
	if err != nil {
		log.Printf("Pagination Seed u3 failed: %v", err)
		return
	}
	defer db.User.Delete(user.Id.EQ(u3.Id)).Exec(ctx)

	u4, err := db.User.Create().SetEmail("pag_delta@example.com").SetPhoneNum("+104").SetLoginCount(10).Exec(ctx)
	if err != nil {
		log.Printf("Pagination Seed u4 failed: %v", err)
		return
	}
	defer db.User.Delete(user.Id.EQ(u4.Id)).Exec(ctx)

	// -------------------------------------------------------------------------
	// SCENARIO 0: No Sorting at all
	// -------------------------------------------------------------------------
	fmt.Println("\n--- SCENARIO 0: No Sorting ---")
	_, err = db.User.FindMany(
		valk.Or(
			user.Email.EQ("pag_alpha@example.com"),
			user.Email.EQ("pag_bravo@example.com"),
			user.Email.EQ("pag_charlie@example.com"),
			user.Email.EQ("pag_delta@example.com"),
		),
	).Exec(ctx)
	if err != nil {
		log.Printf("OrderBy Asc failed: %v", err)
		return
	}
	// printJSON()

	// -------------------------------------------------------------------------
	// SCENARIO 1: Simple Sorting with Single OrderBy (Ascending / Descending)
	// -------------------------------------------------------------------------
	fmt.Println("\n--- SCENARIO 1: OrderBy Email ASC ---")
	_, err = db.User.FindMany(
		valk.Or(
			user.Email.EQ("pag_alpha@example.com"),
			user.Email.EQ("pag_bravo@example.com"),
			user.Email.EQ("pag_charlie@example.com"),
			user.Email.EQ("pag_delta@example.com"),
		),
	).OrderBy(user.Email.Asc()).Exec(ctx)
	if err != nil {
		log.Printf("OrderBy Asc failed: %v", err)
		return
	}
	// printJSON(ascUsers)

	// -------------------------------------------------------------------------
	// SCENARIO 2: Multi-Field Sorting (LoginCount DESC, then Email ASC)
	// Notice: u1 and u3 both have LoginCount = 50, so Email ASC breaks the tie!
	// -------------------------------------------------------------------------
	fmt.Println("\n--- SCENARIO 2: Multi-Field OrderBy (LoginCount DESC, Email ASC) ---")
	_, err = db.User.FindMany(
		valk.Or(
			user.Email.EQ("pag_alpha@example.com"),
			user.Email.EQ("pag_bravo@example.com"),
			user.Email.EQ("pag_charlie@example.com"),
			user.Email.EQ("pag_delta@example.com"),
		),
	).OrderBy(user.LoginCount.Desc(), user.Email.Asc()).Exec(ctx)
	if err != nil {
		log.Printf("Multi-field OrderBy failed: %v", err)
		return
	}
	// printJSON(multiOrderUsers)

	// -------------------------------------------------------------------------
	// SCENARIO 3: Cursor-Based Pagination (Page 1 -> Page 2)
	// -------------------------------------------------------------------------
	fmt.Println("\n--- SCENARIO 3A: Cursor Pagination - Page 1 (Take 2, OrderBy Email ASC) ---")
	page1, err := db.User.FindMany(
		valk.Or(
			user.Email.EQ("pag_alpha@example.com"),
			user.Email.EQ("pag_bravo@example.com"),
			user.Email.EQ("pag_charlie@example.com"),
			user.Email.EQ("pag_delta@example.com"),
		),
	).OrderBy(user.Email.Asc()).Take(2).Exec(ctx)
	if err != nil {
		log.Printf("Cursor Page 1 failed: %v", err)
		return
	}
	// printJSON(page1)

	if len(page1) > 0 {
		lastSeen := page1[len(page1)-1]
		fmt.Printf("\n--- SCENARIO 3B: Cursor Pagination - Page 2 (Cursor after %s) ---\n", lastSeen.Email)
		_, err = db.User.FindMany(
			valk.Or(
				user.Email.EQ("pag_alpha@example.com"),
				user.Email.EQ("pag_bravo@example.com"),
				user.Email.EQ("pag_charlie@example.com"),
				user.Email.EQ("pag_delta@example.com"),
			),
		).OrderBy(user.Email.Asc()).Cursor(user.Email.EQ(lastSeen.Email)).Take(2).Exec(ctx)
		if err != nil {
			log.Printf("Cursor Page 2 failed: %v", err)
			return
		}
		// printJSON(page2)
	}

	// -------------------------------------------------------------------------
	// SCENARIO 4: Cursor Pagination with Multi-Column Sorting & Filtering
	// -------------------------------------------------------------------------
	fmt.Println("\n--- SCENARIO 4: Filter + Multi-Sort + Cursor (LoginCount >= 20, OrderBy LoginCount DESC, Email ASC) ---")
	_, err = db.User.FindMany(
		user.LoginCount.GTE(20),
		valk.Or(
			user.Email.EQ("pag_alpha@example.com"),
			user.Email.EQ("pag_bravo@example.com"),
			user.Email.EQ("pag_charlie@example.com"),
		),
	).
		OrderBy(user.LoginCount.Desc(), user.Email.Asc()).
		Cursor(user.Email.EQ("pag_alpha@example.com")).
		Take(2).
		Exec(ctx)
	if err != nil {
		log.Printf("Scenario 4 failed: %v", err)
		return
	}
	// printJSON(filteredCursor)

	// -------------------------------------------------------------------------
	// SCENARIO 5: Relation Sub-Queries with OrderBy & Take (1-to-Many relation)
	// -------------------------------------------------------------------------
	fmt.Println("\n--- SCENARIO 5: Relation Selection with Sub-Query OrderBy & Take ---")
	post1, _ := db.Post.Create().SetTitle("Zebra Post").SetAuthorId(u1.Id).Exec(ctx)
	defer db.Post.Delete(post.Id.EQ(post1.Id)).Exec(ctx)

	post2, _ := db.Post.Create().SetTitle("Apple Post").SetAuthorId(u1.Id).Exec(ctx)
	defer db.Post.Delete(post.Id.EQ(post2.Id)).Exec(ctx)

	_, err = db.User.FindUnique(user.Id.EQ(u1.Id)).Select(valk.UserSelect{
		Email: true,
		Posts: post.Query().
			OrderBy(post.Title.Asc()).
			Take(5).
			Select(post.Select{
				Id:    true,
				Title: true,
			}),
	}).Exec(ctx)
	if err != nil {
		log.Printf("Scenario 5 failed: %v", err)
		return
	}
	// printJSON(userWithSortedPosts)

	// -------------------------------------------------------------------------
	// SCENARIO 6: Non-Unique OrderBy + Cursor (Triggers Auto-Appended PK Tiebreaker)
	// LoginCount is non-unique (u1 and u3 both have LoginCount = 50).
	// Valkyrie automatically appends "id" ASC to both subquery tuple and ORDER BY!
	// -------------------------------------------------------------------------
	fmt.Println("\n--- SCENARIO 6: Non-Unique OrderBy + Cursor (Auto-Appends PK 'id' Tiebreaker) ---")
	_, err = db.User.FindMany(
		valk.Or(
			user.Email.EQ("pag_alpha@example.com"),
			user.Email.EQ("pag_bravo@example.com"),
			user.Email.EQ("pag_charlie@example.com"),
			user.Email.EQ("pag_delta@example.com"),
		),
	).
		OrderBy(user.LoginCount.Desc()).
		Cursor(user.Email.EQ("pag_alpha@example.com")).
		Take(2).
		Exec(ctx)
	if err != nil {
		log.Printf("Scenario 6 failed: %v", err)
		return
	}
}
