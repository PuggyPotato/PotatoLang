package lexer

import (
	"potatolang/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	// Todo: Add arrow return types to the function
	// i.e. func(x, y) -> number { return x + y }
	input := `
				let x = 5;
				let y = 10.5;

				let add = func(x, y) {
					return x + y;
				}
				
				let z = add(x,y);

				let mut a = 5;
				a = 100.50;
`

	tests := []struct {
		expectedType token.TokenType
		expectedLiteral string
	} {
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "y"},
		{token.ASSIGN, "="},
		{token.NUMBER, "10.5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "func"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.LET, "let"},
		{token.IDENT, "z"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.MUT, "mut"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.NUMBER, "100.50"},
		{token.SEMICOLON, ";"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%s, got=%s", i, tt.expectedLiteral, tok.Literal)
		}
	}
}