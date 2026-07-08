package generator

import (
	"strings"
	"testing"
	"github.com/voidclancy/valk/schema"
)

func TestGenerateClient_NativeDBConstraints(t *testing.T) {
	sch := schema.Schema{
		Models: []*schema.Model{
			{
				Name:      "Item",
				TableName: "items",
				ScalarFields: []*schema.ScalarField{
					{
						Name:   "id",
						Type:   "String",
						GoType: "string",
						IsID:   true,
					},
					{
						Name:   "code",
						Type:   "String",
						GoType: "string",
						NativeType: &schema.NativeType{
							Name: "VarChar",
							Args: []string{"8"},
						},
					},
					{
						Name:   "count",
						Type:   "Int",
						GoType: "int32",
						NativeType: &schema.NativeType{
							Name: "SmallInt",
						},
					},
				},
			},
		},
	}

	outputs, err := GenerateClient(sch, "valk", "github.com/voidclancy/valk", "", "", nil)
	if err != nil {
		t.Fatalf("failed to generate client: %v", err)
	}

	itemCode, ok := outputs["item.go"]
	if !ok {
		t.Fatal("expected item.go in outputs")
	}

	// Verify length checks are generated
	if !strings.Contains(itemCode, "utf8.RuneCountInString(input.Code) > 8") {
		t.Errorf("expected generated code to contain VarChar limit check, got:\n%s", itemCode)
	}

	// Verify SmallInt range checks are generated
	if !strings.Contains(itemCode, "input.Count < -32768 || input.Count > 32767") {
		t.Errorf("expected generated code to contain SmallInt limit check, got:\n%s", itemCode)
	}
}
