---
title: Select
description: Select specific fields and relations to return from a query.
category: Select & Omit
tags:
    [
        select,
        omit,
        scalar selection,
        relational selection,
        nested selection,
        query selection,
        omit fields,
        omit relations,
    ]
categoryOrder: 3
order: 1
---

# Select

`Select` specifies which fields and relations are returned by a query.

Scalar fields and **to-one relations** are selected using a generated `Select` struct. **To-many relations** are selected using their generated query builder.

## Selecting Scalar Fields

Select scalar fields by setting their corresponding field to `true`.

Select to-one relations by assigning their generated `Select` struct.

```go
users, err := db.User.
    FindMany(user.Email.EQ("x@y.com")).
    Select(user.Select{
        Id:    true,
        Email: true,
        Bio:   true,

        Profile: &profile.Select{
            Id:  true,
            Bio: true,
        },
    }).
    Exec(ctx)
```

## Selecting To-Many Relations

Select to-many relations using their generated query builder. This allows filtering, ordering, pagination, and nested selections.

```go
users, err := db.User.
    FindMany(user.Email.EQ("x@y.com")).
    Select(user.Select{
        Id:    true,
        Email: true,
        Bio:   true,

        Posts: post.Query().
            Where(post.Title.Contains("News")).
            OrderBy(post.PublishedAt.Desc()).
            Select(post.Select{
                Id:      true,
                Title:   true,
                Content: true,
            }),
    }).
    Exec(ctx)
```

To-many relation fields also accept the relation's `Select` struct directly. Passing a plain struct - instead of a query builder - selects **all** matching related records, without any filtering, ordering, or pagination:

```go
users, err := db.User.
    FindMany(user.Email.EQ("x@y.com")).
    Select(user.Select{
        Id: true,

        Posts: &post.Select{
            Id:    true,
            Title: true,
        },
    }).
    Exec(ctx)
```

> **Note:** Relation fields are not scalar values. To-one relations hold a pointer to the relation's `Select` struct (for example `*profile.Select`), and to-many relations hold a select query - either the relation's `Select` struct (as shown above) or its query builder. You must assign one of these, never `true`.

Nested selections can be composed to any depth.

```go
users, err := db.User.
    FindMany(user.Email.EQ("x@y.com")).
    Select(user.Select{
        Id:    true,
        Email: true,
        Bio:   true,

        Posts: post.Query().
            Where(post.Title.Contains("News")).
            OrderBy(post.PublishedAt.Desc()).
            Select(post.Select{
                Id:      true,
                Title:   true,
                Content: true,

                Comments: comment.Query().
                    Where(comment.Content.Contains("great")).
                    Select(comment.Select{
                        Id:      true,
                        Content: true,
                    }),
            }),
    }).
    Exec(ctx)
```

### Selection Rules

- Scalar fields are selected using boolean values.
- To-one relations are selected with a pointer to their generated `Select` struct (`*profile.Select`).
- To-many relations are selected with either their generated query builder (for filtering, ordering, and pagination) or their `Select` struct directly (which selects all related records).
- Selections can be nested to any depth.

> **Note:** `Select` and `Omit` are mutually exclusive. Attempting to use both on the same query will result in an error.

> **Note:** If `Select` is omitted (or an empty `Select` struct is provided), Phi returns all scalar fields by default. Relations are never loaded unless explicitly selected.

## Supported By

`Select` can be used with:

- [`Create`](/docs/Create-Records#create)
- [`CreateManyAndReturn`](/docs/Create-Records#createmanyandreturn)
- [`FindUnique`](/docs/Read-Records#findunique)
- [`FindFirst`](/docs/Read-Records#findfirst)
- [`FindMany`](/docs/Read-Records#findmany)
- [`Update`](/docs/Update-Records#update)
- [`UpdateManyAndReturn`](/docs/Update-Records#updatemanyandreturn)
- [`Delete`](/docs/Delete-Records#delete)
