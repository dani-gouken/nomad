package vm

import (
	"errors"
	"fmt"

	"github.com/dani-gouken/nomad/runtime/data"
	"github.com/dani-gouken/nomad/tokenizer"
)

const VM_MAX_CALL_STACK = 16384

type Frame = struct {
	DebugToken  tokenizer.Token
	CurrentFunc *data.RuntimeFunc
	stack       *Stack
	returnAddr  int
}

type CallStack struct {
	data    [VM_MAX_STACK]*Frame
	pointer int
}

func NewFrame(returnAddr int, f *data.RuntimeFunc, t tokenizer.Token) *Frame {
	return &Frame{
		CurrentFunc: f,
		stack:       NewStack(),
		DebugToken:  t,
		returnAddr:  returnAddr,
	}
}

func NewCallStack() *CallStack {
	callStack := &CallStack{
		data: [VM_MAX_CALL_STACK]*Frame{
			NewFrame(-1, nil, tokenizer.Token{}),
		},
		pointer: 1,
	}
	return callStack
}

func (s *CallStack) Push(f *Frame) error {
	if s.pointer >= VM_MAX_STACK {
		return fmt.Errorf("Stack overflow. Maximum stack size of %d reached", VM_MAX_STACK)
	}
	s.data[s.pointer] = f
	s.pointer++
	return nil
}

func (s *CallStack) Pop() (*Frame, error) {
	current, err := s.Current()
	if err != nil {
		return nil, err
	}
	s.pointer--
	return current, nil
}
func (s *CallStack) Current() (*Frame, error) {
	if s.pointer <= 0 {
		return nil, errors.New("stack underflow")
	}
	return s.data[s.pointer-1], nil
}

func (s *CallStack) Get(pointer int) *Frame {
	return s.data[pointer]
}
func (s *CallStack) SetPointer(pointer int) {
	s.pointer = pointer
}
