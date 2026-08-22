package parser

import (
	"potatolang/ast"
	"potatolang/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 67.67;
	let mut z = 50;
	let a,mut b = 50, 59.42;
	let mut c,_ = 90.42, 50.40;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	
	if program == nil {
		t.Fatalf("ParseProgram() returned nil.")
	}
	if len(program.Statements) != 5 {
		t.Fatalf("program.Statements does not contain 5 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier []string
	} {
		{[]string{"x"}},
		{[]string{"y"}},
		{[]string{"z"}},
		{[]string{"a", "b"}},
		{[]string{"c", "_"}},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, names []string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral() not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
	}

	if len(letStmt.Names) != len(names) {
		t.Errorf("expected %d identifiers, got=%d", len(names), len(letStmt.Names))
	}

	for i, name := range names {
		if letStmt.Names[i].Value != name {
			t.Errorf("letStmt.Names[%d].Value not '%s', got=%s", i, name, letStmt.Names[i].Value)
		}

		if letStmt.Names[i].TokenLiteral() != name {
			t.Errorf("letStmt.Names[%d].TokenLiteral() not '%s', got=%s", i, name, letStmt.Names[i].TokenLiteral())
		}
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return 
	}

	t.Errorf("parser had %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	input := `
			return 5;
			return 10;
			return 1, 2;
			return 5, nil;
			`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain 4 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement, got=%T",stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got=%q",returnStmt.TokenLiteral())
		}
	}
}