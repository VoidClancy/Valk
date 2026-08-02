---
title: Create Records
description: Create one or more records.
category: CRUD
tags: [create, createMany, createManyAndReturn, onConflict, conflictAction]
categoryOrder: 2
order: 1
---

# Create Records

Phi provides three methods for creating records:

- `Create`
- `CreateMany`
- `CreateManyAndReturn`

The examples in this guide use the following Prisma schema:

```prisma
model User {
  id        String   @id @default(cuid())
  name      String
  email     String   @unique
  password  String
  bio       String?
  createdAt DateTime @default(now())

  posts     Post[]
}

model Post {
  id        String   @id @default(cuid())
  createdAt DateTime @default(now())
  published Boolean  @default(false)
  title     String
  content   String

  author    User     @relation(fields: [authorId], references: [id])
  authorId  String?
}
```

## Create

Creates a single record and returns the created record.

```go
created, err := db.User.Create().
    SetName("John Doe").
    SetEmail("x@y.com").
    SetPassword("secret").
    SetBio("super cool bio").
    Exec(ctx)
```

### Supported Builder Methods

| Method | Description |
| --- | --- |
| [`Select`](/docs/Select) | Select specific fields and relations. Mutually exclusive with `Omit`. |
| [`Omit`](/docs/Omit) | Omit specific scalar fields. Mutually exclusive with `Select`. |
| [`OnConflict`](/docs/Upsert-Records#onconflict) | Configure behavior when a unique constraint is violated. Supports `.Ignore()`, `.UpdateNewValues()`, and `.Update()`. |

---

## CreateMany

Creates multiple records and returns the number of records created.

Accepts a variadic list of `*CreateBuilder` values.

```go
createdCount, err := db.User.CreateMany(
    db.User.Create().
        SetName("John Doe").
        SetEmail("x@y.com").
        SetPassword("secret"),

    db.User.Create().
        SetName("Jane Doe").
        SetEmail("a@b.com").
        SetPassword("secret"),
).Exec(ctx)
```

Or build the list dynamically:

```go
var usersToCreate []*user.CreateBuilder

for i := range 20 {
    usersToCreate = append(usersToCreate,
        db.User.Create().
            SetName(fmt.Sprintf("user-%d", i)).
            SetEmail(fmt.Sprintf("email-%d@x.com", i)).
            SetPassword(fmt.Sprintf("secret-%d", i)),
    )
}

createdCount, err := db.User.CreateMany(usersToCreate...).Exec(ctx)
```

> **Note:** `user.CreateBuilder` is a type alias for `phi.UserCreateBuilder`. They are identical.

### Supported Builder Methods

| Method | Description |
| --- | --- |
| [`OnConflict`](/docs/Upsert-Records#onconflict) | Configure behavior when a unique constraint is violated. Supports `.Ignore()`, `.UpdateNewValues()`, and `.Update()`. |
| [`SkipDuplicates`](/docs/SkipDuplicates) | Shorthand for `.OnConflict().Ignore()`. |

---

## CreateManyAndReturn

Creates multiple records and returns the created records.

Accepts a variadic list of `*CreateBuilder` values.

```go
createdUsers, err := db.User.CreateManyAndReturn(
    db.User.Create().
        SetName("John Doe").
        SetEmail("x@y.com").
        SetPassword("secret"),

    db.User.Create().
        SetName("Jane Doe").
        SetEmail("a@b.com").
        SetPassword("secret"),
).Exec(ctx)
```

Or build the list dynamically:

```go
var usersToCreate []*user.CreateBuilder

for i := range 20 {
    usersToCreate = append(usersToCreate,
        db.User.Create().
            SetName(fmt.Sprintf("user-%d", i)).
            SetEmail(fmt.Sprintf("email-%d@x.com", i)).
            SetPassword(fmt.Sprintf("secret-%d", i)),
    )
}

createdUsers, err := db.User.CreateManyAndReturn(usersToCreate...).Exec(ctx)
```

### Supported Builder Methods

| Method | Description |
| --- | --- |
| [`Select`](/docs/Select) | Select specific fields and relations. Mutually exclusive with `Omit`. |
| [`Omit`](/docs/Omit) | Omit specific scalar fields. Mutually exclusive with `Select`. |
| [`OnConflict`](/docs/Upsert-Records#onconflict) | Configure behavior when a unique constraint is violated. Supports `.Ignore()`, `.UpdateNewValues()`, and `.Update()`. |
| [`SkipDuplicates`](/docs/SkipDuplicates) | Shorthand for `.OnConflict().Ignore()`. |