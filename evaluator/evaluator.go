package evaluator

import (
	"potatolang/ast"
	"potatolang/object"
)

var (
	TRUE = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NIL = &object.Nil{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

		// Statements
		case *ast.Program:
			return evalStatements(node.Statements)
			
		case *ast.ExpressionStatement:
			return Eval(node.Expression)

		// Expressions
		case *ast.NumberLiteral:
			return &object.Number{Value: node.Value}
		case *ast.Boolean:
			return nativeBoolToBooleanObject(node.Value)
		case *ast.PrefixExpression:
			right := Eval(node.Right)
			return evalPrefixExpression(node.Operator, right)
			
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	} 
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
		case "!":
			return evalBangOperatorExpression(right)
		case "-":
			return evalMinusOperatorExpression(right)
		default:
			return NIL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
		case TRUE:
			return FALSE
		case FALSE:
			return TRUE
		case NIL:
			return TRUE
		default:
			return FALSE
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.NUMBER_OBJ {
		return NIL
	}

	value := right.(*object.Number).Value
	return &object.Number{Value: -value}
}