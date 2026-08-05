# Phi

Phi is a compile-time ORM for Go that mirrors the Prisma developer experience. It parses your `schema.prisma` file and generates a strongly-typed Go database client with zero runtime reflection.

[Documentation Website](https://phi-orm.vercel.app) | [pkg.go.dev](https://pkg.go.dev/github.com/voidclancy/phi)

## Key Features

* **Compile-time Type Safety**: Queries, predicates, and selections fail at build time if invalid.
* **Zero Reflection**: Compiles down to standard `database/sql` calls for maximum query execution speed.
* **Prisma Schema DX**: Use your existing `schema.prisma` models, enums, native attributes, and relations.
* **Declarative Migrations**: Integrated forward-only migration engine powered by Atlas DDL diffing.
* **Extension Hooks**: Intercept and mutate CRUD operations via middleware closures.
* **Multi-Provider**: Native support for PostgreSQL and SQLite.

## Installation

Install the Phi CLI tool:

```bash
go install github.com/voidclancy/phi@latest
```

## Quick Start

1. Initialize your configuration file:

```bash
phi init
```

2. Generate your type-safe Go client:

```bash
phi generate
```

3. Query your database:

```go
db, err := phi.Open("sqlite3",
		"file:memdb1?mode=memory&cache=shared&_pragma=foreign_keys(1)&_time_format=sqlite")

	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

    defer db.Close()

u, err := db.User.FindUnique(
    user.Email.EQ("user@example.com"),
).Exec(ctx)
```

## Documentation

Full guides, API references, CLI commands, and hooks documentation are available on the documentation website:

[https://phi-orm.vercel.app](https://phi-orm.vercel.app)

## License

[Apache 2.0 License](LICENSE)