---
title: Read Hooks
description: Intercept and mutate read queries with extension hooks.
category: Hooks
tags: [hooks, extension, findUnique, findFirst, findMany, count, predicate, middleware]
---

# Read Hooks

Read hooks allow you to intercept `FindUnique`, `FindFirst`, `FindMany`, and `Count` operations before queries hit the database. You can inspect or append predicates, apply soft-delete/tenant filters, modify ordering or pagination, or short-circuit queries with cached data.

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
    FindMany: func(ctx context.Context, args *phi.UserFindManyArgs, next phi.UserFindManyQuery) ([]*phi.User, error) {
        // Pre-hook: append a global tenant filter
        args.Where = append(args.Where, user.LoginCount.GTE(1))

        res, err := next(ctx, args)

        // Post-hook: process or cache results
        return res, err
    },
})
```

---

## Hook Signatures

| Hook Field | Signature | Return Type |
| :--- | :--- | :--- |
| `FindUnique` | `func(ctx context.Context, args *phi.UserFindUniqueArgs, next phi.UserFindUniqueQuery)` | `(*phi.User, error)` |
| `FindFirst` | `func(ctx context.Context, args *phi.UserFindFirstArgs, next phi.UserFindFirstQuery)` | `(*phi.User, error)` |
| `FindMany` | `func(ctx context.Context, args *phi.UserFindManyArgs, next phi.UserFindManyQuery)` | `([]*phi.User, error)` |
| `Count` | `func(ctx context.Context, args *phi.UserCountArgs, next phi.UserCountQuery)` | `(int64, error)` |

---

## 1. Inspecting & Appending Predicates (`Where`)

All read args contain `Where []phi.PredicateOf[User]`. Each predicate exposes `.Column()`, `.Value()`, and `.Children()` for type-safe inspection:

```go
db.User.Use(user.Extension{
    FindMany: func(ctx context.Context, args *phi.UserFindManyArgs, next phi.UserFindManyQuery) ([]*phi.User, error) {
        // Inspect active filters
        for _, pred := range args.Where {
            if children := pred.Children(); len(children) > 0 {
                // Composite key predicate (e.g. @@unique([email, phoneNum]))
                for _, child := range children {
                    fmt.Printf("Composite constituent: %s = %v\n", child.Column, child.Value)
                }
            } else {
                // Standard scalar predicate
                fmt.Printf("Filter applied: %s = %v\n", pred.Column(), pred.Value())
            }
        }

        // Enforce role restriction
        args.Where = append(args.Where, user.Role.EQ(phi.UserRole_STUDENT))

        return next(ctx, args)
    },
})
```

### Combining Predicates with Logical Operators (`Or` / `And`)
Group predicates using logical operators:

```go
args.Where = append(args.Where,
    user.Or(
        user.Role.EQ(phi.UserRole_ADMIN),
        user.Role.EQ(phi.UserRole_TEACHER),
    ),
)
```

---

## 2. Setters vs. Direct Mutation

Every query argument struct provides chainable `Set*` helper methods for replacing values, as well as direct exported fields for appending:

| Setter Method | Target Arguments | Replaces Field |
| :--- | :--- | :--- |
| `SetWhere(...)` | `FindUnique`, `FindFirst`, `FindMany`, `Count` | `Where` |
| `SetOrderBy(...)` | `FindFirst`, `FindMany` | `OrderBy` |
| `SetCursor(...)` | `FindFirst`, `FindMany` | `Cursor` |
| `SetSkip(n int)` | `FindFirst`, `FindMany`, `Count` | `Skip` (`*int`) |
| `SetTake(n int)` | `FindFirst`, `FindMany`, `Count` | `Take` (`*int`) |

### Example: Chainable Setters (`FindMany`)
```go
db.User.Use(user.Extension{
    FindMany: func(ctx context.Context, args *phi.UserFindManyArgs, next phi.UserFindManyQuery) ([]*phi.User, error) {
        args.SetWhere(user.LoginCount.GTE(10)).
            SetOrderBy(user.LoginCount.Desc(), user.Email.Asc()).
            SetSkip(0).
            SetTake(50)

        return next(ctx, args)
    },
})
```

---

## 3. Operations Overview

### `FindUnique`
Receives `*phi.UserFindUniqueArgs` with `Where []phi.PredicateOf[User]` and `Select *phi.UserSelect`.

```go
db.User.Use(user.Extension{
    FindUnique: func(ctx context.Context, args *phi.UserFindUniqueArgs, next phi.UserFindUniqueQuery) (*phi.User, error) {
        // Force select user posts
        args.Select.Posts = post.Query().Where(post.Published.EQ(true))

        return next(ctx, args)
    },
})
```

### `FindFirst` & `FindMany`
Support full pagination, sorting, and relation preloading:

```go
db.User.Use(user.Extension{
    FindFirst: func(ctx context.Context, args *phi.UserFindFirstArgs, next phi.UserFindFirstQuery) (*phi.User, error) {
        // Enforce order by highest login count
        args.SetOrderBy(user.LoginCount.Desc())
        return next(ctx, args)
    },
})
```

### `Count`
Receives `*phi.UserCountArgs` with `Where`, `Skip`, and `Take`:

```go
db.User.Use(user.Extension{
    Count: func(ctx context.Context, args *phi.UserCountArgs, next phi.UserCountQuery) (int64, error) {
        args.Where = append(args.Where, user.LoginCount.GT(0))
        return next(ctx, args)
    },
})
```

---

## 4. Query Short-Circuiting & Caching

Skip database calls by returning cached data directly:

```go
db.User.Use(user.Extension{
    FindUnique: func(ctx context.Context, args *phi.UserFindUniqueArgs, next phi.UserFindUniqueQuery) (*phi.User, error) {
        cacheKey := fmt.Sprintf("user:%v", args.Where[0].Value())

        if cached, found := cache.Get(cacheKey); found {
            return cached.(*phi.User), nil
        }

        res, err := next(ctx, args)
        if err == nil && res != nil {
            cache.Set(cacheKey, res)
        }
        return res, err
    },
})
```
