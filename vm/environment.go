package vm

import (
	"fmt"

	"github.com/dani-gouken/nomad/runtime/data"
	"github.com/dani-gouken/nomad/runtime/types"
)

const ROOT_SCOPE = 0

type Scope struct {
	variables map[string]*data.RuntimeValue
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

func (s *Scope) DeclareVariable(name string, runtimeValue *data.RuntimeValue, declaredType types.RuntimeType) error {
	err := declaredType.Match(runtimeValue.RuntimeType)
	if err != nil {
		return fmt.Errorf("type mismatch, could not assign value of type %s to the variable %s declared as %s", runtimeValue.RuntimeType.GetName(), name, declaredType.GetName())
	}
	s.variables[name] = &data.RuntimeValue{
		RuntimeType: declaredType,
		Value:       runtimeValue.Value,
	}
	return nil
}

func (s *Scope) UnsetVariable(name string) {
	delete(s.variables, name)
}

func (s *Scope) GetVariable(name string) (*data.RuntimeValue, error) {
	value, ok := s.variables[name]
	if !ok {
		return value, fmt.Errorf("undefined variable %s", name)
	}
	return value, nil
}
func (e *Environment) DeclareVariable(
	name string,
	runtimeValue *data.RuntimeValue,
	declaredType types.RuntimeType,
) error {
	scope, err := e.GetCurrentScope()
	if err != nil {
		return err
	}
	return scope.DeclareVariable(name, runtimeValue, declaredType)
}

func (e *Environment) UnsetVariable(name string) error {
	scope, err := e.GetCurrentScope()
	if err != nil {
		return err
	}
	scope.UnsetVariable(name)
	return nil
}

func (e *Environment) GetVariable(name string) (*data.RuntimeValue, error) {
	scope, err := e.GetCurrentScope()
	if err != nil {
		return nil, err
	}
	return scope.GetVariable(name)
}

func NewEnvironment() Environment {
	scopes := make(map[int]Scope)
	scopes[ROOT_SCOPE] = Scope{
		variables: map[string]*data.RuntimeValue{},
	}
	return Environment{
		currentScope: ROOT_SCOPE,
		scopes:       scopes,
	}
}
