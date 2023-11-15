package vm

import "fmt"

const (
	INT_TYPE    = "int"
	BOOL_TYPE   = "bool"
	FLOAT_TYPE  = "float"
	STRING_TYPE = "string"
)

type RuntimeType struct {
	name    string
	methods map[string]Function
}

func (t *RuntimeType) GetName() string {
	return t.name
}

func MakeIntType() *RuntimeType {
	return &RuntimeType{
		name: INT_TYPE,
		methods: map[string]Function{
			"+": {
				Signature: []Parameter{
					{
						Name:     "self",
						TypeName: "int",
						Self:     true,
					},
					{
						Name:     "b",
						TypeName: "int",
					},
				},
				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
					a, _ := ((pv[0]).Value.Value).(int64)
					b, _ := ((pv[1]).Value.Value).(int64)
					return &RuntimeValue{
						TypeName: pv[0].Value.TypeName,
						Value:    a + b,
					}, nil
				},
			},
			"-": {
				Signature: []Parameter{
					{
						Name:     "self",
						TypeName: "int",
						Self:     true,
					},
					{
						Name:     "b",
						TypeName: "int",
					},
				},
				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
					a, _ := ((pv[0]).Value.Value).(int64)
					b, _ := ((pv[1]).Value.Value).(int64)
					return &RuntimeValue{
						TypeName: pv[0].Value.TypeName,
						Value:    a - b,
					}, nil
				},
			},
			"*": {
				Signature: []Parameter{
					{
						Name:     "self",
						TypeName: "int",
						Self:     true,
					},
					{
						Name:     "b",
						TypeName: "int",
					},
				},
				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
					a, _ := ((pv[0]).Value.Value).(int64)
					b, _ := ((pv[1]).Value.Value).(int64)
					return &RuntimeValue{
						TypeName: pv[0].Value.TypeName,
						Value:    a * b,
					}, nil
				},
			},
			"==": {
				Signature: []Parameter{
					{
						Name:     "self",
						TypeName: "int",
						Self:     true,
					},
					{
						Name:     "b",
						TypeName: "int",
					},
				},
				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
					a, _ := ((pv[0]).Value.Value).(int64)
					b, _ := ((pv[1]).Value.Value).(int64)
					return &RuntimeValue{
						TypeName: BOOL_TYPE,
						Value:    a == b,
					}, nil
				},
			},
		},
	}
}
func MakeBoolype() *RuntimeType {
	return &RuntimeType{
		name: BOOL_TYPE,
		methods: map[string]Function{
			"==": {
				Signature: []Parameter{
					{
						Name:     "self",
						TypeName: "bool",
						Self:     true,
					},
					{
						Name:     "b",
						TypeName: "bool",
					},
				},
				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
					a, _ := ((pv[0]).Value.Value).(bool)
					b, _ := ((pv[1]).Value.Value).(bool)
					return &RuntimeValue{
						TypeName: BOOL_TYPE,
						Value:    a == b,
					}, nil
				},
			},
		}}
}

func (t *RuntimeType) GetDeclaredMethods() map[string]Function {
	return t.methods
}

func (t *RuntimeType) GetMethod(name string) (Function, error) {
	res, ok := t.methods[name]
	if !ok {
		return res, fmt.Errorf("method %s::%s not found", t.GetName(), name)
	}
	return res, nil
}
