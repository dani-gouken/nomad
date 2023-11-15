package vm

import (
	"fmt"
)

const ROOT_SCOPE = 0

type Scope struct {
	variables map[string]*RuntimeValue
	parent    int
	id        int
}

type Environment struct {
	scopes       map[int]Scope
	currentScope int
	scopeCounter int
}

func (e *Environment) PushScope() Scope {
	nextScopeId := e.scopeCounter + 1
	scope := Scope{
		parent: e.currentScope,
		id:     nextScopeId,
	}
	e.scopes[nextScopeId] = scope
	return scope
}
func (e *Environment) GetCurrentScope() (*Scope, error) {

	scope, ok := e.scopes[e.currentScope]
	if !ok {
		return &scope, fmt.Errorf("scope with id [%d] not found", e.currentScope)
	}
	return &scope, nil
}

func (e *Scope) isRoot() bool {
	return e.id == ROOT_SCOPE
}

func (e *Environment) PopScope() error {
	scope, err := e.GetCurrentScope()
	if err != nil {
		return err
	}
	if scope.isRoot() {
		return fmt.Errorf("cannot pop root scope")
	}
	parent := scope.parent
	delete(e.scopes, e.currentScope)
	e.currentScope = parent
	return nil
}

func (s *Scope) SetVariable(name string, varType *RuntimeType, value interface{}) {
	s.variables[name] = &RuntimeValue{
		TypeName: varType.name,
		Value:    value,
	}
}

func (s *Scope) UnsetVariable(name string) {
	delete(s.variables, name)
}

func (s *Scope) GetVariable(name string) (*RuntimeValue, error) {
	value, ok := s.variables[name]
	if !ok {
		return value, fmt.Errorf("undefined variable %s", name)
	}
	return value, nil
}
func (e *Environment) SetVariable(name string, declaredType *RuntimeType, constantValue *RuntimeValue) error {
	scope, err := e.GetCurrentScope()
	if err != nil {
		return err
	}
	if declaredType.name != constantValue.TypeName {
		return fmt.Errorf("type mismatch, could not assign value of type [%s(%v)] to the variable [%s(%s)]", constantValue.TypeName, constantValue.Value, name, declaredType.name)
	}
	scope.SetVariable(name, declaredType, constantValue.Value)
	return nil
}

func (e *Environment) UnsetVariable(name string) error {
	scope, err := e.GetCurrentScope()
	if err != nil {
		return err
	}
	scope.UnsetVariable(name)
	return nil
}

func (e *Environment) GetVariable(name string) (*RuntimeValue, error) {
	scope, err := e.GetCurrentScope()
	if err != nil {
		return nil, err
	}
	return scope.GetVariable(name)
}

func NewEnvironment() Environment {
	scopes := make(map[int]Scope)
	scopes[ROOT_SCOPE] = Scope{
		variables: map[string]*RuntimeValue{},
	}
	return Environment{
		currentScope: ROOT_SCOPE,
		scopes:       scopes,
	}
}
