package object

type Environment struct {
	store map[string]Variable
}

type Variable struct {
	Value Object
	IsMut bool
}

func NewEnvironment() *Environment {
	s := make(map[string]Variable)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj.Value, ok
}

func (e *Environment) Set(name string, val Object, isMut bool) Object {
	e.store[name] = Variable{Value: val, IsMut: isMut}
	return val
}
