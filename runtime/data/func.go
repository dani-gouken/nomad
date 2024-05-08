package data

import (
	"fmt"

	"github.com/dani-gouken/nomad/runtime/types"
)

type Parameter struct {
	Name         string
	DefaultValue RuntimeValue
	HasDefault   bool
	RuntimeType  types.RuntimeType
}

type FuncSignature struct {
	ReturnType types.RuntimeType
	Parameters []Parameter
	names      map[string]int
}

type RuntimeFunc struct {
	Begin     int
	Tag       string
	Signature FuncSignature
}

func (s *FuncSignature) AsType() *types.FuncType {
	// TODO: cache this
	funcType := types.NewFuncType()
	for _, param := range s.Parameters {
		funcType.AddParam(param.RuntimeType)
	}
	funcType.SetRet(s.ReturnType)
	return funcType
}

func (f *RuntimeFunc) AddParam(name string, t types.RuntimeType, defaultValue RuntimeValue) error {
	_, ok := f.Signature.names[name]
	if ok {
		return fmt.Errorf("cannot redeclare parameter [%s]", name)
	}
	f.Signature.names[name] = len(f.Signature.names)
	f.Signature.Parameters = append(f.Signature.Parameters, Parameter{
		HasDefault:   defaultValue != RuntimeValue{},
		DefaultValue: defaultValue,
		Name:         name,
		RuntimeType:  t,
	})
	return nil
}

func (f *RuntimeFunc) SetRet(t types.RuntimeType) {
	f.Signature.ReturnType = t
}

func NewRuntimeFunc(t *types.Registrar, beginPtr int) *RuntimeFunc {
	return &RuntimeFunc{
		Begin: beginPtr,
		Tag:   "closure",
		Signature: FuncSignature{
			ReturnType: t.GetOrPanic("void"),
			Parameters: []Parameter{},
			names:      make(map[string]int),
		},
	}
}
