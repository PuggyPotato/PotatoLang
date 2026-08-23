package evaluator

import (
	"fmt"
	"potatolang/ast"
	"potatolang/object"
)

var (
	TRUE = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NIL = &object.Nil{}
	VOID = &object.Void{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

		// Statements
		case *ast.Program:
			return evalProgram(node, env)
			
		case *ast.ExpressionStatement:
			return Eval(node.Expression, env)

		case *ast.LetStatement:
			for i, val := range node.Values {
				evaluated := Eval(val, env)
				if isError(evaluated) {
					return evaluated
				}
				env.Set(node.Names[i].Value, evaluated)
			}

		case *ast.Identifier:
			return evalIdentifier(node, env)

		// Expressions
		case *ast.NumberLiteral:
			return &object.Number{Value: node.Value}
			
		case *ast.Boolean:
			return nativeBoolToBooleanObject(node.Value)
			
		case *ast.PrefixExpression:
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return evalPrefixExpression(node.Operator, right)
			
		case *ast.InfixExpression:
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return evalInfixExpression(node.Operator, left, right)
			
		case *ast.BlockStatement:
			return evalBlockStatement(node, env)
			
		case *ast.IfExpression:
			return evalIfExpression(node, env)

		case *ast.ReturnStatement:
			var evaluatedVals []object.Object

			for _, val := range node.ReturnValues {
				evaluated := Eval(val, env)
				if isError(evaluated) {
					return evaluated
				}
				evaluatedVals = append(evaluatedVals, evaluated)
			}

			return &object.ReturnValue{Values: evaluatedVals}

		case *ast.NilLiteral:
			return NIL
			
	}
	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
			case *object.ReturnValue:
				return result
			case *object.Error:
				return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}

	}
	return VOID
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
			return VOID
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
		case VOID:
			return TRUE
		default:
			return FALSE
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.NUMBER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Number).Value
	return &object.Number{Value: -value}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
		case left.Type() != right.Type():
			return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
		case left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ:
			return evalNumberInfixExpression(operator, left, right)
		case operator == "==":
			return nativeBoolToBooleanObject(left == right)
		case operator == "!=":
			return nativeBoolToBooleanObject(left != right)
		default:
			return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalNumberInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Number).Value
	rightVal := right.(*object.Number).Value

	switch operator {
		case "+":
			return &object.Number{Value: leftVal + rightVal }
		case "-":
			return &object.Number{Value: leftVal - rightVal }
		case "*":
			return &object.Number{Value: leftVal * rightVal }
		case "/":
			return &object.Number{Value: leftVal / rightVal }
		case "<":
			return nativeBoolToBooleanObject(leftVal < rightVal)
		case ">":
			return nativeBoolToBooleanObject(leftVal > rightVal)
		case "==":
			return nativeBoolToBooleanObject(leftVal == rightVal)
		case "!=":
			return nativeBoolToBooleanObject(leftVal != rightVal)
		default:
			return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return VOID
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
		case NIL:
			return false
		case VOID:
			return false
		case TRUE:
			return true
		case FALSE:
			return false
		default: 
			return true
	}
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value )
	}

	return val
}