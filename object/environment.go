package object

type Environment struct {
	store map[string]Variable
	outer *Environment
}

type Variable struct {
	Value Object
	VariableType string
	IsMut bool
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Variable)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj.Value, ok = e.outer.Get(name)
	}
	return obj.Value, ok
}

func (e *Environment) Set(name string, val Object, isMut bool, variableType string) Object {
	e.store[name] = Variable{Value: val, IsMut: isMut, VariableType: variableType }
	return val
}

func (e *Environment) Reassign(name string, val Object, variableType string) (Object, bool, bool, bool)  {// return value, exist, ismut, isSameType
	variable, ok := e.store[name]
	
	if !ok {
		// check for outer scope
		if e.outer != nil {
			return e.outer.Reassign(name, val, variableType)
		}
		
		// did not find the identifier
		return nil, false, false, false
	}

	if !variable.IsMut {
		return variable.Value, true, false, false // found but is not mutable
	}

	if variable.VariableType != variableType {
		return variable.Value, true, true, false
	}

	e.store[name] = Variable{Value: val, IsMut: true, VariableType: variableType}
	return val, true, true, true
}