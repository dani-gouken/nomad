package data

import (
	"fmt"

	"github.com/dani-gouken/nomad/runtime/types"
)

func AddString(a *RuntimeValue, b *RuntimeValue) (*RuntimeValue, error) {
	err := types.ExpectedStringType(a.RuntimeType)
	if err != nil {
		return nil, err
	}
	err = types.ExpectedStringType(b.RuntimeType)
	if err != nil {
		return nil, err
	}

	aString := a.Value.(string)
	bString := b.Value.(string)

	return &RuntimeValue{
		RuntimeType: a.RuntimeType,
		Value:       aString + bString,
	}, nil
}

func ApplyBinaryOpToString(symbol string, lhs *RuntimeValue, rhs *RuntimeValue) (*RuntimeValue, error) {
	switch symbol {
	case "+":
		return AddString(lhs, rhs)
	default:
		return nil, fmt.Errorf("unsupported operand %s for type %s", symbol, lhs.RuntimeType.GetName())
	}
}
