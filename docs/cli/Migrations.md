---
title: Migrations
description: Declarative, forward-only migration workflows and Goose-compatible DDL file generation in Phi.
category: CLI
tags: [migrations, goose, atlas, ddl, schema-diff, phi-migrate]
---

# Migrations

Phi features a production-grade, **forward-only** migration engine (`phi migrate`). Inspired by Prisma's declarative migration workflow and powered by the Atlas DDL calculation engine, Phi diffs your `schema.prisma` against your database state and generates standard, Goose-compatible SQL migration files.

---

## 1. Migration Model: Forward-Only

Phi adopts a **forward-only** migration strategy. In production database engineering, rolling back migrations via down scripts frequently leads to data corruption or accidental column dropping.

When your schema evolves:

1. Update `schema.prisma`.
2. Run `phi migrate <migration_name>`.
3. Phi calculates the exact SQL delta, appends a new versioned `.sql` migration file, and executes it.

---

## 2. CLI Migration Workflows

### Generating & Applying Migrations

Run the CLI command to create and apply a migration:

```bash
phi migrate <migration_name>
```

Or using the short flag:

```bash
phi -m <migration_name>
```

### What Happens During `phi migrate`:

1. **Database Initialization**: Creates the target database automatically if it does not already exist.
2. **Schema Diffing**: Uses Atlas DDL engine to compare your Prisma schema against current database tables.
3. **Goose SQL Generation**: Writes a versioned SQL file (e.g. `./phi/migrations/00001_init.sql`) containing standard Goose migration headers:
    ```sql
    -- +goose Up
    CREATE TABLE "User" (
      "id" TEXT NOT NULL,
      "email" TEXT NOT NULL,
      PRIMARY KEY ("id")
    );
    ```
4. **Execution**: Applies the migration immediately to the database.

---

## 3. Goose & Embed Compatibility

Because Phi migration files use standard Goose `-- +goose Up` headers, you can embed migration scripts directly into your Go binaries using Go 1.16+ `embed.FS`:

```go
package main

import (
    "embed"
    "github.com/pressly/goose/v3"
)

//go:embed phi/migrations/*.sql
var embedMigrations embed.FS

func runEmbeddedMigrations(db *sql.DB) error {
    goose.SetBaseFS(embedMigrations)
    return goose.Up(db, "phi/migrations")
}
```

---

## 4. Dialect-Aware Migration DDL

Phi's migration generator adjusts DDL syntax based on your `schema.prisma` `datasource db { provider = "..." }`:

| Feature               | PostgreSQL DDL                            | SQLite DDL                               |
| :-------------------- | :---------------------------------------- | :--------------------------------------- |
| **Enums**             | Native `CREATE TYPE "Enum" AS ENUM (...)` | Column `TEXT` + `CHECK ("col" IN (...))` |
| **UUID Primary Keys** | `UUID NOT NULL`                           | `TEXT NOT NULL`                          |
| **AutoIncrement**     | `BIGSERIAL` / `SERIAL`                    | `INTEGER PRIMARY KEY AUTOINCREMENT`      |
| **JSON / JSONB**      | `JSONB`                                   | `TEXT` / `BLOB`                          |
