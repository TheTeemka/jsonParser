package lexer

import (
	"jsonParser/token"
	"unicode"
)

type Lexer struct {
	Input   []byte
	char    byte
	pos     int
	nextPos int
	line    int
}

func New(data []byte) *Lexer {
	return &Lexer{
		Input: data,
	}
}
func (l *Lexer) NextToken() token.Token {
	var t token.Token

	l.readChar()
	l.skimWhiteSpaces()

	switch l.char {
	case '{':
		t = newToken(token.LeftBrace, l.line, l.pos, l.pos+1, l.char)
	case '}':
		t = newToken(token.RightBrace, l.line, l.pos, l.pos+1, l.char)
	case '[':
		t = newToken(token.LeftBracket, l.line, l.pos, l.pos+1, l.char)
	case ']':
		t = newToken(token.RightBracket, l.line, l.pos, l.pos+1, l.char)
	case ':':
		t = newToken(token.Colon, l.line, l.pos, l.pos+1, l.char)
	case ',':
		t = newToken(token.Comma, l.line, l.pos, l.pos+1, l.char)
	case '"':
		t.Tp = token.String
		t.Start = l.pos
		t.Literal = l.readString()
		t.End = l.pos + 1
		t.Line = l.line
	case 0:
		t.Tp = token.EOF
		t.Literal = "End"
		t.Line = l.line
	default:
		if unicode.IsDigit(rune(l.char)) || l.char == '-' || l.char == '+' { // number
			t.Tp = token.Number
			t.Start = l.pos
			t.Literal = l.readNumber()
			t.Start = l.pos + 1
			t.Line = l.line
			if !((t.Literal[0] == '-' || t.Literal[0] == '+') && len(t.Literal) == 1) { // check if it is not just "+" or "-"
				break
			}
		} else { // bool or null
			t.Start = l.pos
			t.Literal = l.readLeft()
			t.Start = l.pos + 1
			t.Line = l.line
			if t.Literal == "true" || t.Literal == "false" {
				t.Tp = token.Bool
				break
			} else if t.Literal == "null" {
				t.Tp = token.Null
				break
			}
		}
		return t
	}
	return t
}

func (l *Lexer) readChar() {
	if l.pos+1 >= len(l.Input) {
		l.char = 0
	} else {
		l.char = l.Input[l.nextPos]
	}
	l.pos = l.nextPos
	l.nextPos++
}

func (l *Lexer) skimWhiteSpaces() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' {
		if l.char == '\n' {
			l.line++
		}
		l.readChar()
	}
}

func newToken(tp token.Type, line, start, end int, char ...byte) token.Token {
	return token.Token{
		Tp:      tp,
		Line:    line,
		Start:   start,
		End:     end,
		Literal: string(char),
	}
}

func (l *Lexer) readString() string {
	st := l.pos + 1
	l.readChar()
	for l.char != '"' {
		l.readChar()
	}
	return string(l.Input[st:l.pos])
}

func (l *Lexer) readNumber() string {
	st := l.pos
	var dot bool

	if l.char == '-' {
		l.readChar()
	}

	for unicode.IsDigit(rune(l.char)) || l.char == '.' {
		if l.char == '.' {
			if dot {
				break
			}
			dot = true
		}
		l.readChar()
	}
	l.nextPos--
	return string(l.Input[st:l.pos])
}

func (l *Lexer) readLeft() string {
	st := l.pos
	for l.char != 0 && len(l.Input[st:l.pos+1]) <= 5 {
		s := string(l.Input[st : l.pos+1])
		if s == "true" || s == "false" || s == "null" {
			break
		}

		l.readChar()
	}
	return string(l.Input[st : l.pos+1])
}
