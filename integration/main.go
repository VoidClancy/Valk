package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"integration/phi"
	"integration/phi/categoryToPost"
	"integration/phi/comment"
	"integration/phi/post"
	"integration/phi/profile"
	"integration/phi/user"
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

func dbReset(db *phi.DB) error {
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
	// db := openConn()
	db := openPGConn()
	dbReset(db)
	defer dbReset(db)
	defer db.Close()
	rawDB := db.Raw()
	rawDB.SetMaxOpenConns(10)
	ctx := context.Background()
	runMigrations(db, ctx)

	// seed(db, ctx)
	// runManualTransaction(db, ctx)
	// runBlockBasedTransaction(db, ctx)
	// runPaginationExamples(db, ctx)
	// runExtensionExamples(db, ctx)
	// runCTP(db, ctx)
	//

	usr, err := db.User.Create().SetEmail("xx@yy.com").SetPhoneNum("11").Exec(ctx)
	if err != nil {
		panic(err)
	}
	db.Post.Create().SetTitle("Post 1").SetAuthorId(usr.Id).SetTags([]string{"aa", "bb"}).Exec(ctx)
	posts, err := db.Post.FindMany(post.Tags.HasEvery([]string{"aa", "bb"})).Exec(ctx)

	if err != nil {
		panic(err)
	}
	printJSON(posts)
}

// =============================================================================
// EXTENSION HOOKS , all hook types with inspection & mutation
// =============================================================================
func runExtensionExamples(db *phi.DB, ctx context.Context) {
	_ = ctx
	db.User.Use(user.Extension{
		Create: func(ctx context.Context, args *phi.UserCreateArgs, next phi.UserCreateQuery) (*phi.User, error) {
			return next(ctx, args)
		},
		CreateMany: func(ctx context.Context, args *phi.UserCreateManyArgs, next phi.UserCreateManyQuery) (int64, error) {
			return next(ctx, args)
		},
		CreateManyAndReturn: func(ctx context.Context, args *phi.UserCreateManyAndReturnArgs, next phi.UserCreateManyAndReturnQuery) ([]*phi.User, error) {
			return next(ctx, args)
		},
		FindUnique: func(ctx context.Context, args *phi.UserFindUniqueArgs, next phi.UserFindUniqueQuery) (*phi.User, error) {
			for _, w := range args.Where {
				col, val := w.Column(), w.Value()
				_ = col
				_ = val
			}
			return next(ctx, args)
		},
		FindFirst: func(ctx context.Context, args *phi.UserFindFirstArgs, next phi.UserFindFirstQuery) (*phi.User, error) {
			args.SetOrderBy(user.Email.Asc()).
				SetCursor(user.Email.EQ("x@y.com")).
				SetSkip(10).
				SetTake(20)
			return next(ctx, args)
		},
		FindMany: func(ctx context.Context, args *phi.UserFindManyArgs, next phi.UserFindManyQuery) ([]*phi.User, error) {
			args.Where = append(args.Where, user.LoginCount.GTE(10))
			return next(ctx, args)
		},
		Count: func(ctx context.Context, args *phi.UserCountArgs, next phi.UserCountQuery) (int64, error) {
			return next(ctx, args)
		},
		Delete: func(ctx context.Context, args *phi.UserDeleteArgs, next phi.UserDeleteQuery) (*phi.User, error) {
			return next(ctx, args)
		},
		DeleteMany: func(ctx context.Context, args *phi.UserDeleteManyArgs, next phi.UserDeleteManyQuery) (int64, error) {
			return next(ctx, args)
		},
		Update: func(ctx context.Context, args *phi.UserUpdateArgs, next phi.UserUpdateQuery) (*phi.User, error) {
			if args.Data.Email != nil {
				lower := strings.ToLower(*args.Data.Email)
				args.Data.Email = &lower
			}

			for _, w := range args.Where {
				if w.Column() == user.Email.Column {

				}
				if w.Column() == user.EmailPhone.Column {

				}
			}
			return next(ctx, args)
		},
		UpdateMany: func(ctx context.Context, args *phi.UserUpdateManyArgs, next phi.UserUpdateManyQuery) (int64, error) {

			return next(ctx, args)
		},

		UpdateManyAndReturn: func(ctx context.Context,
			args *phi.UserUpdateManyAndReturnArgs,
			next phi.UserUpdateManyAndReturnQuery) ([]*phi.User, error) {

			return next(ctx, args)

		},
	})
}

// =============================================================================
// EXTENSION API DEEP DIVE , contrasts hook-level vs top-level builder patterns
// =============================================================================
func inconsistency() {
	db := openConn()

	db.CategoryToPost.Use(categoryToPost.Extension{
		FindMany: func(ctx context.Context,
			args *phi.CategoryToPostFindManyArgs,
			next phi.CategoryToPostFindManyQuery) ([]*phi.CategoryToPost, error) {
			for i, w := range args.Where {
				if w.Column() == categoryToPost.PostId.Column {
					args.Where[i] = categoryToPost.CategoryId.EQ(22)
				}

			}
			return next(ctx, args)
		},
	})

	db.CategoryToPost.FindFirst(categoryToPost.PostId_CategoryId.EQ("xx", 12),
		categoryToPost.CategoryId.EQ(32))
	db.User.Use(user.Extension{

		FindUnique: func(ctx context.Context, args *phi.UserFindUniqueArgs, next phi.UserFindUniqueQuery) (*phi.User, error) {
			for _, w := range args.Where {
				col, _ := w.Column(), w.Value()
				if col == user.Email.Column {
					w = user.Email.EQ("usr@example.com")

				}
				if w.Column() == user.Email.Column {
					w = user.EmailPhone.EQ("", "")
				}
			}
			args.SetWhere(user.Email.EQ("unique@example.com"), user.Role.EQ(phi.UserRole_ADMIN))
			args.Select.Email = true
			args.Select.Posts = post.Query().Where(post.Title.Contains("News")).OrderBy(post.Title.Desc())
			return next(ctx, args)
		},

		FindFirst: func(ctx context.Context, args *phi.UserFindFirstArgs, next phi.UserFindFirstQuery) (*phi.User, error) {
			args.SetWhere(user.Email.Like("%@example.com")).
				SetOrderBy(user.LoginCount.Desc(), user.Email.Asc()).
				SetSkip(10).
				SetTake(20).
				SetCursor(user.Email.EQ("cursor@example.com"))
			for i, w := range args.Where {
				if w.Column() == user.Email.Column {
					args.Where[i] = user.Email.Like("%@example.com")
				}
			}
			args.Select.Posts = post.Query().OrderBy(post.Title.Asc())
			return next(ctx, args)
		},

		FindMany: func(ctx context.Context, args *phi.UserFindManyArgs, next phi.UserFindManyQuery) ([]*phi.User, error) {
			args.Where = append(args.Where, user.LoginCount.GTE(10))
			args.SetTake(50)
			args.OrderBy = append(args.OrderBy, user.Email.Asc())
			args.OrderBy = []phi.OrderBy[phi.User]{
				user.Email.Asc(),
				user.LoginCount.Asc(),
			}
			args.SetOrderBy(user.Email.Asc(), user.LoginCount.Asc())
			args.Cursor = user.Email.EQ("X")
			return next(ctx, args)
		},

		Count: func(ctx context.Context, args *phi.UserCountArgs, next phi.UserCountQuery) (int64, error) {
			args.SetWhere(user.LoginCount.GTE(10))
			args.SetTake(1000)
			for i, w := range args.Where {
				if w.Column() == user.LoginCount.Column {
					args.Where[i] = user.LoginCount.GTE(10)
				}
			}
			return next(ctx, args)
		},

		Create: func(ctx context.Context, args *phi.UserCreateArgs, next phi.UserCreateQuery) (*phi.User, error) {
			args.Data.Email = strings.ToLower(args.Data.Email)

			if args.ConflictAction != nil && args.ConflictAction.IsDoNothing() {
				args.ConflictTarget = user.EmailPhone
			}
			return next(ctx, args)
		},

		CreateMany: func(ctx context.Context, args *phi.UserCreateManyArgs, next phi.UserCreateManyQuery) (int64, error) {
			args.AppendData(db.User.Create().SetEmail("xx"), db.User.Create().SetEmail("yy"))
			for _, r := range args.Data {
				fmt.Println(r.Email)
			}
			for i, r := range args.Data {
				if r.Email == "create_many@example" {
					r.Email = fmt.Sprintf("create_many_%d@example", i)
				}
			}

			if args.ConflictAction != nil && args.ConflictAction.IsUpdateNewValues() {
				args.ConflictAction = user.ConflictUpdate(func(u *phi.UserUpsert) {
					u.Role.Set(phi.UserRole_STUDENT)
					u.LoginCount.Increment(1)
				})
			}
			return next(ctx, args)
		},

		CreateManyAndReturn: func(ctx context.Context, args *phi.UserCreateManyAndReturnArgs, next phi.UserCreateManyAndReturnQuery) ([]*phi.User, error) {

			adminRole := phi.UserRole_ADMIN
			args.Data = append(args.Data, &phi.UserCreate{
				Email: "create_many@example",
				Role:  &adminRole,
			})

			for i, r := range args.Data {
				if r.Email == "create_many@example" {
					r.Email = fmt.Sprintf("create_many_%d@example", i)
				}
			}

			args.Select.Posts = post.Query()
			return next(ctx, args)
		},

		Delete: func(ctx context.Context, args *phi.UserDeleteArgs, next phi.UserDeleteQuery) (*phi.User, error) {
			args.SetWhere(user.Email.EQ("delete@example.com"))
			args.Select.Posts = post.Query()
			return next(ctx, args)
		},

		DeleteMany: func(ctx context.Context, args *phi.UserDeleteManyArgs, next phi.UserDeleteManyQuery) (int64, error) {
			args.SetWhere(user.LoginCount.EQ(0))
			return next(ctx, args)
		},
	})

	// Top-level builder equivalents for comparison
	db.User.FindUnique(user.Email.EQ("test@example.com"), user.Role.EQ(phi.UserRole_ADMIN)).
		Select(user.Select{
			Email: true,
			Posts: post.Query().Where(post.Title.Contains("News")).OrderBy(post.Title.Desc()),
		})

	db.User.FindFirst(user.Email.Like("%@example.com")).
		OrderBy(user.LoginCount.Desc(), user.Email.Asc()).
		Skip(10).
		Take(20).
		Cursor(user.Email.EQ("cursor@example.com")).
		Select(user.Select{
			Posts: post.Query().OrderBy(post.Title.Asc()),
		})

	db.User.FindMany(user.LoginCount.GTE(10)).Take(50)
	db.User.Count(user.LoginCount.GTE(10)).Take(1000)

	db.User.Create().
		SetEmail("new@example.com").
		SetPhoneNum("1234567890").
		OnConflict(user.Email).Ignore()

	db.User.CreateMany(
		db.User.Create().SetEmail("batch1@example.com"),
		db.User.Create().SetEmail("batch2@example.com"),
	).OnConflict(user.Email).UpdateNewValues()

	db.User.CreateManyAndReturn(
		db.User.Create().SetEmail("ret1@example.com"),
	).Select(user.Select{
		Posts: post.Query(),
	})

	db.User.Delete(user.Email.EQ("del@example.com")).
		Select(user.Select{
			Posts: post.Query(),
		})

	db.User.DeleteMany(user.LoginCount.EQ(0))

	db.User.Create().
		SetEmail("x").
		SetPhoneNum("x").
		OnConflict(user.Email).
		Ignore().
		Select(
			user.Select{
				Email: true,
				Profile: &profile.Select{
					Id:  true,
					Bio: true,
				},

				Posts: post.Query().Where(post.Id.EQ("xx")),
			},
		)

}

// =============================================================================
// SEED DATA
// =============================================================================
func seed(db *phi.DB, ctx context.Context) *SeedData {
	db.User.Use(user.Extension{
		Create: func(ctx context.Context, args *phi.UserCreateArgs, next phi.UserCreateQuery) (*phi.User, error) {
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
		user.EmailPhone.EQ("x@y.com", "+1111"),
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
	referrer, err := db.User.Create().SetEmail("referrer@example.com").SetPhoneNum("555-0001").SetPassword("pass123").SetRole(phi.UserRole_STUDENT).Select(user.Select{
		Id:    true,
		Email: true,
	}).
		Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create referrer: %v", err)
	}

	referred, err := db.User.Create().SetEmail("referred@example.com").SetPhoneNum("555-0002").SetPassword("pass456").SetRole(phi.UserRole_STUDENT).SetReferredById(referrer.Id).Exec(ctx)
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

	p, err := db.Post.Create().SetTitle("Phi ORM Deep Dive").SetContent("skrrrt").SetAuthorId(referred.Id).Exec(ctx)
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

// =============================================================================
// CONNECTIONS
// =============================================================================
func openConn() *phi.DB {
	db, err := phi.Open("sqlite3", "file:memdb1?mode=memory&cache=shared&_pragma=foreign_keys(1)&_time_format=sqlite")

	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	return db
}

func openPGConn() *phi.DB {
	pgUrl := os.Getenv("DATABASE_DIRECT_URL")
	db, err := phi.Open("postgres", pgUrl)
	_, err = db.Raw().Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")

	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	return db
}

// =============================================================================
// MIGRATIONS
// =============================================================================
func runMigrations(db *phi.DB, ctx context.Context) {
	if err := db.RunMigrations(ctx); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}

// =============================================================================
// COMPOSITE KEY
// =============================================================================
func runCTP(db *phi.DB, ctx context.Context) {
	db.CategoryToPost.Use(categoryToPost.Extension{
		FindUnique: func(ctx context.Context,
			args *phi.CategoryToPostFindUniqueArgs,
			next phi.CategoryToPostFindUniqueQuery) (*phi.CategoryToPost, error) {

			for _, w := range args.Where {
				column, value := w.Column(), w.Value()

				if w.Column() == categoryToPost.PostId.Column {

				}

				if w.Column() == categoryToPost.PostId_CategoryId.Column {
					//since it is a one logical col now, for consistency with scalars col access pattern
				}
				fmt.Println(column, value)
			}
			args.Where = append(args.Where, categoryToPost.PostId_CategoryId.EQ("xc", 22))
			return next(ctx, args)
		},
	})
	db.CategoryToPost.FindUnique(categoryToPost.PostId_CategoryId.EQ("x", 23), categoryToPost.CategoryId.EQ(12)).Exec(ctx)
	db.CategoryToPost.Create().SetCategoryId(22).SetPostId("x").Exec(ctx) //creates with the composite PK constraint
}

// =============================================================================
// MANUAL TRANSACTION
// =============================================================================
func runManualTransaction(db *phi.DB, ctx context.Context) {
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

// =============================================================================
// BLOCK-BASED TRANSACTION
// =============================================================================
func runBlockBasedTransaction(db *phi.DB, ctx context.Context) {
	err := db.Transaction(ctx, func(tx *phi.Tx) error {
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

// =============================================================================
// PAGINATION & ORDERING PLAYGROUND
// =============================================================================
func runPaginationExamples(db *phi.DB, ctx context.Context) {
	fmt.Println("=================== PAGINATION & ORDERBY EXAMPLES ===================")

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

	// Scenario 0: No Sorting
	fmt.Println("\n--- SCENARIO 0: No Sorting ---")
	_, err = db.User.FindMany(
		phi.Or(
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

	// Scenario 1: Single OrderBy (Ascending)
	fmt.Println("\n--- SCENARIO 1: OrderBy Email ASC ---")
	_, err = db.User.FindMany(
		phi.Or(
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

	// Scenario 2: Multi-Field Sorting
	fmt.Println("\n--- SCENARIO 2: Multi-Field OrderBy (LoginCount DESC, Email ASC) ---")
	_, err = db.User.FindMany(
		phi.Or(
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

	// Scenario 3: Cursor Pagination
	fmt.Println("\n--- SCENARIO 3A: Cursor Pagination - Page 1 (Take 2, OrderBy Email ASC) ---")
	page1, err := db.User.FindMany(
		phi.Or(
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

	if len(page1) > 0 {
		lastSeen := page1[len(page1)-1]
		fmt.Printf("\n--- SCENARIO 3B: Cursor Pagination - Page 2 (Cursor after %s) ---\n", lastSeen.Email)
		_, err = db.User.FindMany(
			phi.Or(
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
	}

	// Scenario 4: Filter + Multi-Sort + Cursor
	fmt.Println("\n--- SCENARIO 4: Filter + Multi-Sort + Cursor (LoginCount >= 20, OrderBy LoginCount DESC, Email ASC) ---")
	_, err = db.User.FindMany(
		user.LoginCount.GTE(20),
		phi.Or(
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

	// Scenario 5: Relation Sub-Queries with OrderBy & Take
	fmt.Println("\n--- SCENARIO 5: Relation Selection with Sub-Query OrderBy & Take ---")
	post1, _ := db.Post.Create().SetTitle("Zebra Post").SetAuthorId(u1.Id).Exec(ctx)
	defer db.Post.Delete(post.Id.EQ(post1.Id)).Exec(ctx)

	post2, _ := db.Post.Create().SetTitle("Apple Post").SetAuthorId(u1.Id).Exec(ctx)
	defer db.Post.Delete(post.Id.EQ(post2.Id)).Exec(ctx)

	_, err = db.User.FindUnique(user.Id.EQ(u1.Id)).Select(phi.UserSelect{
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

	// Scenario 6: Non-Unique OrderBy + Cursor
	fmt.Println("\n--- SCENARIO 6: Non-Unique OrderBy + Cursor (Auto-Appends PK 'id' Tiebreaker) ---")
	_, err = db.User.FindMany(
		phi.Or(
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

// =============================================================================
// HELPERS
// =============================================================================
func printJSON(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
