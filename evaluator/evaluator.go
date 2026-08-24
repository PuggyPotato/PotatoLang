package evaluator

import (
	"fmt"
	"potatolang/ast"
	"potatolang/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NIL   = &object.Nil{}
	VOID  = &object.Void{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.LetStatement:
		return EvalLetStatement(node, env)

	case *ast.AssignStatement:
		return EvalReassignStatement(node, env)

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

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		returnTypes := node.ReturnTypes
		return &object.Function{Parameters: params, Env: env, Body: body, ReturnTypes: returnTypes}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

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
		return &object.Number{Value: leftVal + rightVal}
	case "-":
		return &object.Number{Value: leftVal - rightVal}
	case "*":
		return &object.Number{Value: leftVal * rightVal}
	case "/":
		return &object.Number{Value: leftVal / rightVal}
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

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
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
		return newError("identifier not found: %s", node.Value)
	}

	return val
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}
	
	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	if isError(evaluated) {
		return evaluated
	}

	if returnTypes, ok := evaluated.(*object.ReturnValue); ok {
		if len(returnTypes.Values) != len(function.ReturnTypes) {
			return newError("type mismatch: expression expects %d return values, but found %d.", len(function.ReturnTypes), len(returnTypes.Values) )
		}
		for i, returnType := range returnTypes.Values {
			actualType := returnType.Type()
			
			if string(actualType) != function.ReturnTypes[i] {
				isNil := function.ReturnTypes[i] == "error" && string(actualType) == object.NIL_OBJ 
				if !isNil {
					return newError("return type mismatch: function expected to return %s, got %s", function.ReturnTypes[i], string(actualType))
				}
			}
		}
	}
	
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx], param.IsMut)
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValues, ok := obj.(*object.ReturnValue); ok {
		if len(returnValues.Values) == 1 {
			return returnValues.Values[0]
		}
	}
	
	return obj
}

func EvalLetStatement(node *ast.LetStatement, env *object.Environment) object.Object{
	if len(node.Names) > 1 && len(node.Values) == 1 {
		evaluated := Eval(node.Values[0], env)
		if isError(evaluated) {
			return evaluated
		}

		returnVal, ok := evaluated.(*object.ReturnValue)
		if !ok { // returnVal = nil
			// function declaration should not be able to be unpacked e.g. a,b = func() {}
			if _, isFunc := evaluated.(*object.Function); isFunc {
				return newError("function declaration cannot be unpacked")
			}
			// used for common error with 1 value, when returnVal fails  
			return newError("assignment mismatch: %d variables but %d value", len(node.Names), 1)
		}

		if len(returnVal.Values) != len(node.Names) {
			return newError("assignment mismatch: %d variables but %d value", len(node.Names), len(returnVal.Values))
		}
		

		for i, name := range node.Names {
			env.Set(name.Value, returnVal.Values[i], name.IsMut)
		}
		return VOID
	}

	if len(node.Values) != len(node.Names) {
		return newError("assignment mismatch: %d variables but %d value", len(node.Names), len(node.Values))
	}

	for i, val := range node.Values {
		evaluated := Eval(val, env)
		if isError(evaluated) {
			return evaluated
		}
		
		if returnVal, ok := evaluated.(*object.ReturnValue); ok {
			if len(returnVal.Values) == 1 {
				evaluated = returnVal.Values[0]
			}

			if len(returnVal.Values) > 1 {
				return newError("assignment mismatch: 1 variables but %d values.", len(returnVal.Values))
			}
		}
		
		env.Set(node.Names[i].Value, evaluated, node.Names[i].IsMut)
	}
	return VOID
}

func EvalReassignStatement(node *ast.AssignStatement, env *object.Environment) object.Object{
	if len(node.Names) > 1 && len(node.Values) == 1 {
		evaluated := Eval(node.Values[0], env)
		if isError(evaluated) {
			return evaluated
		}

		returnVal, ok := evaluated.(*object.ReturnValue)
		if !ok || len(returnVal.Values) != len(node.Names) {
			return newError("assignment mismatch: %d variables but %d value", len(node.Names), len(returnVal.Values))
		}

		for i, name := range node.Names {
			_, exist, isMut := env.Reassign(name.Value, returnVal.Values[i])
			
			if !exist {
				return newError("undefined: %s",name.Value)
			}

			if !isMut {
				return newError("%s is not mutable.",name.Value)
			}
			
		}
		return VOID
	}

	if len(node.Values) != len(node.Names) {
		return newError("assignment mismatch: %d variables but %d value", len(node.Names), len(node.Values))
	}

	for i, val := range node.Values {
		evaluated := Eval(val, env)
		if isError(evaluated) {
			return evaluated
		}
		_, exist, isMut := env.Reassign(node.Names[i].Value, evaluated)
		
		if !exist {
			return newError("%s is undefined.",node.Names[i].Value)
		}

		if !isMut {
			return newError("%s is not mutable.",node.Names[i].Value)
		}
		
		env.Reassign(node.Names[i].Value, evaluated)
	}
	return VOID
}