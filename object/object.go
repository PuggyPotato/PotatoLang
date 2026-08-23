package object

import "fmt"

const (
	NUMBER_OBJ = "NUMBER"
	BOOLEAN_OBJ = "BOOLEAN"
	NIL_OBJ = "NIL"
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
