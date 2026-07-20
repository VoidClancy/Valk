package generator

import (
	"slices"
	"strings"

	"github.com/voidclancy/valk/schema"
)

var DEFAULT_FUNCS = map[string]string{
	"autoincrement": "",
	"cuid":          "generateCUID()",
	"cuid(1)":       "generateCUID()",
	"cuid(2)":       "generateCUID2()",
	"uuid":          "generateUUID()",
	"uuid(4)":       "generateUUID()",
	"uuid(7)":       "generateUUID7()",
	"ulid":          "generateULID()",
	"nanoid":        "generateNanoID()",
	"now":           "time.Now()",
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	if s == strings.ToUpper(s) {
		s = strings.ToLower(s)
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func lowercase(s string) string {
	if s == "" {
		return ""
	}
	if s == strings.ToUpper(s) {
		s = strings.ToLower(s)
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// returns the relation name if this scalar field is a FK for a relation on the model, empty string if not
func fkForRelation(model *schema.Model, field *schema.ScalarField) string {
	for _, rel := range model.RelationFields {
		for _, fk := range rel.FKFields {
			if fk.Name == field.Name {
				return rel.Name
			}
		}
	}
	return ""
}

func fieldPredType(f *schema.ScalarField, parentPkg string) string {
	if f.EnumRef != nil {
		if f.IsArray {
			return "[]" + parentPkg + "." + f.EnumRef.Name + "Type"
		}
		return parentPkg + "." + f.EnumRef.Name + "Type"
	}
	t := f.GoType
	if f.Optional {
		t = strings.TrimPrefix(t, "*")
	}
	return t
}

func hasJsonField(m *schema.Model) bool {
	for _, sf := range m.ScalarFields {
		if sf.Type == "Json" || strings.Contains(sf.GoType, "json.RawMessage") {
			return true
		}
	}
	return false
}
func hasTimeField(m *schema.Model) bool {
	for _, sf := range m.ScalarFields {
		if sf.Type == "DateTime" || strings.Contains(sf.GoType, "time.Time") {
			return true
		}
	}
	return false
}
func isKnownDefaultFunc(funcName string) bool {
	val, ok := DEFAULT_FUNCS[funcName]
	return ok && val != ""
}

func defaultFuncCall(funcName string) string {
	return DEFAULT_FUNCS[funcName]
}
func hasStringField(m *schema.Model) bool {
	for _, sf := range m.ScalarFields {
		if sf.GoType == "string" || strings.Contains(sf.GoType, "string") {
			return true
		}
	}
	return false
}
func hasNetField(m *schema.Model) bool {
	for _, sf := range m.ScalarFields {
		if sf.NativeType != nil && sf.NativeType.Name == "Inet" {
			return true
		}
	}
	return false
}
func hasHstoreField(m *schema.Model) bool {
	for _, sf := range m.ScalarFields {
		if strings.TrimPrefix(sf.GoType, "*") == "map[string]*string" {
			return true
		}
	}
	return false
}
func hasHstoreAnywhere(sch schema.Schema) bool {
	return slices.ContainsFunc(sch.Models, hasHstoreField)
}
func hasNetAnywhere(sch schema.Schema) bool {
	return slices.ContainsFunc(sch.Models, hasNetField)
}
func hstoreExpr(goType string, expr string) string {
	if strings.TrimPrefix(goType, "*") == "map[string]*string" {
		return "ToHstore(" + expr + ")"
	}
	return expr
}
