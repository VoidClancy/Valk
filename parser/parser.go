package parser

import (
	"fmt"
)

type Parser struct {
	Tokens []Token
	Pos    int
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
		panic(fmt.Sprintf("unexpected EOF, expected %v", t))
	}

	tok := p.current()

	if tok.Type != t {
		panic(fmt.Sprintf(
			"expected %v, got %v (%q) at line %d, col %d",
			t,
			tok.Type,
			tok.Value,
			tok.Line,
			tok.Col,
		))
	}

	p.advance()

	return tok
}

func (p *Parser) skipBlock() {
	for !p.eof() && p.current().Type != LBRACE && p.current().Type != EOF {
		p.advance()
	}
	if p.eof() || p.current().Type == EOF {
		return
	}
	p.expect(LBRACE)
	braceCount := 1
	for !p.eof() && braceCount > 0 && p.current().Type != EOF {
		if p.current().Type == LBRACE {
			braceCount++
		} else if p.current().Type == RBRACE {
			braceCount--
		}
		p.advance()
	}
}

func (p *Parser) ParseModel() Model {
	tok := p.expect(IDENT)
	if tok.Value != "model" {
		panic(fmt.Sprintf("expected 'model' keyword, got %q", tok.Value))
	}

	model := Model{
		Name: p.expect(IDENT).Value,
	}

	p.expect(LBRACE)

	for !p.eof() && p.current().Type != RBRACE && p.current().Type != EOF {
		if p.current().Type == ATAT {
			model.Attributes = append(
				model.Attributes,
				p.parseModelAttribute(),
			)
			continue
		}

		if p.current().Type == IDENT {
			model.Fields = append(model.Fields, p.parseField())
			continue
		}

		tok := p.current()
		panic(fmt.Sprintf("unexpected token %q (%v) at line %d, col %d inside model %s", tok.Value, tok.Type, tok.Line, tok.Col, model.Name))
	}

	p.expect(RBRACE)

	return model
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
	p.expect(prefix)

	attr := Attribute{
		Name: p.parseAttributeName(),
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

func valPtr(v Value) *Value {
	return &v
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
		left = Value{
			Type:   ValBinary,
			Scalar: tok.Value,
			Left:   valPtr(left),
			Right:  valPtr(right),
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
		// Could be a function call or a simple identifier
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
	default:
		panic(fmt.Sprintf("unexpected token %q (%v) at line %d, col %d when parsing value", tok.Value, tok.Type, tok.Line, tok.Col))
	}
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

func (p *Parser) parseField() Field {
	field := Field{
		Name: p.expect(IDENT).Value,
		Type: p.parseFieldType(),
	}

	if p.current().Type == LBRACKET {
		p.expect(LBRACKET)
		p.expect(RBRACKET)
		field.IsArray = true
	}

	if p.current().Type == QUESTION {
		p.expect(QUESTION)
		field.IsOptional = true
	}

	for p.current().Type == AT {
		field.Attributes = append(
			field.Attributes,
			p.parseFieldAttribute(),
		)
	}
	return field
}

func (p *Parser) Parse() Schema {
	var schema Schema

	for !p.eof() && p.current().Type != EOF {
		tok := p.current()
		if tok.Type == IDENT && tok.Value == "model" {
			schema.Models = append(
				schema.Models,
				p.ParseModel(),
			)
		} else {
			// Skip other blocks
			p.skipBlock()
		}
	}

	return schema
}
