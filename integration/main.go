package main

import (
	"context"
	"encoding/json"
	"fmt"
	"integration/valk"
	"integration/valk/post"
	"integration/valk/user"
	"os"

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
	// _, err = db.User.Create().SetEmail("c@y.com").SetPhoneNum("+1111").SetId("1234").Exec(ctx)
	// foundUser, err := db.User.FindUnique(
	// 	user.EmailPhoneUnique("c@y.com", "+1111"),
	// 	user.And(
	// 		user.PhoneNum.Contains("1111"),
	// 		user.Id.Contains("234"),
	// 	),
	// ).Select(user.Select{
	// 	Id:    true,
	// 	Email: true,

	// 	Profile: &profile.Select{
	// 		Id:  true,
	// 		Bio: true,
	// 	},

	// 	Posts: post.Query().
	// 		Where(post.AuthorId.Contains("234")).
	// 		Select(post.Select{
	// 			Id:    true,
	// 			Title: true,
	// 		}),
	// }).
	// 	Exec(ctx)

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// printJSON(foundUser)

	// _, err = db.Post.FindMany(post.Id.Contains("xx")).Select(post.Select{}).Exec(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _, err = db.Post.FindUnique(post.Id.EQ("xxx")).Select(post.Select{}).Exec(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// user, err := db.User.Create().SetId("122").SetEmail("xasx").SetPhoneNum("+122111").Exec(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// printJSON(user)
	// bulkUsers, err := db.User.CreateManyAndReturn(
	// 	db.User.Create().SetId("11").SetEmail("xx").SetPhoneNum("+1111"),
	// 	db.User.Create().SetId("22").SetEmail("xy").SetPhoneNum("+11112").SetRole(valk.UserRole.Admin),
	// 	db.User.Create().SetId("22222").SetEmail("xcy").SetPhoneNum("+ss11112").SetRole(valk.UserRole.Admin),
	// ).Exec(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// printJSON(bulkUsers)

	// var builders []*user.CreateBuilder
	// for i := range 20 {
	// 	builder := db.User.Create().
	// 		SetEmail(fmt.Sprintf("user%d@gmail.com", i)).
	// 		SetPassword(fmt.Sprintf("pass%d", i)).
	// 		SetPhoneNum(fmt.Sprintf("+1111%d", i))
	// 	if i%2 == 0 {
	// 		builder.SetRole(valk.UserRole.Admin)

	// 	}

	// 	builders = append(builders, builder)
	// }
	// count, err := db.User.CreateMany(builders...).
	// 	OnConflict(user.Email).Update(func(u *valk.UserUpsert) {
	// 	u.Role.Set(string(valk.UserRole.Admin))
	// }).SkipDuplicates().
	// 	Exec(ctx)
	// if err != nil {
	// 	log.Fatalf("failed to seed users: %v", err)
	// }

	// printJSON(count)
	db.User.Use(user.Extension{
		Create: func(ctx context.Context, input *valk.UserCreate, next valk.UserCreateQuery) (*valk.User, error) {

			fmt.Println("CREATING USER WITH EMAIl: ", input.Email)
			usr, err := next(ctx, input)
			if err != nil {
				return nil, fmt.Errorf("FAILED TO CREATE USER: %v ", err)
			}

			return usr, nil
		},
	})
	// user1, err := db.User.Create().SetEmail("a").SetPhoneNum("11").Exec(ctx)
	// if err != nil {
	// 	log.Fatalf("failed to seed users: %v", err)
	// }
	// user2, err := db.User.Create().SetEmail("a").SetPhoneNum("11").OnConflict(user.EmailPhone).Ignore().Exec(ctx)
	// if err != nil {
	// 	log.Fatalf("failed to seed users: %v", err)
	// }
	// fmt.Println("USER1:")
	// printJSON(user1)
	// fmt.Println("USER2:")
	// printJSON(user2)
	// db.User.Create().SetEmail("xxxc").SetPhoneNum("6969").SetId("Bleh").Exec(ctx)
	// deletedCnt, err := db.User.DeleteMany(user.Id.EQ("Bleh"), user.Password.Contains("xx")).Exec(ctx)
	// fmt.Printf("DELETED %d USERS \n", deletedCnt)

	// usersCnt, err := db.User.Count(
	// 	user.Id.Contains("x"),
	// 	user.Or(
	// 		user.Password.Contains("y"),
	// 		user.Email.NEQ("c"),
	// 	)).
	// 	Exec(ctx)
	// fmt.Printf("COUNT %d USERS \n", usersCnt)

	// // --- DELETE SCENARIO 1: Simple Delete (No Select -> Returns All Scalar Fields) ---
	// delUser1, err := db.User.Create().
	// 	SetEmail("del1@example.com").
	// 	SetPhoneNum("+10001").
	// 	Exec(ctx)
	// if err != nil {
	// 	log.Fatalf("failed to create user 1: %v", err)
	// }
	// deleted1, err := db.User.Delete(user.Id.EQ(delUser1.Id)).Exec(ctx)
	// if err != nil {
	// 	log.Fatalf("failed to delete user 1: %v", err)
	// }
	// fmt.Println("DELETE SCENARIO 1 (Simple Delete - All Scalars):")
	// printJSON(deleted1)

	// // --- DELETE SCENARIO 2: Delete with Scalar Field Selection ---
	// delUser2, err := db.User.Create().
	// 	SetEmail("del2@example.com").
	// 	SetPhoneNum("+10002").
	// 	Exec(ctx)
	// if err != nil {
	// 	log.Fatalf("failed to create user 2: %v", err)
	// }

	// deleted2, err := db.User.Delete(user.Id.EQ(delUser2.Id)).
	// 	Select(valk.UserSelect{
	// 		Email: true,
	// 		Role:  true,
	// 	}).
	// 	Exec(ctx)
	// if err != nil {
	// 	log.Fatalf("failed to delete user 2: %v", err)
	// }
	// fmt.Println("DELETE SCENARIO 2 (Scalar Field Selection):")
	// printJSON(deleted2)

	// // --- DELETE SCENARIO 3: Delete with Nested Relations Selection ---
	// delUser3, err := db.User.Create().
	// 	SetEmail("del3@example.com").
	// 	SetPhoneNum("+10003").
	// 	Exec(ctx)
	// if err != nil {
	// 	log.Fatalf("failed to create user 3: %v", err)
	// }

	// _, err = db.Profile.Create().
	// 	SetBio("Bio of del3").
	// 	SetUserId(delUser3.Id).
	// 	Exec(ctx)
	// if err != nil {
	// 	log.Fatalf("failed to create profile for user 3: %v", err)
	// }

	// _, err = db.Post.Create().
	// 	SetTitle("First Post of del3").
	// 	SetAuthorId(delUser3.Id).
	// 	Exec(ctx)
	// if err != nil {
	// 	log.Fatalf("failed to create post for user 3: %v", err)
	// }

	// deleted3, err := db.User.Delete(user.Id.EQ(delUser3.Id)).
	// 	Select(valk.UserSelect{
	// 		Email: true,
	// 		Profile: &profile.Select{
	// 			Bio: true,
	// 		},
	// 		Posts: post.Query().Select(post.Select{
	// 			Title: true,
	// 		}),
	// 	}).
	// 	Exec(ctx)
	// if err != nil {
	// 	log.Fatalf("failed to delete user 3: %v", err)
	// }
	runPaginationExamples(db, ctx)
}

// func seed(db *valk.DB, ctx context.Context) *SeedData {

// 	db.User.Use(user.Extension{
// 		Create: func(ctx context.Context, input *valk.UserCreate, next valk.UserCreateQuery) (*valk.User, error) {
// 			return next(ctx, input)
// 		},
// 	})

// 	var usersToCreate []*user.CreateBuilder

// 	for i := range 20 {
// 		usersToCreate = append(usersToCreate, db.User.Create().
// 			SetEmail(fmt.Sprintf("email-%d", i)).
// 			SetPhoneNum(fmt.Sprintf("555-%d", i)).
// 			SetPassword(fmt.Sprintf("password-%d", i)),
// 		)
// 	}

// 	_, err := db.User.FindUnique(
// 		user.EmailPhoneUnique("x@y.com", "+1111"),
// 	).Select(user.Select{
// 		Id:       true,
// 		Email:    true,
// 		PhoneNum: true,
// 		Profile:  &profile.Select{},

// 		Posts: post.Query().Where(post.And(
// 			post.Title.Contains("super-cool-post"),
// 			post.Published.EQ(true),
// 		)).
// 			Select(post.Select{
// 				Id:    true,
// 				Title: true,
// 				Comments: comment.Query().Where(comment.Or(
// 					comment.AuthorId.Contains("xyz"),
// 					comment.AuthorId.Contains("abc"),
// 				)),
// 			}),
// 	}).
// 		Exec(ctx)

// 	users, err := db.User.CreateManyAndReturn(usersToCreate...).Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to create users: %v", err)
// 	}
// 	fmt.Printf("CreateManyAndReturn: %d users returned with auto-generated IDs\n", len(users))

// 	if _, err := db.User.CreateMany(
// 		db.User.Create().
// 			SetEmail("test").
// 			SetPhoneNum("555-test").
// 			SetPassword("passwd"),
// 		db.User.Create().
// 			SetEmail("again").
// 			SetPhoneNum("555-again").
// 			SetPassword("123456"),
// 	).Exec(ctx); err != nil {
// 		log.Fatalf("failed to CreateMany: %v", err)
// 	}
// 	referrer, err := db.User.Create()// 		.SetEmail("referrer@example.com")// 		.SetPhoneNum("555-0001")// 		.SetPassword("pass123")// 		.SetRole(valk.UserRole.Student)//.Select(user.Select{
// 		Id:    true,
// 		Email: true,
// 	}).
// 		Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to create referrer: %v", err)
// 	}

// 	referred, err := db.User.Create()// 		.SetEmail("referred@example.com")// 		.SetPhoneNum("555-0002")// 		.SetPassword("pass456")// 		.SetRole(valk.UserRole.Student)// 		.SetReferredById(referrer.Id)//.Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to create referred: %v", err)
// 	}

// 	prof, err := db.Profile.Create()// 		.SetBio("BLEH")// 		.SetUserId(referred.Id)// 		.SetCreatedAt(time.Now())//.Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to create profile: %v", err)
// 	}
// 	fmt.Println("PROFILE:")
// 	printJSON(prof)

// 	categoryTest, err := db.Category.Create()// 		.SetName("TEST")//.Exec(ctx)

// 	if err != nil {
// 		log.Fatalf("failed to create category: %v", err)
// 	}
// 	fmt.Println("CATEGORY:")
// 	printJSON(categoryTest)

// 	p, err := db.Post.Create()// 		.SetTitle("Valkyrie ORM Deep Dive")// 		.SetContent("skrrrt")// 		.SetAuthorId(referred.Id)//.Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to create post: %v", err)
// 	}

// 	cat, err := db.Category.Create()// 		.SetName("Programming")//.Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to create category: %v", err)
// 	}

// 	_, err = db.CategoryToPost.Create()// 		.SetPostId(p.Id)// 		.SetCategoryId(cat.Id)//.Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to create CategoryToPost: %v", err)
// 	}

// 	meta1 := json.RawMessage(`{"rating":5,"verified":true}`)
// 	_, err = db.Comment.Create()// 		.SetTextify(100)// 		.SetDummy3("dummy_val_1")// 		.SetDummy1(42)// 		.SetDummy2("dummy_val_2")// 		.SetPostId(p.Id)// 		.SetAuthorId(referrer.Id)// 		.SetMeta(meta1)//.Select(comment.Select{
// 		Post: &post.Select{
// 			Id:     true,
// 			Title:  true,
// 			Author: user.Query().Where(user.Id.EQ(p.Id)).OrderBy(user.Id.Asc()),
// 		},
// 	}).
// 		Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to create comment 1: %v", err)
// 	}

// 	meta2 := json.RawMessage(`{"rating":4,"verified":false}`)
// 	_, err = db.Comment.Create()// 		.SetTextify(200)// 		.SetDummy3("dummy_val_3")// 		.SetDummy1(84)// 		.SetDummy2("dummy_val_4")// 		.SetPostId(p.Id)// 		.SetAuthorId(referred.Id)// 		.SetMeta(meta2)//.Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to create comment 2: %v", err)
// 	}

// 	return &SeedData{
// 		ReferrerId: referrer.Id,
// 		ReferredId: referred.Id,
// 		PostId:     p.Id,
// 		Meta1:      meta1,
// 		Meta2:      meta2,
// 	}
// }

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
