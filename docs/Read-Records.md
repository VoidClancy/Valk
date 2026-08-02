---
title: Read Records
description: Learn how to query records using FindUnique, FindFirst, and FindMany.
category: CRUD
tags: [read, query, findUnique, findFirst, findMany]
order: 2
---

# Read Records

Phi provides three methods for reading records:

- `FindUnique`
- `FindFirst`
- `FindMany`

All three methods support the same predicate system and share the same selection API.

## FindUnique

Retrieves a single record by a required unique predicate. Additional predicates may be provided to further constrain the query, but the first argument **must** be a unique predicate.

### Examples

```go
user, err := db.User.
    FindUnique(
        user.Email.EQ("x@y.com"),
    ).
    Exec(ctx)
```

```go
user, err := db.User.
    FindUnique(
        user.Email.EQ("x@y.com"),
        user.Bio.Contains("the"),
    ).
    Exec(ctx)
```

```go
user, err := db.User.
    FindUnique(
        user.Email.EQ("x@y.com"),
        user.Or(
            user.Bio.Contains("the"),
            user.PhoneNum.HasPrefix("+1"),
        ),
    ).
    Exec(ctx)
```

### Supported Builder Methods

| Method | Description |
| --- | --- |
| [`Select`](/docs/Select) | Select specific fields and relations. |
| [`Omit`](/docs/Omit) | Omit specific scalar fields. |

---

## FindFirst

Retrieves the first record that satisfies all supplied predicates.

### Examples

```go
user, err := db.User.
    FindFirst(
        user.Email.EQ("x@y.com"),
    ).
    Exec(ctx)
```

```go
user, err := db.User.
    FindFirst(
        user.Email.EQ("x@y.com"),
        user.Bio.Contains("the"),
    ).
    Exec(ctx)
```

```go
user, err := db.User.
    FindFirst(
        user.Email.EQ("x@y.com"),
        user.Or(
            user.Bio.Contains("the"),
            user.PhoneNum.HasPrefix("+1"),
        ),
    ).
    Exec(ctx)
```

### Supported Builder Methods

| Method | Description |
| --- | --- |
| [`Select`](/docs/Select) | Select specific fields and relations. |
| [`Omit`](/docs/Omit) | Omit specific scalar fields. |
| [`Skip`](/docs/Skip) | Skip a number of matching records. Mirrors SQL's `OFFSET`. |
| [`Take`](/docs/Take) | Limit the number of returned records. Mirrors SQL's `LIMIT`. |
| [`OrderBy`](/docs/OrderBy) | Specify the ordering of results. Mirrors SQL's `ORDER BY`. |
| [`Cursor`](/docs/Cursor) | Continue querying from a specific cursor. |

---

## FindMany

Retrieves all records that satisfy all supplied predicates.

### Examples

```go
users, err := db.User.
    FindMany(
        user.Email.EQ("x@y.com"),
    ).
    Exec(ctx)
```

```go
users, err := db.User.
    FindMany(
        user.Email.EQ("x@y.com"),
        user.Bio.Contains("the"),
    ).
    Exec(ctx)
```

```go
users, err := db.User.
    FindMany(
        user.Email.EQ("x@y.com"),
        user.Or(
            user.Bio.Contains("the"),
            user.PhoneNum.HasPrefix("+1"),
        ),
    ).
    Exec(ctx)
```

### Supported Builder Methods

| Method | Description |
| --- | --- |
| [`Select`](/docs/Select) | Select specific fields and relations. |
| [`Omit`](/docs/Omit) | Omit specific scalar fields. |
| [`Skip`](/docs/Skip) | Skip a number of matching records. Mirrors SQL's `OFFSET`. |
| [`Take`](/docs/Take) | Limit the number of returned records. Mirrors SQL's `LIMIT`. |
| [`OrderBy`](/docs/OrderBy) | Specify the ordering of results. Mirrors SQL's `ORDER BY`. |
| [`Cursor`](/docs/Cursor) | Continue querying from a specific cursor. |