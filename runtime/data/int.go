package data

import (
	"fmt"

	"github.com/dani-gouken/nomad/runtime/types"
)

func AddInt(a *RuntimeValue, b *RuntimeValue) (*RuntimeValue, error) {
	err := types.ExpectedIntType(a.RuntimeType)
	if err != nil {
		return nil, err
	}
	err = types.ExpectedIntType(b.RuntimeType)
	if err != nil {
		return nil, err
	}

	aInt := a.Value.(int64)
	bInt := b.Value.(int64)

	return &RuntimeValue{
		RuntimeType: a.RuntimeType,
		Value:       aInt + bInt,
	}, nil
}

func SubInt(a *RuntimeValue, b *RuntimeValue) (*RuntimeValue, error) {
	err := types.ExpectedIntType(a.RuntimeType)
	if err != nil {
		return nil, err
	}
	err = types.ExpectedIntType(b.RuntimeType)
	if err != nil {
		return nil, err
	}

	aInt := a.Value.(int64)
	bInt := b.Value.(int64)

	return &RuntimeValue{
		RuntimeType: a.RuntimeType,
		Value:       aInt - bInt,
	}, nil
}

func MultInt(a *RuntimeValue, b *RuntimeValue) (*RuntimeValue, error) {
	err := types.ExpectedIntType(a.RuntimeType)
	if err != nil {
		return nil, err
	}
	err = types.ExpectedIntType(b.RuntimeType)
	if err != nil {
		return nil, err
	}
	aInt := a.Value.(int64)
	bInt := b.Value.(int64)

	return &RuntimeValue{
		RuntimeType: a.RuntimeType,
		Value:       aInt * bInt,
	}, nil
}

func DivInt(a *RuntimeValue, b *RuntimeValue) (*RuntimeValue, error) {
	err := types.ExpectedIntType(a.RuntimeType)
	if err != nil {
		return nil, err
	}
	err = types.ExpectedIntType(b.RuntimeType)
	if err != nil {
		return nil, err
	}
	aInt := a.Value.(int64)
	bInt := b.Value.(int64)

	return &RuntimeValue{
		RuntimeType: a.RuntimeType,
		Value:       aInt / bInt,
	}, nil
}

func ApplyBinaryOpToInt(symbol string, lhs *RuntimeValue, rhs *RuntimeValue) (*RuntimeValue, error) {
	switch symbol {
	case "+":
		return AddInt(lhs, rhs)
	case "-":
		return SubInt(lhs, rhs)
	case "*":
		return MultInt(lhs, rhs)
	case "/":
		return DivInt(lhs, rhs)
	default:
		return nil, fmt.Errorf("unsupported operand %s for type %s", symbol, lhs.RuntimeType.GetName())
	}
}
