package types

import (
	"fmt"
)

type FuncType struct {
	anonymous  bool
	generic    bool
	returnType RuntimeType
	parameters []RuntimeType
}

func (f *FuncType) GetName() string {
	name := "func"
	if len(f.parameters) > 0 {
		name = name + "("
	}
	for i, parameter := range f.parameters {
		name = name + parameter.GetName()
		if i == len(f.parameters)-1 {
			name = name + ")"
		} else {
			name = name + ","
		}
	}

	if f.returnType != nil {
		name = name + " -> (" + f.returnType.GetName() + ")"
	}

	return name
}

func (f *FuncType) SetName(name string) {
	f.anonymous = false
}

func (o *FuncType) IsAnonymous() bool {
	return o.anonymous
}

func (o *FuncType) IsGeneric() bool {
	return o.generic
}

func (t *FuncType) Match(t2 RuntimeType) error {
	t2Func, err := ToFuncType(t2)
	if err != nil {
		return err
	}
	err = t.returnType.Match(t2Func.returnType)
	if err != nil {
		return fmt.Errorf("return type mismatch, %s", err.Error())
	}
	if len(t2Func.parameters) != len(t.parameters) {
		return fmt.Errorf("parameter length mismatch, %d expected, got %d", len(t.parameters), len(t2Func.parameters))
	}
	for i := range t.parameters {
		err := t.parameters[i].Match(t2Func.parameters[i])
		if err != nil {
			return fmt.Errorf("parameter %d type mismatch, %s", i, err.Error())
		}
	}
	return nil
}

func (f *FuncType) AddParam(t RuntimeType) {
	f.parameters = append(f.parameters, t)
}

func (f *FuncType) SetRet(t RuntimeType) {
	f.returnType = t
}

func NewFuncType() *FuncType {
	return &FuncType{
		anonymous:  true,
		generic:    true,
		parameters: []RuntimeType{},
	}
}

func IsFuncType(t RuntimeType) bool {
	_, err := ToFuncType(t)
	return err == nil
}

func ToFuncType(t RuntimeType) (*FuncType, error) {
	tObj, ok := t.(*FuncType)
	if !ok {
		return nil, fmt.Errorf("func type expected")
	}
	return tObj, nil
}
