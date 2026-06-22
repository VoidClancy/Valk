package parser

import (
	"fmt"
)

type TokenType int

const (
	EOF TokenType = iota

	IDENT   // e.g. User, String, id, unique, now, db
	STRING  // e.g. "users", "DATABASE_URL"
	NUMBER  // e.g. 123, -45.67
	BOOLEAN // e.g. true, false

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

	EQUAL     // ==
	NOT_EQUAL // !=
	LT        // <
	GT        // >
	LTE       // <=
	GTE       // >=
	AND       // &&
	OR        // ||
	BANG      // !
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

type Lexer struct {
	input  string
	pos    int
	line   int
	col    int
	length int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		line:   1,
		col:    1,
		length: len(input),
	}
}

func (l *Lexer) peek() byte {
	if l.pos >= l.length {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) peekNext() byte {
	if l.pos+1 >= l.length {
		return 0
	}
	return l.input[l.pos+1]
}

func (l *Lexer) advance() byte {
	if l.pos >= l.length {
		return 0
	}
	ch := l.input[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return ch
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespaceAndComments()

	if l.pos >= l.length {
		return Token{Type: EOF, Value: "", Line: l.line, Col: l.col}
	}

	ch := l.peek()
	line := l.line
	col := l.col

	// Identifiers and Keywords
	if isLetter(ch) {
		val := l.readIdentifier()
		// Check for keywords
		t := IDENT
		switch val {
		case "true", "false":
			t = BOOLEAN
		}
		return Token{Type: t, Value: val, Line: line, Col: col}
	}

	// Numbers
	if isDigit(ch) || (ch == '-' && isDigit(l.peekNext())) {
		val := l.readNumber()
		return Token{Type: NUMBER, Value: val, Line: line, Col: col}
	}

	// Strings
	if ch == '"' {
		val := l.readString()
		return Token{Type: STRING, Value: val, Line: line, Col: col}
	}

	// Multi-character and single-character symbols
	l.advance()
	switch ch {
	case '{':
		return Token{Type: LBRACE, Value: "{", Line: line, Col: col}
	case '}':
		return Token{Type: RBRACE, Value: "}", Line: line, Col: col}
	case '(':
		return Token{Type: LPAREN, Value: "(", Line: line, Col: col}
	case ')':
		return Token{Type: RPAREN, Value: ")", Line: line, Col: col}
	case '[':
		return Token{Type: LBRACKET, Value: "[", Line: line, Col: col}
	case ']':
		return Token{Type: RBRACKET, Value: "]", Line: line, Col: col}
	case ',':
		return Token{Type: COMMA, Value: ",", Line: line, Col: col}
	case ':':
		return Token{Type: COLON, Value: ":", Line: line, Col: col}
	case '.':
		return Token{Type: DOT, Value: ".", Line: line, Col: col}
	case '?':
		return Token{Type: QUESTION, Value: "?", Line: line, Col: col}
	case '=':
		if l.peek() == '=' {
			l.advance()
			return Token{Type: EQUAL, Value: "==", Line: line, Col: col}
		}
		return Token{Type: ASSIGN, Value: "=", Line: line, Col: col}
	case '!':
		if l.peek() == '=' {
			l.advance()
			return Token{Type: NOT_EQUAL, Value: "!=", Line: line, Col: col}
		}
		return Token{Type: BANG, Value: "!", Line: line, Col: col}
	case '<':
		if l.peek() == '=' {
			l.advance()
			return Token{Type: LTE, Value: "<=", Line: line, Col: col}
		}
		return Token{Type: LT, Value: "<", Line: line, Col: col}
	case '>':
		if l.peek() == '=' {
			l.advance()
			return Token{Type: GTE, Value: ">=", Line: line, Col: col}
		}
		return Token{Type: GT, Value: ">", Line: line, Col: col}
	case '&':
		if l.peek() == '&' {
			l.advance()
			return Token{Type: AND, Value: "&&", Line: line, Col: col}
		}
		panic(fmt.Sprintf("unknown character %q at line %d, col %d", ch, line, col))
	case '|':
		if l.peek() == '|' {
			l.advance()
			return Token{Type: OR, Value: "||", Line: line, Col: col}
		}
		panic(fmt.Sprintf("unknown character %q at line %d, col %d", ch, line, col))
	case '@':
		if l.peek() == '@' {
			l.advance()
			return Token{Type: ATAT, Value: "@@", Line: line, Col: col}
		}
		return Token{Type: AT, Value: "@", Line: line, Col: col}
	default:
		// Unknown character
		panic(fmt.Sprintf("unknown character %q at line %d, col %d", ch, line, col))
	}
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' || ch == '$'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < l.length {
		ch := l.peek()
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			l.advance()
			continue
		}
		// Comment starting with //
		if ch == '/' && l.peekNext() == '/' {
			// Skip comment until end of line
			l.advance() // first /
			l.advance() // second /
			for l.pos < l.length && l.peek() != '\n' {
				l.advance()
			}
			continue
		}
		break
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.pos
	l.advance() // consume first letter
	for l.pos < l.length {
		ch := l.peek()
		if isLetter(ch) || isDigit(ch) {
			l.advance()
		} else {
			break
		}
	}
	return l.input[start:l.pos]
}

func (l *Lexer) readNumber() string {
	start := l.pos
	if l.peek() == '-' {
		l.advance()
	}
	for l.pos < l.length && isDigit(l.peek()) {
		l.advance()
	}
	if l.pos < l.length && l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance() // consume .
		for l.pos < l.length && isDigit(l.peek()) {
			l.advance()
		}
	}
	return l.input[start:l.pos]
}

func (l *Lexer) readString() string {
	l.advance() // consume open quote
	var result []byte
	for l.pos < l.length {
		ch := l.peek()
		if ch == '"' {
			l.advance() // consume close quote
			return string(result)
		}
		if ch == '\\' {
			l.advance() // consume backslash
			if l.pos < l.length {
				escaped := l.advance()
				switch escaped {
				case 'n':
					result = append(result, '\n')
				case 't':
					result = append(result, '\t')
				case 'r':
					result = append(result, '\r')
				case '\\':
					result = append(result, '\\')
				case '"':
					result = append(result, '"')
				default:
					result = append(result, '\\', escaped)
				}
			}
			continue
		}
		result = append(result, l.advance())
	}
	panic(fmt.Sprintf("unterminated string literal at line %d", l.line))
}

func ExtractTokens(schema string) []Token {
	lex := NewLexer(schema)
	var tokens []Token
	for {
		tok := lex.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}
	return tokens
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
	default:
		return "UNKNOWN"
	}
}
