package migration

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/voidclancy/valk/schema"

	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/sqlite"
)

type SqliteDialect struct{}

func (SqliteDialect) QuoteIdent(name string) string {
	return `"` + name + `"`
}

func (SqliteDialect) GetSQLType(sf *schema.ScalarField) string {
	sqlType := strings.ToUpper(sf.SQLType)
	switch sqlType {
	case "VARCHAR", "TEXT", "UUID":
		return "TEXT"
	case "TIMESTAMP":
		return "TIMESTAMP"
	case "INTEGER", "BIGINT", "BOOLEAN":
		return "INTEGER"
	case "DOUBLE PRECISION", "NUMERIC":
		return "REAL"
	case "JSONB", "BYTEA":
		return "BLOB"
	default:
		return "TEXT"
	}
}

func (SqliteDialect) GetSQLDefault(dv *schema.DefaultValue, pslType string) string {
	switch dv.Kind {
	case schema.DefaultLiteral:
		if pslType == schema.TypeBoolean {
			return strings.ToUpper(dv.Literal)
		}
		if pslType == schema.TypeInt || pslType == schema.TypeBigInt || pslType == schema.TypeFloat || pslType == schema.TypeDecimal {
			return dv.Literal
		}
		return fmt.Sprintf("'%s'", strings.ReplaceAll(dv.Literal, "'", "''"))

	case schema.DefaultEnumValue:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(dv.EnumValue, "'", "''"))

	case schema.DefaultFunc:
		switch dv.FuncName {
		case "now":
			return "CURRENT_TIMESTAMP"
		case "uuid", "uuid(4)", "uuid(7)":
			return ""
		case "cuid", "cuid(1)", "cuid(2)":
			return ""
		case "ulid":
			return ""
		case "nanoid":
			return ""
		}

	case schema.DefaultDBGenerated:
		return dv.DBExpression
	}

	return ""
}

func (SqliteDialect) GenerateEnum(enum *schema.Enum) string {
	return "" //will be inlined using CHECK
}

func (SqliteDialect) FormatAutoIncrement(sqlType string) (string, string) {
	return "INTEGER", "" // SQLite autoincrement requires INTEGER type
}

func (d SqliteDialect) FormatSinglePK(tableName, colName string, isAutoInc bool) (string, string) {
	if isAutoInc {
		return "PRIMARY KEY AUTOINCREMENT", ""
	}
	pkName := tableName + "_pkey"
	return "", fmt.Sprintf("  CONSTRAINT %s PRIMARY KEY (%s)", d.QuoteIdent(pkName), d.QuoteIdent(colName))
}

func (d SqliteDialect) FormatEnumConstraint(colName string, enum *schema.Enum) string {
	var enumVals []string
	for _, ev := range enum.ValueMap {
		enumVals = append(enumVals, fmt.Sprintf("'%s'", ev.DBName))
	}
	return fmt.Sprintf("%s IN (%s)", d.QuoteIdent(colName), strings.Join(enumVals, ", "))
}

func (SqliteDialect) SupportsNativeEnums() bool { return false }

func (SqliteDialect) OpenConn(db *sql.DB) (migrate.Driver, error) {
	return sqlite.Open(db)
}
