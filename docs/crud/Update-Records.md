---
title: Update Records
description: Update one or more records using type-safe update builders.
category: CRUD
tags: [update, updateMany, updateManyAndReturn]
---

# Update Records

Phi provides three update methods:

- `Update`
- `UpdateMany`
- `UpdateManyAndReturn`

All three methods share the same predicate system and assignment API.

---

## Update

Updates a **single record**.

The first predicate **must** uniquely identify a record (for example, `Id.EQ()` or a unique field such as `Email.EQ()`). Additional predicates may be supplied to further constrain the update.

### Basic

```go
user, err := db.User.Update(user.Email.EQ("x@y.com")).
	SetEmail("x@y.com").
	SetPassword("secret").
	SetBio("super cool bio").
	Exec(ctx)
```

### With Additional Predicates

```go
user, err := db.User.Update(
	user.Email.EQ("x@y.com"),
	user.Bio.Contains("the"),
).
	SetPassword("secret").
	SetBio("super cool bio").
	Exec(ctx)
```

### Using Assignments

```go
user, err := db.User.Update(
	user.Email.EQ("x@y.com"),
	user.Or(
		user.Bio.Contains("the"),
		user.PhoneNum.HasPrefix("+1"),
	),
).
	Assignments(
		user.Password.Set("secret"),
		user.Bio.Set("super cool bio"),
	).
	Exec(ctx)
```

### Supported Builder Methods

| Method                        | Description                                         |
| ----------------------------- | --------------------------------------------------- |
| [`Select`](/docs/Select)      | Return only the selected fields or relations.       |
| [`Omit`](/docs/Omit)          | Return all scalar fields except the omitted ones.   |
| `Set{Field}`                  | Update a single field.                              |
| [`Assignments`](#assignments) | Update multiple fields using generated assignments. |

---

## UpdateMany

Updates **all records** matching the supplied predicates.

Unlike `Update`, no unique predicate is required.

### Basic

```go
count, err := db.User.UpdateMany(
	user.Bio.Contains("golang"),
).
	SetBio("Super cool bio").
	Exec(ctx)
```

### Multiple Predicates

```go
count, err := db.User.UpdateMany(
	user.Email.HasSuffix("@example.com"),
	user.Bio.Contains("the"),
).
	SetPassword("secret").
	SetBio("super cool bio").
	Exec(ctx)
```

### Using Assignments

```go
count, err := db.User.UpdateMany(
	user.Or(
		user.Bio.Contains("the"),
		user.PhoneNum.HasPrefix("+1"),
	),
).
	Assignments(
		user.Password.Set("secret"),
		user.Bio.Set("super cool bio"),
	).
	Exec(ctx)
```

### Supported Builder Methods

| Method                        | Description                                         |
| ----------------------------- | --------------------------------------------------- |
| `Set{Field}`                  | Update a single field.                              |
| [`Assignments`](#assignments) | Update multiple fields using generated assignments. |

---

## UpdateManyAndReturn

Updates **all matching records** and returns them.

By default, all scalar fields are returned. You can customize the returned data with `Select` or `Omit`.

### Basic

```go
users, err := db.User.UpdateManyAndReturn(
	user.Bio.Contains("golang"),
).
	SetBio("Super cool bio").
	Exec(ctx)
```

### Multiple Predicates

```go
users, err := db.User.UpdateManyAndReturn(
	user.Email.HasSuffix("@example.com"),
	user.Bio.Contains("the"),
).
	SetPassword("secret").
	SetBio("super cool bio").
	Exec(ctx)
```

### Using Assignments

```go
users, err := db.User.UpdateManyAndReturn(
	user.Or(
		user.Bio.Contains("the"),
		user.PhoneNum.HasPrefix("+1"),
	),
).
	Assignments(
		user.Password.Set("secret"),
		user.Bio.Set("super cool bio"),
	).
	Exec(ctx)
```

### Supported Builder Methods

| Method                        | Description                                         |
| ----------------------------- | --------------------------------------------------- |
| [`Select`](/docs/Select)      | Return only the selected fields or relations.       |
| [`Omit`](/docs/Omit)          | Return all scalar fields except the omitted ones.   |
| `Set{Field}`                  | Update a single field.                              |
| [`Assignments`](#assignments) | Update multiple fields using generated assignments. |

---

## Assignments

`Assignments` accepts a variadic list of generated field assignments.

For most updates, chaining `Set{Field}` methods is the most readable approach:

```go
user, err := db.User.Update(
    user.Email.EQ("x@y.com"),
).
    SetUsername("JohnDoe.sh").
    SetEmail("new@example.com").
    SetBio("super cool bio").
    Exec(ctx)
```

`Assignments` becomes especially useful when the fields being updated are determined at runtime. Because assignments are ordinary values, they can be built, reused, filtered, and passed between functions before executing the query.

```go
assignments := []phi.FieldAssignmentOf[phi.User]{}

if req.Username != nil {
    assignments = append(assignments, user.Username.Set(*req.Username))
}

if req.Bio != nil {
    assignments = append(assignments, user.Bio.Set(*req.Bio))
}

updatedUser, err := db.User.Update(
    user.Id.EQ(id),
).
    Assignments(assignments...).
    Exec(ctx)
```

This is difficult to achieve with chained `Set{Field}` methods, since the chain must be known at compile time.

> **Note:** `Assignments` only accepts assignments generated for the same model. Passing `post.Title.Set(...)` to `db.User.Update(...).Assignments(...)` results in a compile-time type error.
