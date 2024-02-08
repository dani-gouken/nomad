package vm

import (
	"fmt"
	"strconv"

	nomadError "github.com/dani-gouken/nomad/errors"
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
		return fmt.Errorf("runtime error: unable to store value of runtime type %s", runtimeType)
	}
}

func (vm *Vm) Interpret(instructions []Instruction) error {
loop:
	for i := 0; i < len(instructions); i++ {
		instruction := instructions[i]
		// fmt.Println(fmt.Sprintf("executing %s %s %s", instruction.Code, instruction.Arg1, instruction.Arg2))
		switch instruction.Code {
		case OP_HALT:
			break loop
		case OP_PUSH_CONST:
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
				return nomadError.RuntimeErrorUnsupportedOperand("not(!)", value.TypeName, instruction.DebugToken)
			}
			boolValue := value.Value.(bool)
			if err != nil {
				return err
			}
			vm.stack.PushBool(!boolValue)
		case JUMP:
			addr, err := strconv.Atoi(instruction.Arg1)
			if err != nil {
				panic(err)
			}
			i = addr - 1
		case JUMP_NOT, JUMP_IF:
			value, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			if value.TypeName != BOOL_TYPE {
				return nomadError.RuntimeErrorUnsupportedOperand("boolean expected for comparison", value.TypeName, instruction.DebugToken)
			}
			v, ok := value.Value.(bool)
			if !ok {
				panic("boolean expected")
			}
			addr, err := strconv.Atoi(instruction.Arg1)
			if err != nil {
				panic(err)
			}
			if (instruction.Code == JUMP && v) || (instruction.Code == JUMP_NOT && !v) {
				i = addr - 1
			}

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
				return nomadError.RuntimeErrorUnsupportedOperand("negative (-)", value.TypeName, instruction.DebugToken)
			}
		case OP_EQ:
			rhs, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			lhs, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			vm.stack.PushBool(rhs.Value == lhs.Value)
		case OP_ADD, OP_SUB, OP_MULT, OP_DIV, OP_CMP:
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
			t, err := vm.types.Get(lhs.TypeName)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			_, err = t.GetMethod(opSymbol)
			if err != nil {
				return nomadError.RuntimeErrorUnsupportedOperand(opSymbol, t.name, instruction.DebugToken)
			}
			result, err := vm.CallMethod(lhs, opSymbol, instruction.DebugToken, rhs)
			if err != nil {
				return err
			}
			vm.stack.Push(*result)
		case OP_PUSH_SCOPE:
			vm.Env().PushScope()
		case OP_LABEL:
			continue
		case OP_POP_SCOPE:
			vm.Env().PopScope()
		case OP_SET_VAR:
			value, err := vm.stack.Pop()

			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			variableName := instruction.Arg1
			variable, err := vm.env.GetVariable(variableName)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			err = vm.checkTypeCompatibility(variable.TypeName, value.TypeName)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			variable.Value = value.Value
			variable.TypeName = value.TypeName
		case OP_DECL_VAR:
			value, err := vm.stack.Pop()

			if err != nil {
				return err
			}

			declaredTypeName := instruction.Arg1
			declaredType, err := vm.types.Get(declaredTypeName)
			possibleTypes := []*RuntimeType{}

			if err != nil {
				compositeType, err := vm.types.GetComposite(declaredTypeName)
				if err != nil {
					return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
				}
				for i := 0; i < len(compositeType.Cases); i++ {
					t, err := vm.types.Get(compositeType.Cases[i])
					if err != nil {
						return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
					}
					possibleTypes = append(possibleTypes, t)
				}
			} else {
				possibleTypes = append(possibleTypes, declaredType)
			}

			valueType, err := vm.types.Get(value.TypeName)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			name := instruction.Arg2
			err = vm.Env().DeclareVariable(
				name,
				valueType,
				value,
				declaredType,
				possibleTypes,
			)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
		case OP_POP_CONST:
			vm.stack.Pop()
		case OP_LOAD_VAR:
			value, err := vm.Env().GetVariable(instruction.Arg1)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			vm.stack.Push(*value)

		default:
			return nomadError.RuntimeError(fmt.Sprintf("failed to interpret instruction [%s]", instruction.Code), instruction.DebugToken)
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
	case OP_CMP:
		return "<->", nil
	case OP_DIV:
		return "/", nil
	}
	return "", fmt.Errorf("unknown operator %s", op)

}

func (vm *Vm) checkTypeCompatibility(typeA string, typeB string) error {
	possibleTypesA, err := vm.types.GetPossible(typeA)

	if err != nil {
		return err
	}
	possibleTypesB, err := vm.types.GetPossible(typeB)
	if err != nil {
		return err
	}

	for i := 0; i < len(possibleTypesA); i++ {
		possibleA := possibleTypesA[i]
		for j := 0; j < len(possibleTypesB); j++ {
			possibleB := possibleTypesB[j]
			if possibleA == possibleB {
				return nil
			}
		}

	}
	return fmt.Errorf("type %s is not compatible with %s", typeA, typeB)
}
func (vm *Vm) CallMethod(self *RuntimeValue, method string, debugToken tokenizer.Token, parameters ...*RuntimeValue) (*RuntimeValue, error) {
	t, err := vm.types.Get(self.TypeName)
	if err != nil {
		return nil, nomadError.RuntimeError(err.Error(), debugToken)
	}
	function, err := t.GetMethod(method)
	if err != nil {
		return nil, nomadError.RuntimeError(err.Error(), debugToken)
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
		return nil, nomadError.RuntimeError(fmt.Sprintf("invalid parameter when calling %s::%s, expected %d parameters got %d", self.TypeName, method, len(function.Signature), len(parameters)), debugToken)
	}

	for i := 0; i < len(function.Signature)-1; i++ {

		parameter := function.Signature[i+1]
		value := parameters[i]
		err := vm.checkTypeCompatibility(parameter.TypeName, value.TypeName)

		if err != nil {
			return nil, nomadError.RuntimeError(fmt.Sprintf("type mismatch on %s::%s, expected parameter %d to be %s got %s", self.TypeName, method, i+1, parameter.TypeName, value.TypeName), debugToken)
		}
		callParameters = append(callParameters, ParameterValue{
			Parameter: parameter,
			Value:     value,
		})
	}
	result, err := function.Run(callParameters)
	if err != nil {
		return nil, nomadError.RuntimeError(err.Error(), debugToken)
	}
	if result.TypeName != function.ReturnTypeName {
		return nil, nomadError.RuntimeError(fmt.Sprintf("expected return type for %s::%s is %s got %s", self.TypeName, method, function.ReturnTypeName, result.TypeName), debugToken)
	}
	return result, err
}
