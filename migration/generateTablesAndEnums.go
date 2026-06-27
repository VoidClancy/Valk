package migration

import (
	"fmt"
	"strings"
	"valkyrie/schema"
)

type tableProps struct {
	model   *schema.Model
	dialect Dialect
	dbName  string
}

func generateTables(schemaDef *schema.Schema, dialect Dialect, sb *strings.Builder) {
	for _, model := range schemaDef.Models {
		tableName := model.TableName
		if tableName == "" {
			tableName = model.Name
		}

		fmt.Fprintf(sb, "CREATE TABLE %s (\n", dialect.QuoteIdent(tableName))

		// 1. add scalar fields (cols)
		columns, tableConstraints := generateScalarFields(model, dialect, tableName)

		// 2. composite PK table constraint
		if pkConstraint := generateCompositePK(model, tableName, dialect); pkConstraint != "" {
			tableConstraints = append(tableConstraints, pkConstraint)
		}

		// 3. composite unique constraints
		tableConstraints = append(tableConstraints, generateCompositeUniques(model, tableName, dialect)...)

		// 4. FK table constraints
		tableConstraints = append(tableConstraints, generateForeignKeys(model, tableName, dialect)...)

		// Append tableConstraints to cols
		columns = append(columns, tableConstraints...)

		sb.WriteString(strings.Join(columns, ",\n"))
		sb.WriteString("\n);\n\n")

		generateIndexes(model, tableName, dialect, sb)

		if len(model.Indexes) > 0 {
			sb.WriteString("\n")
		}
	}
}

func generateEnums(enums []*schema.Enum, dialect Dialect, sb *strings.Builder) {
	for _, enum := range enums {
		enumDDL := dialect.GenerateEnum(enum)
		if enumDDL != "" {
			sb.WriteString(enumDDL)
		}
	}
}

//----------------------------------------------

func generateCompositePK(model *schema.Model, tableName string, dialect Dialect) string {
	if len(model.CompositePK) == 0 {
		return ""
	}
	var pkCols []string
	for _, pkField := range model.CompositePK {
		cName := pkField
		for _, sf := range model.ScalarFields {
			if sf.Name == pkField {
				if sf.ColName != "" {
					cName = sf.ColName
				}
				break
			}
		}
		pkCols = append(pkCols, dialect.QuoteIdent(cName))
	}
	pkName := tableName + "_pkey"
	return fmt.Sprintf("  CONSTRAINT %s PRIMARY KEY (%s)", dialect.QuoteIdent(pkName), strings.Join(pkCols, ", "))
}

func generateCompositeUniques(model *schema.Model, tableName string, dialect Dialect) []string {
	var constraints []string
	for _, uniq := range model.CompositeUnique {
		var uniqCols []string
		var uniqColNames []string
		for _, uField := range uniq.Fields {
			cName := uField
			for _, sf := range model.ScalarFields {
				if sf.Name == uField {
					if sf.ColName != "" {
						cName = sf.ColName
					}
					break
				}
			}
			uniqCols = append(uniqCols, dialect.QuoteIdent(cName))
			uniqColNames = append(uniqColNames, cName)
		}
		uniqName := uniq.Name
		if uniqName == "" {
			uniqName = tableName + "_" + strings.Join(uniqColNames, "_") + "_key"
		}
		constraints = append(constraints, fmt.Sprintf("  CONSTRAINT %s UNIQUE (%s)", dialect.QuoteIdent(uniqName), strings.Join(uniqCols, ", ")))
	}
	return constraints
}

func generateForeignKeys(model *schema.Model, tableName string, dialect Dialect) []string {
	var constraints []string
	for _, rf := range model.RelationFields {
		if len(rf.FKFields) > 0 && len(rf.RefFields) > 0 {
			var fkCols []string
			var fkColNames []string
			for _, fkField := range rf.FKFields {
				colName := fkField.ColName
				if colName == "" {
					colName = fkField.Name
				}
				fkCols = append(fkCols, dialect.QuoteIdent(colName))
				fkColNames = append(fkColNames, colName)
			}

			var refCols []string
			targetTable := rf.TargetModel.TableName
			if targetTable == "" {
				targetTable = rf.TargetModel.Name
			}
			for _, refField := range rf.RefFields {
				colName := refField.ColName
				if colName == "" {
					colName = refField.Name
				}
				refCols = append(refCols, dialect.QuoteIdent(colName))
			}

			fkName := fmt.Sprintf("%s_%s_fkey", tableName, strings.Join(fkColNames, "_"))
			fkConstraint := fmt.Sprintf("  CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
				dialect.QuoteIdent(fkName),
				strings.Join(fkCols, ", "),
				dialect.QuoteIdent(targetTable),
				strings.Join(refCols, ", "),
			)

			if rf.OnDelete != "" {
				fkConstraint += " ON DELETE " + formatReferentialAction(rf.OnDelete)
			}
			if rf.OnUpdate != "" {
				fkConstraint += " ON UPDATE " + formatReferentialAction(rf.OnUpdate)
			}

			constraints = append(constraints, fkConstraint)
		}
	}
	return constraints
}

func generateIndexes(model *schema.Model, tableName string, dialect Dialect, sb *strings.Builder) {
	for _, idx := range model.Indexes {
		var idxCols []string
		var colNamesForName []string
		for _, iField := range idx.Fields {
			cName := iField
			for _, sf := range model.ScalarFields {
				if sf.Name == iField {
					cName = sf.ColName
					break
				}
			}
			idxCols = append(idxCols, dialect.QuoteIdent(cName))
			colNamesForName = append(colNamesForName, cName)
		}

		idxName := idx.Name
		if idxName == "" {
			idxName = fmt.Sprintf("%s_%s_idx", tableName, strings.Join(colNamesForName, "_"))
		}

		fmt.Fprintf(sb, "CREATE INDEX %s ON %s (%s);\n",
			dialect.QuoteIdent(idxName),
			dialect.QuoteIdent(tableName),
			strings.Join(idxCols, ", "))
	}
}

func generateScalarFields(model *schema.Model, dialect Dialect, tableName string) ([]string, []string) {
	var columns []string
	var tableConstraints []string

	for _, sf := range model.ScalarFields {
		colName := sf.ColName
		if colName == "" {
			colName = sf.Name
		}

		sqlType := dialect.GetSQLType(sf)

		var colParts []string
		colParts = append(colParts, dialect.QuoteIdent(colName))
		colParts = append(colParts, sqlType)

		// Handle Auto Increment / Serial
		if sf.Default != nil && sf.Default.FuncName == "autoincrement" {
			typeOverride, keyword := dialect.FormatAutoIncrement(sqlType)
			if typeOverride != "" {
				sqlType = typeOverride
				colParts[1] = sqlType // Update type in parts
			}
			if keyword != "" {
				colParts = append(colParts, keyword)
			}
		}

		// Nullability
		if sf.Optional {
			colParts = append(colParts, "NULL")
		} else {
			// Serial columns in postgres are implicitly not null, but explicit is good
			colParts = append(colParts, "NOT NULL")
		}

		// Single field Primary Key
		if sf.IsID && len(model.CompositePK) == 0 {
			isAutoInc := sf.Default != nil && sf.Default.FuncName == "autoincrement"
			inlinePK, tablePK := dialect.FormatSinglePK(tableName, colName, isAutoInc)
			if inlinePK != "" {
				colParts = append(colParts, inlinePK)
			}
			if tablePK != "" {
				tableConstraints = append(tableConstraints, tablePK)
			}
		}

		// Single field Unique constraint
		if sf.IsUnique {
			uniqName := tableName + "_" + colName + "_key"
			tableConstraints = append(tableConstraints, fmt.Sprintf("  CONSTRAINT %s UNIQUE (%s)", dialect.QuoteIdent(uniqName), dialect.QuoteIdent(colName)))
		}

		// Default value (except autoincrement which we handled)
		if sf.Default != nil && sf.Default.FuncName != "autoincrement" {
			var defStr string
			if sf.EnumRef != nil && sf.Default.Kind == schema.DefaultEnumValue {
				dbVal := sf.Default.EnumValue
				for _, ev := range sf.EnumRef.ValueMap {
					if ev.Name == sf.Default.EnumValue {
						dbVal = ev.DBName
						break
					}
				}
				defStr = fmt.Sprintf("'%s'", strings.ReplaceAll(dbVal, "'", "''"))
			} else {
				defStr = dialect.GetSQLDefault(sf.Default, sf.Type)
			}
			if defStr != "" {
				colParts = append(colParts, "DEFAULT "+defStr)
			}
		}

		// Inline Enum Constraints (e.g. SQLite CHECK constraints)
		if sf.EnumRef != nil {
			constraint := dialect.FormatEnumConstraint(colName, sf.EnumRef)
			if constraint != "" {
				colParts = append(colParts, constraint)
			}
		}

		columns = append(columns, "  "+strings.Join(colParts, " "))
	}

	return columns, tableConstraints
}

func formatReferentialAction(action string) string {
	switch strings.ToLower(action) {
	case "setnull":
		return "SET NULL"
	case "setdefault":
		return "SET DEFAULT"
	case "noaction":
		return "NO ACTION"
	default:
		return strings.ToUpper(action)
	}
}
