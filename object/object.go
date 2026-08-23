package object

import (
	"fmt"
	"strings"
)

const (
	NUMBER_OBJ = "NUMBER"
	BOOLEAN_OBJ = "BOOLEAN"
	NIL_OBJ = "NIL"
	VOID_OBJ = "VOID"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
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