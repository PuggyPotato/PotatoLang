package token

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	// Identifiers + Literals
	IDENT = "IDENT"
	NUMBER = "NUMBER" // 1, 2, 3, 4, 1.1, 1.2 ...
	STRING = "STRING"

	// Operators
	ASSIGN = "="
	PLUS = "+"
	MINUS = "-"
	BANG = "!"
	ASTERISK = "*"
	SLASH = "/"
	LT = "<"
	GT = ">"
	EQ = "=="
	NOT_EQ = "!="
	RIGHT_ARROW = "->"

	// Delimiters
	COMMA = ","
	SEMICOLON = ";"
	COLON = ":"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET = "LET"
	TRUE = "TRUE"
	FALSE = "FALSE"
	IF = "IF"
	ELSE = "ELSE"
	RETURN = "RETURN"
	MUT = "MUT"

	// Types
	NUMBER_TYPE = "NUMBER_TYPE"
	STRING_TYPE = "STRING_TYPE"
	ERROR_TYPE = "ERROR_TYPE"
	BOOLEAN_TYPE = "BOOLEAN_TYPE"

	VOID_VALUE = "VOID_VALUE"
	NIL_VALUE = "NIL_VALUE"
)

var keywords = map[string] TokenType {
	"func": FUNCTION,
	"let": LET,
	"true": TRUE,
	"false": FALSE,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
	"mut": MUT,
	"number": NUMBER_TYPE,
	"bool": BOOLEAN_TYPE,
	"string": STRING_TYPE,
	"error": ERROR_TYPE,
	"nil": NIL_VALUE,
	"void": VOID_VALUE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}