---
title: Upsert Records
description: Configure conflict handling for create operations.
category: CRUD
tags: [upsert, onConflict, conflictAction, create, createManyAndReturn]
categoryOrder: 2
order: 5
---

# Upsert Records

Phi performs upserts through the `OnConflict` API available on create operations.

Rather than exposing a separate `Upsert` method, Phi allows you to configure how `Create`, `CreateMany`, and `CreateManyAndReturn` behave when a unique constraint is violated.

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
```

## OnConflict

`OnConflict` specifies which unique constraint should be handled if an insert would violate it.

For example, the following handles conflicts on the `email` unique constraint:

```go
created, err := db.User.Create().
    SetName("John Doe").
    SetEmail("x@y.com").
    SetPassword("secret").
    SetBio("super cool bio").
    OnConflict(user.Email).Ignore().
    Exec(ctx)
```

`OnConflict` supports the following actions:

- `Ignore()`
- `UpdateNewValues()`
- `Update(func(u *phi.UserUpsert) { ... })`

---

## Ignore

Ignores records that would violate the specified unique constraint.

```go
created, err := db.User.Create().
    SetName("John Doe").
    SetEmail("x@y.com").
    SetPassword("secret").
    SetBio("super cool bio").
    OnConflict(user.Email).Ignore().
    Exec(ctx)
```

This mirrors SQL's `ON CONFLICT DO NOTHING`.

---

## UpdateNewValues

Updates every writable field using the values supplied to `Create`.

```go
created, err := db.User.Create().
    SetName("John Doe").
    SetEmail("x@y.com").
    SetPassword("secret").
    SetBio("super cool bio").
    OnConflict(user.Email).
    UpdateNewValues().
    Exec(ctx)
```

This mirrors SQL's `ON CONFLICT DO UPDATE SET ...` using the inserted values.

---

## Update

Provides full control over the update performed when a conflict occurs.

```go
created, err := db.User.Create().
    SetName("John Doe").
    SetEmail("x@y.com").
    SetPassword("secret").
    SetBio("super cool bio").
    OnConflict(user.Email).Update(func(u *phi.UserUpsert) {
        u.Bio.Set("new bio")
        u.Email.Set("new-email")
    }).
    Exec(ctx)
```

Only the fields specified inside the callback are updated.

> **Note:** `user.Upsert` is a type alias for `phi.UserUpsert` - they are identical. The callback parameter type is generated per model, so you can also write the callback as `func(u *user.Upsert) { ... }` and it will behave exactly the same.

---

## SkipDuplicates

`SkipDuplicates` is available only on `CreateMany` and `CreateManyAndReturn`.

It is equivalent to:

```go
.OnConflict(user.<unique field>).Ignore()
```

and mirrors SQL's `ON CONFLICT DO NOTHING`.

---

## Supported By

`OnConflict` is supported by:

- [`Create`](/docs/Create-Records#create)
- [`CreateMany`](/docs/Create-Records#createmany)
- [`CreateManyAndReturn`](/docs/Create-Records#createmanyandreturn)

`SkipDuplicates` is supported by:

- [`CreateMany`](/docs/Create-Records#createmany)
- [`CreateManyAndReturn`](/docs/Create-Records#createmanyandreturn)
