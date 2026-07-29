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

func hasFieldWhere(m *schema.Model, pred func(*schema.ScalarField) bool) bool {
	return slices.ContainsFunc(m.ScalarFields, pred)
}

func hasJsonField(m *schema.Model) bool {
	return hasFieldWhere(m, func(sf *schema.ScalarField) bool {
		return sf.Type == "Json" || strings.Contains(sf.GoType, "json.RawMessage")
	})
}
func hasJsonAnywhere(sch schema.Schema) bool {
	return slices.ContainsFunc(sch.Models, hasJsonField)
}
func hasTimeField(m *schema.Model) bool {
	return hasFieldWhere(m, func(sf *schema.ScalarField) bool {
		return sf.Type == "DateTime" || strings.Contains(sf.GoType, "time.Time")
	})
}
func hasTimeAnywhere(sch schema.Schema) bool {
	return slices.ContainsFunc(sch.Models, hasTimeField)
}
func isKnownDefaultFunc(funcName string) bool {
	val, ok := DEFAULT_FUNCS[funcName]
	return ok && val != ""
}

func defaultFuncCall(funcName string) string {
	return DEFAULT_FUNCS[funcName]
}
func hasUuidField(m *schema.Model) bool {
	return hasFieldWhere(m, func(sf *schema.ScalarField) bool {
		return sf.NativeType != nil && sf.NativeType.Name == "Uuid"
	})
}
func hasUuidAnywhere(sch schema.Schema) bool {
	return slices.ContainsFunc(sch.Models, hasUuidField)
}
func hasFloatField(m *schema.Model) bool {
	return hasFieldWhere(m, func(sf *schema.ScalarField) bool {
		return sf.Type == "Float"
	})
}
func hasFloatAnywhere(sch schema.Schema) bool {
	return slices.ContainsFunc(sch.Models, hasFloatField)
}
func hasDecimalField(m *schema.Model) bool {
	return hasFieldWhere(m, func(sf *schema.ScalarField) bool {
		return sf.Type == "Decimal"
	})
}
func hasDecimalAnywhere(sch schema.Schema) bool {
	return slices.ContainsFunc(sch.Models, hasDecimalField)
}
func hasNetField(m *schema.Model) bool {
	return hasFieldWhere(m, func(sf *schema.ScalarField) bool {
		return sf.NativeType != nil && sf.NativeType.Name == "Inet"
	})
}
func hasNetAnywhere(sch schema.Schema) bool {
	return slices.ContainsFunc(sch.Models, hasNetField)
}
func hasHstoreField(m *schema.Model) bool {
	return hasFieldWhere(m, func(sf *schema.ScalarField) bool {
		return strings.TrimPrefix(sf.GoType, "*") == "map[string]*string"
	})
}
func hasHstoreAnywhere(sch schema.Schema) bool {
	return slices.ContainsFunc(sch.Models, hasHstoreField)
}
func hstoreExpr(goType string, expr string) string {
	if strings.TrimPrefix(goType, "*") == "map[string]*string" {
		return "ToHstore(" + expr + ")"
	}
	return expr
}
