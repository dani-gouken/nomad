package vm

import (
	"fmt"

	nomadError "github.com/dani-gouken/nomad/errors"
	"github.com/dani-gouken/nomad/tokenizer"
)

type TypeRegistrar struct {
	data           map[string]*RuntimeType
	compositeTypes map[string]*CompositeType
}

func (r *TypeRegistrar) Add(name string, t *RuntimeType, token tokenizer.Token) error {
	if r.Has(name) {
		return nomadError.RuntimeError(fmt.Sprintf("cannot redeclare type %s", name), token)
	}
	r.data[name] = t
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

	r.Add(INT_TYPE, MakeIntType(), tokenizer.Token{})
	r.Add(BOOL_TYPE, MakeBoolType(), tokenizer.Token{})
	r.Add(FLOAT_TYPE, MakeFloatType(), tokenizer.Token{})
	r.Add(TYPE_TYPE, MakeTypeType(), tokenizer.Token{})
	r.Add(ARRAY_TYPE, MakeArrayType(), tokenizer.Token{})

	num := NewCompositeType(NUM_TYPE, INT_TYPE, FLOAT_TYPE)
	r.AddCompositeType(&num)
	return r
}
