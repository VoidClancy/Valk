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
	user1, err := db.User.Create().SetEmail("a").SetPhoneNum("11").Exec(ctx)
	if err != nil {
		log.Fatalf("failed to seed users: %v", err)
	}
	user2, err := db.User.Create().SetEmail("a").SetPhoneNum("11").OnConflict(user.EmailPhone).Ignore().Exec(ctx)
	if err != nil {
		log.Fatalf("failed to seed users: %v", err)
	}
	fmt.Println("USER1:")
	printJSON(user1)
	fmt.Println("USER2:")
	printJSON(user2)

	db.User.Create().SetEmail("xxxc").SetPhoneNum("6969").SetId("Bleh").Exec(ctx)
	deletedCnt, err := db.User.DeleteMany(user.Id.EQ("Bleh"), user.Password.Contains("xx")).Exec(ctx)
	fmt.Printf("DELETED %d USERS \n", deletedCnt)

	usersCnt, err := db.User.Count(
		user.Id.Contains("x"),
		user.Or(
			user.Password.Contains("y"),
			user.Email.NEQ("c"),
		)).
		Exec(ctx)
	fmt.Printf("COUNT %d USERS \n", usersCnt)

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
// 	referrer, err := db.User.Create(
// 		user.Email.Set("referrer@example.com"),
// 		user.PhoneNum.Set("555-0001"),
// 		user.Password.Set("pass123"),
// 		user.Role.Set(valk.UserRole.Student),
// 	).Select(user.Select{
// 		Id:    true,
// 		Email: true,
// 	}).
// 		Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to create referrer: %v", err)
// 	}

// 	referred, err := db.User.Create(
// 		user.Email.Set("referred@example.com"),
// 		user.PhoneNum.Set("555-0002"),
// 		user.Password.Set("pass456"),
// 		user.Role.Set(valk.UserRole.Student),
// 		user.ReferredById.Set(referrer.Id),
// 	).Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to create referred: %v", err)
// 	}

// 	prof, err := db.Profile.Create(
// 		profile.Bio.Set("BLEH"),
// 		profile.UserId.Set(referred.Id),
// 		profile.CreatedAt.Set(time.Now()),
// 	).Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to create profile: %v", err)
// 	}
// 	fmt.Println("PROFILE:")
// 	printJSON(prof)

// 	categoryTest, err := db.Category.Create(
// 		category.Name.Set("TEST"),
// 	).Exec(ctx)

// 	if err != nil {
// 		log.Fatalf("failed to create category: %v", err)
// 	}
// 	fmt.Println("CATEGORY:")
// 	printJSON(categoryTest)

// 	p, err := db.Post.Create(
// 		post.Title.Set("Valkyrie ORM Deep Dive"),
// 		post.Content.Set("skrrrt"),
// 		post.AuthorId.Set(referred.Id),
// 	).Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to create post: %v", err)
// 	}

// 	cat, err := db.Category.Create(
// 		category.Name.Set("Programming"),
// 	).Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to create category: %v", err)
// 	}

// 	_, err = db.CategoryToPost.Create(
// 		categoryToPost.PostId.Set(p.Id),
// 		categoryToPost.CategoryId.Set(cat.Id),
// 	).Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to create CategoryToPost: %v", err)
// 	}

// 	meta1 := json.RawMessage(`{"rating":5,"verified":true}`)
// 	_, err = db.Comment.Create(
// 		comment.Textify.Set(100),
// 		comment.Dummy3.Set("dummy_val_1"),
// 		comment.Dummy1.Set(42),
// 		comment.Dummy2.Set("dummy_val_2"),
// 		comment.PostId.Set(p.Id),
// 		comment.AuthorId.Set(referrer.Id),
// 		comment.Meta.Set(meta1),
// 	).Select(comment.Select{
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
// 	_, err = db.Comment.Create(
// 		comment.Textify.Set(200),
// 		comment.Dummy3.Set("dummy_val_3"),
// 		comment.Dummy1.Set(84),
// 		comment.Dummy2.Set("dummy_val_4"),
// 		comment.PostId.Set(p.Id),
// 		comment.AuthorId.Set(referred.Id),
// 		comment.Meta.Set(meta2),
// 	).Exec(ctx)
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
	author, err := tx.User.Create(
		user.Email.Set("clancySizer@gmail.com"),
		user.PhoneNum.Set("+1234567890"),
	).Exec(ctx)
	if err != nil {
		fmt.Printf("failed to create user: %+v", err)
		return
	}

	postWithAuthor, err := tx.Post.Create(
		post.Title.Set("A Post"),
		post.AuthorId.Set(author.Id),
	).Select(post.Select{
		Id:    true,
		Title: true,
		Author: user.Query().Select(user.Select{
			Email: true,
		}),
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

		author, err := tx.User.Create(
			user.Email.Set("clancySizer@gmail.com"),
			user.PhoneNum.Set("+1234567890"),
		).Exec(ctx)
		if err != nil {
			return err
		}

		postWithAuthor, err := tx.Post.Create(
			post.Title.Set("A Post"),
			post.AuthorId.Set(author.Id),
		).Select(post.Select{
			Id:    true,
			Title: true,
			Author: user.Query().Select(user.Select{
				Email: true,
			}),
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
