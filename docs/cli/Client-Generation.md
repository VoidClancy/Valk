---
title: Client Generation
description: Generate strongly-typed Go client code from schema.prisma using the phi generate command.
category: CLI
tags: [cli, generate, client, code-generation, phi-g]
---

# Client Generation Command

The `phi generate` command parses your `schema.prisma` file and generates a 100% type-safe, zero-allocation Go ORM client in the configured output directory.

---

## 1. Quick Usage

Run the generate command from your project root:

```bash
phi generate
```

Or using the short flag:

```bash
phi -g
```

---

## 2. Generation Workflow

When `phi generate` is executed, the CLI performs the following steps:

1. **Configuration Resolution**: Reads `phi.yml` or `phi.json` to determine schema input path and client output directory.
2. **Prisma Schema Parsing**: Parses data models, scalar fields, enums, relation constraints, default values, and native `@db.*` types.
3. **Template Rendering**: Executes Phi's internal Go code generation templates:
   * **`client.go`**: Core DB client, connection handles, predicate data types, and dialect queries.
   * **`[model].go`**: Model structs, query builders, create/update/upsert inputs, select/omit maps, and model delegates.
   * **`[model]/[model].go`**: Isolated model sub-packages exporting type-safe field constants (`user.Email`, `user.Role`, etc.).
   * **`errors.go` & `validation.go`**: Sentinel error definitions, driver error translation engines, and in-memory input validation guards.
4. **Formatting & Cleanup**: Ensures generated Go source files are cleanly formatted and ready for instant compilation.

---

## 3. When to Re-Run `phi generate`

Re-run `phi generate` whenever you update your Prisma schema:
* Adding or renaming models or fields.
* Modifying field types or adding `@unique` / `@id` constraints.
* Defining or updating enum values.
* Changing relation rules (`onDelete`).

> [!TIP]
> Add `phi generate` to your project Makefile or CI build script to ensure generated client code is always in sync with your `schema.prisma`.

---

## 4. Generated Package Architecture

For a schema containing a `User` model, `phi generate` outputs the following package structure under `./phi`:

```text
phi/
├── client.go           # Core Client instance, DB handle, raw SQL helpers
├── runtime.go          # Dialect execution engine & statement caching
├── user.go             # User model struct, UserQueryBuilder, UserCreateBuilder, UserUpdateBuilder
└── user/
    └── user.go         # Model package exporting field constants (user.Email, user.Id, etc.)
```
