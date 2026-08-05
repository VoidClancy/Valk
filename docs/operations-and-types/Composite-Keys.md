---
title: Composite Keys
description: Working with multi-column primary keys (@@id) and compound unique constraints (@@unique).
category: Schema & Types
tags: [composite, primary-key, unique, predicates, multi-column]
---

# Composite Keys

Phi provides full support for multi-column primary keys (`@@id`) and compound unique constraints (`@@unique`). Composite keys are exposed as first-class generated objects with multi-argument `.EQ(...)` predicates, full CRUD support, `.Children()` inspection methods, and seamless integration with extension hooks and upserts.

---

## Schema Definitions

### 1. Composite Unique Constraint (`@@unique`)
```prisma
model User {
  id       String @id @default(cuid())
  email    String
  phoneNum String

  @@unique([email, phoneNum])
}
```
* **Generated Target**: `user.EmailPhone`

### 2. Composite Primary Key (`@@id`)
```prisma
model CategoryToPost {
  postId     String
  categoryId Int

  @@id([postId, categoryId])
}
```
* **Generated Target**: `categoryToPost.PostId_CategoryId`

---

## Access Pattern & Predicates (`EQ`)

For composite keys, Phi generates a composite helper struct in the model package (e.g. `user.EmailPhone` or `categoryToPost.PostId_CategoryId`).

### Multi-Argument `EQ(...)` Predicate
Call `.EQ(...)` passing values in the exact order specified in the Prisma schema:

```go
// Matches (email = 'a@b.com' AND phoneNum = '+1000')
p1 := user.EmailPhone.EQ("a@b.com", "+1000")

// Matches (postId = 'post-1' AND categoryId = 42)
p2 := categoryToPost.PostId_CategoryId.EQ("post-1", 42)
```

### Predicate Inspection & `.Children()` Method
When inspecting a composite predicate:
* **`.Column()`**: Returns the logical composite column name (e.g. `"emailPhone"` or `"PostId_CategoryId"`).
* **`.Children()`**: Returns a slice of `phi.ChildPredicate` containing constituent child columns and values (`[]ChildPredicate{{Column: "email", Value: "a@b.com"}, ...}}`).
* **`.Value()`**: Returns a `map[string]any` holding constituent column keys and values.

```go
p := user.EmailPhone.EQ("a@b.com", "+1000")

fmt.Println(p.Column()) // "emailPhone"

// Inspect constituent columns using Children()
for _, child := range p.Children() {
    fmt.Printf("%s = %v\n", child.Column, child.Value)
    // email = a@b.com
    // phoneNum = +1000
}
```

---

## CRUD Operations with Composite Keys

### 1. `FindUnique` by Composite Key

Fetch a single record targeting a compound unique or composite primary key:

```go
// Unique lookup via composite @@unique
u, err := db.User.FindUnique(
    user.EmailPhone.EQ("user@example.com", "+1000"),
).Exec(ctx)

// Unique lookup via composite @@id
ctp, err := db.CategoryToPost.FindUnique(
    categoryToPost.PostId_CategoryId.EQ("post-10", 42),
).Exec(ctx)
```

### 2. Combining Additional Predicates
You can pass additional scalar predicates alongside a composite key in `FindUnique`:

```go
u, err := db.User.FindUnique(
    user.EmailPhone.EQ("user@example.com", "+1000"),
    user.Role.EQ(phi.UserRole_STUDENT),
).Exec(ctx)
```

### 3. Updating Records by Composite Key

Update records matching a composite key:

```go
u, err := db.User.Update(
    user.EmailPhone.EQ("user@example.com", "+1000"),
).SetPassword("new-secret").Exec(ctx)
```

### 4. Deleting Records by Composite Key

Delete join table or composite records:

```go
deleted, err := db.CategoryToPost.Delete(
    categoryToPost.PostId_CategoryId.EQ("post-10", 42),
).Exec(ctx)
```

---

## Conflict Resolution (`OnConflict`) with Composite Keys

Use composite keys directly as `OnConflict` targets for upsert queries:

```go
// Ignore duplicate composite inserts
affected, err := db.CategoryToPost.CreateMany(
    db.CategoryToPost.Create().SetPostId("post-1").SetCategoryId(42),
).OnConflict(categoryToPost.PostId_CategoryId).Ignore().Exec(ctx)

// Update existing matching composite record
affected, err := db.User.CreateMany(
    db.User.Create().SetEmail("user@example.com").SetPhoneNum("+1000").SetPassword("updated"),
).OnConflict(user.EmailPhone).UpdateNewValues().Exec(ctx)
```

---

## Extension Hooks & Composite Keys

In extension hooks, a composite predicate appears as a single logical column in `args.Where`. Use `.Children()` to iterate over its constituent fields without type casting:

```go
db.User.Use(user.Extension{
    FindUnique: func(ctx context.Context, args *phi.UserFindUniqueArgs, next phi.UserFindUniqueQuery) (*phi.User, error) {
        for _, w := range args.Where {
            switch w.Column() {
            case user.Email.Column:
                // Scalar unique predicate
            case user.EmailPhone.Column:
                // Iterate constituent fields with .Children()
                for _, child := range w.Children() {
                    fmt.Printf("Composite field: %s = %v\n", child.Column, child.Value)
                }
            }
        }
        return next(ctx, args)
    },
})
```
