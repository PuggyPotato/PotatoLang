package parser

import (
	"fmt"
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

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s, got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() not %s, got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestNumberLiteralExpression(t *testing.T) {
	input := "6.7;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.NumberLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 6.7 {
		t.Errorf("literal.Value is not %f, got=%f", 6.7, literal.Value)
	}
	if literal.TokenLiteral() != "6.7" {
		t.Errorf("literal.TokenLiteral() is not %s, got=%s", "6.7", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input string
		operator string
		value float64
	} {
		{"!5;", "!", 5},
		{"!50.678", "!", 50.678},
		{"-15;", "-", 15},
		{"-69.42", "-", 69.420},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program does not have 1 statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression) 
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression, got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testNumberLiteral(t, exp.Right, tt.value) {
			return
		}
	}
}

func testNumberLiteral(t *testing.T, nl ast.Expression, value float64) bool {
	num, ok := nl.(*ast.NumberLiteral)
	if !ok {
		t.Errorf("nl not *ast.NumberLiteral. got=%T", nl)
		return false
	}

	if num.Value != value {
		t.Errorf("num.Value not %g, got=%g", value, num.Value)
		return false
	}

	if num.TokenLiteral() != fmt.Sprintf("%g",value) {
		t.Errorf("num.TokenLiteral not %g, got=%s", value, num.TokenLiteral())
		return false
	}

	return true
}

func TestParseInfixExpression(t *testing.T) {
	infixTests := []struct{
		input string
		leftValue float64
		operator string
		rightValue float64
	} {
		{"5.6 + 5.2;", 5.6, "+", 5.2},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testNumberLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testNumberLiteral(t, exp.Right, tt.rightValue) {
			return 
		}
		
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct{
		input string
		expected string
	} {
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
 			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3.69 > 4",
			"((5 < 4) != (3.69 > 4))",
		},
		{
			"3.5 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3.5 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
  		"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type)	{
		case float64:
			return testNumberLiteral(t, exp, v)
		case int:
			return testNumberLiteral(t, exp, float64(v))
		case string:
			return testIdentifier(t, exp, v)
		case bool:
			return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func TestBooleanExpression(t *testing.T) {
	booleanTest := []struct {
		input string
		expected bool
	} {
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range booleanTest {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
	
		if len(program.Statements) != 1 {
			t.Fatalf("program does not have enough statements. got=%d", len(program.Statements))
		}
	
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}
	
		literal, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("stmt.Expression not *ast.Boolean. got=%T", stmt.Expression)
		}
		if literal.Value != tt.expected {
			t.Errorf("literal.Value is not %t, got=%t", tt.expected, literal.Value)
		}
		if literal.TokenLiteral() != fmt.Sprintf("%t",tt.expected){
			t.Errorf("literal.TokenLiteral() is not %t, got=%s", tt.expected, literal.TokenLiteral())
		}
		
	}
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok :=exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
		return false
	}

	return true
}