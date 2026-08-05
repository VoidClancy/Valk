---
title: Query Filters
description: The predicate system for filtering records with EQ, NEQ, GT, LT, In, Contains, Like, and logical operators.
category: Query Filters
tags: [query, filters, eq, neq, gt, lt, in, contains, like, and, or, not, predicates]
---

# Query Filters

Every read, update, and delete operation accepts zero or more _predicates_ that narrow down which records are affected. Predicates are generated per model and live in each model's package (for example `user.Email.EQ(...)`).

The available operators depend on the field's scalar type:

- Every scalar field exposes the comparison and set operators.
- String fields additionally expose the string-matching operators.
- Optional fields expose null checks via `IsNull` and `IsNotNull`.
- Array fields expose `Has`, `HasEvery`, and `HasSome`.

The same predicate types are passed to `FindUnique`, `FindFirst`, `FindMany`, `Update`, `UpdateMany`, `Delete`, `DeleteMany`, and `Count`.

Use the following Prisma schema as reference:

```prisma
enum UserRole {
  ADMIN
  STUDENT
  TEACHER
}

model User {
  id        String    @id @default(cuid())
  email     String    @unique
  phoneNum  String?
  loginCount Int
  bio       String?
  role      UserRole?
  tags      String[]
}
```

## EQ

Matches records where a field equals a value. On unique fields, `EQ` returns a `UniquePredicate` and can also be used as a cursor.

```go
users, err := db.User.FindMany(user.Email.EQ("x@y.com")).Exec(ctx)
```

```sql
SELECT * FROM "users" WHERE "email" = 'x@y.com';
```

## NEQ

Matches records where a field is not equal to a value.

```go
users, err := db.User.FindMany(user.Bio.NEQ("spam")).Exec(ctx)
```

```sql
SELECT * FROM "users" WHERE "bio" != 'spam';
```

## GT / GTE

Matches records where a field is greater than - or greater than or equal to - a value.

```go
users, err := db.User.FindMany(user.LoginCount.GT(10)).Exec(ctx)
users, err := db.User.FindMany(user.LoginCount.GTE(10)).Exec(ctx)
```

```sql
SELECT * FROM "users" WHERE "loginCount" > 10;
SELECT * FROM "users" WHERE "loginCount" >= 10;
```

## LT / LTE

Matches records where a field is less than - or less than or equal to - a value.

```go
users, err := db.User.FindMany(user.LoginCount.LT(10)).Exec(ctx)
users, err := db.User.FindMany(user.LoginCount.LTE(10)).Exec(ctx)
```

```sql
SELECT * FROM "users" WHERE "loginCount" < 10;
SELECT * FROM "users" WHERE "loginCount" <= 10;
```

## In / NotIn

Matches records where a field is (or is not) one of a list of values.

```go
users, err := db.User.FindMany(
    user.Role.In([]phi.UserRoleType{phi.UserRole_ADMIN, phi.UserRole_TEACHER}),
).Exec(ctx)
users, err := db.User.FindMany(user.Email.NotIn([]string{"a@b.com", "c@d.com"})).Exec(ctx)
```

```sql
SELECT * FROM "users" WHERE "role" IN ('ADMIN', 'TEACHER');
SELECT * FROM "users" WHERE "email" NOT IN ('a@b.com', 'c@d.com');
```

## Between

Matches records where a field falls within an inclusive range.

```go
users, err := db.User.FindMany(user.LoginCount.Between(5, 100)).Exec(ctx)
```

```sql
SELECT * FROM "users" WHERE "loginCount" BETWEEN 5 AND 100;
```

## Contains

Matches string fields that contain the given substring. This maps to SQL's `LIKE '%value%'`.

```go
users, err := db.User.FindMany(user.Bio.Contains("golang")).Exec(ctx)
```

```sql
SELECT * FROM "users" WHERE "bio" LIKE '%golang%';
```

## HasPrefix / HasSuffix

Matches string fields that begin or end with a given substring.

```go
users, err := db.User.FindMany(user.Email.HasPrefix("admin")).Exec(ctx)
users, err := db.User.FindMany(user.Email.HasSuffix("@example.com")).Exec(ctx)
```

```sql
SELECT * FROM "users" WHERE "email" LIKE 'admin%';
SELECT * FROM "users" WHERE "email" LIKE '%@example.com';
```

## Like / ILike

Provides raw `LIKE` pattern matching. Use `%` as a wildcard. `ILike` performs a case-insensitive match on the underlying database.

```go
users, err := db.User.FindMany(user.Email.Like("%@example.com")).Exec(ctx)
users, err := db.User.FindMany(user.Bio.ILike("%golang%")).Exec(ctx)
```

```sql
SELECT * FROM "users" WHERE "email" LIKE '%@example.com';
-- Case-insensitive, e.g. uses "email" ILIKE '%golang%' on PostgreSQL
SELECT * FROM "users" WHERE "bio" ILIKE '%golang%';
```

> **Note:** Dialects vary. PostgreSQL supports a native `ILIKE`, so `ILike` maps directly to `col ILIKE val`. SQLite has no `ILIKE` operator, so Phi emulates it with `LOWER(col) LIKE LOWER(val)` to get the same case-insensitive behavior.

## IsNull / IsNotNull

Match records where an optional field is null or not null.

```go
users, err := db.User.FindMany(user.Bio.IsNull()).Exec(ctx)
users, err := db.User.FindMany(user.Bio.IsNotNull()).Exec(ctx)
```

```sql
SELECT * FROM "users" WHERE "bio" IS NULL;
SELECT * FROM "users" WHERE "bio" IS NOT NULL;
```

## And

Combines predicates so that all of them must match; it compiles to the same `AND` in SQL.

```go
users, err := db.User.FindMany(
    user.And(
        user.LoginCount.GTE(10),
        user.Email.HasSuffix("@example.com"),
    ),
).Exec(ctx)
```

```sql
SELECT * FROM "users" WHERE "loginCount" >= 10 AND "email" LIKE '%@example.com';
```

## Multiple Predicates Are AND-ed

Query methods accept zero or more predicates, so you can usually pass conditions directly without `And`. Every predicate passed to a method is combined with `AND` in the generated SQL:

```go
users, err := db.User.FindMany(
    user.LoginCount.GTE(10),
    user.Email.HasSuffix("@example.com"),
).Exec(ctx)
```

```sql
SELECT * FROM "users" WHERE "loginCount" >= 10 AND "email" LIKE '%@example.com';
```

The two produce identical SQL. Choose based on what you are doing:

- **Method arguments** (`FindMany(p1, p2)`) - the simplest, most readable way to flatten a set of `AND` conditions. Use this for everyday filtering.
- **`user.And(p1, p2, ...)`** - bundles several predicates into a single value you can store, reuse, or pass as one argument. Use this when a condition needs to be shared across queries, built at runtime, passed into a helper, or embedded inside a larger boolean expression such as `user.Or(user.And(a, b), user.And(c, d))`. Because it returns one predicate, you can nest it inside `And`, `Or`, or `Not`, which require individual predicate values.

## Or

Combines predicates so that any one of them may match.

```go
users, err := db.User.FindMany(
    user.Or(
        user.Role.EQ(phi.UserRole_ADMIN),
        user.LoginCount.GTE(100),
    ),
).Exec(ctx)
```

```sql
SELECT * FROM "users" WHERE "role" = 'ADMIN' OR "loginCount" >= 100;
```

## Not

Negates a single predicate.

```go
users, err := db.User.FindMany(
    user.Not(user.Role.EQ(phi.UserRole_ADMIN)),
).Exec(ctx)
```

```sql
SELECT * FROM "users" WHERE NOT ("role" = 'ADMIN');
```

## Array Fields

Array fields (such as `tags String[]`) expose containment operators:

- `Has(value)` - the array contains the value.
- `HasEvery(values)` - the array contains every value.
- `HasSome(values)` - the array contains at least one value.

```go
users, err := db.User.FindMany(user.Tags.HasEvery([]string{"go", "grpc"})).Exec(ctx)
```

> **Note:** Predicates are compiled after validation. Any predicate that fails validation produces an error instead of running an invalid query.
