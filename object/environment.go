package object

type Environment struct {
	store map[string]Variable
	outer *Environment
}

type Variable struct {
	Value Object
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

func (e *Environment) Set(name string, val Object, isMut bool) Object {
	e.store[name] = Variable{Value: val, IsMut: isMut}
	return val
}

func (e *Environment) Reassign(name string, val Object) (Object, bool, bool)  {// return value, exist, ismut
	variable, ok := e.store[name]
	
	if !ok {
		// check for outer scope
		if e.outer != nil {
			return e.outer.Reassign(name, val)
		}
		
		// did not find the identifier
		return nil, false, false
	}

	if !variable.IsMut {
		return nil, true, false
	}

	e.store[name] = Variable{Value: val, IsMut: true}
	return val, true, true
}