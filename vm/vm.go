package vm

import (
	"fmt"
	"strconv"

	nomadError "github.com/dani-gouken/nomad/errors"
	"github.com/dani-gouken/nomad/runtime/data"
	"github.com/dani-gouken/nomad/runtime/types"
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
	types types.Registrar
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
		types: types.NewRegistrar(),
	}
}

func (vm *Vm) pushConst(runtimeType string, value string) error {
	switch runtimeType {
	case types.BOOL_TYPE:
		return vm.stack.Push(data.RuntimeValue{
			Value:       value == OP_CONST_TRUE,
			RuntimeType: vm.types.GetOrPanic(runtimeType),
		})
	case types.INT_TYPE:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		return vm.stack.Push(data.RuntimeValue{
			Value:       int64(intVal),
			RuntimeType: vm.types.GetOrPanic(runtimeType),
		})
	case types.FLOAT_TYPE:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		return vm.stack.Push(data.RuntimeValue{
			Value:       float64(floatVal),
			RuntimeType: vm.types.GetOrPanic(runtimeType),
		})
	case types.STRING_TYPE:
		return vm.stack.Push(data.RuntimeValue{
			Value:       value,
			RuntimeType: vm.types.GetOrPanic(runtimeType),
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
			if err != nil {
				return err
			}
			_, err = types.ToObjectType(value.RuntimeType)
			if err != nil {
				fmt.Printf("<%s> %v\n", value.RuntimeType.GetName(), value.Value)
			} else {
				vObj := value.Value.(*data.RuntimeObject)
				fmt.Print(value.RuntimeType.GetName())
				fmt.Print("{")
				lastKey := ""

				for k1 := range vObj.GetFields() {
					lastKey = k1
				}
				for k, v := range vObj.GetFields() {
					fmt.Print(v.RuntimeType.GetName())
					fmt.Print(" ")
					fmt.Print(k)
					fmt.Print(" ")
					fmt.Print(v.Value)
					if k != lastKey {
						fmt.Print(", ")
					}
				}
				fmt.Print("}")
				fmt.Print("\n")
			}
		case OP_NOT:
			value, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			err = types.ExpectedBoolType(value.RuntimeType)
			if err != nil {
				return nomadError.RuntimeErrorUnsupportedOperand("not(!)", value.RuntimeType.GetName(), instruction.DebugToken)
			}
			boolValue := value.Value.(bool)
			if err != nil {
				return err
			}
			vm.stack.PushBool(vm.types, !boolValue)
		case OP_JUMP:
			addr, err := strconv.Atoi(instruction.Arg1)
			if err != nil {
				panic(err)
			}
			i = addr - 1
		case OP_JUMP_NOT, OP_JUMP_IF:
			value, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			err = types.ExpectedBoolType(value.RuntimeType)
			if err != nil {
				return nomadError.RuntimeErrorUnsupportedOperand("boolean expected for comparison", value.RuntimeType.GetName(), instruction.DebugToken)
			}
			v, ok := value.Value.(bool)
			if !ok {
				panic("boolean expected")
			}
			addr, err := strconv.Atoi(instruction.Arg1)
			if err != nil {
				panic(err)
			}
			if (instruction.Code == OP_JUMP && v) || (instruction.Code == OP_JUMP_NOT && !v) {
				i = addr - 1
			}

		case OP_NEGATIVE:
			value, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			switch value.RuntimeType.GetName() {
			case types.INT_TYPE:
				intValue := value.Value.(int64)
				vm.stack.PushInt(vm.types, -intValue)
			case types.FLOAT_TYPE:
				floatValue := value.Value.(float64)
				vm.stack.PushFloat(vm.types, -floatValue)
			default:
				return nomadError.RuntimeErrorUnsupportedOperand("negative (-)", value.RuntimeType.GetName(), instruction.DebugToken)
			}
		case OP_LEN:
			value, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			_, arrayTypeErr := types.ToArrayType(value.RuntimeType)
			scalarType, scalarTypeErr := types.ToScalarType(value.RuntimeType)

			if (arrayTypeErr != nil && scalarTypeErr != nil) || (scalarTypeErr == nil && !scalarType.IsString()) {
				return nomadError.RuntimeErrorUnsupportedOperand("len", value.RuntimeType.GetName(), instruction.DebugToken)
			}
			if arrayTypeErr == nil {
				arrValue := value.Value.(data.RuntimeArray)
				vm.stack.PushInt(vm.types, int64(len(arrValue.Values)))
			}
			if scalarTypeErr == nil {
				stringValue := value.Value.(string)
				vm.stack.PushInt(vm.types, int64(len(stringValue)))
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
			vm.stack.PushBool(vm.types, rhs.Value == lhs.Value)
		case OP_EQ_2:
			rhs, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			lhs1, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			lhs2, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			vm.stack.PushBool(vm.types, (rhs.Value == lhs1.Value) || (rhs.Value == lhs2.Value))
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
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			result, err := data.ApplyBinaryOp(vm.types, opSymbol, lhs, rhs)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			vm.stack.Push(*result)
		case OP_OR, OP_AND:
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
			if lhs.RuntimeType.GetName() != types.BOOL_TYPE {
				return nomadError.RuntimeError(fmt.Sprintf("cannot call operator %s on value of type %s", opSymbol, lhs.RuntimeType.GetName()), instruction.DebugToken)
			}
			if rhs.RuntimeType.GetName() != types.BOOL_TYPE {
				return nomadError.RuntimeError(fmt.Sprintf("cannot call operator %s on value of type %s", opSymbol, rhs.RuntimeType.GetName()), instruction.DebugToken)
			}
			lhsValue := lhs.Value.(bool)
			rhsValue := rhs.Value.(bool)
			var res bool
			if instruction.Code == OP_OR {
				res = lhsValue || rhsValue
			} else {
				res = lhsValue && rhsValue
			}
			vm.stack.PushBool(vm.types, res)
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
			err = variable.RuntimeType.Match(value.RuntimeType)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			variable.Value = value.Value
		case OP_DECL_VAR:
			value, err := vm.stack.Pop()
			t, err := vm.stack.Pop()

			tScalar, err := types.ToScalarType(t.RuntimeType)
			if err != nil {
				return err
			}

			if !tScalar.IsType() {
				return nomadError.RuntimeError(fmt.Sprintf("type value expected, got %s", t.RuntimeType.GetName()), instruction.DebugToken)
			}

			declaredType := t.Value.(types.RuntimeType)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			name := instruction.Arg2
			err = vm.Env().DeclareVariable(
				name,
				value,
				declaredType,
			)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
		case OP_ARR_INIT:
			t, err := vm.stack.Pop()
			if err != nil {
				return err
			}

			tScalar, err := types.ToScalarType(t.RuntimeType)
			if err != nil {
				return err
			}
			err = types.ExpectedTypeType(tScalar)
			if err != nil {
				return err
			}

			arraySubType := t.Value.(types.RuntimeType)
			arrayType := types.NewArrayType(arraySubType)
			value := data.RuntimeValue{
				Value:       data.RuntimeArray{},
				RuntimeType: arrayType,
			}
			vm.stack.Push(value)
		case OP_ARR_TYPE:
			t, err := vm.stack.Pop()
			if err != nil {
				return err
			}

			tScalar, err := types.ToScalarType(t.RuntimeType)
			if err != nil {
				return err
			}
			err = types.ExpectedTypeType(tScalar)
			if err != nil {
				return err
			}

			arraySubType := t.Value.(types.RuntimeType)
			arrayType := types.NewArrayType(arraySubType)
			vm.stack.PushType(vm.types, arrayType)

		case OP_ARR_PUSH:
			value, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			array, err := vm.stack.Current()
			if err != nil {
				return err
			}
			t, err := types.ToArrayType(array.RuntimeType)
			if err != nil {
				return nomadError.RuntimeError("cannot push to non-array types", instruction.DebugToken)
			}
			runtimeArray, _ := array.Value.(data.RuntimeArray)
			err = t.MatchSubtype(value.RuntimeType)
			if err != nil {
				return nomadError.RuntimeError(fmt.Sprintf("type mismatch, %s expected, %s given", t.GetSubtype().GetName(), value.RuntimeType.GetName()), instruction.DebugToken)
			}
			runtimeArray.Values = append(runtimeArray.Values, *value)
			array.Value = runtimeArray
		case OP_ARR_LOAD:
			index, err := vm.stack.Pop()
			array, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			_, err = types.ToArrayType(array.RuntimeType)
			if err != nil {
				return nomadError.RuntimeError("cannot push to non-array types", instruction.DebugToken)
			}

			err = types.ExpectedIntType(index.RuntimeType)
			if err != nil {
				return nomadError.RuntimeError("index should be an integer", instruction.DebugToken)
			}

			runtimeArray, _ := array.Value.(data.RuntimeArray)
			i, _ := index.Value.(int64)
			if err != nil {
				return nomadError.RuntimeError("index should be an integer", instruction.DebugToken)
			}
			vm.stack.Push(runtimeArray.Values[i])
		case OP_POP_CONST:
			vm.stack.Pop()
		case OP_LOAD_VAR:
			value, err := vm.Env().GetVariable(instruction.Arg1)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			vm.stack.Push(*value)
		case OP_LOAD_TYPE:
			value, err := vm.types.Get(instruction.Arg1)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			vm.stack.PushType(vm.types, value)
		case OP_DECL_TYPE:
			value, err := vm.stack.Pop()
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			err = types.ExpectedTypeType(value.RuntimeType)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			vType := value.Value.(types.RuntimeType)
			objectType, err := types.ToObjectType(vType)
			if err == nil && objectType.IsAnonymous() {
				objectType.SetName(instruction.Arg1)
				vm.types.Add(objectType, instruction.DebugToken)
			} else {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
		case OP_OBJ_TYPE:
			obj := types.NewObjectType()
			vm.stack.PushType(vm.types, obj)
		case OP_OBJ_INIT:
			typeName := instruction.Arg1
			t, err := vm.types.Get(typeName)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}

			tObj, err := types.ToObjectType(t)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}

			obj := data.NewRuntimeObject()
			for k, v := range tObj.GetDefaults() {
				vValue, ok := v.(data.RuntimeValue)
				if !ok {
					return nomadError.RuntimeError("object default is expected to be a runtime value", instruction.DebugToken)
				}
				obj.SetField(k, vValue)
			}

			vm.stack.Push(data.RuntimeValue{
				RuntimeType: t,
				Value:       obj,
			})

		case OP_OBJ_TYPE_SET_FIELD:
			fieldDefaultValue, err := vm.stack.Pop()
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}

			fieldTypeValue, err := vm.stack.Pop()
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}

			err = types.ExpectedTypeType(fieldTypeValue.RuntimeType)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}

			fieldType := fieldTypeValue.Value.(types.RuntimeType)

			object, err := vm.stack.Current()
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}

			err = types.ExpectedTypeType(object.RuntimeType)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}

			objectTypeValue := object.Value.(types.RuntimeType)

			objectType, err := types.ToObjectType(objectTypeValue)

			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			objectType.AddField(instruction.Arg1, fieldType, *fieldDefaultValue)
		case OP_OBJ_SET_FIELD:
			value, err := vm.stack.Pop()
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			field := instruction.Arg1

			objectValue, err := vm.stack.Current()
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}

			t, err := types.ToObjectType(objectValue.RuntimeType)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}

			object := objectValue.Value.(*data.RuntimeObject)
			fieldType, err := t.GetFieldType(field)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}

			err = fieldType.Match(value.RuntimeType)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}
			object.SetField(field, *value)
		case OP_OBJ_LOAD:
			objectValue, err := vm.stack.Pop()
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}

			_, err = types.ToObjectType(objectValue.RuntimeType)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}

			object := objectValue.Value.(*data.RuntimeObject)
			field := instruction.Arg1
			v, err := object.GetField(field)
			if err != nil {
				return nomadError.RuntimeError(err.Error(), instruction.DebugToken)
			}

			vm.stack.Push(*v)
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
	case OP_AND:
		return "&", nil
	case OP_OR:
		return "|", nil
	case OP_DIV:
		return "/", nil
	}
	return "", fmt.Errorf("unknown operator %s", op)

}
