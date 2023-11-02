package vm

import (
	"fmt"
	"strconv"
)

const (
	OP_CONST_TRUE  = "TRUE"
	OP_CONST_FALSE = "FALSE"
)

type Vm struct {
	fp    int
	stack Stack
	env   Environment
}

func (vm *Vm) Env() *Environment {
	return &vm.env
}

type Instruction struct {
	Code string
	Arg1 string
	Arg2 string
}

func New() *Vm {
	return &Vm{
		stack: Stack{
			pointer: 1,
		},
		env: NewEnvironment(),
	}
}

func (vm *Vm) pushConst(runtimeType string, value string) error {
	switch runtimeType {
	case VM_RUNTIME_BOOL:
		return vm.stack.PushBool(value == OP_CONST_TRUE)

	case VM_RUNTIME_INT:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		return vm.stack.PushInt(int64(intVal))

	default:
		return fmt.Errorf("runtime error: Unable to store value of runtime type %s", runtimeType)
	}
}

func SupportAddOp(val *RuntimeValue) bool {
	return (*val).GetType() == VM_RUNTIME_INT
}

func (vm *Vm) Interpret(instructions []Instruction) error {
loop:
	for i := 0; i < len(instructions); i++ {
		instruction := instructions[i]
		switch instruction.Code {
		case OP_HALT:
			break loop
		case OP_STORE_CONST:
			t := instruction.Arg1
			v := instruction.Arg2
			vm.pushConst(t, v)
		case OP_DEBUG_PRINT:
			value, err := vm.stack.Current()
			if err == nil {
				fmt.Println(*value)
			}
		case OP_NOT:
			value, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			boolValue, err := RuntimeValueAsBool(value)
			if err != nil {
				return err
			}
			vm.stack.PushBool(!boolValue.Value)
		case OP_NEGATIVE:
			value, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			boolValue, err := RuntimeValueAsInt(value)
			if err != nil {
				return err
			}
			vm.stack.PushInt(-boolValue.Value)
		case OP_ADD:
			v1, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			v2, err := vm.stack.Pop()
			if err != nil {
				return err
			}

			if !SupportAddOp(v1) {
				return fmt.Errorf("runtime error: Type error, cannot add value of type [%s]", (*v1).GetType())
			}
			if !SupportAddOp(v2) {
				return fmt.Errorf("runtime error: Type error, cannot add value of type [%s]", (*v2).GetType())
			}
			if (*v1).GetType() != (*v2).GetType() {
				return fmt.Errorf("runtime error: Type mismatch, cannot add value of type [%s] to value of type [%s]", (*v1).GetType(), (*v2).GetType())
			}
			v1Int, _ := RuntimeValueAsInt(v1)
			v2Int, _ := RuntimeValueAsInt(v2)

			vm.stack.PushInt(v1Int.Value + v2Int.Value)
		case OP_PUSH_SCOPE:
			vm.Env().PushScope()
		case OP_POP_SCOPE:
			vm.Env().PopScope()
		case OP_STORE_VAR:
			value, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			vm.Env().SetVariable(
				instruction.Arg1,
				value,
			)
		case OP_POP_CONST:
			vm.stack.Pop()
		case OP_LOAD_VAR:
			value, err := vm.Env().GetVariable(instruction.Arg1)
			if err != nil {
				return err
			}
			vm.stack.Push(*value)
		case OP_SUB:
			v2, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			v1, err := vm.stack.Pop()
			if err != nil {
				return err
			}

			if !SupportAddOp(v1) {
				return fmt.Errorf("runtime error: Type error, cannot add value of type [%s]", (*v1).GetType())
			}
			if !SupportAddOp(v2) {
				return fmt.Errorf("runtime error: Type error, cannot add value of type [%s]", (*v2).GetType())
			}
			if (*v1).GetType() != (*v2).GetType() {
				return fmt.Errorf("runtime error: Type mismatch, cannot add value of type [%s] to value of type [%s]", (*v1).GetType(), (*v2).GetType())
			}
			v1Int, _ := RuntimeValueAsInt(v1)
			v2Int, _ := RuntimeValueAsInt(v2)
			vm.stack.PushInt(v1Int.Value - v2Int.Value)
		case OP_EQ:
			v2, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			v1, err := vm.stack.Pop()
			if err != nil {
				return err
			}

			if (*v1).GetType() != (*v2).GetType() {
				vm.stack.PushBool(false)
			} else if (*v1).GetType() == VM_RUNTIME_INT {
				v1Int, _ := RuntimeValueAsInt(v1)
				v2Int, err := RuntimeValueAsInt(v2)
				if err != nil {
					return vm.stack.PushBool(false)
				}
				vm.stack.PushBool(v1Int.Value == v2Int.Value)
			} else if (*v1).GetType() == VM_RUNTIME_BOOL {
				v1Bool, _ := RuntimeValueAsBool(v1)
				v2Bool, err := RuntimeValueAsBool(v2)
				if err != nil {
					return vm.stack.PushBool(false)
				}
				vm.stack.PushBool(v1Bool.Value == v2Bool.Value)
			} else {
				vm.stack.PushBool(false)
			}
		case OP_MULT:
			v1, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			v2, err := vm.stack.Pop()
			if err != nil {
				return err
			}

			if !SupportAddOp(v1) {
				return fmt.Errorf("runtime error: Type error, cannot add value of type [%s]", (*v1).GetType())
			}
			if !SupportAddOp(v2) {
				return fmt.Errorf("runtime error: Type error, cannot add value of type [%s]", (*v2).GetType())
			}
			if (*v1).GetType() != (*v2).GetType() {
				return fmt.Errorf("runtime error: Type mismatch, cannot add value of type [%s] to value of type [%s]", (*v1).GetType(), (*v2).GetType())
			}
			v1Int, _ := RuntimeValueAsInt(v1)
			v2Int, _ := RuntimeValueAsInt(v2)

			vm.stack.PushInt(v1Int.Value * v2Int.Value)
		default:
			return fmt.Errorf("runtime error: Failed to interpret instruction [%s]", instruction.Code)
		}
	}
	return nil
}
