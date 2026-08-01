package generator

import (
	"slices"
	"strings"

	"github.com/voidclancy/phi/schema"
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

func hasModelType(m *schema.Model, targetTypes ...string) bool {
	for _, sf := range m.ScalarFields {
		for _, t := range targetTypes {
			switch t {
			case "Array":
				if sf.IsArray {
					return true
				}
			case "Hstore":
				if strings.Contains(sf.GoType, "map[string]*string") {
					return true
				}
			case "DateTime", "Time":
				if sf.Type == "DateTime" || strings.Contains(sf.GoType, "time.Time") {
					return true
				}
			case "Json":
				if (sf.Type == "Json" || strings.Contains(sf.GoType, "json.RawMessage")) && !sf.IsArray {
					return true
				}
			default:
				if sf.Type == t || (sf.NativeType != nil && sf.NativeType.Name == t) {
					return true
				}
			}
		}
	}
	return false
}

func hasType(sch schema.Schema, targetTypes ...string) bool {
	for _, m := range sch.Models {
		if hasModelType(m, targetTypes...) {
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

func hasDefaultFunc(sch schema.Schema, names ...string) bool {
	for _, m := range sch.Models {
		for _, sf := range m.ScalarFields {
			if sf.Default != nil && sf.Default.Kind == schema.DefaultFunc {
				if slices.Contains(names, sf.Default.FuncName) {
					return true
				}
			}
			if sf.IsID && sf.GoType == "string" && sf.Default == nil && slices.Contains(names, "cuid") {
				return true
			}
		}
	}
	return false
}
func hstoreExpr(goType string, expr string) string {
	if strings.TrimPrefix(goType, "*") == "map[string]*string" {
		return "ToHstore(" + expr + ")"
	}
	return expr
}
