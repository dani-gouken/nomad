package vm

import (
	"fmt"
)

type TypeRegistrar struct {
	data map[string]*RuntimeType
}

func (r *TypeRegistrar) Add(t *RuntimeType) error {
	if r.Has(t.GetName()) {
		return fmt.Errorf("cannot redeclare type  %s", t.GetName())
	}
	r.data[t.GetName()] = t
	return nil
}

func (r *TypeRegistrar) Get(name string) (*RuntimeType, error) {
	if !r.Has(name) {
		return nil, fmt.Errorf(fmt.Sprintf("unknown type [%s]", name))
	}
	return r.data[name], nil
}

func (r *TypeRegistrar) Has(name string) bool {
	_, ok := r.data[name]
	return ok
}

func NewTypeRegistrar() TypeRegistrar {
	r := TypeRegistrar{
		data: make(map[string]*RuntimeType),
	}
	r.Add(MakeIntType())
	r.Add(MakeBoolype())
	return r
}
