package schema

import (
	"fmt"
)

type astAST struct {
	Datasources []astDatasourceDecl

	Models []astModelDecl
	Enums  []astEnumDecl
}

type astDatasourceDecl struct {
	Name       string
	Properties []astKeyValue
	Line, Col  int
}

type astKeyValue struct {
	Key       string
	Value     Value
	Line, Col int
}

type astEnumDecl struct {
	Name       string
	Values     []astEnumValueDecl
	Attributes []Attribute
	Line, Col  int
}

type astEnumValueDecl struct {
	Name       string
	Attributes []Attribute
	Line, Col  int
}

type astModelDecl struct {
	Name       string
	Fields     []astFieldDecl
	Attributes []Attribute
	Line, Col  int
}

type astFieldDecl struct {
	Name       string
	TypeName   string
	IsArray    bool
	Optional   bool
	Attributes []Attribute
	Line, Col  int
}

type delimiter struct {
	typ  TokenType
	line int
	col  int
}

func delimChar(t TokenType) string {
	switch t {
	case LPAREN, RPAREN:
		return "("
	case LBRACKET, RBRACKET:
		return "["
	case LBRACE, RBRACE:
		return "{"
	default:
		return ""
	}
}

func isMatching(open, close TokenType) bool {
	return (open == LPAREN && close == RPAREN) ||
		(open == LBRACKET && close == RBRACKET) ||
		(open == LBRACE && close == RBRACE)
}

type Parser struct {
	Tokens     []Token
	Pos        int
	errors     DiagnosticList
	openDelims []delimiter
}

func (p *Parser) pushDelim(tok Token) {
	p.openDelims = append(p.openDelims, delimiter{
		typ:  tok.Type,
		line: tok.Line,
		col:  tok.Col,
	})
}

func (p *Parser) popDelim(expectedClosing TokenType) {
	if len(p.openDelims) == 0 {
		return
	}
	last := p.openDelims[len(p.openDelims)-1]
	if isMatching(last.typ, expectedClosing) {
		p.openDelims = p.openDelims[:len(p.openDelims)-1]
	}
}

func (p *Parser) clearDelimsToLastLBrace() {
	for len(p.openDelims) > 0 {
		last := p.openDelims[len(p.openDelims)-1]
		p.openDelims = p.openDelims[:len(p.openDelims)-1]
		if last.typ == LBRACE {
			break
		}
	}
}

func (p *Parser) clearAllDelims() {
	p.openDelims = nil
}

func (p *Parser) eof() bool {
	return p.Pos >= len(p.Tokens)
}

func (p *Parser) current() Token {
	if p.eof() {
		return Token{Type: EOF}
	}
	return p.Tokens[p.Pos]
}

func (p *Parser) lastTokenLine() int {
	if len(p.Tokens) > 0 {
		return p.Tokens[len(p.Tokens)-1].Line
	}
	return 1
}

func (p *Parser) lastTokenCol() int {
	if len(p.Tokens) > 0 {
		return p.Tokens[len(p.Tokens)-1].Col
	}
	return 1
}

func (p *Parser) errorf(format string, args ...any) {
	tok := p.current()
	line := tok.Line
	col := tok.Col
	if tok.Type == EOF {
		line = p.lastTokenLine()
		col = p.lastTokenCol()
	}
	panic(Diagnostic{
		Severity: SevError,
		Message:  fmt.Sprintf(format, args...),
		Pos:      Position{Line: line, Col: col},
		Source:   "parser",
	})
}

func (p *Parser) advance() {
	if !p.eof() {
		p.Pos++
	}
}

func (p *Parser) next() Token {
	if p.Pos+1 >= len(p.Tokens) {
		return Token{Type: EOF}
	}
	return p.Tokens[p.Pos+1]
}

func (p *Parser) expect(t TokenType) Token {
	if p.eof() {
		tok := p.current()
		if t == RBRACKET || t == RPAREN || t == RBRACE {
			if len(p.openDelims) > 0 {
				last := p.openDelims[len(p.openDelims)-1]
				if isMatching(last.typ, t) {
					panic(Diagnostic{
						Severity: SevError,
						Message:  fmt.Sprintf("unterminated %q starting at line %d, col %d", delimChar(last.typ), last.line, last.col),
						Pos:      Position{Line: tok.Line, Col: tok.Col},
						Source:   "parser",
					})
				}
			}
		}
		panic(Diagnostic{
			Severity: SevError,
			Message:  fmt.Sprintf("unexpected EOF, expected %v", t),
			Pos:      Position{Line: p.lastTokenLine(), Col: p.lastTokenCol()},
			Source:   "parser",
		})
	}

	tok := p.current()

	if tok.Type != t {
		if t == RBRACKET || t == RPAREN || t == RBRACE {
			if len(p.openDelims) > 0 {
				last := p.openDelims[len(p.openDelims)-1]
				if isMatching(last.typ, t) {
					panic(Diagnostic{
						Severity: SevError,
						Message:  fmt.Sprintf("unterminated %q starting at line %d, col %d", delimChar(last.typ), last.line, last.col),
						Pos:      Position{Line: tok.Line, Col: tok.Col},
						Source:   "parser",
					})
				}
			}
		}
		panic(Diagnostic{
			Severity: SevError,
			Message:  fmt.Sprintf("expected %v, got %v (%q)", t, tok.Type, tok.Value),
			Pos:      Position{Line: tok.Line, Col: tok.Col},
			Source:   "parser",
		})
	}

	switch tok.Type {
	case LPAREN, LBRACKET, LBRACE:
		p.pushDelim(tok)
	case RPAREN, RBRACKET, RBRACE:
		p.popDelim(tok.Type)
	}

	p.advance()

	return tok
}

func (p *Parser) parseDatasourceDecl() astDatasourceDecl {
	startTok := p.expect(IDENT)
	name := p.expect(IDENT).Value
	p.expect(LBRACE)

	var props []astKeyValue
	for !p.eof() && p.current().Type != RBRACE && p.current().Type != EOF {
		keyTok := p.expect(IDENT)
		p.expect(ASSIGN)
		val := p.parseValue()
		props = append(props, astKeyValue{
			Key:   keyTok.Value,
			Value: val,
			Line:  keyTok.Line,
			Col:   keyTok.Col,
		})
	}
	p.expect(RBRACE)

	return astDatasourceDecl{
		Name:       name,
		Properties: props,
		Line:       startTok.Line,
		Col:        startTok.Col,
	}
}

func (p *Parser) safeParseDatasourceDecl(ast *astAST) {
	defer func() {
		if r := recover(); r != nil {
			if diag, ok := r.(Diagnostic); ok {
				p.errors = append(p.errors, diag)
				p.recoverToTopLevel()
			} else {
				panic(r)
			}
		}
	}()
	ast.Datasources = append(ast.Datasources, p.parseDatasourceDecl())
}

func (p *Parser) parseEnumDecl() astEnumDecl {
	startTok := p.expect(IDENT) // "enum"
	name := p.expect(IDENT).Value
	p.expect(LBRACE)

	var values []astEnumValueDecl
	var attrs []Attribute

	for !p.eof() && p.current().Type != RBRACE && p.current().Type != EOF {
		if p.current().Type == ATAT {
			attrs = append(attrs, p.parseModelAttribute())
			continue
		}

		if p.current().Type == IDENT {
			val := p.current().Value
			if val == "model" || val == "enum" || val == "datasource" {
				if len(p.openDelims) > 0 {
					last := p.openDelims[len(p.openDelims)-1]
					if last.typ == LBRACE {
						p.errorf("unterminated %q starting at line %d, col %d", delimChar(last.typ), last.line, last.col)
					}
				}
				p.errorf("expected RBRACE, got %s", val)
			}
		}

		valTok := p.expect(IDENT)
		var valAttrs []Attribute
		for p.current().Type == AT {
			valAttrs = append(valAttrs, p.parseFieldAttribute())
		}
		values = append(values, astEnumValueDecl{
			Name:       valTok.Value,
			Attributes: valAttrs,
			Line:       valTok.Line,
			Col:        valTok.Col,
		})
	}
	p.expect(RBRACE)

	return astEnumDecl{
		Name:       name,
		Values:     values,
		Attributes: attrs,
		Line:       startTok.Line,
		Col:        startTok.Col,
	}
}

func (p *Parser) safeParseEnumDecl(ast *astAST) {
	defer func() {
		if r := recover(); r != nil {
			if diag, ok := r.(Diagnostic); ok {
				p.errors = append(p.errors, diag)
				p.recoverToTopLevel()
			} else {
				panic(r)
			}
		}
	}()
	ast.Enums = append(ast.Enums, p.parseEnumDecl())
}

func (p *Parser) parseModelDecl() astModelDecl {
	startTok := p.expect(IDENT)
	name := p.expect(IDENT).Value
	p.expect(LBRACE)

	var fields []astFieldDecl
	var attrs []Attribute

	for !p.eof() && p.current().Type != RBRACE && p.current().Type != EOF {
		if p.current().Type == ATAT {
			attrs = append(attrs, p.parseModelAttribute())
			continue
		}

		if p.current().Type == IDENT {
			val := p.current().Value
			if val == "model" || val == "enum" || val == "datasource" {
				if len(p.openDelims) > 0 {
					last := p.openDelims[len(p.openDelims)-1]
					if last.typ == LBRACE {
						p.errorf("unterminated %q starting at line %d, col %d", delimChar(last.typ), last.line, last.col)
					}
				}
				p.errorf("expected RBRACE, got %s", val)
			}
			fields = append(fields, p.parseFieldDecl())
			continue
		}

		tok := p.current()
		p.errorf("unexpected token %q (%v) inside model %s", tok.Value, tok.Type, name)
	}
	p.expect(RBRACE)

	return astModelDecl{
		Name:       name,
		Fields:     fields,
		Attributes: attrs,
		Line:       startTok.Line,
		Col:        startTok.Col,
	}
}

func (p *Parser) safeParseModelDecl(ast *astAST) {
	defer func() {
		if r := recover(); r != nil {
			if diag, ok := r.(Diagnostic); ok {
				p.errors = append(p.errors, diag)
				p.recoverFromModelError()
			} else {
				panic(r)
			}
		}
	}()
	ast.Models = append(ast.Models, p.parseModelDecl())
}

func (p *Parser) parseFieldType() string {
	ident := p.expect(IDENT).Value
	if ident == "Unsupported" && p.current().Type == LPAREN {
		p.expect(LPAREN)
		strVal := p.expect(STRING).Value
		p.expect(RPAREN)
		return fmt.Sprintf("Unsupported(%q)", strVal)
	}
	return ident
}

func (p *Parser) parseFieldDecl() astFieldDecl {
	nameTok := p.expect(IDENT)
	fieldType := p.parseFieldType()

	var isArray bool
	if p.current().Type == LBRACKET {
		p.expect(LBRACKET)
		p.expect(RBRACKET)
		isArray = true
	}

	var optional bool
	if p.current().Type == QUESTION {
		p.expect(QUESTION)
		optional = true
	}

	var attrs []Attribute
	for p.current().Type == AT {
		attrs = append(attrs, p.parseFieldAttribute())
	}

	return astFieldDecl{
		Name:       nameTok.Value,
		TypeName:   fieldType,
		IsArray:    isArray,
		Optional:   optional,
		Attributes: attrs,
		Line:       nameTok.Line,
		Col:        nameTok.Col,
	}
}

func (p *Parser) parseAttributeName() string {
	name := p.expect(IDENT).Value
	for p.current().Type == DOT {
		p.expect(DOT)
		name += "." + p.expect(IDENT).Value
	}
	return name
}

func (p *Parser) parseAttribute(prefix TokenType) Attribute {
	tok := p.expect(prefix)

	if p.current().Type != IDENT {
		p.errorf("missing attribute name after %q", tok.Value)
	}

	attr := Attribute{
		Name: p.parseAttributeName(),
		Line: tok.Line,
		Col:  tok.Col,
	}

	if p.current().Type == LPAREN {
		attr.Args = p.parseAttributeArgs()
	}

	return attr
}

func (p *Parser) parseModelAttribute() Attribute {
	return p.parseAttribute(ATAT)
}

func (p *Parser) parseFieldAttribute() Attribute {
	return p.parseAttribute(AT)
}

func (p *Parser) parseAttributeArgs() []Argument {
	p.expect(LPAREN)
	var args []Argument

	for p.current().Type != RPAREN && p.current().Type != EOF {
		args = append(args, p.parseArgument())

		if p.current().Type == COMMA {
			p.expect(COMMA)
		}
	}

	p.expect(RPAREN)

	return args
}

func (p *Parser) parseArgument() Argument {
	var name string
	if p.current().Type == IDENT && p.next().Type == COLON {
		name = p.expect(IDENT).Value
		p.expect(COLON)
	}

	val := p.parseValue()
	return Argument{
		Name:  name,
		Value: val,
	}
}

func (p *Parser) parseValue() Value {
	return p.parseExpression(0)
}

func (p *Parser) parseExpression(precedence int) Value {
	left := p.parsePrimaryValue()

	for {
		tok := p.current()
		opPrecedence := getPrecedence(tok.Type)
		if opPrecedence < precedence {
			break
		}

		p.advance() // consume operator
		right := p.parseExpression(opPrecedence + 1)
		lCopy := left
		rCopy := right
		left = Value{
			Type:   ValBinary,
			Scalar: tok.Value,
			Left:   &lCopy,
			Right:  &rCopy,
		}
	}

	return left
}

func getPrecedence(t TokenType) int {
	switch t {
	case OR:
		return 1
	case AND:
		return 2
	case EQUAL, NOT_EQUAL:
		return 3
	case LT, GT, LTE, GTE:
		return 4
	default:
		return -1
	}
}

func (p *Parser) parsePrimaryValue() Value {
	tok := p.current()
	switch tok.Type {
	case STRING:
		p.advance()
		return Value{
			Type:   ValLiteral,
			Scalar: tok.Value,
		}
	case NUMBER:
		p.advance()
		return Value{
			Type:   ValLiteral,
			Scalar: tok.Value,
		}
	case BOOLEAN:
		p.advance()
		return Value{
			Type:   ValLiteral,
			Scalar: tok.Value,
		}
	case LBRACKET:
		p.expect(LBRACKET)
		var list []Value
		for p.current().Type != RBRACKET && p.current().Type != EOF {
			list = append(list, p.parseValue())
			if p.current().Type == COMMA {
				p.expect(COMMA)
			}
		}
		p.expect(RBRACKET)
		return Value{
			Type:  ValArray,
			Array: list,
		}
	case IDENT:
		name := p.parseAttributeName()
		if p.current().Type == LPAREN {
			p.expect(LPAREN)
			var args []Argument
			for p.current().Type != RPAREN && p.current().Type != EOF {
				args = append(args, p.parseArgument())
				if p.current().Type == COMMA {
					p.expect(COMMA)
				}
			}
			p.expect(RPAREN)
			return Value{
				Type:   ValFunc,
				Scalar: name,
				Args:   args,
			}
		}
		return Value{
			Type:   ValIdent,
			Scalar: name,
		}
	case ILLEGAL:
		p.advance()
		msg := fmt.Sprintf("illegal character %q", tok.Value)
		if tok.Value == "unterminated string literal" {
			msg = "unterminated string literal"
		}
		panic(Diagnostic{
			Severity: SevError,
			Message:  msg,
			Pos:      Position{Line: tok.Line, Col: tok.Col},
			Source:   "lexer",
		})
	default:
		if len(p.openDelims) > 0 {
			last := p.openDelims[len(p.openDelims)-1]
			if last.typ == LPAREN || last.typ == LBRACKET {
				p.errorf("unterminated %q starting at line %d, col %d", delimChar(last.typ), last.line, last.col)
			}
		}
		p.errorf("unexpected token %q (%v) when parsing value", tok.Value, tok.Type)
		return Value{}
	}
}

func (p *Parser) recoverFromModelError() {
	p.clearDelimsToLastLBrace()
	braceCount := 0
	for !p.eof() && p.current().Type != EOF {
		tok := p.current()
		if tok.Type == LBRACE {
			braceCount++
		} else if tok.Type == RBRACE {
			braceCount--
			if braceCount <= 0 {
				p.advance()
				break
			}
		} else if braceCount <= 0 && tok.Type == IDENT && tok.Value == "model" {
			break
		}
		p.advance()
	}
}

func (p *Parser) recoverToTopLevel() {
	p.clearAllDelims()
	for !p.eof() && p.current().Type != EOF {
		tok := p.current()
		if tok.Type == IDENT {
			switch tok.Value {
			case "model", "enum", "datasource":
				return
			}
		}
		p.advance()
	}
}

func (p *Parser) ParseAST() *astAST {
	var ast astAST

	for !p.eof() && p.current().Type != EOF {
		tok := p.current()
		if tok.Type == ILLEGAL {
			msg := fmt.Sprintf("illegal character %q", tok.Value)
			if tok.Value == "unterminated string literal" {
				msg = "unterminated string literal"
			}
			p.errors = append(p.errors, Diagnostic{
				Severity: SevError,
				Message:  msg,
				Pos:      Position{Line: tok.Line, Col: tok.Col},
				Source:   "lexer",
			})
			p.advance()
			continue
		}

		if tok.Type != IDENT {
			p.errors = append(p.errors, Diagnostic{
				Severity: SevError,
				Message:  fmt.Sprintf("unexpected token %q (%v) at top-level", tok.Value, tok.Type),
				Pos:      Position{Line: tok.Line, Col: tok.Col},
				Source:   "parser",
			})
			p.advance()
			continue
		}

		switch tok.Value {
		case "model":
			p.safeParseModelDecl(&ast)
		case "enum":
			p.safeParseEnumDecl(&ast)
		case "datasource":
			p.safeParseDatasourceDecl(&ast)

		default:
			p.errors = append(p.errors, Diagnostic{
				Severity: SevError,
				Message:  fmt.Sprintf("unexpected top-level identifier %q", tok.Value),
				Pos:      Position{Line: tok.Line, Col: tok.Col},
				Source:   "parser",
			})
			p.advance()
		}
	}

	return &ast
}
