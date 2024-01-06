package vm

import (
	"fmt"
)

type TypeRegistrar struct {
	data           map[string]*RuntimeType
	compositeTypes map[string]*CompositeType
}

func (r *TypeRegistrar) Add(t *RuntimeType) error {
	if r.Has(t.GetName()) {
		return fmt.Errorf("cannot redeclare type  %s", t.GetName())
	}
	r.data[t.GetName()] = t
	return nil
}

func (r *TypeRegistrar) GetPossible(typeName string) ([]*RuntimeType, error) {
	possibleTypes := []*RuntimeType{}
	runtimeType, err := r.Get(typeName)
	if err == nil {
		possibleTypes = append(possibleTypes, runtimeType)
		return possibleTypes, nil
	}

	composite, err := r.GetComposite(typeName)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(composite.Cases); i++ {
		possibleTypeName := composite.Cases[i]
		possibleType, err := r.Get(possibleTypeName)
		if err != nil {
			return nil, err
		}
		possibleTypes = append(possibleTypes, possibleType)

	}
	return possibleTypes, nil
}

func (r *TypeRegistrar) AddCompositeType(t *CompositeType) {
	r.compositeTypes[t.Alias] = t
}

func (r *TypeRegistrar) Get(name string) (*RuntimeType, error) {
	if !r.Has(name) {
		return nil, fmt.Errorf(fmt.Sprintf("unknown type [%s]", name))
	}
	return r.data[name], nil
}

func (r *TypeRegistrar) GetComposite(name string) (*CompositeType, error) {
	if !r.HasComposite(name) {
		return nil, fmt.Errorf(fmt.Sprintf("unknown type [%s]", name))
	}
	return r.compositeTypes[name], nil
}

func (r *TypeRegistrar) Has(name string) bool {
	_, ok := r.data[name]
	return ok
}
func (r *TypeRegistrar) HasComposite(name string) bool {
	_, ok := r.compositeTypes[name]
	return ok
}

func NewTypeRegistrar() TypeRegistrar {
	r := TypeRegistrar{
		data:           make(map[string]*RuntimeType),
		compositeTypes: make(map[string]*CompositeType),
	}

	r.Add(MakeIntType())
	r.Add(MakeBoolype())
	r.Add(MakeFloatType())

	num := NewCompositeType(NUM_TYPE, INT_TYPE, FLOAT_TYPE)
	r.AddCompositeType(&num)
	return r
}
