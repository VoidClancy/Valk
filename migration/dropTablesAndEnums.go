package migration

import (
	"fmt"
	"strings"
	"valkyrie/schema"
)

type dropProps struct {
	schemaDef   schema.Schema
	dialect     Dialect
	downBuilder *strings.Builder
}

func dropTables(schemaDef *schema.Schema, dialect Dialect, downBuilder *strings.Builder) {
	for i := len(schemaDef.Models) - 1; i >= 0; i-- {
		model := schemaDef.Models[i]
		tableName := model.TableName
		if tableName == "" {
			tableName = model.Name
		}
		fmt.Fprintf(downBuilder, "DROP TABLE IF EXISTS %s;\n", dialect.QuoteIdent(tableName))
	}
}

func dropEnums(schemaDef *schema.Schema, dialect Dialect, downBuilder *strings.Builder) {
	for i := len(schemaDef.Enums) - 1; i >= 0; i-- {
		enum := schemaDef.Enums[i]
		name := enum.Name
		if enum.TableMapName != "" {
			name = enum.TableMapName
		}
		if dialect.GenerateEnum(enum) != "" {
			fmt.Fprintf(downBuilder, "DROP TYPE IF EXISTS %s;\n", dialect.QuoteIdent(name))
		}
	}
}
