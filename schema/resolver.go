package schema

import (
	"fmt"
)

type TypeMap struct {
	GoType, SQLType string
}

var scalarTypeInfo = map[string]TypeMap{
	"String":   {"string", "TEXT"},
	"Int":      {"int32", "INTEGER"},
	"BigInt":   {"int64", "BIGINT"},
	"Float":    {"float64", "DOUBLE PRECISION"},
	"Decimal":  {"string", "NUMERIC"},
	"Boolean":  {"bool", "BOOLEAN"},
	"DateTime": {"time.Time", "TIMESTAMP"},
	"Json":     {"json.RawMessage", "JSONB"},
	"Bytes":    {"[]byte", "BYTEA"},
}

type Resolver struct {
	ast        *astAST
	errors     DiagnosticList
	enums      map[string]*Enum
	models     map[string]*Model
	modelOrder []*Model
}

func NewResolver(ast *astAST) *Resolver {
	return &Resolver{
		ast:    ast,
		enums:  map[string]*Enum{},
		models: map[string]*Model{},
	}
}

func (r *Resolver) errorf(line, col int, format string, args ...any) {
	r.errors = append(r.errors, Diagnostic{
		Severity: SevError,
		Message:  fmt.Sprintf(format, args...),
		Pos:      Position{Line: line, Col: col},
		Source:   "resolver",
	})
}

func (r *Resolver) Resolve() *Schema {
	schema := &Schema{}

	r.resolveDatasource(schema)
	r.registerEnums()
	r.registerModels()
	r.resolveScalarFields()
	r.resolveRelationFields()
	r.linkInverseRelations()
	r.inferRelationKinds()
	r.resolveModelAttributes()

	for _, decl := range r.ast.Enums {
		if enum, ok := r.enums[decl.Name]; ok {
			schema.Enums = append(schema.Enums, enum)
		}
	}

	schema.Models = r.modelOrder
	schema.Errors = r.errors

	return schema
}

func (r *Resolver) resolveDatasource(schema *Schema) {
	if len(r.ast.Datasources) == 0 {
		return
	}
	ds := r.ast.Datasources[0]
	if len(r.ast.Datasources) > 1 {
		r.errorf(ds.Line, ds.Col, "multiple datasource blocks found; only one is supported")
	}
	out := Datasource{Name: ds.Name}
	for _, kv := range ds.Properties {
		switch kv.Key {
		case "provider":
			if kv.Value.Type == ValLiteral {
				out.Provider = kv.Value.Scalar
			}

		}
	}
	schema.Datasource = out
}

func (r *Resolver) registerEnums() {
	for _, decl := range r.ast.Enums {
		if _, exists := r.enums[decl.Name]; exists {
			r.errorf(decl.Line, decl.Col, "enum %q is declared more than once", decl.Name)
			continue
		}
		e := &Enum{Name: decl.Name}
		for _, v := range decl.Values {
			ev := EnumValue{Name: v.Name, DBName: v.Name}
			for _, a := range v.Attributes {
				if a.Name == "map" && len(a.Args) == 1 && a.Args[0].Value.Type == ValLiteral {
					ev.DBName = a.Args[0].Value.Scalar
				}
			}
			e.Values = append(e.Values, ev.Name)
			e.ValueMap = append(e.ValueMap, ev)
		}
		for _, a := range decl.Attributes {
			if a.Name == "map" && len(a.Args) == 1 && a.Args[0].Value.Type == ValLiteral {
				e.TableMapName = a.Args[0].Value.Scalar
			}
		}
		r.enums[decl.Name] = e
	}
}

func (r *Resolver) registerModels() {
	for _, decl := range r.ast.Models {
		if _, exists := r.models[decl.Name]; exists {
			r.errorf(decl.Line, decl.Col, "model %q is declared more than once", decl.Name)
			continue
		}
		m := &Model{Name: decl.Name, TableName: decl.Name}
		r.models[decl.Name] = m
		r.modelOrder = append(r.modelOrder, m)
	}
}

func (r *Resolver) resolveScalarFields() {
	for _, decl := range r.ast.Models {
		m, ok := r.models[decl.Name]
		if !ok {
			continue
		}
		for _, fd := range decl.Fields {
			if info, ok := scalarTypeInfo[fd.TypeName]; ok {
				m.ScalarFields = append(m.ScalarFields, r.buildScalarField(fd, info.GoType, info.SQLType, nil))
				continue
			}
			if enum, ok := r.enums[fd.TypeName]; ok {
				m.ScalarFields = append(m.ScalarFields, r.buildScalarField(fd, "string", "TEXT", enum))
				continue
			}
			if len(fd.TypeName) > 13 && fd.TypeName[:12] == "Unsupported(" {
				// E.g. Unsupported("geometry(Point)")
				sqlType := fd.TypeName[13 : len(fd.TypeName)-2]
				m.ScalarFields = append(m.ScalarFields, r.buildScalarField(fd, "any", sqlType, nil))
				continue
			}
			if _, isModel := r.models[fd.TypeName]; !isModel {
				r.errorf(fd.Line, fd.Col, "field %q has unknown type %q (not a scalar, enum, or model)", fd.Name, fd.TypeName)
			}
		}
	}
}

func (r *Resolver) buildScalarField(fd astFieldDecl, goType, sqlType string, enum *Enum) *ScalarField {
	sf := &ScalarField{
		Name:       fd.Name,
		Type:       fd.TypeName,
		ColName:    fd.Name,
		GoType:     goType,
		SQLType:    sqlType,
		IsArray:    fd.IsArray,
		Optional:   fd.Optional,
		EnumRef:    enum,
		Attributes: fd.Attributes,
	}
	if fd.IsArray {
		sf.GoType = "[]" + sf.GoType
	}
	if fd.Optional && !fd.IsArray {
		sf.GoType = "*" + sf.GoType
	}

	for _, a := range fd.Attributes {
		switch a.Name {
		case "id":
			sf.IsID = true
		case "unique":
			sf.IsUnique = true
		case "updatedAt":
			sf.IsUpdatedAt = true
		case "default":
			sf.Default = r.resolveDefaultValue(a, enum)
		case "map":
			if len(a.Args) == 1 && a.Args[0].Value.Type == ValLiteral {
				sf.ColName = a.Args[0].Value.Scalar
			}
		default:
			if isNativeDBType(a.Name) {
				sf.NativeType = &NativeType{
					Name: stripDBPrefix(a.Name),
					Args: stringifyArgs(a.Args),
				}
				if mapped := nativeTypeToSQL(sf.NativeType.Name); mapped != "" {
					sf.SQLType = mapped
				}
			}
		}
	}
	return sf
}

func isNativeDBType(attrName string) bool {
	return len(attrName) > 3 && attrName[:3] == "db."
}

func stripDBPrefix(attrName string) string { return attrName[3:] }

func stringifyArgs(args []Argument) []string {
	out := make([]string, 0, len(args))
	for _, a := range args {
		switch a.Value.Type {
		case ValLiteral:
			out = append(out, a.Value.Scalar)
		case ValIdent:
			out = append(out, a.Value.Scalar)
		}
	}
	return out
}

func nativeTypeToSQL(name string) string {
	switch name {
	case "VarChar", "Char":
		return "VARCHAR"
	case "Text":
		return "TEXT"
	case "Decimal", "Numeric":
		return "NUMERIC"
	case "Uuid":
		return "UUID"
	case "Timestamptz":
		return "TIMESTAMPTZ"
	case "Date":
		return "DATE"
	default:
		return ""
	}
}

func (r *Resolver) resolveDefaultValue(a Attribute, enum *Enum) *DefaultValue {
	if len(a.Args) == 0 {
		return nil
	}
	v := a.Args[0].Value

	switch v.Type {
	case ValFunc:
		if v.Scalar == "dbgenerated" {
			if len(v.Args) == 1 && v.Args[0].Value.Type == ValLiteral {
				dbExpr := v.Args[0].Value.Scalar
				return &DefaultValue{
					Value:        fmt.Sprintf("dbgenerated(%q)", dbExpr),
					Type:         "DBGenerated",
					Kind:         DefaultDBGenerated,
					DBExpression: dbExpr,
				}
			}
			return nil
		}
		funcCallStr := v.Scalar + "()"
		return &DefaultValue{
			Value:    funcCallStr,
			Type:     "Func",
			Kind:     DefaultFunc,
			FuncName: v.Scalar,
		}

	case ValIdent:
		if enum != nil {
			return &DefaultValue{
				Value:     v.Scalar,
				Type:      "EnumValue",
				Kind:      DefaultEnumValue,
				EnumValue: v.Scalar,
			}
		}
		return &DefaultValue{
			Value:   v.Scalar,
			Type:    "Literal",
			Kind:    DefaultLiteral,
			Literal: v.Scalar,
		}

	case ValLiteral:
		return &DefaultValue{
			Value:   v.Scalar,
			Type:    "Literal",
			Kind:    DefaultLiteral,
			Literal: v.Scalar,
		}
	}
	return nil
}

func (r *Resolver) resolveRelationFields() {
	for _, decl := range r.ast.Models {
		m, ok := r.models[decl.Name]
		if !ok {
			continue
		}
		for _, fd := range decl.Fields {
			if _, isScalar := scalarTypeInfo[fd.TypeName]; isScalar {
				continue
			}
			if _, isEnum := r.enums[fd.TypeName]; isEnum {
				continue
			}
			target, isModel := r.models[fd.TypeName]
			if !isModel {
				continue
			}
			m.RelationFields = append(m.RelationFields, r.buildRelationField(fd, m, target))
		}
	}
}

func (r *Resolver) buildRelationField(fd astFieldDecl, owner, target *Model) *RelationField {
	rf := &RelationField{
		Name:            fd.Name,
		Type:            fd.TypeName,
		TargetModelName: fd.TypeName,
		TargetModel:     target,
		Optional:        fd.Optional,
		IsArray:         fd.IsArray,
	}

	for _, a := range fd.Attributes {
		if a.Name != "relation" {
			r.errorf(a.Line, a.Col, "attribute '@%s' is not allowed on relation fields", a.Name)
		}
		for _, arg := range a.Args {
			switch {
			case arg.Name == "" && arg.Value.Type == ValLiteral:
				rf.RelationName = arg.Value.Scalar
			case arg.Name == "fields":
				rf.FKFields = r.resolveFieldRefs(owner, arg.Value, fd.Line, fd.Col)
			case arg.Name == "references":
				rf.RefFields = r.resolveFieldRefs(target, arg.Value, fd.Line, fd.Col)
			case arg.Name == "onDelete":
				rf.OnDelete = identOrLiteral(arg.Value)
			case arg.Name == "onUpdate":
				rf.OnUpdate = identOrLiteral(arg.Value)
			case arg.Name == "name":
				rf.RelationName = identOrLiteral(arg.Value)
			}
		}
	}

	if rf.RelationName == "" {
		rf.RelationName = synthesizeRelationName(owner.Name, target.Name)
	}

	return rf
}

func identOrLiteral(v Value) string {
	if v.Type == ValIdent || v.Type == ValLiteral {
		return v.Scalar
	}
	return ""
}

func synthesizeRelationName(a, b string) string {
	if a <= b {
		return a + "To" + b
	}
	return b + "To" + a
}

func (r *Resolver) resolveFieldRefs(m *Model, v Value, line, col int) []*ScalarField {
	if v.Type != ValArray {
		r.errorf(line, col, "expected an array of field names")
		return nil
	}
	var out []*ScalarField
	for _, item := range v.Array {
		if item.Type != ValIdent {
			r.errorf(line, col, "expected a field name")
			continue
		}
		sf := findScalarField(m, item.Scalar)
		if sf == nil {
			r.errorf(line, col, "relation references field %q which does not exist on model %q", item.Scalar, m.Name)
			continue
		}
		out = append(out, sf)
	}
	return out
}

func findScalarField(m *Model, name string) *ScalarField {
	for _, sf := range m.ScalarFields {
		if sf.Name == name {
			return sf
		}
	}
	return nil
}

func (r *Resolver) linkInverseRelations() {
	type key struct {
		a, b string
		name string
	}
	type entry struct {
		owner *Model
		field *RelationField
	}
	index := map[key][]entry{}

	for _, m := range r.modelOrder {
		for _, rf := range m.RelationFields {
			k := key{a: m.Name, b: rf.TargetModelName, name: rf.RelationName}
			index[k] = append(index[k], entry{owner: m, field: rf})
		}
	}

	seen := map[*RelationField]bool{}
	for _, m := range r.modelOrder {
		for _, rf := range m.RelationFields {
			if seen[rf] || rf.Inverse != nil {
				continue
			}
			reverseKey := key{a: rf.TargetModelName, b: m.Name, name: rf.RelationName}
			partners := index[reverseKey]

			var partner *RelationField
			for _, cand := range partners {
				if cand.field == rf {
					continue
				}
				if cand.field.Inverse == nil && !seen[cand.field] {
					partner = cand.field
					break
				}
			}

			if partner == nil {
				continue
			}

			rf.Inverse = partner
			partner.Inverse = rf
			seen[rf] = true
			seen[partner] = true
		}
	}
}

func (r *Resolver) inferRelationKinds() {
	decided := map[*RelationField]bool{}
	for _, m := range r.modelOrder {
		for _, rf := range m.RelationFields {
			if decided[rf] {
				continue
			}
			decided[rf] = true

			switch {
			case rf.IsArray && rf.Inverse != nil && rf.Inverse.IsArray:
				rf.Kind = RelManyToMany
				rf.JoinTableName = "_" + rf.RelationName
				decided[rf.Inverse] = true
				rf.Inverse.Kind = RelManyToMany
				rf.Inverse.JoinTableName = rf.JoinTableName

			case rf.IsArray:
				rf.Kind = RelOneToMany
				if rf.Inverse != nil {
					decided[rf.Inverse] = true
					rf.Inverse.Kind = RelManyToOne
				}

			case len(rf.FKFields) > 0:
				if rf.Inverse != nil && rf.Inverse.IsArray {
					rf.Kind = RelManyToOne
					decided[rf.Inverse] = true
					rf.Inverse.Kind = RelOneToMany
				} else {
					rf.Kind = RelOneToOne
					if rf.Inverse != nil {
						decided[rf.Inverse] = true
						rf.Inverse.Kind = RelOneToOne
					}
				}

			default:
				rf.Kind = RelOneToOne
				if rf.Inverse != nil {
					decided[rf.Inverse] = true
					rf.Inverse.Kind = RelOneToOne
				}
			}
		}
	}
}

func (r *Resolver) resolveModelAttributes() {
	for _, decl := range r.ast.Models {
		m, ok := r.models[decl.Name]
		if !ok {
			continue
		}
		m.Attributes = decl.Attributes

		for _, a := range decl.Attributes {
			switch a.Name {
			case "map":
				if len(a.Args) == 1 && a.Args[0].Value.Type == ValLiteral {
					m.TableName = a.Args[0].Value.Scalar
				}
			case "id":
				m.CompositePK = r.resolveFieldNameList(m, a, "@@id")
			case "unique":
				fields, name := r.resolveFieldNameListWithName(m, a, "@@unique")
				if fields != nil {
					m.CompositeUnique = append(m.CompositeUnique, UniqueConstraint{
						Fields: fields,
						Name:   name,
					})
				}
			case "index":
				idx := Index{}
				for _, arg := range a.Args {
					if arg.Name == "" && arg.Value.Type == ValArray {
						idx.Fields = r.fieldNamesFromArray(m, arg.Value, a.Line, a.Col)
					}
					if arg.Name == "name" {
						idx.Name = identOrLiteral(arg.Value)
					}
				}
				if len(idx.Fields) > 0 {
					m.Indexes = append(m.Indexes, idx)
				}
			}
		}
	}
}

func (r *Resolver) resolveFieldNameList(m *Model, a Attribute, attrLabel string) []string {
	if len(a.Args) == 0 || a.Args[0].Value.Type != ValArray {
		r.errorf(a.Line, a.Col, "%s expects an array of field names", attrLabel)
		return nil
	}
	return r.fieldNamesFromArray(m, a.Args[0].Value, a.Line, a.Col)
}

func (r *Resolver) resolveFieldNameListWithName(m *Model, a Attribute, attrLabel string) ([]string, string) {
	var fields []string
	var name string
	foundArray := false
	for _, arg := range a.Args {
		if arg.Name == "" && arg.Value.Type == ValArray {
			fields = r.fieldNamesFromArray(m, arg.Value, a.Line, a.Col)
			foundArray = true
		}
		if arg.Name == "name" {
			name = identOrLiteral(arg.Value)
		}
	}
	if !foundArray {
		r.errorf(a.Line, a.Col, "%s expects an array of field names", attrLabel)
	}
	return fields, name
}

func (r *Resolver) fieldNamesFromArray(m *Model, v Value, line, col int) []string {
	var out []string
	for _, item := range v.Array {
		if item.Type != ValIdent {
			continue
		}
		if findScalarField(m, item.Scalar) == nil {
			r.errorf(line, col, "field %q does not exist on model %q", item.Scalar, m.Name)
			continue
		}
		out = append(out, item.Scalar)
	}
	return out
}
