package vm

import (
	"errors"
	"fmt"

	"github.com/dani-gouken/nomad/runtime/data"
	"github.com/dani-gouken/nomad/runtime/types"
)

const VM_MAX_STACK = 16384

type Stack struct {
	data    [VM_MAX_STACK]data.RuntimeValue
	pointer int
}

func (s *Stack) Push(value data.RuntimeValue) error {
	if s.pointer >= VM_MAX_STACK {
		return fmt.Errorf("Stack overflow. Maximum stack size of %d reached", VM_MAX_STACK)
	}
	s.data[s.pointer] = value
	s.pointer++
	return nil
}

func (s *Stack) Pop() (*data.RuntimeValue, error) {
	current, err := s.Current()
	if err != nil {
		return nil, err
	}
	s.pointer--
	return current, nil
}
func (s *Stack) Current() (*data.RuntimeValue, error) {
	if s.pointer <= 0 {
		return nil, errors.New("stack underflow")
	}
	return &s.data[s.pointer-1], nil
}

func (s *Stack) PushBool(t types.Registrar, value bool) error {
	return s.Push(data.RuntimeValue{
		RuntimeType: t.GetOrPanic(types.BOOL_TYPE),
		Value:       value,
	})
}

func (s *Stack) PushType(t types.Registrar, typeValue types.RuntimeType) error {
	return s.Push(data.RuntimeValue{
		RuntimeType: t.GetOrPanic(types.TYPE_TYPE),
		Value:       typeValue,
	})
}
func (s *Stack) PushInt(t types.Registrar, value int64) error {
	return s.Push(data.RuntimeValue{
		Value:       value,
		RuntimeType: t.GetOrPanic(types.INT_TYPE),
	})
}
func (s *Stack) PushFloat(t types.Registrar, value float64) error {
	return s.Push(data.RuntimeValue{
		Value:       value,
		RuntimeType: t.GetOrPanic(types.FLOAT_TYPE),
	})
}

func (s *Stack) Get(pointer int) *data.RuntimeValue {
	return &s.data[pointer]
}

func NewStack() *Stack {
	return &Stack{
		data:    [16384]data.RuntimeValue{},
		pointer: 1,
	}
}
