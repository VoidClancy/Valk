package migration

import (
	"fmt"
	"strings"
	"valkyrie/schema"
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
	return strings.ToUpper(sf.SQLType)
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
		case "uuid":
			return "gen_random_uuid()"
		case "cuid":
			return ""
		}

	case schema.DefaultDBGenerated:
		return dv.DBExpression
	}

	return ""
}
