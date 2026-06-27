package migration

import (
	"fmt"
	"strings"
	"valkyrie/schema"
)

type SqliteDialect struct{}

func (SqliteDialect) QuoteIdent(name string) string {
	return `"` + name + `"`
}

func (SqliteDialect) GetSQLType(sf *schema.ScalarField) string {
	sqlType := strings.ToUpper(sf.SQLType)
	switch sqlType {
	case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
		return "TEXT"
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
		case "uuid":
			return ""
		case "cuid":
			return ""
		}

	case schema.DefaultDBGenerated:
		return dv.DBExpression
	}

	return ""
}
