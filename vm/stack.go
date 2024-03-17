package vm

import (
	"errors"
	"fmt"
)

const VM_MAX_STACK = 16384

type Stack struct {
	data    [VM_MAX_STACK]RuntimeValue
	pointer int
}

func (s *Stack) Push(value RuntimeValue) error {
	if s.pointer >= VM_MAX_STACK {
		return fmt.Errorf("Stack overflow. Maximum stack size of %d reached", VM_MAX_STACK)
	}
	s.data[s.pointer] = value
	s.pointer++
	return nil
}

func (s *Stack) Pop() (*RuntimeValue, error) {
	current, err := s.Current()
	if err != nil {
		return nil, err
	}
	s.pointer--
	return current, nil
}
func (s *Stack) Current() (*RuntimeValue, error) {
	if s.pointer <= 0 {
		return nil, errors.New("the stack is empty")
	}
	return &s.data[s.pointer-1], nil
}

func (s *Stack) PushBool(value bool) error {
	return s.Push(RuntimeValue{
		TypeName: BOOL_TYPE,
		Value:    value,
	})
}

func (s *Stack) PushType(typeValue *RuntimeType) error {
	return s.Push(RuntimeValue{
		TypeName: TYPE_TYPE,
		Value:    typeValue,
	})
}
func (s *Stack) PushInt(value int64) error {
	return s.Push(RuntimeValue{
		Value:    value,
		TypeName: INT_TYPE,
	})
}
func (s *Stack) PushFloat(value float64) error {
	return s.Push(RuntimeValue{
		Value:    value,
		TypeName: FLOAT_TYPE,
	})
}

func (s *Stack) Get(pointer int) *RuntimeValue {
	return &s.data[pointer]
}
