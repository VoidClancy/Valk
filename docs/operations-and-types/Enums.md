---
title: Enums
description: Define, query, migrate, and enforce type-safe enum values across PostgreSQL and SQLite.
category: Schema & Types
tags: [enum, postgresql, sqlite, migrations, ddl, type-safety, validation]
---

# Enums

Phi generates strongly-typed Go constants, validation routines, and dialect-specific database DDL for Prisma `enum` definitions. Enums provide compile-time type safety across your client builders, input structs, and query predicates while mapping cleanly to PostgreSQL and SQLite migrations.

---

## Defining Enums in Prisma Schema

Define an enum in your `schema.prisma`:

```prisma
enum UserRole {
  ADMIN
  STUDENT
  TEACHER
}

model User {
  id    String   @id @default(cuid())
  email String   @unique
  role  UserRole @default(STUDENT)
}
```

---

## Dialect Storage & DDL Generation

Phi's migration engine (`phi migrate`) tailors DDL generation based on the database provider:

### 1. PostgreSQL (Native DDL Enums)

PostgreSQL natively supports custom enum types. During migration generation, Phi generates a dedicated `CREATE TYPE` statement followed by an enum-typed column:

#### Generated DDL:
```sql
CREATE TYPE "UserRole" AS ENUM (
  'ADMIN',
  'STUDENT',
  'TEACHER'
);

CREATE TABLE "User" (
  "id" TEXT NOT NULL,
  "email" TEXT NOT NULL,
  "role" "UserRole" NOT NULL DEFAULT 'STUDENT',
  PRIMARY KEY ("id")
);
```

* **Storage**: Stored using PostgreSQL's internal 4-byte enum OID storage engine.
* **Enforcement**: Database engine strictly rejects non-enum string values at the SQL driver boundary.

---

### 2. SQLite (TEXT Storage + Inline CHECK Constraints)

SQLite does not have native enum types. Phi creates the column as `TEXT` and appends an inline SQL `CHECK` constraint:

#### Generated DDL:
```sql
CREATE TABLE "User" (
  "id" TEXT NOT NULL,
  "email" TEXT NOT NULL,
  "role" TEXT NOT NULL DEFAULT 'STUDENT',
  PRIMARY KEY ("id"),
  CONSTRAINT "user_role_check" CHECK ("role" IN ('ADMIN', 'STUDENT', 'TEACHER'))
);
```

* **Storage**: Stored as plain UTF-8 `TEXT`.
* **Enforcement**:
  1. **Database DDL**: SQL `CHECK` constraint prevents invalid strings from being inserted directly via raw SQL.
  2. **Application Level**: Phi generates `UserRoleType` Go types with `.IsValid()` checks, guaranteeing strict type safety before queries hit SQLite.

---

## Default Value Formatting

When an `@default(...)` decorator is added to an enum field in your Prisma schema:

```prisma
role UserRole @default(STUDENT)
```

Phi's migration engine formats default values for the database DDL:

| Provider | Generated Default SQL | Behavior |
| :--- | :--- | :--- |
| **PostgreSQL** | `DEFAULT 'STUDENT'` | Coerced to `"UserRole"` enum type by Postgres. |
| **SQLite** | `DEFAULT 'STUDENT'` | Inserted as string default matching `TEXT` column. |

In Go client builders (`Create()`), if an optional or defaulted enum field is left unspecified (`nil`), Phi automatically applies the schema default value during record map computation.

---

## Access Pattern & Naming Conventions

All enum types, constants, and validation methods live in the root generated package (`phi`):

### Root Client Package (`phi`)
* **Enum Type**: `type UserRoleType string`
* **Enum Values**: `phi.UserRole_ADMIN`, `phi.UserRole_STUDENT`, `phi.UserRole_TEACHER`
* **Validation Method**: `(e UserRoleType).IsValid() bool`

```go
role := phi.UserRole_ADMIN
if role.IsValid() {
    fmt.Println("Role is valid:", role)
}
```

### Model Package (`user`)
The model package (e.g. `user`) exports predicate builders typed against `phi.UserRoleType`:
* **Field Predicates**: `user.Role` (type `phi.Field[phi.User, phi.UserRoleType]`)

---

## Usage Examples

### 1. Inserting Records (`Create` / `CreateMany`)

Pass root enum constants directly to builder `Set*` methods or struct fields:

```go
u, err := db.User.Create().
    SetEmail("student@example.com").
    SetPhoneNum("+1001").
    SetRole(phi.UserRole_STUDENT).
    Exec(ctx)

u2, err := db.User.Create().
    SetEmail("admin@example.com").
    SetPhoneNum("+1002").
    SetRole(phi.UserRole_ADMIN).
    Exec(ctx)
```

### 2. Querying with Enum Predicates

Filter records by enum values using `EQ()`, `NEQ()`, `In()`, and `NotIn()`:

```go
// Find all Teachers or Admins
users, err := db.User.FindMany(
    user.Or(
        user.Role.EQ(phi.UserRole_TEACHER),
        user.Role.EQ(phi.UserRole_ADMIN),
    ),
).Exec(ctx)

// Using In operator
students, err := db.User.FindMany(
    user.Role.In([]phi.UserRoleType{phi.UserRole_STUDENT}),
).Exec(ctx)
```

### 3. Updating Enum Fields

Update enum values on existing records:

```go
updatedUser, err := db.User.Update(user.Id.EQ("user-123")).
    SetRole(phi.UserRole_ADMIN).
    Exec(ctx)
```

### 4. Conflict Updates in Upserts

Set or update enum values during `OnConflict` upsert execution:

```go
affected, err := db.User.CreateMany(
    db.User.Create().SetEmail("user@example.com").SetPhoneNum("+1000").SetRole(phi.UserRole_STUDENT),
).OnConflict(user.Email).Update(func(u *phi.UserUpsert) {
    u.Role.Set(phi.UserRole_ADMIN)
}).Exec(ctx)
```
