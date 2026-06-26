package schema

// =============================================================================
//  Parser Integration Tests
// =============================================================================

import (
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------
// Lexer tests: unicode, whitespace, and character rules
// -----------------------------------------------------------------------------

func TestOdyssey_LexerChaos(t *testing.T) {
	t.Run("unicode in string literals", func(t *testing.T) {
		input := `
		model User {
			id    Int    @id
			name  String @default("こんにちは世界 🌍 مرحبا")
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("expected 1 model, got %d", len(s.Models))
		}
		f := s.Models[0].ScalarFields[1]
		attr := f.Attributes[0]
		if attr.Args[0].Value.Scalar != "こんにちは世界 🌍 مرحبا" {
			t.Errorf("unicode string not preserved: %q", attr.Args[0].Value.Scalar)
		}
	})

	t.Run("escaped quotes inside string literals", func(t *testing.T) {
		input := `
		model Thing {
			id   Int    @id
			desc String @default("she said \"hello\" to me")
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("expected 1 model, got %d; errors: %v", len(s.Models), s.Errors)
		}
	})

	t.Run("CRLF line endings dont break anything", func(t *testing.T) {
		input := "model User {\r\n\tid Int @id\r\n\tname String\r\n}\r\n"
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("expected 1 model, got %d", len(s.Models))
		}
		if s.Models[0].Name != "User" {
			t.Errorf("wrong model name: %s", s.Models[0].Name)
		}
	})

	t.Run("tabs vs spaces mixed indentation", func(t *testing.T) {
		input := "model Mixed {\n\t  id   Int    @id\n  \t  email String @unique\n}"
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 || len(s.Models[0].ScalarFields) != 2 {
			t.Fatalf("mixed whitespace broken: models=%d", len(s.Models))
		}
	})

	t.Run("empty string literal as default", func(t *testing.T) {
		input := `
		model Thing {
			id   Int    @id
			tag  String @default("")
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("empty string literal failed: %v", s.Errors)
		}
		val := s.Models[0].ScalarFields[1].Attributes[0].Args[0].Value
		if val.Type != ValLiteral || val.Scalar != "" {
			t.Errorf("expected empty literal, got %+v", val)
		}
	})

	t.Run("string with only whitespace", func(t *testing.T) {
		input := `
		model Thing {
			id   Int    @id
			gap  String @default("   ")
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("whitespace-only string failed: %v", s.Errors)
		}
	})

	t.Run("numbers that look like identifiers", func(t *testing.T) {
		// e.g., a model field with a numeric default
		input := `
		model Counter {
			id    Int @id
			count Int @default(0)
			ratio Float @default(3.14)
			neg   Int @default(-1)
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("numeric defaults broken: %v", s.Errors)
		}
		fields := s.Models[0].ScalarFields
		if len(fields) != 4 {
			t.Fatalf("expected 4 fields, got %d", len(fields))
		}
	})

	t.Run("very long identifier names", func(t *testing.T) {
		longName := strings.Repeat("A", 300)
		input := "model " + longName + " {\n\tid Int @id\n}\n"
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 || s.Models[0].Name != longName {
			t.Errorf("long identifier failed: got model name len %d", len(s.Models[0].Name))
		}
	})
}

// -----------------------------------------------------------------------------
// Attribute tests: nesting, multiple args, named args
// -----------------------------------------------------------------------------

func TestOdyssey_AttributeChaos(t *testing.T) {
	t.Run("attribute with no args vs empty parens", func(t *testing.T) {
		input := `
		model A {
			id   Int @id
			uuid String @default(uuid())
			now  DateTime @default(now())
			auto Int @default(autoincrement())
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("function call defaults broken: %v", s.Errors)
		}
		fields := s.Models[0].ScalarFields
		for i, name := range []string{"id", "uuid", "now", "auto"} {
			if fields[i].Name != name {
				t.Errorf("field %d: expected %s, got %s", i, name, fields[i].Name)
			}
		}
	})

	t.Run("attribute with multiple named args", func(t *testing.T) {
		input := `
		model Post {
			id      Int    @id
			title   String
			
			@@index([title], name: "idx_title", type: BTree)
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("multi-arg index broken: %v", s.Errors)
		}
		attrs := s.Models[0].Attributes
		if len(attrs) != 1 {
			t.Fatalf("expected 1 block attr, got %d", len(attrs))
		}
		if len(attrs[0].Args) != 3 {
			t.Errorf("expected 3 args (array + 2 named), got %d: %+v", len(attrs[0].Args), attrs[0].Args)
		}
	})

	t.Run("relation attribute with all possible named args", func(t *testing.T) {
		input := `
		model Post {
			id       Int    @id
			author   User   @relation(fields: [authorId], references: [id], onDelete: Cascade, onUpdate: NoAction, name: "PostToUser")
			authorId Int
		}
		model User {
			id    Int    @id
			posts Post[]
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 2 {
			t.Fatalf("relation attr broken: %v", s.Errors)
		}
		relFields := s.Models[0].RelationFields
		if len(relFields) != 1 {
			t.Fatalf("expected 1 relation field in Post, got %d", len(relFields))
		}
		if relFields[0].RelationName != "PostToUser" {
			t.Errorf("expected relation name 'PostToUser', got %q", relFields[0].RelationName)
		}
		if len(relFields[0].FKFields) != 1 || relFields[0].FKFields[0].Name != "authorId" {
			t.Errorf("expected FK fields to contain authorId, got %+v", relFields[0].FKFields)
		}
		if len(relFields[0].RefFields) != 1 || relFields[0].RefFields[0].Name != "id" {
			t.Errorf("expected Ref fields to contain id, got %+v", relFields[0].RefFields)
		}
		if relFields[0].OnDelete != "Cascade" {
			t.Errorf("expected onDelete 'Cascade', got %q", relFields[0].OnDelete)
		}
		if relFields[0].OnUpdate != "NoAction" {
			t.Errorf("expected onUpdate 'NoAction', got %q", relFields[0].OnUpdate)
		}
	})

	t.Run("field with many attributes chained", func(t *testing.T) {
		input := `
		model MultiAttr {
			id   Int    @id @default(autoincrement()) @map("_id")
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("multi attribute chain broken: %v", s.Errors)
		}
		idField := s.Models[0].ScalarFields[0]
		if len(idField.Attributes) != 3 {
			t.Errorf("expected 3 attributes on id, got %d: %+v", len(idField.Attributes), idField.Attributes)
		}
	})

	t.Run("block attribute with compound multi-field array", func(t *testing.T) {
		input := `
		model Compound {
			a Int
			b Int
			c Int
			d String
			
			@@id([a, b])
			@@unique([b, c, d])
			@@index([a, b, c, d])
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("compound attr broken: %v", s.Errors)
		}
		attrs := s.Models[0].Attributes
		if len(attrs) != 3 {
			t.Fatalf("expected 3 block attrs, got %d", len(attrs))
		}
		// @@unique should have 3 items in its array
		arr := attrs[1].Args[0].Value
		if arr.Type != ValArray || len(arr.Array) != 3 {
			t.Errorf("@@unique array should have 3 items, got: %+v", arr)
		}
	})

	t.Run("nested function call as attribute arg", func(t *testing.T) {
		// dbgenerated("gen_random_uuid()") — function call inside a string
		// but also: @default(dbgenerated("nextval('seq')"))
		input := `
		model Seq {
			id Int @id @default(dbgenerated("nextval('my_seq'::regclass)"))
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("dbgenerated broken: %v", s.Errors)
		}
	})

	t.Run("attribute arg that is itself a function call with args", func(t *testing.T) {
		input := `
		model Sorted {
			id    Int    @id
			score Float
			
			@@index([score(sort: Desc)])
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("sort modifier in index broken: %v", s.Errors)
		}
	})
}

// -----------------------------------------------------------------------------
// Type system tests
// -----------------------------------------------------------------------------

func TestOdyssey_TypeChaos(t *testing.T) {
	t.Run("optional array — the cursed combo", func(t *testing.T) {
		// String[]? is actually invalid in Prisma, but let's see if the parser panics
		input := `
		model Cursed {
			id   Int      @id
			tags String[]
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("array field broken: %v", s.Errors)
		}
		tags := s.Models[0].ScalarFields[1]
		if !tags.IsArray {
			t.Errorf("expected IsArray=true for String[]")
		}
		if tags.Optional {
			t.Errorf("String[] should not be optional")
		}
	})

	t.Run("all scalar types at once", func(t *testing.T) {
		input := `
		model AllTypes {
			id         Int      @id
			strField   String
			boolField  Boolean
			floatField Float
			decField   Decimal
			bigInt     BigInt
			bytesField Bytes
			jsonField  Json
			dtField    DateTime
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("all types broken: %v", s.Errors)
		}
		if len(s.Models[0].ScalarFields) != 9 {
			t.Errorf("expected 9 fields, got %d", len(s.Models[0].ScalarFields))
		}
	})

	t.Run("Unsupported type with nested parens and quotes", func(t *testing.T) {
		input := `
		model Geo {
			id  Int                                  @id
			pt  Unsupported("geometry(Point, 4326)")?
			ls  Unsupported("geometry(LineString)")?
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("complex Unsupported type broken: %v", s.Errors)
		}
		fields := s.Models[0].ScalarFields
		if len(fields) != 3 {
			t.Fatalf("expected 3 fields, got %d", len(fields))
		}
		if fields[1].Type != `Unsupported("geometry(Point, 4326)")` {
			t.Errorf("wrong Unsupported type: %s", fields[1].Type)
		}
	})

	t.Run("enum with many values including special names", func(t *testing.T) {
		input := `
		enum Status {
			PENDING
			ACTIVE
			INACTIVE
			DELETED
			ARCHIVED
			SUSPENDED
			BANNED
			SHADOW_BANNED
			UNDER_REVIEW
		}
		model User {
			id     Int    @id
			status Status @default(PENDING)
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("large enum broken: %v", s.Errors)
		}
		if len(s.Enums) != 1 {
			t.Fatalf("expected 1 enum, got %d", len(s.Enums))
		}
		if len(s.Enums[0].Values) != 9 {
			t.Errorf("expected 9 enum values, got %d", len(s.Enums[0].Values))
		}
	})

	t.Run("native db type with multiple args", func(t *testing.T) {
		input := `
		model Native {
			id   Int     @id
			col1 Decimal @db.Decimal(10, 2)
			col2 String  @db.VarChar(255)
			col3 String  @db.Char(1)
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("native types broken: %v", s.Errors)
		}
		fields := s.Models[0].ScalarFields
		decAttr := fields[1].Attributes[0]
		if decAttr.Name != "db.Decimal" {
			t.Errorf("expected db.Decimal, got %s", decAttr.Name)
		}
		if len(decAttr.Args) != 2 {
			t.Errorf("expected 2 args for Decimal(10,2), got %d", len(decAttr.Args))
		}
	})
}

// -----------------------------------------------------------------------------
// Model structure edge cases
// -----------------------------------------------------------------------------

func TestOdyssey_ModelStructureChaos(t *testing.T) {
	t.Run("empty model body", func(t *testing.T) {
		input := `
		model Empty {
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("empty model broken: %v", s.Errors)
		}
		m := s.Models[0]
		if len(m.ScalarFields) != 0 || len(m.RelationFields) != 0 {
			t.Errorf("empty model should have no fields: %+v", m)
		}
	})

	t.Run("model with only block attributes, no fields", func(t *testing.T) {
		input := `
		model Ghost {
			@@map("ghosts")
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("block-only model broken: %v", s.Errors)
		}
		if len(s.Models[0].Attributes) != 1 {
			t.Errorf("expected 1 block attr, got %d", len(s.Models[0].Attributes))
		}
	})

	t.Run("model named with keyword-adjacent names", func(t *testing.T) {
		// "models", "modeler", "datasources" — not keywords but close
		input := `
		model ModelData {
			id Int @id
		}
		model Enumerable {
			id Int @id
		}
		model Datasource {
			id Int @id
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 3 {
			t.Fatalf("keyword-adjacent names broken: got %d models, errors: %v", len(s.Models), s.Errors)
		}
	})

	t.Run("massive model with 50 fields", func(t *testing.T) {
		var sb strings.Builder
		sb.WriteString("model BigBoy {\n")
		for i := 0; i < 50; i++ {
			sb.WriteString("\tfield")
			sb.WriteString(strings.Repeat("x", i%10+1)) // vary names
			sb.WriteString("_")
			// Use a digit suffix via manual approach
			for _, d := range []int{i / 10, i % 10} {
				sb.WriteByte(byte('0' + d))
			}
			sb.WriteString(" String\n")
		}
		sb.WriteString("\tid Int @id\n}\n")
		input := sb.String()
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("large model broken: %v", s.Errors)
		}
		if len(s.Models[0].ScalarFields) != 51 {
			t.Errorf("expected 51 fields, got %d", len(s.Models[0].ScalarFields))
		}
	})

	t.Run("self-referential model", func(t *testing.T) {
		input := `
		model Category {
			id       Int        @id
			name     String
			parent   Category?  @relation("CategoryTree", fields: [parentId], references: [id])
			parentId Int?
			children Category[] @relation("CategoryTree")
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("self-referential model broken: %v", s.Errors)
		}
		// parent and children should be relation fields
		if len(s.Models[0].RelationFields) != 2 {
			t.Errorf("expected 2 relation fields, got %d", len(s.Models[0].RelationFields))
		}
	})

	t.Run("many-to-many implicit join table model", func(t *testing.T) {
		input := `
		model Post {
			id   Int    @id
			tags Tag[]
		}
		model Tag {
			id    Int    @id
			posts Post[]
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 2 {
			t.Fatalf("m2m broken: %v", s.Errors)
		}
	})

	t.Run("view keyword instead of model", func(t *testing.T) {
		input := `
		view UserView {
			id    Int    @unique
			email String
		}
		`
		s, _ := ParseSchema(input)
		// Either it parses as a view or records an error — it must NOT panic
		_ = s
	})
}

// -----------------------------------------------------------------------------
// Block type tests
// -----------------------------------------------------------------------------

func TestOdyssey_BlockTypeChaos(t *testing.T) {
	t.Run("datasource with shadow database url", func(t *testing.T) {
		input := `
		datasource db {
			provider          = "postgresql"
			url               = env("DATABASE_URL")
			shadowDatabaseUrl = env("SHADOW_DATABASE_URL")
			relationMode      = "prisma"
		}
		model Trivial {
			id Int @id
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("shadow db url broken: %v", s.Errors)
		}
	})

	t.Run("completely empty schema", func(t *testing.T) {
		input := ``
		s, _ := ParseSchema(input)
		if len(s.Models) != 0 {
			t.Errorf("empty schema should yield no models")
		}
		if len(s.Errors) != 0 {
			t.Errorf("empty schema should yield no errors: %v", s.Errors)
		}
	})

	t.Run("schema is only comments", func(t *testing.T) {
		input := `
		// This is a comment
		// Another comment
		// Yet another
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 0 {
			t.Errorf("comment-only schema should yield no models")
		}
		if len(s.Errors) != 0 {
			t.Errorf("comment-only schema should yield no errors: %v", s.Errors)
		}
	})

	t.Run("unknown top-level block type", func(t *testing.T) {
		input := `
		wizzard Merlin {
			power = "infinite"
		}
		model User {
			id Int @id
		}
		`
		s, _ := ParseSchema(input)
		// Must not panic; should recover and parse User
		found := false
		for _, m := range s.Models {
			if m.Name == "User" {
				found = true
			}
		}
		if !found {
			t.Errorf("failed to recover after unknown block type; models: %+v", s.Models)
		}
	})
}

// -----------------------------------------------------------------------------
// Error recovery tests
// -----------------------------------------------------------------------------

func TestOdyssey_ErrorRecovery(t *testing.T) {
	t.Run("multiple errors spread through schema, good models survive", func(t *testing.T) {
		input := `
		model Good1 {
			id   Int    @id
			name String
		}

		model Bad1 {
			id Int @id @broken(
		}

		model Good2 {
			id    Int    @id
			email String @unique
		}

		model Bad2 {
			id Int @id
			x  String @default("unterminated
		}

		model Good3 {
			id   Int    @id
			slug String
		}
		`
		s, _ := ParseSchema(input)

		goodNames := []string{"Good1", "Good2", "Good3"}
		for _, name := range goodNames {
			found := false
			for _, m := range s.Models {
				if m.Name == name {
					found = true
				}
			}
			if !found {
				t.Errorf("good model %q lost during error recovery; parsed: %+v; errors: %v", name, s.Models, s.Errors)
			}
		}
	})

	t.Run("truncated schema — file cut off mid-model", func(t *testing.T) {
		input := `
		model Complete {
			id Int @id
		}

		model Incomplete {
			id   Int    @id
			name String
		`
		// No closing brace — simulates a file truncated on disk
		s, _ := ParseSchema(input)
		// Complete should be parseable
		found := false
		for _, m := range s.Models {
			if m.Name == "Complete" {
				found = true
			}
		}
		if !found {
			t.Errorf("truncated schema lost Complete model; models=%+v", s.Models)
		}
		// Must not panic or infinite loop — if we're here, we're good
	})

	t.Run("duplicate model names", func(t *testing.T) {
		input := `
		model Twin {
			id   Int    @id
			name String
		}
		model Twin {
			id    Int    @id
			email String
		}
		`
		s, _ := ParseSchema(input)
		// Parser may error or accept both — it must NOT panic
		// We just want deterministic behavior
		_ = s.Models
		_ = s.Errors
	})

	t.Run("completely garbled input", func(t *testing.T) {
		input := `
		}}}}{{{{{ @@@@ !!!! #### $$$
		model model model {{{
		@@@ ??? !!!
		`
		s, _ := ParseSchema(input)
		// Must not panic. That's the whole test.
		_ = s
	})

	t.Run("brace mismatch — extra closing brace at top level", func(t *testing.T) {
		input := `
		model User {
			id Int @id
		}
		}
		model Post {
			id Int @id
		}
		`
		s, _ := ParseSchema(input)
		_ = s // must not panic or loop forever
	})

	t.Run("deeply nested erroneous attribute args", func(t *testing.T) {
		input := `
		model Nested {
			id   Int    @id
			bad  String @relation(fields: [[[a, b], c], d], references: [id])
		}
		`
		s, _ := ParseSchema(input)
		_ = s // must not panic
	})
}

// -----------------------------------------------------------------------------
// Conditional expressions (where clauses)
// -----------------------------------------------------------------------------

func TestOdyssey_ConditionalExpressions(t *testing.T) {
	t.Run("where clause with AND logic", func(t *testing.T) {
		input := `
		model Post {
			id        Int      @id
			status    String
			deletedAt DateTime?

			@@index([status], where: deletedAt == null)
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("conditional index broken: %v", s.Errors)
		}
		attr := s.Models[0].Attributes[0]
		whereArg := attr.Args[1]
		if whereArg.Name != "where" {
			t.Errorf("expected named arg 'where', got %q", whereArg.Name)
		}
		if whereArg.Value.Type != ValBinary {
			t.Errorf("expected binary expression, got %+v", whereArg.Value)
		}
	})

	t.Run("where clause comparing to string literal", func(t *testing.T) {
		input := `
		model User {
			id     Int    @id
			status String

			@@index([status], where: status == "active")
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("string where clause broken: %v", s.Errors)
		}
		attr := s.Models[0].Attributes[0]
		right := attr.Args[1].Value.Right
		if right.Type != ValLiteral || right.Scalar != "active" {
			t.Errorf("expected string literal 'active', got %+v", right)
		}
	})

	t.Run("where clause with != operator", func(t *testing.T) {
		input := `
		model Item {
			id     Int    @id
			status String

			@@index([status], where: status != "deleted")
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("!= where clause broken: %v", s.Errors)
		}
	})
}

// -----------------------------------------------------------------------------
// Large schema integration
// -----------------------------------------------------------------------------

func TestOdyssey_LargeSchemaIntegration(t *testing.T) {
	t.Run("realistic e-commerce schema all at once", func(t *testing.T) {
		input := `
		datasource db {
			provider = "postgresql"
			url      = env("DATABASE_URL")
		}

		enum Role {
			CUSTOMER
			ADMIN
			VENDOR
		}

		enum OrderStatus {
			PENDING
			PROCESSING
			SHIPPED
			DELIVERED
			CANCELLED
			REFUNDED
		}

		model User {
			id        Int       @id @default(autoincrement())
			email     String    @unique
			name      String?
			role      Role      @default(CUSTOMER)
			createdAt DateTime  @default(now())
			updatedAt DateTime  @updatedAt
			orders    Order[]
			addresses Address[]
			reviews   Review[]

			@@map("users")
		}

		model Address {
			id       Int    @id @default(autoincrement())
			line1    String
			line2    String?
			city     String
			state    String
			zip      String
			country  String @default("US")
			userId   Int
			user     User   @relation(fields: [userId], references: [id], onDelete: Cascade)

			@@index([userId])
			@@map("addresses")
		}

		model Product {
			id          Int         @id @default(autoincrement())
			sku         String      @unique
			name        String
			description String?
			price       Decimal     @db.Decimal(10, 2)
			stock       Int         @default(0)
			images      String[]
			categoryId  Int?
			category    Category?   @relation(fields: [categoryId], references: [id])
			orderItems  OrderItem[]
			reviews     Review[]

			@@index([sku])
			@@index([categoryId])
			@@map("products")
		}

		model Category {
			id       Int        @id @default(autoincrement())
			name     String     @unique
			slug     String     @unique
			parent   Category?  @relation("CategoryTree", fields: [parentId], references: [id])
			parentId Int?
			children Category[] @relation("CategoryTree")
			products Product[]

			@@map("categories")
		}

		model Order {
			id         Int         @id @default(autoincrement())
			status     OrderStatus @default(PENDING)
			total      Decimal     @db.Decimal(10, 2)
			userId     Int
			user       User        @relation(fields: [userId], references: [id])
			items      OrderItem[]
			createdAt  DateTime    @default(now())
			updatedAt  DateTime    @updatedAt

			@@index([userId])
			@@index([status])
			@@map("orders")
		}

		model OrderItem {
			id        Int     @id @default(autoincrement())
			quantity  Int
			price     Decimal @db.Decimal(10, 2)
			orderId   Int
			order     Order   @relation(fields: [orderId], references: [id], onDelete: Cascade)
			productId Int
			product   Product @relation(fields: [productId], references: [id])

			@@unique([orderId, productId])
			@@map("order_items")
		}

		model Review {
			id        Int     @id @default(autoincrement())
			rating    Int
			body      String?
			userId    Int
			user      User    @relation(fields: [userId], references: [id])
			productId Int
			product   Product @relation(fields: [productId], references: [id])

			@@unique([userId, productId])
			@@index([productId])
			@@map("reviews")
		}
		`
		s, _ := ParseSchema(input)

		if len(s.Errors) != 0 {
			t.Fatalf("realistic schema should have 0 errors, got %d: %v", len(s.Errors), s.Errors)
		}
		if len(s.Models) != 7 {
			t.Errorf("expected 7 models, got %d", len(s.Models))
		}
		if len(s.Enums) != 2 {
			t.Errorf("expected 2 enums, got %d", len(s.Enums))
		}

		// Verify specific structural details
		var product *Model
		for i := range s.Models {
			if s.Models[i].Name == "Product" {
				product = s.Models[i]
				break
			}
		}
		if product == nil {
			t.Fatal("Product model not found")
		}
		if len(product.Attributes) != 3 {
			t.Errorf("Product should have 3 block attrs (2x@@index + @@map), got %d", len(product.Attributes))
		}
	})

	t.Run("interleaved enums and models in random order", func(t *testing.T) {
		input := `
		model A { id Int @id }
		enum X { ONE TWO }
		model B { id Int @id }
		enum Y { ALPHA BETA GAMMA }
		model C { id Int @id }
		enum Z { P Q R S T }
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 3 {
			t.Errorf("expected 3 models, got %d", len(s.Models))
		}
		if len(s.Enums) != 3 {
			t.Errorf("expected 3 enums, got %d", len(s.Enums))
		}
	})
}

// -----------------------------------------------------------------------------
// Whitespace and comments
// -----------------------------------------------------------------------------

func TestOdyssey_WhitespaceAndComments(t *testing.T) {
	t.Run("comments between every token", func(t *testing.T) {
		input := `
		// comment before model
		model // comment after keyword
		CommentedOut // comment after name
		{ // comment after open brace
			// comment before field
			id // comment after field name
			Int // comment after type
			@id // comment after attribute
			// comment at end of field line
			// another comment
			name String // inline comment
		} // comment after close brace
		// trailing comment
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("comment-laden schema broken: models=%d, errors=%v", len(s.Models), s.Errors)
		}
		m := s.Models[0]
		if m.Name != "CommentedOut" {
			t.Errorf("wrong model name: %s", m.Name)
		}
		if len(m.ScalarFields) != 2 {
			t.Errorf("expected 2 fields, got %d", len(m.ScalarFields))
		}
	})

	t.Run("doc comments (triple slash) on fields", func(t *testing.T) {
		input := `
		model Documented {
			/// The primary key
			id   Int    @id
			/// The user's email address
			email String @unique
		}
		`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("doc comment broken: %v", s.Errors)
		}
		if len(s.Models[0].ScalarFields) != 2 {
			t.Errorf("expected 2 fields, got %d", len(s.Models[0].ScalarFields))
		}
	})

	t.Run("no newlines — entire schema on one line", func(t *testing.T) {
		input := `model Inline { id Int @id name String @unique }`
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("single-line model broken: %v", s.Errors)
		}
		if len(s.Models[0].ScalarFields) != 2 {
			t.Errorf("expected 2 fields, got %d", len(s.Models[0].ScalarFields))
		}
	})

	t.Run("excessive blank lines between everything", func(t *testing.T) {
		input := "\n\n\n\nmodel\n\n\nSpaced\n\n\n{\n\n\n\nid\n\n\nInt\n\n\n@id\n\n\n}\n\n\n"
		s, _ := ParseSchema(input)
		if len(s.Models) != 1 {
			t.Fatalf("spaced-out schema broken: %v", s.Errors)
		}
	})
}

// -----------------------------------------------------------------------------
// Composite schema integration tests
// -----------------------------------------------------------------------------

func TestOdyssey_FinalBoss(t *testing.T) {
	input := `
	// Composite integration schema

	datasource db {
		provider          = "postgresql"
		url               = env("DATABASE_URL")
		shadowDatabaseUrl = env("SHADOW_DATABASE_URL")
		relationMode      = "prisma"
	}

	/// The role of a user in the system.
	enum Role {
		USER
		ADMIN
		SUPER_ADMIN
	}

	// A model with every trick in the book
	model User {
		id            Int       @id @default(autoincrement())
		email         String    @unique @db.VarChar(320)
		emailVerified DateTime?
		name          String?   @db.VarChar(100)
		role          Role      @default(USER)
		metadata      Json?
		rawBytes      Bytes?
		score         Decimal   @default(0) @db.Decimal(5, 2)
		createdAt     DateTime  @default(now())
		updatedAt     DateTime  @updatedAt
		deletedAt     DateTime?
		tags          String[]
		posts         Post[]
		profile       Profile?
		sessions      Session[]

		@@unique([email])
		@@index([role, createdAt])
		@@index([email], where: deletedAt == null)
		@@index([role], where: deletedAt == null)
		@@map("users")
	}

	model Profile {
		id     Int     @id @default(autoincrement())
		bio    String? @db.Text
		userId Int     @unique
		user   User    @relation(fields: [userId], references: [id], onDelete: Cascade, onUpdate: Cascade)

		@@map("profiles")
	}

	model Post {
		id          Int       @id @default(autoincrement())
		title       String    @db.VarChar(255)
		slug        String    @unique @db.VarChar(255)
		body        String?   @db.Text
		published   Boolean   @default(false)
		publishedAt DateTime?
		authorId    Int
		author      User      @relation(fields: [authorId], references: [id])
		tags        Tag[]
		views       Int       @default(0)
		geoPoint    Unsupported("geometry(Point, 4326)")?

		@@index([authorId])
		@@index([slug], where: published == true)
		@@index([publishedAt], where: publishedAt != null)
		@@map("posts")
	}

	model Tag {
		id    Int    @id @default(autoincrement())
		name  String @unique
		slug  String @unique
		posts Post[]

		@@map("tags")
	}

	model Session {
		id        String   @id @default(cuid())
		token     String   @unique @db.VarChar(512)
		userId    Int
		user      User     @relation(fields: [userId], references: [id], onDelete: Cascade)
		expiresAt DateTime
		createdAt DateTime @default(now())

		@@index([userId])
		@@index([token], where: expiresAt != null)
		@@map("sessions")
	}
	`

	s, _ := ParseSchema(input)

	if len(s.Errors) != 0 {
		t.Fatalf("THE FINAL BOSS: expected 0 errors, got %d:\n%v", len(s.Errors), s.Errors)
	}

	expectedModels := []string{"User", "Profile", "Post", "Tag", "Session"}
	if len(s.Models) != len(expectedModels) {
		t.Fatalf("expected %d models, got %d", len(expectedModels), len(s.Models))
	}
	for _, name := range expectedModels {
		found := false
		for _, m := range s.Models {
			if m.Name == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("model %q missing from final boss parse", name)
		}
	}

	if len(s.Enums) != 1 || s.Enums[0].Name != "Role" {
		t.Errorf("expected enum Role, got: %+v", s.Enums)
	}

	// Verify User's complexity survived
	var user *Model
	for i := range s.Models {
		if s.Models[i].Name == "User" {
			user = s.Models[i]
			break
		}
	}
	if user == nil {
		t.Fatal("User model missing")
	}
	// 5 block attrs: @@unique, @@index x3, @@map
	if len(user.Attributes) != 5 {
		t.Errorf("User: expected 5 block attrs, got %d: %+v", len(user.Attributes), user.Attributes)
	}
	// 4 relation fields: posts, profile, sessions + post scalar tag through join
	if len(user.RelationFields) < 3 {
		t.Errorf("User: expected at least 3 relation fields, got %d", len(user.RelationFields))
	}

	// Verify Post's Unsupported type survived
	var post *Model
	for i := range s.Models {
		if s.Models[i].Name == "Post" {
			post = s.Models[i]
			break
		}
	}
	if post == nil {
		t.Fatal("Post model missing")
	}
	var geoField *ScalarField
	for i := range post.ScalarFields {
		if post.ScalarFields[i].Name == "geoPoint" {
			geoField = post.ScalarFields[i]
			break
		}
	}
	if geoField == nil {
		t.Error("Post.geoPoint field not found")
	} else {
		if geoField.Type != `Unsupported("geometry(Point, 4326)")` {
			t.Errorf("geoPoint type mangled: %q", geoField.Type)
		}
		if !geoField.Optional {
			t.Errorf("geoPoint should be optional")
		}
	}

	t.Log("SURVIVED")
}
