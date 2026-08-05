---
title: Delete Records
description: Delete one or more records using type-safe delete builders.
category: CRUD
tags: [delete, deleteMany, relations, cascade, foreign-key]
---

# Delete Records

Phi provides two methods for deleting records:

- `Delete`
- `DeleteMany`

> **Note:** `Delete` always returns the deleted record, regardless of the underlying database. On databases that support `DELETE ... RETURNING`, Phi performs the deletion in a single query. On databases that don't, Phi transparently executes the operation inside a transaction by fetching the record before deleting it, ensuring consistent behavior across all supported SQL dialects.

---

## Relation Deletion Behavior (`onDelete`)

Deletion behavior for foreign key relations is controlled by your `schema.prisma` `@relation(onDelete: ...)` rules:

```prisma
model Post {
  id       String @id @default(cuid())
  authorId String
  author   User   @relation(fields: [authorId], references: [id], onDelete: Cascade)
}
```

Phi translates `@relation(onDelete: ...)` rules into database foreign key constraints in DDL:

| Prisma `onDelete` Action | SQL Foreign Key DDL Action | Database Behavior on Deleting Parent Record                                            |
| :----------------------- | :------------------------- | :------------------------------------------------------------------------------------- |
| `Cascade`                | `ON DELETE CASCADE`        | Child records referencing the parent are automatically deleted by the database engine. |
| `Restrict`               | `ON DELETE RESTRICT`       | Prevents deletion of the parent record if dependent child records exist.               |
| `NoAction`               | `ON DELETE NO ACTION`      | Prevents deletion unless constraints are satisfied within the transaction.             |
| `SetNull`                | `ON DELETE SET NULL`       | Sets foreign key columns on referencing child records to `NULL`.                       |
| `SetDefault`             | `ON DELETE SET DEFAULT`    | Resets foreign key columns on child records to their schema default values.            |

---

## Delete

Deletes a **single record**.

The first predicate **must** uniquely identify a record (for example, `Id.EQ()` or another unique field). Additional predicates may be supplied to further constrain the deletion.

By default, all scalar fields are returned. Use `Select` or `Omit` to customize the returned data.

### Basic

```go
user, err := db.User.Delete(
    user.Email.EQ("x@y.com"),
).Exec(ctx)
```

### With Additional Predicates

```go
user, err := db.User.Delete(
    user.Email.EQ("x@y.com"),
    user.Bio.Contains("golang"),
).Exec(ctx)
```

### Returning Selected Fields

```go
user, err := db.User.Delete(
    user.Id.EQ(id),
).
    Select(user.Select{
        Id:       true,
        Username: true,
    }).
    Exec(ctx)
```

### Omitting Fields

```go
user, err := db.User.Delete(
    user.Id.EQ(id),
).
    Omit(user.Omit{
        Password: true,
    }).
    Exec(ctx)
```

### Supported Builder Methods

| Method                   | Description                                       |
| :----------------------- | :------------------------------------------------ |
| [`Select`](/docs/Select) | Return only the selected fields and relations.    |
| [`Omit`](/docs/Omit)     | Return all scalar fields except the omitted ones. |

---

## DeleteMany

Deletes **all records** matching the supplied predicates.

Unlike `Delete`, no unique predicate is required.

`DeleteMany` returns the number of rows deleted.

### Basic

```go
deleted, err := db.User.DeleteMany(
    user.Bio.Contains("inactive"),
).Exec(ctx)
```

### Multiple Predicates

```go
deleted, err := db.User.DeleteMany(
    user.Email.HasSuffix("@example.com"),
    user.LoginCount.LT(5),
).Exec(ctx)
```

### Using Logical Predicates

```go
deleted, err := db.User.DeleteMany(
    user.Or(
        user.Bio.Contains("spam"),
        user.PhoneNum.HasPrefix("+999"),
    ),
).Exec(ctx)
```

### Return Value

```go
deleted, err := db.User.DeleteMany(
    user.Bio.Contains("inactive"),
).Exec(ctx)

fmt.Printf("Deleted %d users\n", deleted)
```

### Supported Builder Methods

`DeleteMany` exposes no additional builder methods beyond its predicates.

> **Note:** `DeleteMany` only returns the number of deleted rows. If you need the deleted records themselves, query them before deleting.
