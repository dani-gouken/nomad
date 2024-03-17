package vm

import "fmt"

const (
	INT_TYPE    = "int"
	BOOL_TYPE   = "bool"
	TYPE_TYPE   = "type"
	FLOAT_TYPE  = "float"
	NUM_TYPE    = "num"
	STRING_TYPE = "string"
	ARRAY_TYPE  = "array"
)

type RuntimeType struct {
	name    string
	subtype string
	methods map[string]Function
}

func (t *RuntimeType) GetName() string {
	return t.name
}

func (t *RuntimeType) GetSubtype() string {
	return t.subtype
}

func MakeIntType() *RuntimeType {
	return &RuntimeType{
		name: INT_TYPE,
		methods: map[string]Function{
			"+": {
				Signature: []Parameter{
					{
						Name:     "self",
						TypeName: INT_TYPE,
						Self:     true,
					},
					{
						Name:     "b",
						TypeName: INT_TYPE,
					},
				},
				ReturnTypeName: INT_TYPE,
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
						TypeName: INT_TYPE,
						Self:     true,
					},
					{
						Name:     "b",
						TypeName: INT_TYPE,
					},
				},
				ReturnTypeName: INT_TYPE,
				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
					a, _ := ((pv[0]).Value.Value).(int64)
					b, _ := ((pv[1]).Value.Value).(int64)
					return &RuntimeValue{
						TypeName: pv[0].Value.TypeName,
						Value:    a - b,
					}, nil
				},
			},
			"<->": {
				Signature: []Parameter{
					{
						Name:     "self",
						TypeName: INT_TYPE,
						Self:     true,
					},
					{
						Name:     "b",
						TypeName: NUM_TYPE,
					},
				},
				ReturnTypeName: INT_TYPE,
				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
					a, _ := ((pv[0]).Value.Value).(int64)
					b, _ := ((pv[1]).Value.Value).(int64)
					var result int64 = 0
					if a < b {
						result = -1
					}
					if a > b {
						result = 1
					}
					return &RuntimeValue{
						TypeName: INT_TYPE,
						Value:    result,
					}, nil
				},
			},
			"/": {
				Signature: []Parameter{
					{
						Name:     "self",
						TypeName: INT_TYPE,
						Self:     true,
					},
					{
						Name:     "b",
						TypeName: INT_TYPE,
					},
				},
				ReturnTypeName: FLOAT_TYPE,
				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
					a, _ := ((pv[0]).Value.Value).(int64)
					b, _ := ((pv[1]).Value.Value).(int64)

					return &RuntimeValue{
						TypeName: FLOAT_TYPE,
						Value:    float64(a) / float64(b),
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
				ReturnTypeName: "int",
				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
					a, _ := ((pv[0]).Value.Value).(int64)
					b, _ := ((pv[1]).Value.Value).(int64)
					return &RuntimeValue{
						TypeName: pv[0].Value.TypeName,
						Value:    a * b,
					}, nil
				},
			},
		},
	}
}
func MakeFloatType() *RuntimeType {
	return &RuntimeType{
		name: FLOAT_TYPE,
		methods: map[string]Function{
			"+": {
				Signature: []Parameter{
					{
						Name:     "self",
						TypeName: FLOAT_TYPE,
						Self:     true,
					},
					{
						Name:     "b",
						TypeName: FLOAT_TYPE,
					},
				},
				ReturnTypeName: FLOAT_TYPE,
				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
					a, _ := ((pv[0]).Value.Value).(float64)
					b, _ := ((pv[1]).Value.Value).(float64)
					return &RuntimeValue{
						TypeName: pv[0].Value.TypeName,
						Value:    a + b,
					}, nil
				},
			},
			"/": {
				Signature: []Parameter{
					{
						Name:     "self",
						TypeName: FLOAT_TYPE,
						Self:     true,
					},
					{
						Name:     "b",
						TypeName: FLOAT_TYPE,
					},
				},
				ReturnTypeName: FLOAT_TYPE,
				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
					a, _ := ((pv[0]).Value.Value).(float64)
					b, _ := ((pv[1]).Value.Value).(float64)
					return &RuntimeValue{
						TypeName: FLOAT_TYPE,
						Value:    a / b,
					}, nil
				},
			},
			"-": {
				Signature: []Parameter{
					{
						Name:     "self",
						TypeName: FLOAT_TYPE,
						Self:     true,
					},
					{
						Name:     "b",
						TypeName: FLOAT_TYPE,
					},
				},
				ReturnTypeName: FLOAT_TYPE,
				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
					a, _ := ((pv[0]).Value.Value).(float64)
					b, _ := ((pv[1]).Value.Value).(float64)
					return &RuntimeValue{
						TypeName: pv[0].Value.TypeName,
						Value:    a - b,
					}, nil
				},
			},
			"<->": {
				Signature: []Parameter{
					{
						Name:     "self",
						TypeName: FLOAT_TYPE,
						Self:     true,
					},
					{
						Name:     "b",
						TypeName: NUM_TYPE,
					},
				},
				ReturnTypeName: INT_TYPE,
				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
					a, _ := ((pv[0]).Value.Value).(float64)
					b, _ := ((pv[1]).Value.Value).(float64)
					var result int64 = 0
					if a < b {
						result = -1
					}
					if a > b {
						result = 1
					}
					return &RuntimeValue{
						TypeName: INT_TYPE,
						Value:    result,
					}, nil
				},
			},
			"*": {
				Signature: []Parameter{
					{
						Name:     "self",
						TypeName: FLOAT_TYPE,
						Self:     true,
					},
					{
						Name:     "b",
						TypeName: FLOAT_TYPE,
					},
				},
				ReturnTypeName: FLOAT_TYPE,
				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
					a, _ := ((pv[0]).Value.Value).(float64)
					b, _ := ((pv[1]).Value.Value).(float64)
					return &RuntimeValue{
						TypeName: pv[0].Value.TypeName,
						Value:    a * b,
					}, nil
				},
			},
		},
	}
}
func MakeBoolType() *RuntimeType {
	return &RuntimeType{
		name:    BOOL_TYPE,
		methods: map[string]Function{}}
}

func MakeTypeType() *RuntimeType {
	return &RuntimeType{
		name:    TYPE_TYPE,
		methods: map[string]Function{}}
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

func MakeArrayType() *RuntimeType {
	return &RuntimeType{
		name:    ARRAY_TYPE,
		methods: map[string]Function{}}
}
