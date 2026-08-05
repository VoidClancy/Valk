---
title: Raw SQL
description: Accessing standard database/sql handles for custom queries and vendor commands.
category: Operations
tags: [raw-sql, sql-db, sql-tx, fallback, database-handle]
---

# Raw SQL Fallback

When you need to execute complex custom SQL, vendor-specific functions, or raw database commands, Phi provides direct access to Go's standard `database/sql` driver handles.

---

## 1. Accessing `*sql.DB` (`db.Raw()`)

Call `db.Raw()` on your client instance to retrieve the underlying `*sql.DB` handle:

```go
// Execute raw query directly on *sql.DB
rows, err := db.Raw().QueryContext(ctx, `SELECT count(*) FROM "User" WHERE active = $1`, true)
if err != nil {
    return err
}
defer rows.Close()
```

---

## 2. Accessing `*sql.Tx` (`tx.Raw()`)

Call `tx.Raw()` inside any transaction to retrieve the underlying `*sql.Tx` handle:

```go
err := db.Transaction(ctx, func(tx *phi.Tx) error {
    // Execute raw SQL on the active transaction handle
    _, err := tx.Raw().ExecContext(ctx, `UPDATE "User" SET loginCount = loginCount + 1 WHERE id = $1`, userID)
    return err
})
```
