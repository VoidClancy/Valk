---
title: Error Handling
description: Sentinel errors, driver error translation, constraint inspection, and client-side validation rules in Phi.
category: Errors
tags: [errors, validation, constraints, translation, postgresql, sqlite, sentinel]
---

# Error Handling & Validation

Phi features a robust, two-tiered error architecture:

1. **Normalized Domain & Database Errors**: Standardized sentinel errors and inspection helpers that mask vendor differences across PostgreSQL and SQLite.
2. **Client-Side Validation**: Pre-query in-memory checks that catch malformed strings, invalid UUIDs, out-of-bounds numbers, and bad JSON before contacting the database.

---

## 1. Domain Sentinel Errors & Inspectors

Phi normalizes database errors into domain sentinel errors. Inspect them using `errors.Is(err, target)` or Ent-style inspector functions:

| Sentinel Error             | Ent-Style Inspector            | Description                                                     |
| :------------------------- | :----------------------------- | :-------------------------------------------------------------- |
| `phi.ErrNotFound`          | `phi.IsNotFound(err)`          | Query expected a record but found none (`sql.ErrNoRows`).       |
| `phi.ErrNoRowsAffected`    | `phi.IsNoRowsAffected(err)`    | Update or delete operation affected 0 rows.                     |
| `phi.ErrConstraint`        | `phi.IsConstraintError(err)`   | Base error for any database constraint violation.               |
| `phi.ErrUniqueConstraint`  | `phi.IsUniqueConstraint(err)`  | Unique index or primary key collision.                          |
| `phi.ErrFKConstraint`      | `phi.IsFKConstraint(err)`      | Foreign key reference violation.                                |
| `phi.ErrNotNullConstraint` | `phi.IsNotNullConstraint(err)` | NOT NULL column violation.                                      |
| `phi.ErrCheckConstraint`   | `phi.IsCheckConstraint(err)`   | SQL `CHECK` constraint failure.                                 |
| `phi.ErrDeadlock`          | `phi.IsDeadlock(err)`          | Database deadlock detected.                                     |
| `phi.ErrLockTimeout`       | `phi.IsLockTimeout(err)`       | Lock wait timeout exceeded.                                     |
| `phi.ErrSerialization`     | `phi.IsSerialization(err)`     | Transaction serialization / concurrent update conflict.         |
| `phi.ErrTxDone`            | `phi.IsTxDone(err)`            | Operation attempted on completed transaction (`sql.ErrTxDone`). |
| `phi.ErrConnClosed`        | `phi.IsConnClosed(err)`        | Database connection closed (`sql.ErrConnDone`).                 |

### Usage Example

```go
u, err := db.User.FindUnique(user.Email.EQ("user@example.com")).Exec(ctx)
if phi.IsNotFound(err) {
    // Return HTTP 404 Not Found
    return
}

if phi.IsUniqueConstraint(err) {
    // Return HTTP 409 Conflict
    return
}
```

---

## 2. Concrete Error Structs

For advanced debugging, Phi wraps errors in detailed concrete structs:

### `*phi.NotFoundError`

Carries the target model name alongside the underlying cause:

```go
var nf *phi.NotFoundError
if errors.As(err, &nf) {
    fmt.Printf("Model %s was not found: %v\n", nf.Model, nf.Cause)
}
```

### `*phi.ConstraintError`

Carries the constraint kind, table name, constraint identifier, and driver cause:

```go
var ce *phi.ConstraintError
if errors.As(err, &ce) {
    fmt.Printf("Constraint %q violated on table %q (Kind: %v)\n", ce.Constraint, ce.Table, ce.Kind)
}
```

---

## 3. Driver Translation Engine (`TranslateDBError`)

Phi normalizes database driver errors via `TranslateDBError(err)`:

### Idempotency

If `err` is already a normalized Phi domain error, `TranslateDBError` returns it unchanged to prevent redundant error wrapping.

### Dialect Mapping Rules

#### Standard Library (`database/sql`)

- `sql.ErrNoRows` -> `phi.ErrNotFound`
- `sql.ErrTxDone` -> `phi.ErrTxDone`
- `sql.ErrConnDone` -> `phi.ErrConnClosed`

#### PostgreSQL (`*pq.Error` & `SQLState`)

- Code `23505` -> `ErrUniqueConstraint`
- Code `23503` -> `ErrFKConstraint`
- Code `23502` -> `ErrNotNullConstraint`
- Code `23514` -> `ErrCheckConstraint`
- Code `40001` -> `ErrSerialization`
- Code `40P01` -> `ErrDeadlock`

#### SQLite (`ExtendedCode` & `Code`)

- Extended Codes `2067`, `1555` -> `ErrUniqueConstraint`
- Extended Code `787` -> `ErrFKConstraint`
- Extended Code `1299` -> `ErrNotNullConstraint`
- Extended Code `275` -> `ErrCheckConstraint`

---

## 4. Client-Side Validation (`ValidationError`)

Phi runs in-memory validation on query inputs **before** making database roundtrips.

### Inspecting Validation Errors

Use `phi.IsValidationError(err)` or `errors.As`:

```go
_, err := db.User.Create().SetEmail("invalid\x00user").Exec(ctx)
if phi.IsValidationError(err) {
    var ve *phi.ValidationError
    if errors.As(err, &ve) {
        for _, fe := range ve.Errors {
            fmt.Printf("Field: %s, Rule: %s, Message: %s\n", fe.Field, fe.Rule, fe.Msg)
        }
    }
}
```

### Built-in Validation Rules

| Field Type        | Rule     | Validation Description                                                                               |
| :---------------- | :------- | :--------------------------------------------------------------------------------------------------- |
| **String**        | `safety` | Rejects strings containing null bytes (`\x00`) or invalid UTF-8 sequences.                           |
| **String**        | `length` | Enforces `@db.VarChar(n)` maximum rune limits.                                                       |
| **UUID**          | `format` | Validates standard 36-character UUID regex syntax.                                                   |
| **Decimal**       | `format` | Validates numeric string syntax and enforces scale precision.                                        |
| **Float**         | `range`  | Rejects `NaN` and infinite (`Inf`) values.                                                           |
| **JSON**          | `format` | Enforces `json.Valid(val)` syntax on `json.RawMessage`.                                              |
| **Bit String**    | `format` | Enforces string contains only `'0'` and `'1'`.                                                       |
| **Inet / CIDR**   | `format` | Enforces valid IP or CIDR address syntax via Go `net` package.                                       |
| **Integer Types** | `range`  | Enforces range bounds (`SmallInt`: $-32768$ to $32767$, `TinyInt`: $-128$ to $127$, `Oid`: $\ge 0$). |
