package schema

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

	if isLetter(ch) {
		val := l.readIdentifier()

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

	if ch == '"' {
		val, ok := l.readString()
		if !ok {
			return Token{Type: ILLEGAL, Value: "unterminated string literal", Line: line, Col: col}
		}
		return Token{Type: STRING, Value: val, Line: line, Col: col}
	}

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
		return Token{Type: ILLEGAL, Value: string(ch), Line: line, Col: col}
	case '|':
		if l.peek() == '|' {
			l.advance()
			return Token{Type: OR, Value: "||", Line: line, Col: col}
		}
		return Token{Type: ILLEGAL, Value: string(ch), Line: line, Col: col}
	case '@':
		if l.peek() == '@' {
			l.advance()
			return Token{Type: ATAT, Value: "@@", Line: line, Col: col}
		}
		return Token{Type: AT, Value: "@", Line: line, Col: col}
	default:
		return Token{Type: ILLEGAL, Value: string(ch), Line: line, Col: col}
	}
}
func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < l.length {
		ch := l.peek()
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			l.advance()
			continue
		}

		if ch == '/' && l.peekNext() == '/' {

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

func (l *Lexer) readString() (string, bool) {
	l.advance() //  open quote
	var result []byte
	for l.pos < l.length {
		ch := l.peek()
		if ch == '"' {
			l.advance() //  close quote
			return string(result), true
		}
		if ch == '\n' || ch == '\r' {
			return string(result), false
		}
		if ch == '\\' {
			l.advance() //  backslash escape
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
	return string(result), false
}
func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func ExtractTokens(file string) []Token {
	lex := NewLexer(file)
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
