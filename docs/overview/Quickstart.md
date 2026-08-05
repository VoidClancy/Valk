---
title: Quickstart
description: Install the Phi CLI, create a project, define a schema, migrate, and generate the type-safe Go ORM client.
category: Overview
tags: [quickstart, install, setup]
---

# Phi Documentation

## Quickstart

### Install the Phi CLI

```bash
go install github.com/voidclancy/phi@latest
```

---

## Create a New Project

Initialize a Go module if you don't already have one.

```bash
mkdir demo
cd demo
go mod init demo
```

---

## Initialize Phi

Generate a configuration file:

```bash
phi init [directory]
```

By default, this creates a **`phi.yml`** file.

To generate a specific format:

**JSON**

```bash
phi init json [directory]
```

**YAML**

```bash
phi init yml [directory]
```

Example `phi.yml`:

```yaml
database:
    url_env: DATABASE_URL # Environment variable for the client connection.
    direct_url_env: DATABASE_DIRECT_URL # Environment variable used for migrations.

schema: ./schema.prisma # Path to your Prisma schema.

client_name: phi # Name of the generated Go package.

output:
    client: ./phi
    migrations: ./phi/migrations # Keep this inside the client directory if you plan to embed migrations.

log:
    - none # Available: none, all, query, warn, error
```

---

## Define Your Schema

Create a `schema.prisma` file.

```prisma
datasource db {
  provider = "sqlite"
}

model User {
  id        String   @id @default(cuid())
  username  String   @unique
  email     String   @unique
  bio       String?
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
}
```

> **Tip:** Install the Prisma Language Server extension for syntax highlighting and autocompletion.

---

## Create a Migration

Generate and apply a migration:

```bash
phi migrate <migration_name>
```

or

```bash
phi -m <migration_name>
```

This command will:

- Create the database if it doesn't exist.
- Generate a migration.
- Apply the migration.

---

## Generate the Client

Generate a type-safe client from your schema:

```bash
phi generate
```

or

```bash
phi -g
```

The client is generated in the directory specified by your configuration file.

Run this command whenever your schema changes.

---

# Usage

Create a `main.go` file.

### `main.go`

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"demo/phi"
	"demo/phi/user"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	db, err := phi.Open(
		"sqlite3",
		"file:memdb1?mode=memory&cache=shared&_pragma=foreign_keys(1)&_time_format=sqlite",
	)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	db.Raw().SetMaxOpenConns(10)

	ctx := context.Background()

	// Create a user.
	createdUser, err := db.User.Create().
		SetUsername("VoidClancy").
		SetEmail("x@y.com").
		SetBio("super cool bio").
		Exec(ctx)
	if err != nil {
		return err
	}

	result, _ := json.MarshalIndent(createdUser, "", "  ")
	fmt.Printf("created user:\n%s\n", result)

	// Find the user.
	foundUser, err := db.User.FindUnique(
		user.Id.EQ(createdUser.Id),
	).Exec(ctx)
	if err != nil {
		return err
	}

	result, _ = json.MarshalIndent(foundUser, "", "  ")
	fmt.Printf("retrieved user:\n%s\n", result)

	if foundUser.Bio != nil {
		fmt.Printf("User bio: %s\n", *foundUser.Bio)
	}

	return nil
}
```

### Example output

```text
created user:
{
  "id": "cmf8z6qp40000a1b2c3d4e5f6",
  "username": "VoidClancy",
  "email": "x@y.com",
  "bio": "super cool bio",
  "createdAt": "2026-07-30T12:48:16.314Z",
  "updatedAt": "2026-07-30T12:48:16.314Z"
}

retrieved user:
{
  "id": "cmf8z6qp40000a1b2c3d4e5f6",
  "username": "VoidClancy",
  "email": "x@y.com",
  "bio": "super cool bio",
  "createdAt": "2026-07-30T12:48:16.314Z",
  "updatedAt": "2026-07-30T12:48:16.314Z"
}

User bio: super cool bio
```
