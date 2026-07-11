package migration

import (
	"database/sql"
	"fmt"
	"strings"
	"github.com/voidclancy/valk/schema"

	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/postgres"
)

type PostgresDialect struct{}

func (d PostgresDialect) QuoteIdent(name string) string {
	return `"` + name + `"`
}

func (d PostgresDialect) GetSQLType(sf *schema.ScalarField) string {
	// If its a custom PG enum, return the quoted enum name
	if sf.EnumRef != nil {
		enumName := sf.EnumRef.Name
		if sf.EnumRef.TableMapName != "" {
			enumName = sf.EnumRef.TableMapName
		}
		return d.QuoteIdent(enumName)
	}
	if sf.NativeType != nil && len(sf.NativeType.Args) > 0 {
		return fmt.Sprintf("%s(%s)", strings.ToLower(sf.SQLType), strings.Join(sf.NativeType.Args, ", "))
	}
	return strings.ToLower(sf.SQLType)
}

func (d PostgresDialect) GetSQLDefault(dv *schema.DefaultValue, pslType string) string {
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

func (d PostgresDialect) GenerateEnum(enum *schema.Enum) string {
	name := enum.Name
	if enum.TableMapName != "" {
		name = enum.TableMapName
	}
	var quotedValues []string
	for _, val := range enum.ValueMap {
		quotedValues = append(quotedValues, fmt.Sprintf("  '%s'", val.DBName))
	}
	return fmt.Sprintf("CREATE TYPE %s AS ENUM (\n%s\n);\n\n",
		d.QuoteIdent(name), strings.Join(quotedValues, ",\n"))
}

func (PostgresDialect) FormatAutoIncrement(sqlType string) (string, string) {
	if sqlType == "bigint" {
		return "bigserial", ""
	}
	return "serial", ""
}

func (d PostgresDialect) FormatSinglePK(tableName, colName string, isAutoInc bool) (string, string) {
	pkName := tableName + "_pkey"
	return "", fmt.Sprintf("  CONSTRAINT %s PRIMARY KEY (%s)", d.QuoteIdent(pkName), d.QuoteIdent(colName))
}

func (d PostgresDialect) FormatEnumConstraint(colName string, enum *schema.Enum) string { return "" }

func (PostgresDialect) SupportsNativeEnums() bool { return true }

func (PostgresDialect) OpenConn(db *sql.DB) (migrate.Driver, error) {
	return postgres.Open(db)
}
