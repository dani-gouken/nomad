package types

import (
	"fmt"
)

type VoidType struct {
}

func (f *VoidType) GetName() string {
	return "void"
}

func (t *VoidType) Match(t2 RuntimeType) error {
	_, err := ToVoidType(t2)
	if err != nil {
		return err
	}
	return nil
}

func MakeVoidType() *VoidType {
	return &VoidType{}
}

func IsVoidType(t RuntimeType) bool {
	_, err := ToVoidType(t)
	return err == nil
}

func ToVoidType(t RuntimeType) (*VoidType, error) {
	tVoid, ok := t.(*VoidType)
	if !ok {
		return nil, fmt.Errorf("object type expected")
	}
	return tVoid, nil
}
