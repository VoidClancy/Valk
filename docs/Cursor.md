---
title: Cursor
description: Fetch pages of results relative to a known record, identified by a unique predicate.
category: Pagination & Sort
tags: [cursor, pagination, keyset]
categoryOrder: 5
order: 4
---

# Cursor

`Cursor` implements cursor-based (keyset) pagination. Instead of skipping an arbitrary number of rows like [`Skip`](/docs/Skip), you position the query _after_ a specific record, identified by a unique predicate such as `user.Email.EQ(...)` or `user.Id.EQ(...)`.

This approach stays stable under concurrent inserts and updates and is far more efficient on large datasets than offset pagination.

`Cursor` can be used with `FindFirst`, `FindMany`, and nested relation query builders. It is typically combined with [`OrderBy`](/docs/OrderBy) and [`Take`](/docs/Take).

## Fetching the Next Page

Combine a cursor with a positive [`Take`](/docs/Take) to fetch the page following a known record:

```go
nextPage, err := db.User.
    FindMany(user.Email.Contains("@example.com")).
    OrderBy(user.Email.Asc()).
    Cursor(user.Email.EQ(lastEmail)).
    Take(20).
    Exec(ctx)
```

This returns the next 20 records after `lastEmail`.

## Fetching the Previous Page

Combine a cursor with a negative [`Take`](/docs/Take) to fetch the records that precede it:

```go
prevPage, err := db.User.
    FindMany(user.Email.Contains("@example.com")).
    OrderBy(user.Email.Asc()).
    Cursor(user.Email.EQ(firstEmailOfCurrentPage)).
    Take(-20).
    Exec(ctx)
```

> **Note:** The cursor value must reference a unique field, such as the record's `Id` or another unique column.

## Usage Tips

- Always provide a stable [`OrderBy`](/docs/OrderBy); it must be consistent from page to page.
- Pass the last record of the current page as the cursor for the next page, and the first record for the previous page.
- Use [`Skip`](/docs/Skip) together with a cursor when you want to jump a fixed number of rows past the cursor.

## Supported By

- [`FindMany`](/docs/Read-Records#findmany)
- [`FindFirst`](/docs/Read-Records#findfirst)
- Nested relation queries via [`Select`](/docs/Select)

## Resulting SQL

A cursor query compiles into a key-based comparison on the cursor's column, combined with the ordering and limit:

```sql
SELECT * FROM "users"
WHERE "email" > 'last@example.com'
ORDER BY "email" ASC
LIMIT 20;
```

A negative [`Take`](/docs/Take) swaps the comparison direction to fetch preceding rows.
