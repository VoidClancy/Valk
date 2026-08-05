---
title: Order By
description: Sort query results using Asc or Desc on any scalar field. Mirrors SQL's ORDER BY.
category: Pagination & Sort
tags: [order, sort, asc, desc]
---

# Order By

`OrderBy` sorts the returned records using the generated `.Asc()` and `.Desc()` methods available on every scalar field. It mirrors SQL's `ORDER BY`.

`OrderBy` can be used with `FindFirst`, `FindMany`, and nested relation query builders.

## Sorting Ascending

```go
users, err := db.User.
    FindMany(user.Bio.Contains("golang")).
    OrderBy(user.CreatedAt.Asc()).
    Exec(ctx)
```

## Sorting Descending

```go
users, err := db.User.
    FindMany(user.Bio.Contains("golang")).
    OrderBy(user.CreatedAt.Desc()).
    Exec(ctx)
```

## Sort by Multiple Columns

Pass multiple columns to `OrderBy`. The first column takes precedence, then the next, and so on:

```go
users, err := db.User.
    FindMany(user.Bio.Contains("golang")).
    OrderBy(
        user.Role.Asc(),
        user.CreatedAt.Desc(),
    ).
    Exec(ctx)
```

## Ordering Nested Relations

`OrderBy` is also available on nested to-many relation query builders within a [`Select`](/docs/Select):

```go
users, err := db.User.
    FindMany(user.Email.Contains("@example.com")).
    Select(user.Select{
        Id: true,
        Posts: post.Query().
            Where(post.Published.EQ(true)).
            OrderBy(post.CreatedAt.Desc()).
            Select(post.Select{
                Id:    true,
                Title: true,
            }),
    }).
    Exec(ctx)
```

> **Note:** When using cursor-based pagination, `OrderBy` must be supplied and should be consistent across pages for stable results.

> **Note:** On cursor queries, Phi appends the model's primary key columns to the `ORDER BY` clause when no unique column is present. If you provide no `OrderBy` at all, it orders solely by the primary key. If the provided columns aren't unique, the primary key is appended as a tie-breaker to keep pagination deterministic. Appended columns follow the same direction as the query (ascending for forward pages, descending for negative takes).

## Supported By

- [`FindMany`](/docs/Read-Records#findmany)
- [`FindFirst`](/docs/Read-Records#findfirst)
- Nested relation queries via [`Select`](/docs/Select)

## Resulting SQL

Each `.Asc()` / `.Desc()` maps to a column in the `ORDER BY` clause:

```sql
SELECT * FROM "users"
ORDER BY "role" ASC, "createdAt" DESC;
```
