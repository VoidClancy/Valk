package schema

func ParseSchema(input string) (*Schema, DiagnosticList) {

	tokens := ExtractTokens(input)

	parser := &Parser{Tokens: tokens}
	ast := parser.ParseAST()

	resolver := NewResolver(ast)
	schema := resolver.Resolve()
	schema.Errors = append(parser.errors, schema.Errors...)

	return schema, schema.Errors
}
