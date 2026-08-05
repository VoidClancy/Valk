---
title: Count Records
description: Count the number of records matching a set of predicates.
category: CRUD
tags: [count, count records]
---

# Count Records

`Count` returns the number of records that match the supplied predicates.

## Count

To count all records matching one or more predicates, call `Count` and execute the query.

```go
usersCount, err := db.User.
    Count(user.LoginCount.GTE(10)).
    Exec(ctx)
```

You can also combine `Count` with pagination methods to count only a subset of the matching records.

```go
usersCount, err := db.User.
    Count(user.LoginCount.GTE(10)).
    Take(10).
    Skip(2).
    Exec(ctx)
```

### Supported Builder Methods

| Method               | Description                                                                |
| -------------------- | -------------------------------------------------------------------------- |
| [`Take`](/docs/Take) | Limit the number of records included in the count. Mirrors SQL's `LIMIT`.  |
| [`Skip`](/docs/Skip) | Skip a number of matching records before counting. Mirrors SQL's `OFFSET`. |

> **Note:** `Take` and `Skip` affect the result of the count. For example, if 100 records match a predicate and you call `Take(10)`, the returned count will be `10`, not `100`.
