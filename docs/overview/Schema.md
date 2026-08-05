---
title: Schema Walkthrough
description: Prisma schema definitions, supported scalar/native types, and dialect mapping rules in Phi.
category: Schema & Types
tags: [schema, prisma, types, postgresql, sqlite, native-types, element-type]
---

# Schema Walkthrough & Types

Phi uses Prisma Schema language (`schema.prisma`) to define data models, relationships, field constraints, and database providers.

---

## 1. Schema Definition Walkthrough

A minimal `schema.prisma` file consists of a **datasource** block and one or more **models**:

```prisma
datasource db {
  provider = "postgres"
}

model User {
  id         String    @id @default(cuid())
  email      String    @unique
  phoneNum   String    @unique
  role       UserRole  @default(STUDENT)
  loginCount Int       @default(0)
  createdAt  DateTime  @default(now())
  updatedAt  DateTime  @updatedAt

  posts      Post[]

  @@unique([email, phoneNum])
}

model Post {
  id        String   @id @default(cuid())
  title     String
  content   String
  published Boolean  @default(false)
  authorId  String
  author    User     @relation(fields: [authorId], references: [id], onDelete: Cascade)
}

enum UserRole {
  ADMIN
  STUDENT
  TEACHER
}
```

---

## 2. Supported Scalar Types & Go Mappings

| Prisma Scalar Type | Go Client Type     | Struct Field (Required) | Struct Field (Optional `?`) |
| :----------------- | :----------------- | :---------------------- | :-------------------------- |
| `String`           | `string`           | `string`                | `*string`                   |
| `Boolean`          | `bool`             | `bool`                  | `*bool`                     |
| `Int`              | `int32`            | `int32`                 | `*int32`                    |
| `BigInt`           | `int64`            | `int64`                 | `*int64`                    |
| `Float`            | `float64`          | `float64`               | `*float64`                  |
| `Decimal`          | `string`           | `string`                | `*string`                   |
| `DateTime`         | `time.Time`        | `time.Time`             | `*time.Time`                |
| `Json`             | `json.RawMessage`  | `json.RawMessage`       | `*json.RawMessage`          |
| `Bytes`            | `[]byte`           | `[]byte`                | `*[]byte`                   |
| `Enum`             | `phi.UserRoleType` | `phi.UserRoleType`      | `*phi.UserRoleType`         |

---

## 3. Native Attributes (`@db.*`) & Provider Mappings

Phi adheres strictly to Prisma's provider-specific feature specifications:

### PostgreSQL vs. SQLite Column Type Mapping

| Prisma Definition                         | PostgreSQL SQL DDL              | SQLite SQL DDL                      |
| :---------------------------------------- | :------------------------------ | :---------------------------------- |
| `id String @id @default(cuid())`          | `TEXT NOT NULL`                 | `TEXT NOT NULL`                     |
| `id String @id @default(uuid())`          | `TEXT` or `UUID`                | `TEXT NOT NULL`                     |
| `field String @db.VarChar(255)`           | `VARCHAR(255)`                  | Not supported by Prisma schema      |
| `field String @db.Text`                   | `TEXT`                          | Not supported by Prisma schema      |
| `field Int @id @default(autoincrement())` | `SERIAL` / `BIGSERIAL`          | `INTEGER PRIMARY KEY AUTOINCREMENT` |
| `createdAt DateTime @default(now())`      | `TIMESTAMP`                     | `TIMESTAMP`                         |
| `data Json`                               | `JSONB`                         | `TEXT` / `BLOB`                     |
| `role UserRole`                           | Native `CREATE TYPE "UserRole"` | `TEXT` + `CHECK ("role" IN (...))`  |

---

## 4. Scalar Array Fields in SQLite (`/// @elementType`)

PostgreSQL natively supports scalar array types (e.g. `String[]` -> `text[]`). Because SQLite does not natively support scalar array types, Prisma schema restricts `String[]` on SQLite providers.

To use typed scalar arrays in SQLite models, Phi supports the `/// @elementType <Type>` doc-comment decorator placed above a `Json` field defaulting to `"[]"`:

```prisma
model User {
  id   String @id @default(cuid())

  /// @elementType String
  tags Json   @default("[]")
}
```

- **Go Client Generation**: Phi generates `Tags []string` on the model struct and handles JSON serialization and deserialization transparently.
- **SQLite DDL**: Generated as a `TEXT` or `BLOB` JSON column.

---

## 5. SQL Indexes (`@@index`)

Phi supports single-column and multi-column index declarations in `schema.prisma`:

```prisma
model User {
  id        String   @id @default(cuid())
  email     String   @unique
  role      UserRole @default(STUDENT)
  createdAt DateTime @default(now())

  @@index([email])
  @@index([role, createdAt])
}
```

### Migration DDL Generation
When `phi migrate` is executed, index definitions are calculated by the migration engine and output as standard DDL statements:

```sql
CREATE INDEX "User_email_idx" ON "User" ("email");
CREATE INDEX "User_role_createdAt_idx" ON "User" ("role", "createdAt");
```

### Query Performance vs Unique Predicates
Non-unique `@@index` declarations speed up query execution (`FindMany`, `FindFirst`, `Count`) directly on the database engine. Because non-unique indexes can match multiple rows, they optimize SQL performance without generating `FindUnique` predicate handles (only `@id`, `@unique`, `@@id`, and `@@unique` generate `FindUnique` targets).
