package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

type rawUserRow struct {
	id           string
	email        string
	phoneNum     string
	password     sql.NullString
	role         string
	roleOptional sql.NullString
	loginCount   int32
	referredById sql.NullString
}

func (r *rawUserRow) fields() []any {
	return []any{&r.id, &r.email, &r.phoneNum, &r.password, &r.role, &r.roleOptional, &r.loginCount, &r.referredById}
}

type rawNullableUserRow struct {
	id           sql.NullString
	email        sql.NullString
	phoneNum     sql.NullString
	password     sql.NullString
	role         sql.NullString
	roleOptional sql.NullString
	loginCount   sql.NullInt64
	referredById sql.NullString
}

func (r *rawNullableUserRow) fields() []any {
	return []any{&r.id, &r.email, &r.phoneNum, &r.password, &r.role, &r.roleOptional, &r.loginCount, &r.referredById}
}

type rawDeepUserRow struct {
	u  rawUserRow
	r  rawNullableUserRow
	rr rawNullableUserRow
}

func (r *rawDeepUserRow) fields() []any {
	fields := make([]any, 0, 24)
	fields = append(fields, r.u.fields()...)
	fields = append(fields, r.r.fields()...)
	fields = append(fields, r.rr.fields()...)
	return fields
}

func benchRawCreate(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)

	query := rawQueryCreate

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, err := db.ExecContext(ctx, query,
			fmt.Sprintf("raw-create-%d", i),
			fmt.Sprintf("raw-create-%d@example.com", i),
			fmt.Sprintf("raw-phone-%d", i),
			"STUDENT",
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchRawCreateMany(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)

	query := rawQueryCreateMany

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		args := make([]any, 0, 50)
		for j := range 10 {
			n := i*10 + j
			args = append(args,
				fmt.Sprintf("raw-cmany-%d", n),
				fmt.Sprintf("raw-cmany-%d@example.com", n),
				fmt.Sprintf("raw-cmany-phone-%d", n),
				"STUDENT",
				0,
			)
		}
		_, err := db.ExecContext(ctx, query, args...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchRawCreateManyAndReturn(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)

	query := rawQueryCreateManyAndReturn

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		args := make([]any, 0, 50)
		for j := 0; j < 10; j++ {
			n := i*10 + j
			args = append(args,
				fmt.Sprintf("raw-cmar-%d", n),
				fmt.Sprintf("raw-cmar-%d@example.com", n),
				fmt.Sprintf("raw-cmar-phone-%d", n),
				"STUDENT",
				0,
			)
		}
		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			b.Fatal(err)
		}
		if rows.Err() != nil {
			b.Fatal(rows.Err())
		}
		for rows.Next() {
			var id, email, phoneNum, role string
			var loginCount int32
			if err := rows.Scan(&id, &email, &phoneNum, &role, &loginCount); err != nil {
				b.Fatal(err)
			}
		}
		rows.Close()
	}
}

func benchRawFindUnique(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)
	seedData(db, "raw-fu")

	query := rawQueryFindUnique

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		email := fmt.Sprintf("raw-fu-user-%d@example.com", i%seedCount)
		var u rawUserRow
		err := db.QueryRowContext(ctx, query, email).Scan(u.fields()...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchRawFindFirst(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)
	seedData(db, "raw-ff")

	query := rawQueryFindFirst

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		email := fmt.Sprintf("raw-ff-user-%d@example.com", i%seedCount)
		var u rawUserRow
		err := db.QueryRowContext(ctx, query, email).Scan(u.fields()...)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			b.Fatal(err)
		}
	}
}

func benchRawFindMany(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)
	seedData(db, "raw-fm")

	query := rawQueryFindMany

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		offset := i % seedCount
		rows, err := db.QueryContext(ctx, query, offset)
		if err != nil {
			b.Fatal(err)
		}
		if rows.Err() != nil {
			b.Fatal(rows.Err())
		}
		for rows.Next() {
			var u rawUserRow
			if err := rows.Scan(u.fields()...); err != nil {
				b.Fatal(err)
			}
		}
		rows.Close()
	}
}

func benchRawUpsert(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)
	seedData(db, "raw-ups")

	query := rawQueryUpsert

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		email := fmt.Sprintf("raw-ups-user-%d@example.com", i%seedCount)
		phone := fmt.Sprintf("raw-ups-phone-%d", i)
		_, err := db.ExecContext(ctx, query,
			fmt.Sprintf("raw-ups-id-%d", i),
			email, phone, "STUDENT", int32(i),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchRawReadDeepRelation(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)
	seedRelations(db, "raw-rdr")

	query := rawQueryDeepRelation

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		email := fmt.Sprintf("raw-rdr-grand-%d@example.com", i%500)
		var row rawDeepUserRow
		err := db.QueryRowContext(ctx, query, email).Scan(row.fields()...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchRawCreateWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)
	seedRelations(db, "raw-cwds")

	queryInsert := fmt.Sprintf(
		`INSERT INTO %s (%s, %s, %s, %s, %s) VALUES (%s, %s, %s, %s, %s)`,
		activeDialect.Quote("User"),
		activeDialect.Quote("id"), activeDialect.Quote("email"), activeDialect.Quote("phoneNum"), activeDialect.Quote("role"), activeDialect.Quote("referredById"),
		activeDialect.BindVar(1), activeDialect.BindVar(2), activeDialect.BindVar(3), activeDialect.BindVar(4), activeDialect.BindVar(5),
	)
	querySelectUserByID := fmt.Sprintf(
		`SELECT %s, %s, %s, %s, %s, %s, %s, %s FROM %s WHERE %s = %s`,
		activeDialect.Quote("id"), activeDialect.Quote("email"), activeDialect.Quote("phoneNum"),
		activeDialect.Quote("password"), activeDialect.Quote("role"), activeDialect.Quote("roleOptional"),
		activeDialect.Quote("loginCount"), activeDialect.Quote("referredById"),
		activeDialect.Quote("User"),
		activeDialect.Quote("id"),
		activeDialect.BindVar(1),
	)

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("raw-cwds-parent-id-%d", i%500)
		email := fmt.Sprintf("raw-cwds-new-%d@example.com", i)
		_, err := db.ExecContext(ctx, queryInsert,
			fmt.Sprintf("raw-cwds-new-id-%d", i),
			email,
			fmt.Sprintf("raw-cwds-new-phone-%d", i),
			"STUDENT",
			parentID,
		)
		if err != nil {
			b.Fatal(err)
		}
		var p rawUserRow
		err = db.QueryRowContext(ctx, querySelectUserByID, parentID).Scan(p.fields()...)
		if err != nil {
			b.Fatal(err)
		}
		if p.referredById.Valid {
			var gp rawUserRow
			err = db.QueryRowContext(ctx, querySelectUserByID, p.referredById.String).Scan(gp.fields()...)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func benchRawCreateManyAndReturnWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)
	seedRelations(db, "raw-cmwds")

	// Build insertion query for 10 records
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s, %s) VALUES `,
		activeDialect.Quote("User"),
		activeDialect.Quote("id"), activeDialect.Quote("email"), activeDialect.Quote("phoneNum"), activeDialect.Quote("role"), activeDialect.Quote("referredById"),
	))
	for j := 0; j < 10; j++ {
		if j > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(activeDialect.Placeholders(j*5+1, 5))
	}
	queryInsert := sb.String()

	// Build select query for 10 records
	userCols := []string{"id", "email", "phoneNum", "password", "role", "roleOptional", "loginCount", "referredById"}
	var selectCols []string
	for _, col := range userCols {
		selectCols = append(selectCols, fmt.Sprintf("u.%s", activeDialect.Quote(col)))
	}
	for _, col := range userCols {
		selectCols = append(selectCols, fmt.Sprintf("r.%s", activeDialect.Quote(col)))
	}
	for _, col := range userCols {
		selectCols = append(selectCols, fmt.Sprintf("rr.%s", activeDialect.Quote(col)))
	}
	querySelect := fmt.Sprintf(
		`SELECT %s FROM %s u LEFT JOIN %s r ON u.%s = r.%s LEFT JOIN %s rr ON r.%s = rr.%s WHERE u.%s IN (%s)`,
		strings.Join(selectCols, ", "),
		activeDialect.Quote("User"), activeDialect.Quote("User"), activeDialect.Quote("referredById"), activeDialect.Quote("id"),
		activeDialect.Quote("User"), activeDialect.Quote("referredById"), activeDialect.Quote("id"),
		activeDialect.Quote("id"),
		activeDialect.Placeholders(1, 10)[1:len(activeDialect.Placeholders(1, 10))-1],
	)

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("raw-cmwds-parent-id-%d", i%500)
		insertArgs := make([]any, 50)
		ids := make([]any, 10)
		for j := 0; j < 10; j++ {
			id := fmt.Sprintf("raw-cmwds-new-id-%d-%d", i, j)
			ids[j] = id
			insertArgs[j*5] = id
			insertArgs[j*5+1] = fmt.Sprintf("raw-cmwds-new-%d-%d@example.com", i, j)
			insertArgs[j*5+2] = fmt.Sprintf("raw-cmwds-new-phone-%d-%d", i, j)
			insertArgs[j*5+3] = "STUDENT"
			insertArgs[j*5+4] = parentID
		}
		_, err := db.ExecContext(ctx, queryInsert, insertArgs...)
		if err != nil {
			b.Fatal(err)
		}

		rows, err := db.QueryContext(ctx, querySelect, ids...)
		if err != nil {
			b.Fatal(err)
		}

		if rows.Err() != nil {
			b.Fatal(rows.Err())
		}
		for rows.Next() {
			var row rawDeepUserRow
			if err := rows.Scan(row.fields()...); err != nil {
				rows.Close()
				b.Fatal(err)
			}
		}
		rows.Close()
	}
}

func benchRawUpsertWithDeepSelect(b *testing.B) {
	ctx := context.Background()
	db := openDB(b)
	defer db.Close()
	createSchema(db)
	seedRelations(db, "raw-uwds")

	// Preseed
	for i := 0; i < seedCount; i++ {
		_, err := db.ExecContext(ctx, rawQueryCreate,
			fmt.Sprintf("raw-uwds-id-%d", i),
			fmt.Sprintf("raw-uwds-%d@example.com", i),
			fmt.Sprintf("raw-uwds-phone-%d", i),
			"STUDENT",
		)
		if err != nil {
			b.Fatal(err)
		}
	}

	queryUpsert := fmt.Sprintf(
		`INSERT INTO %s (%s, %s, %s, %s, %s, %s) VALUES (%s, %s, %s, %s, %s, %s) `+
			`ON CONFLICT(%s) DO UPDATE SET %s = EXCLUDED.%s, %s = EXCLUDED.%s, %s = EXCLUDED.%s, %s = EXCLUDED.%s`,
		activeDialect.Quote("User"),
		activeDialect.Quote("id"), activeDialect.Quote("email"), activeDialect.Quote("phoneNum"), activeDialect.Quote("role"), activeDialect.Quote("loginCount"), activeDialect.Quote("referredById"),
		activeDialect.BindVar(1), activeDialect.BindVar(2), activeDialect.BindVar(3), activeDialect.BindVar(4), activeDialect.BindVar(5), activeDialect.BindVar(6),
		activeDialect.Quote("email"),
		activeDialect.Quote("phoneNum"), activeDialect.Quote("phoneNum"),
		activeDialect.Quote("role"), activeDialect.Quote("role"),
		activeDialect.Quote("loginCount"), activeDialect.Quote("loginCount"),
		activeDialect.Quote("referredById"), activeDialect.Quote("referredById"),
	)
	querySelect := rawQueryDeepRelation

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		parentID := fmt.Sprintf("raw-uwds-parent-id-%d", i%500)
		email := fmt.Sprintf("raw-uwds-%d@example.com", i%seedCount)
		_, err := db.ExecContext(ctx, queryUpsert,
			fmt.Sprintf("raw-uwds-id-new-%d", i),
			email,
			fmt.Sprintf("raw-uwds-phone-new-%d", i),
			"STUDENT",
			int32(i),
			parentID,
		)
		if err != nil {
			b.Fatal(err)
		}

		var row rawDeepUserRow
		err = db.QueryRowContext(ctx, querySelect, email).Scan(row.fields()...)
		if err != nil {
			b.Fatal(err)
		}
	}
}
