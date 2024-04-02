package data

import (
	"fmt"

	"github.com/dani-gouken/nomad/runtime/types"
)

func OrBool(a *RuntimeValue, b *RuntimeValue) (*RuntimeValue, error) {
	err := types.ExpectedBoolType(a.RuntimeType)
	if err != nil {
		return nil, err
	}
	err = types.ExpectedBoolType(b.RuntimeType)
	if err != nil {
		return nil, err
	}

	aBool := a.Value.(bool)
	bBool := b.Value.(bool)
	return &RuntimeValue{
		RuntimeType: a.RuntimeType,
		Value:       aBool || bBool,
	}, nil
}

func AndBool(a *RuntimeValue, b *RuntimeValue) (*RuntimeValue, error) {
	err := types.ExpectedBoolType(a.RuntimeType)
	if err != nil {
		return nil, err
	}
	err = types.ExpectedBoolType(b.RuntimeType)
	if err != nil {
		return nil, err
	}

	aBool := a.Value.(bool)
	bBool := b.Value.(bool)
	return &RuntimeValue{
		RuntimeType: a.RuntimeType,
		Value:       aBool && bBool,
	}, nil
}

func ApplyBinaryOpToBool(symbol string, lhs *RuntimeValue, rhs *RuntimeValue) (*RuntimeValue, error) {
	switch symbol {
	case "&":
		return AndBool(lhs, rhs)
	case "|":
		return OrBool(lhs, rhs)
	default:
		return nil, fmt.Errorf("unsupported operand %s for type %s", symbol, lhs.RuntimeType.GetName())
	}
}
