---
title: Update Hooks
description: Intercept and mutate record updates with extension hooks.
category: Hooks
tags: [hooks, extension, update, updateMany, updateManyAndReturn, predicate, middleware]
---

# Update Hooks

Update hooks allow you to intercept `Update`, `UpdateMany`, and `UpdateManyAndReturn` operations before queries hit the database. You can validate or hash input fields, auto-stamp audit timestamps, mutate target predicates, or pre-load relations on returned records.

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
    Update: func(ctx context.Context, args *phi.UserUpdateArgs, next phi.UserUpdateQuery) (*phi.User, error) {
        // Pre-hook: normalize email if being updated
        if args.Data.Email != nil {
            lower := strings.ToLower(*args.Data.Email)
            args.Data.Email = &lower
        }

        res, err := next(ctx, args)

        // Post-hook: runs after query execution
        return res, err
    },
})
```

---

## Hook Signatures

| Hook Field | Signature | Return Type |
| :--- | :--- | :--- |
| `Update` | `func(ctx context.Context, args *phi.UserUpdateArgs, next phi.UserUpdateQuery)` | `(*phi.User, error)` |
| `UpdateMany` | `func(ctx context.Context, args *phi.UserUpdateManyArgs, next phi.UserUpdateManyQuery)` | `(int64, error)` |
| `UpdateManyAndReturn` | `func(ctx context.Context, args *phi.UserUpdateManyAndReturnArgs, next phi.UserUpdateManyAndReturnQuery)` | `([]*phi.User, error)` |

---

## 1. Mutating Update Data (`args.Data`)

All update structs use **optional pointer fields** (`*string`, `*int32`, etc.).
* `nil`: The field is **unmodified** by the query.
* `non-nil`: The field **will be updated** to the dereferenced value.

```go
db.User.Use(user.Extension{
    Update: func(ctx context.Context, args *phi.UserUpdateArgs, next phi.UserUpdateQuery) (*phi.User, error) {
        // Lowercase email if updated
        if args.Data.Email != nil {
            lower := strings.ToLower(*args.Data.Email)
            args.Data.Email = &lower
        }

        // Hash password if updated
        if args.Data.Password != nil {
            hashed := hashPassword(*args.Data.Password)
            args.Data.Password = &hashed
        }

        return next(ctx, args)
    },
})
```

---

## 2. Targeting & Where Clauses (`Where`)

Inspect or append predicates to restrict which records can be updated:

```go
db.User.Use(user.Extension{
    UpdateMany: func(ctx context.Context, args *phi.UserUpdateManyArgs, next phi.UserUpdateManyQuery) (int64, error) {
        // Restrict bulk updates to STUDENT accounts only
        args.Where = append(args.Where, user.Role.EQ(phi.UserRole_STUDENT))
        return next(ctx, args)
    },
})
```

### Setters for Retargeting
`SetWhere` replaces the active target filter:

```go
db.User.Use(user.Extension{
    Update: func(ctx context.Context, args *phi.UserUpdateArgs, next phi.UserUpdateQuery) (*phi.User, error) {
        args.SetWhere(user.Email.EQ("target@example.com"))
        return next(ctx, args)
    },
})
```

---

## 3. Relation Selection on Returned Records

`Update` and `UpdateManyAndReturn` support `args.Select` for preloading relations on updated records:

```go
db.User.Use(user.Extension{
    UpdateManyAndReturn: func(ctx context.Context, args *phi.UserUpdateManyAndReturnArgs, next phi.UserUpdateManyAndReturnQuery) ([]*phi.User, error) {
        // Return updated users along with their published posts
        args.Select.Posts = post.Query().Where(post.Published.EQ(true))
        return next(ctx, args)
    },
})
```
