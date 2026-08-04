# Phi

Phi is a **compile-time ORM for Go** that mirrors the Prisma developer experience. It parses your `schema.prisma` file and generates a type-safe client with zero reflection at runtime.

## Features

- **Compile-time type safety** - queries, predicates, and results are generated from your schema, so mistakes fail at build time, not runtime.
- **Zero runtime overhead** - no reflection, no sidecars. Everything compiles down to standard `database/sql`.
- **Prisma-style schema** - keep writing `schema.prisma` and get types, builders, and migrations generated from it.
- **Multi-provider** - PostgreSQL, SQLite (for now).
- **Schema migrations** - generate and apply migrations straight from the CLI, with optional embedded migrations.

## Documentation

[https://phi-orm.vercel.app](https://phi-orm.vercel.app)