---
title: Transactions
description: Closure-based transactions, automatic rollback, panic recovery, and manual transaction control.
category: Operations
tags: [transactions, rollback, panic-recovery, begin-tx, commit]
---

# Transactions

Phi provides closure-based transactions with automatic rollback and panic recovery, as well as manual transaction handles.

---

## 1. Closure-Based Transactions (`db.Transaction`)

`db.Transaction(ctx, fn)` manages the entire lifecycle of a transaction automatically:

```go
err := db.Transaction(ctx, func(tx *phi.Tx) error {
    // 1. Create a user within the transaction
    u, err := tx.User.Create().
        SetEmail("tx-user@example.com").
        SetPhoneNum("+1999").
        Exec(ctx)
    if err != nil {
        return err // Triggers automatic tx.Rollback()
    }

    // 2. Create a post associated with the new user
    _, err = tx.Post.Create().
        SetTitle("First Post").
        SetContent("Hello World").
        SetAuthorId(u.Id).
        Exec(ctx)
    if err != nil {
        return err // Triggers automatic tx.Rollback()
    }

    return nil // Triggers automatic tx.Commit()
})
```

---

## 2. Lifecycle & Execution Rules

1. **Automatic Begin**: Executes `db.BeginTx(ctx, nil)` to start the transaction.
2. **Commit on Success**: If `fn(tx)` returns `nil`, `tx.Commit()` is called automatically.
3. **Rollback on Error**: If `fn(tx)` returns a non-nil `error`, `tx.Rollback()` is automatically executed, and the original error is returned.
4. **Panic Protection & Recovery**:
   If a panic occurs anywhere inside `fn(tx)`, Phi's internal `defer` block catches the panic, rolls back the transaction to prevent database lockups, and then re-throws (`repanic`) the original panic:

   ```go
   // Internal Phi Panic-Safety Guard:
   defer func() {
       if p := recover(); p != nil {
           _ = tx.Rollback() // Guarantees database rollback before repanicking
           panic(p)          // Re-throws original panic
       }
   }()
   ```

---

## 3. Manual Transactions (`db.BeginTx`)

For workflows requiring manual control across function boundaries:

```go
// 1. Begin manual transaction
tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return err
}

// 2. Execute operations on tx
u, err := tx.User.Create().
    SetEmail("manual-tx@example.com").
    SetPhoneNum("+2000").
    Exec(ctx)

if err != nil {
    _ = tx.Rollback()
    return err
}

// 3. Commit manually
if err := tx.Commit(); err != nil {
    return err
}
```
