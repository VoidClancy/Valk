package migration

import (
	"fmt"
	"strings"
	providers "valkyrie/dbProviders"
	vs "valkyrie/schema"

	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlite"
)

type atlasSchemaBuilder struct {
	targetSchema *schema.Schema
	schemaDef    *vs.Schema
	provider     providers.DbProvider
	dialect      Dialect
	enumsMap     map[string]*schema.EnumType
	tablesMap    map[string]*schema.Table
}

func ConvertToAtlasSchema(schemaDef *vs.Schema, provider providers.DbProvider, schemaName string) (*schema.Schema, error) {
	builder := &atlasSchemaBuilder{
		targetSchema: &schema.Schema{
			Name: schemaName,
		},
		schemaDef: schemaDef,
		provider:  provider,
		dialect:   GetDialect(provider),
		enumsMap:  make(map[string]*schema.EnumType),
		tablesMap: make(map[string]*schema.Table),
	}

	builder.buildEnumsMap()

	if err := builder.buildTablesMap(); err != nil {
		return nil, err
	}

	if err := builder.buildForeignKeys(); err != nil {
		return nil, err
	}

	return builder.targetSchema, nil
}

func findColumn(table *schema.Table, name string) *schema.Column {
	for _, c := range table.Columns {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (b *atlasSchemaBuilder) buildEnumsMap() {
	if b.dialect != nil && b.dialect.SupportsNativeEnums() {
		for _, enum := range b.schemaDef.Enums {
			enumName := enum.Name
			if enum.TableMapName != "" {
				enumName = enum.TableMapName
			}

			var values []string
			for _, val := range enum.ValueMap {
				values = append(values, val.DBName)
			}

			enumType := &schema.EnumType{
				T:      enumName,
				Values: values,
				Schema: b.targetSchema,
			}
			b.targetSchema.Objects = append(b.targetSchema.Objects, enumType)
			b.enumsMap[enumName] = enumType
		}
	}
}

func (b *atlasSchemaBuilder) buildTablesMap() error {
	for _, model := range b.schemaDef.Models {
		tableName := model.EffectiveTableName()

		table := &schema.Table{
			Name:   tableName,
			Schema: b.targetSchema,
		}

		if err := b.buildColumns(model, table); err != nil {
			return err
		}

		b.buildPrimaryKey(model, table)
		b.buildUniqueConstraints(model, table)
		b.buildIndexes(model, table)
		b.buildEnumCheckConstraints(model, table)

		b.targetSchema.Tables = append(b.targetSchema.Tables, table)
		b.tablesMap[tableName] = table
	}

	return nil
}

func (b *atlasSchemaBuilder) buildColumns(model *vs.Model, table *schema.Table) error {
	for _, sf := range model.ScalarFields {
		colName := sf.EffectiveColName()

		var typ schema.Type
		var err error
		var sqlTypeStr string

		if sf.EnumRef != nil && b.dialect != nil && b.dialect.SupportsNativeEnums() {
			enumName := sf.EnumRef.Name
			if sf.EnumRef.TableMapName != "" {
				enumName = sf.EnumRef.TableMapName
			}
			sqlTypeStr = enumName
			var exists bool
			typ, exists = b.enumsMap[enumName]
			if !exists {
				typ = &schema.EnumType{T: enumName}
			}
		} else {
			sqlTypeStr = getSQLType(sf, b.provider)
			if sf.Default != nil && sf.Default.FuncName == "autoincrement" && b.dialect != nil {
				typeOverride, _ := b.dialect.FormatAutoIncrement(sqlTypeStr)
				if typeOverride != "" {
					sqlTypeStr = typeOverride
				}
			}
			if b.provider == providers.Postgres || b.provider == providers.Postgresql {
				typ, err = postgres.ParseType(sqlTypeStr)
			} else {
				typ, err = sqlite.ParseType(sqlTypeStr)
			}
			if err != nil {
				return fmt.Errorf("failed to parse type %q for field %s.%s: %w", sqlTypeStr, model.Name, sf.Name, err)
			}
		}

		column := &schema.Column{
			Name: colName,
			Type: &schema.ColumnType{
				Type: typ,
				Null: sf.Optional,
				Raw:  sqlTypeStr,
			},
		}

		if sf.Default != nil && sf.Default.FuncName == "autoincrement" {
			if b.provider == providers.Sqlite {
				column.Attrs = append(column.Attrs, &sqlite.AutoIncrement{})
			}
		}

		// Handle other defaults
		if sf.Default != nil && sf.Default.FuncName != "autoincrement" {
			var defaultVal string
			if sf.EnumRef != nil && sf.Default.Kind == vs.DefaultEnumValue {
				dbVal := sf.Default.EnumValue
				for _, ev := range sf.EnumRef.ValueMap {
					if ev.Name == sf.Default.EnumValue {
						dbVal = ev.DBName
						break
					}
				}
				defaultVal = fmt.Sprintf("'%s'", strings.ReplaceAll(dbVal, "'", "''"))
			} else {
				defaultVal = getSQLDefault(sf.Default, sf.Type, b.provider)
			}
			if defaultVal != "" {
				column.Default = &schema.RawExpr{X: defaultVal}
			}
		}

		table.Columns = append(table.Columns, column)
	}
	return nil
}

func (b *atlasSchemaBuilder) buildPrimaryKey(model *vs.Model, table *schema.Table) {
	tableName := table.Name
	if len(model.CompositePK) > 0 {
		pk := &schema.Index{
			Name:   tableName + "_pkey",
			Unique: true,
			Table:  table,
		}
		for _, pkField := range model.CompositePK {
			targetColName := getColName(model, pkField)
			if col := findColumn(table, targetColName); col != nil {
				pk.Parts = append(pk.Parts, &schema.IndexPart{C: col})
			}
		}
		table.PrimaryKey = pk
	} else {
		for _, sf := range model.ScalarFields {
			if sf.IsID {
				colName := sf.EffectiveColName()
				if col := findColumn(table, colName); col != nil {
					table.PrimaryKey = &schema.Index{
						Name:   tableName + "_pkey",
						Unique: true,
						Table:  table,
						Parts:  []*schema.IndexPart{{C: col}},
					}
				}
				break
			}
		}
	}
}

func (b *atlasSchemaBuilder) buildUniqueConstraints(model *vs.Model, table *schema.Table) {
	tableName := table.Name
	for _, sf := range model.ScalarFields {
		if sf.IsUnique {
			colName := sf.EffectiveColName()
			if col := findColumn(table, colName); col != nil {
				table.Indexes = append(table.Indexes, &schema.Index{
					Name:   tableName + "_" + colName + "_key",
					Unique: true,
					Table:  table,
					Parts:  []*schema.IndexPart{{C: col}},
				})
			}
		}
	}

	for _, uniq := range model.CompositeUnique {
		uniqName := uniq.Name
		var uniqColNames []string
		var indexParts []*schema.IndexPart
		for _, uField := range uniq.Fields {
			targetColName := getColName(model, uField)
			uniqColNames = append(uniqColNames, targetColName)
			if col := findColumn(table, targetColName); col != nil {
				indexParts = append(indexParts, &schema.IndexPart{C: col})
			}
		}
		if uniqName == "" {
			uniqName = tableName + "_" + strings.Join(uniqColNames, "_") + "_key"
		}
		table.Indexes = append(table.Indexes, &schema.Index{
			Name:   uniqName,
			Unique: true,
			Table:  table,
			Parts:  indexParts,
		})
	}
}

func (b *atlasSchemaBuilder) buildIndexes(model *vs.Model, table *schema.Table) {
	tableName := table.Name
	for _, idx := range model.Indexes {
		idxName := idx.Name
		var colNamesForName []string
		var indexParts []*schema.IndexPart
		for _, iField := range idx.Fields {
			targetColName := getColName(model, iField)
			colNamesForName = append(colNamesForName, targetColName)
			if col := findColumn(table, targetColName); col != nil {
				indexParts = append(indexParts, &schema.IndexPart{C: col})
			}
		}
		if idxName == "" {
			idxName = tableName + "_" + strings.Join(colNamesForName, "_") + "_idx"
		}
		table.Indexes = append(table.Indexes, &schema.Index{
			Name:   idxName,
			Unique: false,
			Table:  table,
			Parts:  indexParts,
		})
	}
}

func (b *atlasSchemaBuilder) buildEnumCheckConstraints(model *vs.Model, table *schema.Table) {
	if b.dialect == nil {
		return
	}
	tableName := table.Name
	for _, sf := range model.ScalarFields {
		if sf.EnumRef != nil {
			colName := sf.EffectiveColName()
			expr := b.dialect.FormatEnumConstraint(colName, sf.EnumRef)
			if expr != "" {
				checkConstraint := &schema.Check{
					Name: tableName + "_" + colName + "_check",
					Expr: expr,
				}
				table.Attrs = append(table.Attrs, checkConstraint)
			}
		}
	}
}

func (b *atlasSchemaBuilder) buildForeignKeys() error {
	for _, model := range b.schemaDef.Models {
		tableName := model.EffectiveTableName()
		table := b.tablesMap[tableName]

		for _, rf := range model.RelationFields {
			if len(rf.FKFields) > 0 && len(rf.RefFields) > 0 {
				var fkColNames []string
				var localCols []*schema.Column
				for _, fkField := range rf.FKFields {
					colName := fkField.EffectiveColName()
					fkColNames = append(fkColNames, colName)
					if col := findColumn(table, colName); col != nil {
						localCols = append(localCols, col)
					}
				}

				targetTable := rf.TargetModel.EffectiveTableName()
				refTable := b.tablesMap[targetTable]
				if refTable == nil {
					return fmt.Errorf("referenced table %q not found for relation in table %q", targetTable, tableName)
				}

				var refCols []*schema.Column
				for _, refField := range rf.RefFields {
					colName := refField.EffectiveColName()
					if col := findColumn(refTable, colName); col != nil {
						refCols = append(refCols, col)
					}
				}

				fkName := fmt.Sprintf("%s_%s_fkey", tableName, strings.Join(fkColNames, "_"))
				fk := &schema.ForeignKey{
					Symbol:     fkName,
					Table:      table,
					RefTable:   refTable,
					Columns:    localCols,
					RefColumns: refCols,
					OnDelete:   mapRefOption(rf.OnDelete),
					OnUpdate:   mapRefOption(rf.OnUpdate),
				}
				table.ForeignKeys = append(table.ForeignKeys, fk)
				for _, col := range localCols {
					col.ForeignKeys = append(col.ForeignKeys, fk)
				}
			}
		}
	}
	return nil
}

func getColName(model *vs.Model, fieldName string) string {
	for _, sf := range model.ScalarFields {
		if sf.Name == fieldName {
			return sf.EffectiveColName()
		}
	}
	return fieldName
}

func mapRefOption(opt string) schema.ReferenceOption {
	switch strings.ToLower(opt) {
	case "cascade":
		return schema.Cascade
	case "restrict":
		return schema.Restrict
	case "set_null", "setnull":
		return schema.SetNull
	case "set_default", "setdefault":
		return schema.SetDefault
	case "no_action", "noaction":
		return schema.NoAction
	default:
		return schema.NoAction
	}
}
