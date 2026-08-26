package object

import (
	"bytes"
	"fmt"
	"potatolang/ast"
	"strings"
)

const (
	NUMBER_OBJ = "number"
	BOOLEAN_OBJ = "boolean"
	NIL_OBJ = "nil"
	VOID_OBJ = "void"
	RETURN_VALUE_OBJ = "return_value"
	ERROR_OBJ = "error"
	FUNCTION_OBJ = "function"
	STRING_OBJ = "string"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Number struct {
	Value float64
}
func (n *Number) Inspect() string { return fmt.Sprintf("%g", n.Value)}
func (n *Number) Type() ObjectType { return NUMBER_OBJ }


type Boolean struct {
	Value bool
}
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value)}
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }


type Nil struct{}
func (n *Nil) Inspect() string { return "nil" }
func (n *Nil) Type() ObjectType { return NIL_OBJ }

type Void struct{}
func (v *Void) Inspect() string { return "void" }
func (v *Void) Type() ObjectType { return VOID_OBJ }

type ReturnValue struct {
	Values []Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string { 
	values := []string{}
	for _, val := range rv.Values {
		values = append(values, val.Inspect())
	}
	return strings.Join(values, ", ")
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

type Function struct {
	Parameters []*ast.TypedIdentifier
	ReturnTypes []string 
	Body *ast.BlockStatement
	Env *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	out.WriteString("func")
	out.WriteString("(")
	for i, param := range f.Parameters {
		if f.Parameters[i].IsMut {
			out.WriteString("mut ")
		}
		out.WriteString(param.String())
		out.WriteString(": ")
		out.WriteString(f.Parameters[i].Type)
		out.WriteString(", ")
	}

	out.WriteString(")")
	
	if len(f.ReturnTypes) >= 1 {
		returnTypes := []string{}
		out.WriteString("-> ")
		for _, returnType := range f.ReturnTypes {
			returnTypes = append(returnTypes, returnType)
		}
		out.WriteString(strings.Join(returnTypes, ", "))
	}
	out.WriteString("{\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string { return s.Value }