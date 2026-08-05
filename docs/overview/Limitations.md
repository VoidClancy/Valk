---
title: Current Limitations
description: Known boundaries and architectural scope of Phi.
category: Overview
tags: [limitations, roadmap, nested-creates, joins, aggregations]
---

# Current Limitations

Phi focuses on high-performance, zero-allocation Prisma-style query building and code generation. To maintain speed, clarity, and predictable SQL generation, certain features are currently out of scope or planned for future releases.

---

## 1. No Inline Nested Creates (For Now)

Phi does not currently support nested inline record creation inside a single `Create()` builder call (e.g. creating a parent record and its children in one nested builder input).

### Current Pattern:

Use sequential creation inside a closure transaction:

```go
err := db.Transaction(ctx, func(tx *phi.Tx) error {
    user, err := tx.User.Create().SetEmail("user@example.com").Exec(ctx)
    if err != nil {
        return err
    }

    _, err = tx.Post.Create().SetTitle("First Post").SetAuthorId(user.Id).Exec(ctx)
    return err
})
```

---

## 2. No Inline Nested Relation Edges Update (For Now)

Updating foreign key relation links or edges directly through nested relation mutations (such as inline `connect`, `disconnect`, or `update` on relation fields) is not currently supported in `Update()` builders.

### Current Pattern:

Update foreign key fields directly or execute updates on target relation models:

```go
// Update foreign key column directly
updatedPost, err := db.Post.Update(post.Id.EQ("post-1")).
    SetAuthorId("new-user-id").
    Exec(ctx)
```

---

## 3. No Arbitrary SQL Joins (For Now)

Phi does not expose an explicit SQL `.Join()` query builder API.

### Relation Preloading Strategy:

Phi uses Prisma-style **2-pass batch selection** for loading relations:

1. Primary records are queried using clean, index-optimized SQL predicates.
2. Selected relations (`args.Select.Posts = post.Query()`) are fetched via secondary batched `IN (...)` queries and linked in memory.

If your use case requires custom multi-table SQL `JOIN` aggregations, use the raw SQL fallback (`db.Raw()`).

---

## 4. Aggregations & GroupBy

Phi provides `.Count()` for total record counts. Complex aggregations (`SUM`, `AVG`, `MIN`, `MAX`) and `GROUP BY / HAVING` queries are not built into the builder DSL.

### Current Pattern:

Use `db.Raw()` for custom aggregations:

```go
var total int64
err := db.Raw().QueryRowContext(ctx, `SELECT SUM(loginCount) FROM "User"`).Scan(&total)
```

---

## 5. Single-Level Relation Preloading

Relation selection (`args.Select.Posts = post.Query()`) preloads immediate 1-level relations. Deep multi-level nested graph preloading (e.g. `User` -> `Posts` -> `Comments`) is not currently exposed in a single builder tree.

---

---

## 7. Compile-Time Type Safety & Raw Expression Boundaries

Phi intentionally restricts passing raw, un-sanitized SQL string snippets into `.Where(...)` predicate chains. Every builder predicate is strongly typed to its model column to guarantee compile-time type safety and prevent SQL injection vulnerabilities.

If your query requires custom raw SQL expressions (such as complex subqueries or custom `WHERE` clauses), execute them using the raw SQL fallback:

```go
rows, err := db.Raw().QueryContext(ctx, `SELECT * FROM "User" WHERE custom_func(loginCount) = $1`, 5)
```
