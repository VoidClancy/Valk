package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

type EvolutionStep struct {
	Name    string
	Schema  string
	Imports []string
	Code    string
}

var evolutionSteps = []EvolutionStep{
	{
		Name: "001_init",
		Schema: `
model User {
  id    String @id
  email String @unique
  age   Int
}
`,
		Imports: []string{},
		Code: `
			
			u, err := db.User.Create().SetId("user-1").SetEmail("user1@example.com").SetAge(30).Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to seed user-1: %w", err)
			}
			if u.Id != "user-1" {
				return fmt.Errorf("expected user-1, got %q", u.Id)
			}
		`,
	},
	{
		Name: "002_optional_fields",
		Schema: `
model User {
  id    String  @id
  email String  @unique
  age   Int?
  role  String?
}
`,
		Imports: []string{
			`user "integration/sandbox/valk/user"`,
		},
		Code: `
			// Verify user-1 is intact
			u1, err := db.User.FindUnique(user.Id.EQ("user-1")).Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to find user-1: %w", err)
			}
			if u1.Age == nil || *u1.Age != 30 {
				return fmt.Errorf("invalid age: expected 30, got %v", u1.Age)
			}
			if u1.Role != nil {
				return fmt.Errorf("expected role to be nil, got %q", *u1.Role)
			}

			// Create user-2 with nil age and role admin
			u2, err := db.User.Create().SetId("user-2").SetEmail("user2@example.com").SetRole("admin").Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to create user-2: %w", err)
			}
			if u2.Age != nil {
				return fmt.Errorf("expected age to be nil, got %v", *u2.Age)
			}
			if u2.Role == nil || *u2.Role != "admin" {
				return fmt.Errorf("expected role to be admin, got %v", u2.Role)
			}
		`,
	},
	{
		Name: "003_defaults_and_uniques",
		Schema: `
model User {
  id     String  @id
  email  String  @unique
  age    Int?
  role   String?
  phone  String? @unique
  status String? @default("active")
}
`,
		Imports: []string{
			`user "integration/sandbox/valk/user"`,
		},
		Code: `
			// Verify old records received the default values for the new status column
			u1, err := db.User.FindUnique(user.Id.EQ("user-1")).Exec(ctx)
			if err != nil {
				return fmt.Errorf("user-1 not found: %w", err)
			}
			if u1.Status == nil || *u1.Status != "active" {
				return fmt.Errorf("expected status 'active' for user-1, got %v", u1.Status)
			}

			// Create user-3 with phone "+12345" and no status (should use default "active")
			u3, err := db.User.Create().SetId("user-3").SetEmail("user3@example.com").SetPhone("+12345").Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to create user-3: %w", err)
			}
			if u3.Status == nil || *u3.Status != "active" {
				return fmt.Errorf("expected default status 'active', got %v", u3.Status)
			}
			if u3.Phone == nil || *u3.Phone != "+12345" {
				return fmt.Errorf("expected phone '+12345', got %v", u3.Phone)
			}

			// Test unique constraint on phone
			_, err = db.User.Create().SetId("user-4").SetEmail("user4@example.com").SetPhone("+12345").Exec(ctx) // duplicate phone!
			if err == nil {
				return fmt.Errorf("expected unique constraint error on phone, got nil")
			}
		`,
	},
	{
		Name: "004_relations",
		Schema: `
model User {
  id     String  @id
  email  String  @unique
  age    Int?
  role   String?
  phone  String?
  status String? @default("active")
  posts  Post[]
}

model Post {
  id       String @id
  title    String
  authorId String
  author   User   @relation(fields: [authorId], references: [id])
}
`,
		Imports: []string{
			`post "integration/sandbox/valk/post"`,
		},
		Code: `
			// Create a post for user-1
			p, err := db.Post.Create().SetId("post-1").SetTitle("ORM Evolution").SetAuthorId("user-1").Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to create post: %w", err)
			}
			if p.Id != "post-1" {
				return fmt.Errorf("expected post-1, got %q", p.Id)
			}

			// Retrieve post with author relationship
			pLoaded, err := db.Post.FindUnique(post.Id.EQ("post-1")).Select(valk.PostSelect{
				Id:    true,
				Title: true,
				Author: &valk.UserSelect{
					Id:    true,
					Email: true,
				},
			}).Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to query post with relations: %w", err)
			}
			if pLoaded.Author == nil {
				return fmt.Errorf("expected author to be loaded, got nil")
			}
			if pLoaded.Author.Email != "user1@example.com" {
				return fmt.Errorf("expected author email 'user1@example.com', got %q", pLoaded.Author.Email)
			}
		`,
	},
}

func TestSchemaEvolution(t *testing.T) {
	provider := getActiveProvider()

	// 1. Setup clean sandbox directory
	sandboxDir, err := filepath.Abs("./sandbox")
	if err != nil {
		t.Fatalf("failed to get absolute path for sandbox: %v", err)
	}

	_ = os.RemoveAll(sandboxDir)
	err = os.MkdirAll(sandboxDir, 0755)
	if err != nil {
		t.Fatalf("failed to create sandbox dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(sandboxDir)
	}()

	// Determine connection DSN and prepare temporary databases/schemas
	var dsn string
	if provider == "postgres" {
		mainDsn := getPostgresDSN()
		db, err := sql.Open("postgres", mainDsn)
		if err != nil {
			t.Fatalf("failed to connect to main pg: %v", err)
		}
		_, err = db.Exec("DROP SCHEMA IF EXISTS ephemeral_evolution CASCADE; CREATE SCHEMA ephemeral_evolution;")
		db.Close()
		if err != nil {
			t.Fatalf("failed to recreate schema: %v", err)
		}

		defer func() {
			db, err := sql.Open("postgres", mainDsn)
			if err == nil {
				_, _ = db.Exec("DROP SCHEMA IF EXISTS ephemeral_evolution CASCADE;")
				db.Close()
			}
		}()

		if strings.Contains(mainDsn, "?") {
			dsn = mainDsn + "&search_path=ephemeral_evolution"
		} else {
			dsn = mainDsn + "?search_path=ephemeral_evolution"
		}
	} else {
		dsn = "file:" + filepath.Join(sandboxDir, "evolution.db")
	}

	// Write static valk.json configuration
	valkJson := `{
  "database": {
    "url_env": "DATABASE_URL"
  },
  "schema": "./schema.prisma",
  "output": {
    "client": "./valk",
    "migrations": "./valk/migrations"
  }
}`
	err = os.WriteFile(filepath.Join(sandboxDir, "valk.json"), []byte(valkJson), 0644)
	if err != nil {
		t.Fatalf("failed to write valk.json: %v", err)
	}

	valkBin, err := filepath.Abs("../bin/valk")
	if err != nil {
		t.Fatalf("failed to get absolute path for valk: %v", err)
	}

	// Helper to execute valk binary
	runValk := func(args ...string) {
		cmd := exec.Command(valkBin, args...)
		cmd.Dir = sandboxDir
		cmd.Env = append(os.Environ(),
			"DATABASE_URL="+dsn,
			"DATABASE_DIRECT_URL="+dsn,
		)
		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			t.Fatalf("valk %v failed: %v\nstdout: %s\nstderr: %s", args, err, outBuf.String(), errBuf.String())
		}
	}

	// Run through the evolution steps pipeline
	for _, step := range evolutionSteps {
		t.Logf("Running evolution step: %s", step.Name)

		// A. Write schema file for this step
		schemaFileContent := fmt.Sprintf(`
datasource db {
  provider = "%s"
  url      = env("DATABASE_URL")
}

%s
`, provider, step.Schema)

		err = os.WriteFile(filepath.Join(sandboxDir, "schema.prisma"), []byte(schemaFileContent), 0644)
		if err != nil {
			t.Fatalf("[%s] failed to write schema: %v", step.Name, err)
		}

		// Ensure migrations output folder exists
		err = os.MkdirAll(filepath.Join(sandboxDir, "valk/migrations"), 0755)
		if err != nil {
			t.Fatalf("[%s] failed to create migrations folder: %v", step.Name, err)
		}

		// B. Regenerate client and plan/apply migrations
		runValk("generate")
		runValk("migrate", step.Name)

		// C. Generate and compile main.go execution script for this step
		importsStr := strings.Join(step.Imports, "\n\t")
		goCode := fmt.Sprintf(`package main

import (
	"context"
	"fmt"
	"os"
	"integration/sandbox/valk"
	%s

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %%v\n", err)
		os.Exit(1)
	}
	fmt.Println("SUCCESS")
}

func run() error {
	provider := os.Getenv("PROVIDER")
	dsn := os.Getenv("DATABASE_URL")
	db, err := valk.Open(provider, dsn)
	if err != nil {
		return fmt.Errorf("failed to open client: %%w", err)
	}
	defer db.Close()
	ctx := context.Background()

	%s

	return nil
}
`, importsStr, step.Code)

		err = os.WriteFile(filepath.Join(sandboxDir, "main.go"), []byte(goCode), 0644)
		if err != nil {
			t.Fatalf("[%s] failed to write main.go script: %v", step.Name, err)
		}

		// D. Execute step code and verify success
		cmd := exec.Command("go", "run", "main.go")
		cmd.Dir = sandboxDir
		cmd.Env = append(os.Environ(),
			"PROVIDER="+provider,
			"DATABASE_URL="+dsn,
		)
		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			t.Fatalf("[%s] go run main.go failed: %v\nstdout: %s\nstderr: %s", step.Name, err, outBuf.String(), errBuf.String())
		}

		if !strings.Contains(outBuf.String(), "SUCCESS") {
			t.Fatalf("[%s] execution check failed: %s", step.Name, outBuf.String())
		}
	}
}
