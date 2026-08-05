---
title: Create Hooks
description: Intercept and mutate record creation with extension hooks.
category: Hooks
tags: [hooks, extension, create, createMany, createManyAndReturn, middleware]
---

# Create Hooks

Create hooks allow you to intercept `Create`, `CreateMany`, and `CreateManyAndReturn` operations before queries hit the database. You can validate inputs, mutate fields, handle upsert conflicts, or short-circuit query execution entirely.

---

## Schema Context

The examples in this guide reference the following Prisma schema:

```prisma
model User {
  id         String   @id @default(cuid())
  email      String   @unique
  phoneNum   String   @unique
  password   String?
  role       UserRole @default(STUDENT)
  loginCount Int      @default(0)
  createdAt  DateTime @default(now())

  posts      Post[]
}

model Post {
  id        String   @id @default(cuid())
  title     String
  content   String
  published Boolean  @default(false)
  authorId  String?
}
```

---

## Registration & Chaining

Register hooks on a model delegate using `.Use()`, passing a model extension (e.g. `user.Extension`):

```go
db.User.Use(user.Extension{
    Create: func(ctx context.Context, args *phi.UserCreateArgs, next phi.UserCreateQuery) (*phi.User, error) {
        // Pre-hook: runs before database execution
        args.Data.Email = strings.ToLower(args.Data.Email)

        res, err := next(ctx, args)

        // Post-hook: runs after database execution
        return res, err
    },
})
```

### Chaining Order & Context Flow

Multiple `.Use()` calls stack in middleware order (outermost to innermost):

```
Request → Hook A (Pre) → Hook B (Pre) → Database → Hook B (Post) → Hook A (Post) → Response
```

> **Note:** Context values set via `context.WithValue` flow inward to subsequent hooks and query execution, but are isolated from outer hooks during unwinding.

---

## Hook Signatures

| Hook Field            | Signature                                                                                                 | Return Type            |
| :-------------------- | :-------------------------------------------------------------------------------------------------------- | :--------------------- |
| `Create`              | `func(ctx context.Context, args *phi.UserCreateArgs, next phi.UserCreateQuery)`                           | `(*phi.User, error)`   |
| `CreateMany`          | `func(ctx context.Context, args *phi.UserCreateManyArgs, next phi.UserCreateManyQuery)`                   | `(int64, error)`       |
| `CreateManyAndReturn` | `func(ctx context.Context, args *phi.UserCreateManyAndReturnArgs, next phi.UserCreateManyAndReturnQuery)` | `([]*phi.User, error)` |

---

## 1. Single Record (`Create`)

`*phi.UserCreateArgs` exposes the following query properties:

| Field            | Type                         | Description                                                    |
| :--------------- | :--------------------------- | :------------------------------------------------------------- |
| `Data`           | `*phi.UserCreate`            | The record fields being inserted.                              |
| `Select`         | `*phi.UserSelect`            | Selected scalar and relation fields to return.                 |
| `ConflictTarget` | `phi.UniqueConstraintTarget` | Unique column or composite key for upserts.                    |
| `ConflictAction` | `*phi.ConflictAction`        | Resolution action (`DoNothing`, `UpdateNewValues`, or custom). |

### Mutating Inputs

Modify `args.Data` directly before invoking `next`:

```go
db.User.Use(user.Extension{
    Create: func(ctx context.Context, args *phi.UserCreateArgs, next phi.UserCreateQuery) (*phi.User, error) {
        // Normalize email
        args.Data.Email = strings.ToLower(strings.TrimSpace(args.Data.Email))

        // Hash password if present
        if args.Data.Password != nil {
            hashed := hashPassword(*args.Data.Password)
            args.Data.Password = &hashed
        }

        return next(ctx, args)
    },
})
```

> **Note:** Required fields are concrete types, as they must be provided, no need for `nil` check, optional fields are always pointers, as they may be `nil` to indicate "unset".

### Short-Circuiting Execution

Return early without calling `next(ctx, args)` to bypass database insertion:

```go
db.User.Use(user.Extension{
    Create: func(ctx context.Context, args *phi.UserCreateArgs, next phi.UserCreateQuery) (*phi.User, error) {
        if args.Data.Email == "" {
            return nil, errors.New("email is required")
        }
        return next(ctx, args)
    },
})
```

### Relation Sub-queries & Selections

Customize returned scalar fields and pre-load relations on `args.Select`:

```go
db.User.Use(user.Extension{
    Create: func(ctx context.Context, args *phi.UserCreateArgs, next phi.UserCreateQuery) (*phi.User, error) {
        args.Select.Email = true
        args.Select.Posts = post.Query().Where(post.Published.EQ(true))
        return next(ctx, args)
    },
})
```

---

## 2. Bulk Insert (`CreateMany`)

`*phi.UserCreateManyArgs` exposes:

- `Data []*phi.UserCreate`: Slice of record inputs.
- `ConflictTarget` & `ConflictAction`: Upsert rules for bulk conflicts.

### Batch Input Mutation & Validation

Iterate over `args.Data` to validate or modify entries in-place:

```go
db.User.Use(user.Extension{
    CreateMany: func(ctx context.Context, args *phi.UserCreateManyArgs, next phi.UserCreateManyQuery) (int64, error) {
        for _, record := range args.Data {
            record.Email = strings.ToLower(record.Email)
        }
        return next(ctx, args)
    },
})
```

### Appending Records

Use `args.AppendData(...)` to dynamically inject extra records into the batch:

```go
db.User.Use(user.Extension{
    CreateMany: func(ctx context.Context, args *phi.UserCreateManyArgs, next phi.UserCreateManyQuery) (int64, error) {
        args.AppendData(
            db.User.Create().SetEmail("audit1@example.com").SetPhoneNum("+1001"),
            db.User.Create().SetEmail("audit2@example.com").SetPhoneNum("+1002"),
        )
        return next(ctx, args)
    },
})
```

---

## 3. Bulk Insert and Return (`CreateManyAndReturn`)

`*phi.UserCreateManyAndReturnArgs` combines batch data manipulation with relation selection:

```go
db.User.Use(user.Extension{
    CreateManyAndReturn: func(ctx context.Context, args *phi.UserCreateManyAndReturnArgs, next phi.UserCreateManyAndReturnQuery) ([]*phi.User, error) {
        // Enforce default role across all batch entries
        defaultRole := phi.UserRole_STUDENT
        for _, record := range args.Data {
            if record.Role == nil {
                record.Role = &defaultRole
            }
        }

        // Return nested posts for all created users
        args.Select.Posts = post.Query()

        return next(ctx, args)
    },
})
```

---

## 4. Upsert & Conflict Interception

Intercept conflict resolution strategies (`OnConflict`) across single and bulk creation:

```go
db.User.Use(user.Extension{
    CreateMany: func(ctx context.Context, args *phi.UserCreateManyArgs, next phi.UserCreateManyQuery) (int64, error) {
        if args.ConflictAction != nil && args.ConflictAction.IsUpdateNewValues() {
            // Override conflict action with a custom update builder
            args.ConflictAction = user.ConflictUpdate(func(u *phi.UserUpsert) {
                u.Role.Set(phi.UserRole_STUDENT)
                u.LoginCount.Increment(1)
            })
        }
        return next(ctx, args)
    },
})
```
