---
title: Introduction
description: Overview, architecture, and design philosophy of the Phi Go ORM framework.
category: Overview
tags: [introduction, overview, architecture, prisma, go-orm]
---

# Introduction

**Phi** is a high-performance, type-safe ORM framework for Go. Inspired by Prisma, Phi uses your `schema.prisma` file as the single source of truth to generate strongly-typed Go clients, database migrations, and index-optimized SQL queries.

---

## Core Features & Philosophy

### 1. Prisma Schema as Single Source of Truth
Define your data models, enums, indices, and relation constraints in standard Prisma Schema language (`schema.prisma`). Phi parses your schema and generates all Go structs, field helpers, and DDL migrations automatically.

### 2. 100% Compile-Time Type Safety
Say goodbye to string-based column names and dynamic maps in query logic. Phi generates explicit field accessors for every model attribute:

```go
// Fully type-safe builder: invalid field types or misspelled columns fail at compile-time!
users, err := db.User.FindMany(
    user.Role.EQ(phi.UserRole_ADMIN),
    user.LoginCount.GTE(10),
).Exec(ctx)
```

### 3. Zero-Allocation Query Building
Phi is built from the ground up for speed. Query builders construct AST nodes directly into SQL strings without runtime reflection or heavy memory allocations during query assembly.

### 4. Dual Dialect Support (PostgreSQL & SQLite)
Write one unified Go codebase that works seamlessly across **PostgreSQL** and **SQLite**. Phi handles dialect differences (such as native PostgreSQL ENUMs vs. SQLite `TEXT` + `CHECK` constraints) transparently.

### 5. Production-Ready Tooling
* **Extension Hooks**: Middleware hooks for intercepting, mutating, or caching queries (`Create`, `Read`, `Update`, `Delete`).
* **Transactions**: Closure-based transactions with panic recovery and automatic rollback.
* **Forward-Only Migrations**: Goose-compatible migration generator powered by the Atlas DDL engine.
