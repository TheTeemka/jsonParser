package token

type Type string

const (
	LeftBrace    Type = "{"
	RightBrace   Type = "}"
	LeftBracket  Type = "]"
	RightBracket Type = "["

	Colon Type = ":"
	Comma Type = ","

	String Type = "STRING"
	Number Type = "NUMBER"
	Bool   Type = "BOOL"
	Null   Type = "NULL"

	EOF     Type = "EOF"
	ILLEGAL Type = "ILLEGAL"
)

type Token struct {
	Tp      Type
	Literal string
	Line    int
	Start   int
	End     int
}
