package data

import (
	"fmt"

	"github.com/dani-gouken/nomad/runtime/types"
)

type RuntimeValue struct {
	RuntimeType types.RuntimeType
	Value       interface{}
}

type RuntimeArray struct {
	Values []RuntimeValue
}

type RuntimeObject struct {
	fields map[string]*RuntimeValue
}

func (o *RuntimeObject) GetFields() map[string]*RuntimeValue {
	return o.fields
}

func (o *RuntimeObject) GetField(name string) (*RuntimeValue, error) {
	v, ok := o.fields[name]
	if !ok {
		return v, fmt.Errorf("trying to access undefined field [%s]", name)
	}
	return v, nil
}

func (o *RuntimeObject) SetField(name string, value RuntimeValue) error {
	o.fields[name] = &value
	return nil
}

func ApplyBinaryOp(t types.Registrar, symbol string, lhs *RuntimeValue, rhs *RuntimeValue) (*RuntimeValue, error) {
	err := lhs.RuntimeType.Match(rhs.RuntimeType)
	if err != nil {
		return nil, err
	}
	lhsType, err := types.ToScalarType(lhs.RuntimeType)
	if err != nil {
		return nil, err
	}

	if lhsType.IsFloat() {
		return ApplyBinaryOpToFloat(t, symbol, lhs, rhs)
	}

	if lhsType.IsInt() {
		return ApplyBinaryOpToInt(symbol, lhs, rhs)
	}

	if lhsType.IsString() {
		return ApplyBinaryOpToString(symbol, lhs, rhs)
	}

	return nil, fmt.Errorf("unsupported operand %s for type %s", symbol, lhs.RuntimeType.GetName())
}

func NewRuntimeObject() *RuntimeObject {
	return &RuntimeObject{
		fields: make(map[string]*RuntimeValue),
	}
}
