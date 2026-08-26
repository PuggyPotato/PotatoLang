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
		{"let mut x = 5; if true { x = 10;}", []any{VOID}},
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
			checkExpected(t, evaluatedVal.Values[i], expectedVal)
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
			checkExpected(t, evaluatedVal.Values[i], expectedVal)
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
			"type mismatch: number + boolean",
		},
		{
			"5 + true; 5;",
			"type mismatch: number + boolean",
		},
		{
			"-true",
			"unknown operator: -boolean",
		},
		{
			"true + false;",
			"unknown operator: boolean + boolean",
		},
		{
			"5; true + false; 5",
			"unknown operator: boolean + boolean",
		},
		{
			"if 10 > 1 { true + false; }",
			"unknown operator: boolean + boolean",
		},
		{
			`if 10 > 1 {
				if 10 > 1 {
					return true + false;
				}
				return 1;
			}`,
			"unknown operator: boolean + boolean",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			"x = 10;",
			"x is undefined.",
		},
		{
			`
				let x = 10;
				x = 5;
			`,
			"x is not mutable.",
		},
		{
			"let x = func() -> number { return nil; }; x()",
			"return type mismatch: function expected to return number, got nil",
		},
		{
			"let x = func() { return 5; } x()",
			"type mismatch: expression expects 0 return values, but found 1.",
		},
		{
			"let a = func() -> number, number { return 5,10; }; let x = a();",
			"assignment mismatch: 1 variables but 2 values.",
		},
		{
			"let a, b = func() {}",
			"function declaration cannot be unpacked",
		},
		{
			"let mut x = 10; x = nil;",
			"TypeError: cannot assign 'nil' to variable of type 'number'",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: string - string",
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

func TestFunctionObject(t *testing.T) {
	input := "func(mut x: number) -> number { return x + 2; };"

	evaluated := testEval(input)
	function, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(function.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", function.Parameters)
	}

	if function.Parameters[0].Value != "x" {
		t.Fatalf("value is not 'x'. got=%q", function.Parameters[0].Value)
	}

	if function.Parameters[0].Type != "number" {
		t.Fatalf("type is not 'number'. got=%q", function.Parameters[0].Type)
	}

	if function.Parameters[0].IsMut != true {
		t.Fatalf("isMut is not true. got=%t", function.Parameters[0].IsMut)
	}

	if len(function.ReturnTypes) != 1 {
		t.Fatalf("function has wrong returnTypes. returnTypes=%+v", function.ReturnTypes)
	}

	if function.ReturnTypes[0] != "number" {
		t.Fatalf("returnType is not 'number'. got=%s", function.ReturnTypes[0])
	}

	expectedBody := "return (x + 2);"

	if function.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, function.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input string
		expected []any
	} {
		
		{"let identity = func(x: number) { x; }; identity(5);", []any{VOID}},
		{"let identity = func(x: number) -> number { return x; }; identity(5);", []any{5}},
		{"let double = func(x: number) -> (number, error) { return x * 2, nil; }; double(5);", []any{10, nil}},
		{"let add = func(x: number, y: number) -> number { return x + y; }; add(5, 5);", []any{10}},
		{"let add = func(x: number, y: number) -> number { return x + y; }; add(5 + 5, add(5, 5));", []any{20}},
		{"func(x: number) -> number { return x; }(5)", []any{5}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		if len(tt.expected) == 1 && tt.expected[0] == VOID {
			testVoidObject(t, evaluated)
			continue
		}

		if len(tt.expected) == 1 {
			checkExpected(t, evaluated, tt.expected[0])
			continue
		}

		evaluatedVal, ok := evaluated.(*object.ReturnValue)
		if !ok {
			t.Errorf("evaluated not *object.ReturnValue, got=%T (%+v)", evaluated, evaluated)
			continue
		}

		for i, expectedVal := range tt.expected {
			checkExpected(t, evaluatedVal.Values[i], expectedVal)
		}
	}
}

func checkExpected(t *testing.T, obj object.Object, expectedVal any) {
	switch expected := expectedVal.(type) {
		case int:
			testNumberObject(t, obj, float64(expected))
		case float64:
			testNumberObject(t, obj, expected)
		case *object.Nil:
			testNilObject(t, obj)
		case *object.Void:
			testVoidObject(t, obj)
	}
}

func TestReassignment(t *testing.T) {
	tests := []struct {
		input string
		expected float64
	} {
		{"let mut x = 5; x = 10; x;", 10},
		{"let mut x, mut y = 5, 10; x = y; x;", 10},
		{"let mut x = 5; if true { x = 10.69; } x;", 10.69},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input string
		expected any
	} {
		{`len("")`, 0},
		{`len("four")`, 4},
		{"len(1)", "argument to `len` not supported, got number"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
			 case int:
			 	testNumberObject(t, evaluated, float64(expected))
			 case string:
			 	errObj, ok := evaluated.(*object.Error)
				if !ok {
					 t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
					 continue
				}
				if errObj.Message != expected {
					 t.Errorf("wrong error message. expected=%q, got=%q",
						 expected, errObj.Message)
				}
			
		}

	}
}