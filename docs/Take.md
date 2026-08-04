---
title: Take
description: Limit the number of records returned by a query. Mirrors SQL's LIMIT.
category: Pagination & Sort
tags: [take, limit, pagination]
categoryOrder: 5
order: 1
---

# Take

`Take` limits the number of records returned by a query. It mirrors SQL's `LIMIT`.

A positive value returns the first `N` matching records. A negative value combined with a [`Cursor`](/docs/Cursor) walks backward from that cursor, returning the `N` records that precede it.

`Take` can be used with `FindFirst`, `FindMany`, and `Count`.

## Limiting Results

Pass a positive value to cap the number of rows returned:

```go
users, err := db.User.
    FindMany(user.Bio.Contains("golang")).
    Take(10).
    Exec(ctx)
```

## Walking Backward With a Negative Take

Pass a negative value to return the records _before_ a cursor. This is useful for implementing "previous page" navigation in cursor-based pagination.

```go
prevPage, err := db.User.
    FindMany(user.Email.Contains("@example.com")).
    OrderBy(user.Email.Asc()).
    Cursor(user.Email.EQ(lastEmail)).
    Take(-10).
    Exec(ctx)
```

The results are returned in the configured ordering direction (ascending order in this example).

> **Note:** For `FindFirst`, any negative `Take` behaves like `Take(-1)` - it returns the single record immediately preceding the cursor.

## Supported By

- [`FindMany`](/docs/Read-Records#findmany)
- [`FindFirst`](/docs/Read-Records#findfirst)
- [`Count`](/docs/Count-Records#count)

## Resulting SQL

A positive `Take(10)` maps to SQL's `LIMIT`:

```sql
SELECT * FROM "users" LIMIT 10;
```

When combined with [`Skip`](/docs/Skip), the two become `LIMIT ... OFFSET ...`:

```sql
SELECT * FROM "users" LIMIT 10 OFFSET 20;
```

Negative takes are handled internally with `LIMIT -1` so each database merely returns records preceding the cursor; Phi reverses the result to restore ordering.
