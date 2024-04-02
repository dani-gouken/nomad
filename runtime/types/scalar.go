package types

import (
	"fmt"
)

type ScalarType struct {
	name string
}

func (t *ScalarType) GetName() string {
	return t.name
}

func (t *ScalarType) Match(t2 RuntimeType) error {
	_, ok := t2.(*ScalarType)

	if !ok || (t.GetName() != t2.GetName()) {
		return fmt.Errorf("expected type %s, got %s", t.GetName(), t2.GetName())
	}

	return nil
}

func (t *ScalarType) IsInt() bool {
	return t.GetName() == "int"
}

func (t *ScalarType) IsFloat() bool {
	return t.GetName() == "float"
}

func (t *ScalarType) IsType() bool {
	return t.GetName() == "type"
}

func (t *ScalarType) IsString() bool {
	return t.GetName() == "string"
}

func (t *ScalarType) IsBoolean() bool {
	return t.GetName() == "bool"
}

func New(name string) *ScalarType {
	return &ScalarType{
		name: name,
	}
}

func ToScalarType(t RuntimeType) (*ScalarType, error) {
	tScalar, ok := t.(*ScalarType)
	if !ok {
		return nil, fmt.Errorf("scalar type expected")
	}
	return tScalar, nil
}

func MakeBoolType() *ScalarType {
	return New(BOOL_TYPE)
}

func ExpectedBoolType(t RuntimeType) error {
	tScalar, err := ToScalarType(t)
	if err != nil {
		return fmt.Errorf("expected bool type got %s", t.GetName())
	}
	if !tScalar.IsBoolean() {
		return fmt.Errorf("expected bool type got %s", t.GetName())
	}
	return nil
}

func MakeFloatType() *ScalarType {
	return New(FLOAT_TYPE)
}

func MakeStringType() *ScalarType {
	return New(STRING_TYPE)
}

func ExpectedFloatType(t RuntimeType) error {
	tScalar, err := ToScalarType(t)
	if err != nil {
		return fmt.Errorf("expected float type got %s", t.GetName())
	}
	if !tScalar.IsFloat() {
		return fmt.Errorf("expected float type got %s", t.GetName())
	}
	return nil
}

func MakeIntType() *ScalarType {
	return New(INT_TYPE)
}

func ExpectedIntType(t RuntimeType) error {
	tScalar, err := ToScalarType(t)
	if err != nil {
		return fmt.Errorf("expected int type got %s", t.GetName())
	}
	if !tScalar.IsInt() {
		return fmt.Errorf("expected int type got %s", t.GetName())
	}
	return nil
}

func ExpectedStringType(t RuntimeType) error {
	tScalar, err := ToScalarType(t)
	if err != nil {
		return fmt.Errorf("expected string type got %s", t.GetName())
	}
	if !tScalar.IsString() {
		return fmt.Errorf("expected int type got %s", t.GetName())
	}
	return nil
}

func MakeTypeType() *ScalarType {
	return New(TYPE_TYPE)
}

func ExpectedTypeType(t RuntimeType) error {
	tScalar, err := ToScalarType(t)
	if err != nil {
		return fmt.Errorf("expected type type got %s", t.GetName())
	}
	if !tScalar.IsType() {
		return fmt.Errorf("expected type type got %s", t.GetName())
	}
	return nil
}
