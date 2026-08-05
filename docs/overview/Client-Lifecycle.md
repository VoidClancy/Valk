---
title: Client Lifecycle
description: Connection pooling, statement caching, extension registration, and graceful shutdown in Phi.
category: Overview
tags: [client, lifecycle, connection-pool, caching, open, close]
---

# Client Lifecycle

This guide covers initializing the Phi client, configuring connection pool bounds, registering extension hooks during application startup, statement caching, and executing graceful shutdowns.

---

## 1. Opening a Client Connection

Initialize the generated client using `phi.Open(driverName, dataSourceName)`:

```go
package main

import (
    "log"
    _ "github.com/lib/pq" // Import SQL driver
    "myproject/phi"
)

func main() {
    // Open connection pool
    db, err := phi.Open("postgres", "postgres://user:pass@localhost:5432/dbname?sslmode=disable")
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }
    defer db.Close() // Guarantee clean shutdown
}
```

---

## 2. Connection Pool Configuration

`db.Raw()` returns the underlying `*sql.DB` handle. Use it during startup to configure standard Go database connection pool parameters:

```go
db.Raw().SetMaxOpenConns(25)
db.Raw().SetMaxIdleConns(5)
db.Raw().SetConnMaxLifetime(5 * time.Minute)
db.Raw().SetConnMaxIdleTime(1 * time.Minute)
```

---

## 3. Registering Extension Hooks

Register middleware extensions (`.Use(...)`) during application initialization **before** serving queries. Registered extensions are stored in thread-safe delegates and inherited by transactions:

```go
func initDB() (*phi.DB, error) {
    db, err := phi.Open("sqlite3", "file:app.db?cache=shared&_pragma=foreign_keys(1)")
    if err != nil {
        return nil, err
    }

    // Register global audit hook on user delegate
    db.User.Use(user.Extension{
        Create: func(ctx context.Context, args *phi.UserCreateArgs, next phi.UserCreateQuery) (*phi.User, error) {
            args.Data.Email = strings.ToLower(args.Data.Email)
            return next(ctx, args)
        },
    })

    return db, nil
}
```

---

## 4. Prepared Statement Caching

Phi includes a thread-safe internal prepared statement cache (`stmtCache`). Frequently executed query shapes reuse prepared statements automatically, minimizing database parsing overhead and boosting throughput.

---

## 5. Graceful Shutdown (`db.Close()`)

When your application terminates, call `db.Close()` to flush prepared statement caches, release active database connections, and shut down the connection pool cleanly:

```go
func main() {
    db, err := phi.Open("postgres", connString)
    if err != nil {
        log.Fatal(err)
    }

    // Graceful shutdown on OS signals
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigChan
        log.Println("Shutting down database client...")
        if err := db.Close(); err != nil {
            log.Printf("Error closing database connection: %v", err)
        }
        os.Exit(0)
    }()
}
```
