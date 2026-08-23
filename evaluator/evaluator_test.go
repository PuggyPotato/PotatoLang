package evaluator

import (
	"potatolang/lexer"
	"potatolang/object"
	"potatolang/parser"
	"testing"
)

func TestEvalNumberExpression(t *testing.T) {
	tests := []struct {
		input string
		expected float64
	} {
		{"5", 5},
		{"10", 10},
		{"69.54", 69.54},
		{"-5.12", -5.12},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50.67 + 100.67 + -50", 0},
		{"5 * 2 + 10.23", 20.23},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10.92", 60.92},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 15.67", 42.67},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testNumberObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Number)
	if !ok {
		t.Errorf("object is not number. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%g, want=%g", result.Value, expected)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	} {
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	} {
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input string
		expected []any
	} {
		{"if true { return 10; }", []any{10}},
		{"if true { x = 10;}", []any{VOID}},
		{"if false { return 10; }", []any{VOID}},
		{"if 1 { return 10, 50; }", []any{10, 50}},
		{"if 1 < 2 { return 10; }", []any{10}},
		{"if 1 > 2 { return 10.50; }", []any{VOID}},
		{"if 1 > 2 { return 10; } else { return 20.50; }", []any{20.50}},
		{"if 1 < 2 { return 50.50; } else { return 20; }", []any{50.50}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		if len(tt.expected) == 1 && tt.expected[0] == VOID {
			testVoidObject(t, evaluated)
			continue
		}

		evaluatedVal, ok := evaluated.(*object.ReturnValue)
		if !ok {
			t.Fatalf("evaluated not *object.ReturnValue, got=%T (%+v)", evaluated, evaluated)
		}

		for i, expectedVal := range tt.expected {
			switch expected := expectedVal.(type) {
				case int:
					testNumberObject(t, evaluatedVal.Values[i], float64(expected))
				case float64:
					testNumberObject(t, evaluatedVal.Values[i], expected)
				case *object.Nil:
					testNilObject(t, evaluatedVal.Values[i])
				case *object.Void:
					testVoidObject(t, evaluatedVal.Values[i])
			}
		}
	}
}

func testVoidObject(t *testing.T, obj object.Object) bool {
	if obj != VOID {
		t.Errorf("object is not VOID. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}

func testNilObject(t *testing.T, obj object.Object) bool {
	if obj != NIL {
		t.Errorf("object is not Nil. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input string 
		expected []any
	} {
		{"return 10.6, 50;", []any{10.6, 50}},
		{"return 102; return 9;", []any{102}},
		{"return 2 * 5, 100.50; return 9;", []any{10, 100.50}},
		{"return 2 * 7, nil; return 9;", []any{14, NIL}},

		{
			`if 10 > 1 {
				if 10 > 1 {
					return 10;
				}
				return 1;
			}`, []any{10},
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		if len(tt.expected) == 1 && tt.expected[0] == VOID {
			testVoidObject(t, evaluated)
			continue
		}

		evaluatedVal, ok := evaluated.(*object.ReturnValue)
		if !ok {
			t.Fatalf("evaluated not *object.ReturnValue, got=%T (%+v)", evaluated, evaluated)
		}

		for i, expectedVal := range tt.expected {
			switch expected := expectedVal.(type) {
				case int:
					testNumberObject(t, evaluatedVal.Values[i], float64(expected))
				case float64:
					testNumberObject(t, evaluatedVal.Values[i], expected)
				case *object.Nil:
					testNilObject(t, evaluatedVal.Values[i])
				case *object.Void:
					testVoidObject(t, evaluatedVal.Values[i])
			}
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input string
		expectedMessage string
	} {
		{
			"5 + true;",
			"type mismatch: NUMBER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: NUMBER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if 10 > 1 { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`if 10 > 1 {
				if 10 > 1 {
					return true + false;
				}
				return 1;
			}`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T (%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input string
		expected float64
	} {
		{"let a = 5.67; a;", 5.67},
		{"let a = 2 * 2.5; a;", 5},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
		{"let mut a,b = 5, 10; b;", 10},
	}

	for _, tt := range tests {
		testNumberObject(t, testEval(tt.input), tt.expected)
	}
}