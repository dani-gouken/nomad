package types

import "fmt"

type ArrayType struct {
	subtype RuntimeType
}

func (t *ArrayType) GetName() string {
	return fmt.Sprintf("[%s]", t.subtype.GetName())
}

func (t *ArrayType) GetSubtype() RuntimeType {
	return t.subtype
}

func (t *ArrayType) Match(t2 RuntimeType) error {
	t2Arr, err := ToArrayType(t2)
	if err != nil {
		return err
	}
	err = t.GetSubtype().Match(t2Arr.GetSubtype())
	if err != nil {
		return t.GetSubtype().Match(t2Arr.GetSubtype())
	}

	return nil
}

func (t *ArrayType) MatchSubtype(t2 RuntimeType) error {
	return t.GetSubtype().Match(t2)
}

func NewArrayType(t RuntimeType) *ArrayType {
	return &ArrayType{
		subtype: t,
	}
}

func IsArrayType(t RuntimeType) bool {
	_, err := ToArrayType(t)
	return err == nil
}

func ToArrayType(t RuntimeType) (*ArrayType, error) {
	tArr, ok := t.(*ArrayType)
	if !ok {
		return nil, fmt.Errorf("array type expected")
	}
	return tArr, nil
}
