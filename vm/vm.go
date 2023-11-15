package vm

import (
	"fmt"
	"strconv"

	"github.com/dani-gouken/nomad/tokenizer"
)

const (
	OP_CONST_TRUE  = "TRUE"
	OP_CONST_FALSE = "FALSE"
)

type Vm struct {
	fp    int
	stack Stack
	env   Environment
	types TypeRegistrar
}

func (vm *Vm) Env() *Environment {
	return &vm.env
}

type Instruction struct {
	Code       string
	Arg1       string
	Arg2       string
	DebugToken tokenizer.Token
}

func New() *Vm {
	return &Vm{
		stack: Stack{
			pointer: 1,
		},
		env:   NewEnvironment(),
		types: NewTypeRegistrar(),
	}
}

func (vm *Vm) pushConst(runtimeType string, value string) error {
	switch runtimeType {
	case BOOL_TYPE:
		return vm.stack.Push(RuntimeValue{
			Value:    value == OP_CONST_TRUE,
			TypeName: runtimeType,
		})
	case INT_TYPE:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		return vm.stack.Push(RuntimeValue{
			Value:    int64(intVal),
			TypeName: runtimeType,
		})
	case FLOAT_TYPE:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		return vm.stack.Push(RuntimeValue{
			Value:    float64(floatVal),
			TypeName: runtimeType,
		})

	default:
		return fmt.Errorf("runtime error: Unable to store value of runtime type %s", runtimeType)
	}
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
			err := vm.pushConst(t, v)
			if err != nil {
				return err
			}
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
			if value.TypeName != BOOL_TYPE {
				return RuntimeErrorUnsupportedOperand("not(!)", value.TypeName, instruction.DebugToken)
			}
			boolValue := value.Value.(bool)
			if err != nil {
				return err
			}
			vm.stack.PushBool(!boolValue)
		case OP_NEGATIVE:
			value, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			switch value.TypeName {
			case INT_TYPE:
				intValue := value.Value.(int64)
				vm.stack.PushInt(-intValue)
			case FLOAT_TYPE:
				floatValue := value.Value.(float64)
				vm.stack.PushFloat(-floatValue)
			default:
				return RuntimeErrorUnsupportedOperand("Negative (-)", value.TypeName, instruction.DebugToken)
			}
		case OP_ADD, OP_SUB, OP_EQ, OP_MULT:
			rhs, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			lhs, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			opSymbol, err := OpToSymbol(instruction.Code)
			if err != nil {
				return err
			}
			result, err := vm.CallMethod(lhs, opSymbol, instruction.DebugToken, rhs)
			if err != nil {
				return err
			}
			vm.stack.Push(*result)
		case OP_PUSH_SCOPE:
			vm.Env().PushScope()
		case OP_POP_SCOPE:
			vm.Env().PopScope()
		case OP_STORE_VAR:
			value, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			declaredTypeName := instruction.Arg1
			declaredType, err := vm.types.Get(declaredTypeName)
			if err != nil {
				return RuntimeError(err.Error(), instruction.DebugToken)
			}
			name := instruction.Arg2
			err = vm.Env().SetVariable(
				name,
				declaredType,
				value,
			)
			if err != nil {
				return RuntimeError(err.Error(), instruction.DebugToken)
			}
		case OP_POP_CONST:
			vm.stack.Pop()
		case OP_LOAD_VAR:
			value, err := vm.Env().GetVariable(instruction.Arg1)
			if err != nil {
				return RuntimeError(err.Error(), instruction.DebugToken)
			}
			vm.stack.Push(*value)

		default:
			return fmt.Errorf("runtime error: Failed to interpret instruction [%s]", instruction.Code)
		}
	}
	return nil
}
func OpToSymbol(op string) (string, error) {
	switch op {
	case OP_ADD:
		return "+", nil
	case OP_SUB:
		return "-", nil
	case OP_MULT:
		return "*", nil
	case OP_EQ:
		return "==", nil
	}
	return "", fmt.Errorf("unknown operator %s", op)

}
func (vm *Vm) CallMethod(self *RuntimeValue, method string, debugToken tokenizer.Token, parameters ...*RuntimeValue) (*RuntimeValue, error) {
	t, err := vm.types.Get(self.TypeName)
	if err != nil {
		return nil, RuntimeError(err.Error(), debugToken)
	}
	function, err := t.GetMethod(method)
	if err != nil {
		return nil, RuntimeError(err.Error(), debugToken)
	}
	callParameters := []ParameterValue{
		{
			Value: self,
			Parameter: Parameter{
				Name:     "self",
				Self:     true,
				TypeName: self.TypeName,
			},
		},
	}
	if len(parameters) != (len(function.Signature) - 1) {
		return nil, RuntimeError(fmt.Sprintf("invalid parameter when calling %s::%s, expected %d parameters got %d", self.TypeName, method, len(function.Signature), len(parameters)), debugToken)
	}
	for i := 0; i < len(function.Signature)-1; i++ {
		parameter := function.Signature[i]
		value := parameters[i]
		if parameter.TypeName != value.TypeName {
			return nil, RuntimeError(fmt.Sprintf("type mismatch on %s::%s, expected parameter %d to be %s got %s", self.TypeName, method, i+1, parameter.TypeName, value.TypeName), debugToken)
		}
		callParameters = append(callParameters, ParameterValue{
			Parameter: parameter,
			Value:     value,
		})
	}
	return function.Run(callParameters)
}
