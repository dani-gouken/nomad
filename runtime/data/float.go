package data

import (
	"fmt"

	"github.com/dani-gouken/nomad/runtime/types"
)

func AddFloat(a *RuntimeValue, b *RuntimeValue) (*RuntimeValue, error) {
	err := types.ExpectedFloatType(a.RuntimeType)
	if err != nil {
		return nil, err
	}
	err = types.ExpectedFloatType(b.RuntimeType)
	if err != nil {
		return nil, err
	}

	aFloat := a.Value.(float64)
	bFloat := b.Value.(float64)

	return &RuntimeValue{
		RuntimeType: a.RuntimeType,
		Value:       aFloat + bFloat,
	}, nil
}

func SubFloat(a *RuntimeValue, b *RuntimeValue) (*RuntimeValue, error) {
	err := types.ExpectedFloatType(a.RuntimeType)
	if err != nil {
		return nil, err
	}
	err = types.ExpectedFloatType(b.RuntimeType)
	if err != nil {
		return nil, err
	}

	aFloat := a.Value.(float64)
	bFloat := b.Value.(float64)

	return &RuntimeValue{
		RuntimeType: a.RuntimeType,
		Value:       aFloat - bFloat,
	}, nil
}

func MultFloat(a *RuntimeValue, b *RuntimeValue) (*RuntimeValue, error) {
	err := types.ExpectedFloatType(a.RuntimeType)
	if err != nil {
		return nil, err
	}
	err = types.ExpectedFloatType(b.RuntimeType)
	if err != nil {
		return nil, err
	}
	aFloat := a.Value.(float64)
	bFloat := b.Value.(float64)

	return &RuntimeValue{
		RuntimeType: a.RuntimeType,
		Value:       aFloat * bFloat,
	}, nil
}

func DivFloat(a *RuntimeValue, b *RuntimeValue) (*RuntimeValue, error) {
	err := types.ExpectedFloatType(a.RuntimeType)
	if err != nil {
		return nil, err
	}
	err = types.ExpectedFloatType(b.RuntimeType)
	if err != nil {
		return nil, err
	}
	aFloat := a.Value.(float64)
	bFloat := b.Value.(float64)

	return &RuntimeValue{
		RuntimeType: a.RuntimeType,
		Value:       aFloat / bFloat,
	}, nil
}

func CmpFloat(t types.Registrar, a *RuntimeValue, b *RuntimeValue) (*RuntimeValue, error) {
	err := types.ExpectedFloatType(a.RuntimeType)
	if err != nil {
		return nil, err
	}
	err = types.ExpectedFloatType(b.RuntimeType)
	if err != nil {
		return nil, err
	}
	aFloat := a.Value.(float64)
	bFloat := b.Value.(float64)
	var result int64 = 0

	if aFloat < bFloat {
		result = -1
	}

	if aFloat > bFloat {
		result = 1
	}
	return &RuntimeValue{
		RuntimeType: t.GetOrPanic(types.INT_TYPE),
		Value:       result,
	}, nil
}

func ApplyBinaryOpToFloat(t types.Registrar, symbol string, lhs *RuntimeValue, rhs *RuntimeValue) (*RuntimeValue, error) {
	switch symbol {
	case "+":
		return AddFloat(lhs, rhs)
	case "-":
		return SubFloat(lhs, rhs)
	case "*":
		return MultFloat(lhs, rhs)
	case "/":
		return DivFloat(lhs, rhs)
	case "<->":
		return CmpFloat(t, lhs, rhs)
	default:
		return nil, fmt.Errorf("unsupported operand %s for type %s", symbol, lhs.RuntimeType.GetName())
	}
}
