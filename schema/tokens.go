package schema

import (
	"fmt"
)

type TokenType int

const (
	EOF TokenType = iota

	IDENT  //  Identifiers, whether user defined or reserved keywords
	STRING //  String Literals
	NUMBER //  123, -45.67
	BOOLEAN

	// Punctuations
	LBRACE   // {
	RBRACE   // }
	LPAREN   // (
	RPAREN   // )
	LBRACKET // [
	RBRACKET // ]
	COMMA    // ,
	COLON    // :
	DOT      // .
	QUESTION // ?
	AT       // @
	ATAT     // @@
	ASSIGN   // =

	//FOR THE FUTURE, WHEN IT'S WORTH MAKING MY OWN LSP

	EQUAL     // ==
	NOT_EQUAL // !=
	LT        // <
	GT        // >
	LTE       // <=
	GTE       // >=
	AND       // &&
	OR        // ||
	BANG      // !

	ILLEGAL // Unrecognized characters or unterminated literals
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

func LogTokens(tokens []Token) {
	for i, tok := range tokens {
		fmt.Printf("%03d %-12v %q (line: %d, col: %d)\n", i, tok.Type, tok.Value, tok.Line, tok.Col)
	}
}

func (t TokenType) String() string {
	switch t {
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case BOOLEAN:
		return "BOOLEAN"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case COMMA:
		return "COMMA"
	case COLON:
		return "COLON"
	case DOT:
		return "DOT"
	case QUESTION:
		return "QUESTION"
	case AT:
		return "AT"
	case ATAT:
		return "ATAT"
	case ASSIGN:
		return "ASSIGN"
	case EQUAL:
		return "EQUAL"
	case NOT_EQUAL:
		return "NOT_EQUAL"
	case LT:
		return "LT"
	case GT:
		return "GT"
	case LTE:
		return "LTE"
	case GTE:
		return "GTE"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case BANG:
		return "BANG"
	case ILLEGAL:
		return "ILLEGAL"
	default:
		return "UNKNOWN"
	}
}
