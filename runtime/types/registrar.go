package types

import (
	"fmt"

	nomadError "github.com/dani-gouken/nomad/errors"
	"github.com/dani-gouken/nomad/tokenizer"
)

type Registrar struct {
	data map[string]RuntimeType
}

func (r *Registrar) Add(t RuntimeType, token tokenizer.Token) error {
	if r.Has(t.GetName()) {
		return nomadError.RuntimeError(fmt.Sprintf("cannot redeclare type %s", t.GetName()), token)
	}
	r.data[t.GetName()] = t
	return nil
}

func (r *Registrar) Get(name string) (RuntimeType, error) {
	if !r.Has(name) {
		return nil, fmt.Errorf(fmt.Sprintf("unknown type [%s]", name))
	}
	return r.data[name], nil
}

func (r *Registrar) GetOrPanic(name string) RuntimeType {
	t, err := r.Get(name)
	if err != nil {
		panic(err)
	}
	return t
}

func (r *Registrar) Has(name string) bool {
	_, ok := r.data[name]
	return ok
}

func NewRegistrar() Registrar {
	r := Registrar{
		data: make(map[string]RuntimeType),
	}

	r.Add(MakeVoidType(), tokenizer.Token{})
	r.Add(MakeIntType(), tokenizer.Token{})
	r.Add(MakeFloatType(), tokenizer.Token{})
	r.Add(MakeBoolType(), tokenizer.Token{})
	r.Add(MakeTypeType(), tokenizer.Token{})
	r.Add(MakeStringType(), tokenizer.Token{})

	return r
}
