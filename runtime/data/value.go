package data

import (
	"fmt"

	"github.com/dani-gouken/nomad/runtime/types"
)

type Function struct {
	Name           string
	Signature      []Parameter
	ReturnTypeName string
	Run            func([]ParameterValue) (*RuntimeValue, error)
}

type Parameter struct {
	Name     string
	Self     bool
	TypeName string
}
type ParameterValue struct {
	Parameter
	Value *RuntimeValue
}

type RuntimeValue struct {
	RuntimeType types.RuntimeType
	Value       interface{}
}

type RuntimeArray struct {
	Values []RuntimeValue
}

func ApplyBinaryOp(symbol string, lhs *RuntimeValue, rhs *RuntimeValue) (*RuntimeValue, error) {
	err := lhs.RuntimeType.Match(rhs.RuntimeType)
	if err != nil {
		return nil, err
	}
	lhsType, err := types.ToScalarType(lhs.RuntimeType)
	if err != nil {
		return nil, err
	}

	if lhsType.IsFloat() {
		return ApplyBinaryOpToFloat(symbol, lhs, rhs)
	}

	if lhsType.IsInt() {
		return ApplyBinaryOpToInt(symbol, lhs, rhs)
	}

	if lhsType.IsString() {
		return ApplyBinaryOpToString(symbol, lhs, rhs)
	}

	return nil, fmt.Errorf("unsupported operand %s for type %s", symbol, lhs.RuntimeType.GetName())

}
