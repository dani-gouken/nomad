package vm

import "fmt"

const (
	VM_RUNTIME_BOOL = "BOOL"
	VM_RUNTIME_INT  = "INT"
)

type RuntimeValue interface {
	GetType() string
}

func RuntimeValueAsBool(v *RuntimeValue) (*BoolRuntimeValue, error) {
	vbool, ok := (*v).(*BoolRuntimeValue)
	if !ok {
		return nil, fmt.Errorf("could not convert runtime value of type [%s] to boolean", (*v).GetType())
	}
	return vbool, nil
}

func RuntimeValueAsInt(v *RuntimeValue) (*IntRuntimeValue, error) {
	vint, ok := (*v).(*IntRuntimeValue)
	if !ok {
		return nil, fmt.Errorf("could not convert runtime value of type [%s] to boolean", (*v).GetType())
	}
	return vint, nil
}

type BoolRuntimeValue struct {
	Value bool
}

type IntRuntimeValue struct {
	Value int64
}

func (v *BoolRuntimeValue) GetType() string {
	return VM_RUNTIME_BOOL
}
func (v *IntRuntimeValue) GetType() string {
	return VM_RUNTIME_INT
}
