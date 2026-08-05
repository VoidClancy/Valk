---
title: Omit
description: Omit specific scalar fields from a query result.
category: Select & Omit
tags: [omit, omit fields, scalar selection]
---

# Omit

`Omit` specifies which scalar fields should be excluded from a query result.

Unlike `Select`, `Omit` only applies to scalar fields. Relations are never loaded unless explicitly selected with `Select`.

## Omitting Scalar Fields

Omit scalar fields by setting their corresponding field to `true`. Set a field to `false` (or leave it unset) to include it in the result.

```go
users, err := db.User.
    FindMany(user.Email.EQ("x@y.com")).
    Omit(user.Omit{
        Id:    true,
        Email: true,
        Bio:   true,
    }).
    Exec(ctx)
```

> **Note:** `Select` and `Omit` are mutually exclusive. Attempting to use both on the same query will result in an error.

## Supported By

`Omit` is supported by:

- [`Create`](/docs/Create-Records#create)
- [`CreateManyAndReturn`](/docs/Create-Records#createmanyandreturn)
- [`FindUnique`](/docs/Read-Records#findunique)
- [`FindFirst`](/docs/Read-Records#findfirst)
- [`FindMany`](/docs/Read-Records#findmany)
- [`Update`](/docs/Update-Records#update)
- [`UpdateManyAndReturn`](/docs/Update-Records#updatemanyandreturn)
- [`Delete`](/docs/Delete-Records#delete)