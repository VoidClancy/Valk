package schema

import (
	"testing"
)

func TestParseBasicModel(t *testing.T) {
	input := `
	// This is a comment
	model User {
		id    Int    @id @default(autoincrement())
		email String @unique
		name  String?
		role  Role   @default(USER)
		
		@@unique([email])
		@@map("users")
	}

	datasource db {
		provider = "postgresql"
		url      = env("DATABASE_URL")
	}

	enum Role {
		USER
		ADMIN
	}

	model Post {
		id        String   @id @default(cuid())
		title     String
		published Boolean  @default(false)
		author    User     @relation(fields: [authorId], references: [id])
		authorId  Int
		tags      String[]
		
		@@index([title])
	}
	`

	schema, _ := ParseSchema(input)

	if len(schema.Models) != 2 {
		t.Fatalf("expected 2 models, got %d", len(schema.Models))
	}

	// Verify User model
	user := schema.Models[0]
	if user.Name != "User" {
		t.Errorf("expected model name User, got %s", user.Name)
	}
	if len(user.ScalarFields) != 4 {
		t.Errorf("expected 4 scalar fields in User, got %d", len(user.ScalarFields))
	}
	if len(user.Attributes) != 2 {
		t.Errorf("expected 2 block attributes in User, got %d", len(user.Attributes))
	}

	// Verify User fields
	idField := user.ScalarFields[0]
	if idField.Name != "id" || idField.Type != "Int" || idField.IsArray || idField.Optional {
		t.Errorf("id field parsed incorrectly: %+v", idField)
	}
	if len(idField.Attributes) != 2 {
		t.Errorf("expected 2 attributes on id, got %d", len(idField.Attributes))
	}

	nameField := user.ScalarFields[2]
	if nameField.Name != "name" || nameField.Type != "String" || !nameField.Optional {
		t.Errorf("name field parsed incorrectly: %+v", nameField)
	}

	roleField := user.ScalarFields[3]
	if roleField.Name != "role" || roleField.Type != "Role" || len(roleField.Attributes) != 1 {
		t.Errorf("role field parsed incorrectly: %+v", roleField)
	}
	defaultAttr := roleField.Attributes[0]
	if defaultAttr.Name != "default" || len(defaultAttr.Args) != 1 {
		t.Errorf("default attr parsed incorrectly: %+v", defaultAttr)
	}
	if defaultAttr.Args[0].Value.Type != ValIdent || defaultAttr.Args[0].Value.Scalar != "USER" {
		t.Errorf("expected USER arg, got %+v", defaultAttr.Args[0].Value)
	}

	// Verify Block attributes in User
	uniqueBlock := user.Attributes[0]
	if uniqueBlock.Name != "unique" || len(uniqueBlock.Args) != 1 {
		t.Errorf("unique block attribute parsed incorrectly: %+v", uniqueBlock)
	}
	val := uniqueBlock.Args[0].Value
	if val.Type != ValArray || len(val.Array) != 1 || val.Array[0].Type != ValIdent || val.Array[0].Scalar != "email" {
		t.Errorf("unique block arguments parsed incorrectly: %+v", val)
	}

	mapBlock := user.Attributes[1]
	if mapBlock.Name != "map" || len(mapBlock.Args) != 1 {
		t.Errorf("map block attribute parsed incorrectly: %+v", mapBlock)
	}
	mapVal := mapBlock.Args[0].Value
	if mapVal.Type != ValLiteral || mapVal.Scalar != "users" {
		t.Errorf("map block value parsed incorrectly: %+v", mapVal)
	}

	// Verify Post model
	post := schema.Models[1]
	if post.Name != "Post" {
		t.Errorf("expected model name Post, got %s", post.Name)
	}
	if len(post.ScalarFields) != 5 {
		t.Errorf("expected 5 scalar fields in Post, got %d", len(post.ScalarFields))
	}
	if len(post.RelationFields) != 1 {
		t.Errorf("expected 1 relation field in Post, got %d", len(post.RelationFields))
	}

	tagsField := post.ScalarFields[4]
	if tagsField.Name != "tags" || tagsField.Type != "String" || !tagsField.IsArray || tagsField.Optional {
		t.Errorf("tags field parsed incorrectly: %+v", tagsField)
	}
}

func TestUnsupportedType(t *testing.T) {
	input := `
	model Spatial {
		id       Int                           @id
		location Unsupported("geometry(Point)")?
	}
	`
	schema, _ := ParseSchema(input)

	if len(schema.Models) != 1 {
		t.Fatalf("expected 1 model, got %d", len(schema.Models))
	}

	spatial := schema.Models[0]
	if len(spatial.ScalarFields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(spatial.ScalarFields))
	}

	locField := spatial.ScalarFields[1]
	if locField.Name != "location" {
		t.Errorf("expected location, got %s", locField.Name)
	}
	expectedType := `Unsupported("geometry(Point)")`
	if locField.Type != expectedType {
		t.Errorf("expected type %s, got %s", expectedType, locField.Type)
	}
	if !locField.Optional {
		t.Errorf("expected location to be optional")
	}
}

func TestNativeTypes(t *testing.T) {
	input := `
	model Native {
		id   Int    @id
		desc String @db.VarChar(255)
	}
	`
	schema, _ := ParseSchema(input)

	if len(schema.Models) != 1 {
		t.Fatalf("expected 1 model, got %d", len(schema.Models))
	}

	model := schema.Models[0]
	descField := model.ScalarFields[1]
	if len(descField.Attributes) != 1 {
		t.Fatalf("expected 1 attribute, got %d", len(descField.Attributes))
	}

	attr := descField.Attributes[0]
	if attr.Name != "db.VarChar" {
		t.Errorf("expected db.VarChar, got %s", attr.Name)
	}

	if len(attr.Args) != 1 || attr.Args[0].Value.Type != ValLiteral || attr.Args[0].Value.Scalar != "255" {
		t.Errorf("expected single literal argument 255, got %+v", attr.Args)
	}
}

func TestConditionalIndexes(t *testing.T) {
	input := `
	model User {
		id        Int       @id
		email     String
		status    String
		userId    Int
		createdAt DateTime
		deletedAt DateTime?

		@@index([email], where: deletedAt == null)
		@@index([status], where: status == "active")
		@@index([userId, createdAt], where: deletedAt == null)
	}
	`
	schema, _ := ParseSchema(input)

	if len(schema.Models) != 1 {
		t.Fatalf("expected 1 model, got %d", len(schema.Models))
	}

	model := schema.Models[0]
	if len(model.Attributes) != 3 {
		t.Fatalf("expected 3 model attributes, got %d", len(model.Attributes))
	}

	// 1st index
	idx1 := model.Attributes[0]
	if idx1.Name != "index" || len(idx1.Args) != 2 {
		t.Fatalf("expected index with 2 arguments, got %+v", idx1)
	}
	if idx1.Args[0].Name != "" || idx1.Args[0].Value.Type != ValArray {
		t.Errorf("expected positional array, got %+v", idx1.Args[0])
	}
	whereArg1 := idx1.Args[1]
	if whereArg1.Name != "where" || whereArg1.Value.Type != ValBinary || whereArg1.Value.Scalar != "==" {
		t.Errorf("expected binary where expression, got %+v", whereArg1)
	}
	if whereArg1.Value.Left.Type != ValIdent || whereArg1.Value.Left.Scalar != "deletedAt" {
		t.Errorf("expected Left to be deletedAt, got %+v", whereArg1.Value.Left)
	}
	if whereArg1.Value.Right.Type != ValIdent || whereArg1.Value.Right.Scalar != "null" {
		t.Errorf("expected Right to be null, got %+v", whereArg1.Value.Right)
	}

	// 2nd index
	idx2 := model.Attributes[1]
	whereArg2 := idx2.Args[1]
	if whereArg2.Name != "where" || whereArg2.Value.Type != ValBinary || whereArg2.Value.Scalar != "==" {
		t.Errorf("expected binary where expression, got %+v", whereArg2)
	}
	if whereArg2.Value.Left.Type != ValIdent || whereArg2.Value.Left.Scalar != "status" {
		t.Errorf("expected Left to be status, got %+v", whereArg2.Value.Left)
	}
	if whereArg2.Value.Right.Type != ValLiteral || whereArg2.Value.Right.Scalar != "active" {
		t.Errorf("expected Right to be active, got %+v", whereArg2.Value.Right)
	}
}

func TestParserErrorHandlingAndRecovery(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedModel string
		expectErrors  int
		errorSubstr   string
	}{
		{
			name: "unrecognized lexer character outside comment/string",
			input: `
			model User {
				id Int @id
			}
			# invalid character here
			model Post {
				id String @id
			}
			`,
			expectedModel: "User",
			expectErrors:  4,
		},
		{
			name: "unterminated string literal",
			input: `
			model User {
				id Int @id @default("unterminated)
			}
			model Post {
				id String @id
			}
			`,
			expectedModel: "Post",
			expectErrors:  1,
			errorSubstr:   "unterminated string literal",
		},
		{
			name: "parser syntax error in model block",
			input: `
			model BadModel {
				id Int @id @invalidAttribute(
			}
			model GoodModel {
				id String @id
			}
			`,
			expectedModel: "GoodModel",
			expectErrors:  1,
			errorSubstr:   "unterminated \"(\"",
		},
		{
			name: "schema valk test",
			input: `model Clancy {
    id String @id
    email String @unique
    
}`,
			expectedModel: "Clancy",
			expectErrors:  0,
		},
		{
			name: "unterminated brackets",
			input: `
			model BadModel {
				id Int @id
				roles String[
			}
			model GoodModel {
				id String @id
			}
			`,
			expectedModel: "GoodModel",
			expectErrors:  1,
			errorSubstr:   "unterminated \"[\"",
		},
		{
			name: "missing attribute name",
			input: `
			model BadModel {
				id Int @id @
			}
			model GoodModel {
				id String @id
			}
			`,
			expectedModel: "GoodModel",
			expectErrors:  1,
			errorSubstr:   "missing attribute name after \"@\"",
		},
		{
			name: "unterminated model brace",
			input: `
			model BadModel {
				id Int @id
			
			model GoodModel {
				id String @id
			}
			`,
			expectedModel: "GoodModel",
			expectErrors:  1,
			errorSubstr:   "unterminated \"{\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema, _ := ParseSchema(tt.input)

			if len(schema.Errors) != tt.expectErrors {
				t.Errorf("expected %d errors, got %d. Errors: %v", tt.expectErrors, len(schema.Errors), schema.Errors)
			}

			if tt.errorSubstr != "" && len(schema.Errors) > 0 {
				foundSubstr := false
				for _, err := range schema.Errors {
					if contains(err.Message, tt.errorSubstr) {
						foundSubstr = true
						break
					}
				}
				if !foundSubstr {
					t.Errorf("expected error message to contain %q, but got errors: %v", tt.errorSubstr, schema.Errors)
				}
			}

			// Verify we could still parse the valid model in recovery
			found := false
			for _, m := range schema.Models {
				if m.Name == tt.expectedModel {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("failed to recover and parse expected model %q. parsed: %+v", tt.expectedModel, schema.Models)
			}
		})
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
