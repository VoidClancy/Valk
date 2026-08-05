---
title: Delete Hooks
description: Intercept and mutate record deletion with extension hooks.
category: Hooks
tags: [hooks, extension, delete, deleteMany, predicate, middleware]
---

# Delete Hooks

Delete hooks allow you to intercept `Delete` and `DeleteMany` operations before queries hit the database. You can enforce safety rules (preventing root user or mass table deletion), restrict deletion targets, or return preloaded relations on deleted records.

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
    Delete: func(ctx context.Context, args *phi.UserDeleteArgs, next phi.UserDeleteQuery) (*phi.User, error) {
        // Pre-hook: prevent deletion of protected user
        for _, p := range args.Where {
            if p.Column() == user.Email.Column && p.Value() == "admin@example.com" {
                return nil, errors.New("cannot delete primary admin user")
            }
        }

        res, err := next(ctx, args)

        // Post-hook: runs after successful deletion
        return res, err
    },
})
```

---

## Hook Signatures

| Hook Field | Signature | Return Type |
| :--- | :--- | :--- |
| `Delete` | `func(ctx context.Context, args *phi.UserDeleteArgs, next phi.UserDeleteQuery)` | `(*phi.User, error)` |
| `DeleteMany` | `func(ctx context.Context, args *phi.UserDeleteManyArgs, next phi.UserDeleteManyQuery)` | `(int64, error)` |

---

## 1. Safety Guards & Aborting Deletion

To abort a deletion operation, return an error before invoking `next(ctx, args)`. The database is never queried:

### Guarding Against Mass Deletion (`DeleteMany`)
```go
db.User.Use(user.Extension{
    DeleteMany: func(ctx context.Context, args *phi.UserDeleteManyArgs, next phi.UserDeleteManyQuery) (int64, error) {
        // Refuse empty WHERE clause (prevent truncating entire table)
        if len(args.Where) == 0 {
            return 0, errors.New("unbounded bulk delete is rejected for safety")
        }
        return next(ctx, args)
    },
})
```

---

## 2. Restricting Target Predicates (`Where`)

Inspect `args.Where` or append filters to enforce scoping (e.g., tenant boundaries):

```go
db.User.Use(user.Extension{
    DeleteMany: func(ctx context.Context, args *phi.UserDeleteManyArgs, next phi.UserDeleteManyQuery) (int64, error) {
        // Only allow deleting STUDENT accounts
        args.Where = append(args.Where, user.Role.EQ(phi.UserRole_STUDENT))
        return next(ctx, args)
    },
})
```

---

## 3. Preloading Relations on Deleted Records (`Delete`)

`Delete` carries `args.Select`, allowing you to retrieve selected scalar fields and preloaded relations of the record being deleted:

```go
db.User.Use(user.Extension{
    Delete: func(ctx context.Context, args *phi.UserDeleteArgs, next phi.UserDeleteQuery) (*phi.User, error) {
        // Return deleted user along with their posts
        args.Select.Email = true
        args.Select.Posts = post.Query()

        return next(ctx, args)
    },
})
```
