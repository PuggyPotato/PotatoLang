package evaluator

import (
	"potatolang/ast"
	"potatolang/object"
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