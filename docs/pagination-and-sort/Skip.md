---
title: Skip
description: Skip a number of matching records before returning results. Mirrors SQL's OFFSET.
category: Pagination & Sort
tags: [skip, offset, pagination]
---

# Skip

`Skip` discards a number of matching records before returning the remaining ones. It mirrors SQL's `OFFSET`.

`Skip` can be used with `FindFirst`, `FindMany`, and `Count`.

## Skipping Records

Pass a non-negative offset to ignore the first `N` matching records:

```go
users, err := db.User.
    FindMany(user.Bio.Contains("golang")).
    Skip(20).
    Take(10).
    Exec(ctx)
```

This returns records 21 through 30 of the matching set.

## Combining With a Cursor

When combined with a [`Cursor`](/docs/Cursor), `Skip` applies _after_ the cursor has positioned the query, discarding the next `N` records:

```go
users, err := db.User.
    FindMany(user.Email.Contains("@example.com")).
    OrderBy(user.Email.Asc()).
    Cursor(user.Email.EQ(lastEmail)).
    Skip(10).
    Take(10).
    Exec(ctx)
```

This returns the 10 records that follow the cursor after skipping the next 10.

## Supported By

- [`FindMany`](/docs/Read-Records#findmany)
- [`FindFirst`](/docs/Read-Records#findfirst)
- [`Count`](/docs/Count-Records#count)

## Resulting SQL

`Skip(20)` maps to SQL's `OFFSET`. On databases that require a `LIMIT` alongside it, Phi emits `LIMIT -1`:

```sql
SELECT * FROM "users" LIMIT -1 OFFSET 20;
```

With [`Take`](/docs/Take) it becomes a standard `LIMIT ... OFFSET` pair:

```sql
SELECT * FROM "users" LIMIT 10 OFFSET 20;
```
