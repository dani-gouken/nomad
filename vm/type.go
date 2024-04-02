package vm

// func MakeIntType(registrar *TypeRegistrar) *RuntimeType {
// 	return &RuntimeType{
// 		name: INT_TYPE,
// 		methods: map[string]Function{
// 			"+": {
// 				Signature: []Parameter{
// 					{
// 						Name:     "self",
// 						TypeName: INT_TYPE,
// 						Self:     true,
// 					},
// 					{
// 						Name:     "b",
// 						TypeName: INT_TYPE,
// 					},
// 				},
// 				ReturnTypeName: INT_TYPE,
// 				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
// 					a, _ := ((pv[0]).Value.Value).(int64)
// 					b, _ := ((pv[1]).Value.Value).(int64)
// 					intType, err := registrar.Get(INT_TYPE)
// 					if err != nil {
// 						panic("int type not registered")
// 					}
// 					return &RuntimeValue{
// 						RuntimeType: intType,
// 						Value:       a + b,
// 					}, nil
// 				},
// 			},
// 			"-": {
// 				Signature: []Parameter{
// 					{
// 						Name:     "self",
// 						TypeName: INT_TYPE,
// 						Self:     true,
// 					},
// 					{
// 						Name:     "b",
// 						TypeName: INT_TYPE,
// 					},
// 				},
// 				ReturnTypeName: INT_TYPE,
// 				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
// 					a, _ := ((pv[0]).Value.Value).(int64)
// 					b, _ := ((pv[1]).Value.Value).(int64)
// 					intType, err := registrar.Get(INT_TYPE)
// 					if err != nil {
// 						panic("int type not registered")
// 					}
// 					return &RuntimeValue{
// 						RuntimeType: intType,
// 						Value:       a - b,
// 					}, nil
// 				},
// 			},
// 			"<->": {
// 				Signature: []Parameter{
// 					{
// 						Name:     "self",
// 						TypeName: INT_TYPE,
// 						Self:     true,
// 					},
// 					{
// 						Name:     "b",
// 						TypeName: NUM_TYPE,
// 					},
// 				},
// 				ReturnTypeName: INT_TYPE,
// 				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
// 					a, _ := ((pv[0]).Value.Value).(int64)
// 					b, _ := ((pv[1]).Value.Value).(int64)
// 					var result int64 = 0
// 					if a < b {
// 						result = -1
// 					}
// 					if a > b {
// 						result = 1
// 					}
// 					intType, err := registrar.Get(INT_TYPE)
// 					if err != nil {
// 						panic("int type not registered")
// 					}
// 					return &RuntimeValue{
// 						RuntimeType: intType,
// 						Value:       result,
// 					}, nil
// 				},
// 			},
// 			"/": {
// 				Signature: []Parameter{
// 					{
// 						Name:     "self",
// 						TypeName: INT_TYPE,
// 						Self:     true,
// 					},
// 					{
// 						Name:     "b",
// 						TypeName: INT_TYPE,
// 					},
// 				},
// 				ReturnTypeName: FLOAT_TYPE,
// 				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
// 					a, _ := ((pv[0]).Value.Value).(int64)
// 					b, _ := ((pv[1]).Value.Value).(int64)

// 					intType := registrar.GetOrPanic(INT_TYPE)

// 					return &RuntimeValue{
// 						RuntimeType: intType,
// 						Value:       float64(a) / float64(b),
// 					}, nil
// 				},
// 			},
// 			"*": {
// 				Signature: []Parameter{
// 					{
// 						Name:     "self",
// 						TypeName: "int",
// 						Self:     true,
// 					},
// 					{
// 						Name:     "b",
// 						TypeName: "int",
// 					},
// 				},
// 				ReturnTypeName: "int",
// 				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
// 					a, _ := ((pv[0]).Value.Value).(int64)
// 					b, _ := ((pv[1]).Value.Value).(int64)
// 					intType := registrar.GetOrPanic(INT_TYPE)

// 					return &RuntimeValue{
// 						RuntimeType: intType,
// 						Value:       a * b,
// 					}, nil
// 				},
// 			},
// 		},
// 	}
// }
// func MakeFloatType(registrar *TypeRegistrar) *RuntimeType {
// 	return &RuntimeType{
// 		name: FLOAT_TYPE,
// 		methods: map[string]Function{
// 			"+": {
// 				Signature: []Parameter{
// 					{
// 						Name:     "self",
// 						TypeName: FLOAT_TYPE,
// 						Self:     true,
// 					},
// 					{
// 						Name:     "b",
// 						TypeName: FLOAT_TYPE,
// 					},
// 				},
// 				ReturnTypeName: FLOAT_TYPE,
// 				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
// 					a, _ := ((pv[0]).Value.Value).(float64)
// 					b, _ := ((pv[1]).Value.Value).(float64)
// 					t := registrar.GetOrPanic(FLOAT_TYPE)
// 					return &RuntimeValue{
// 						RuntimeType: t,
// 						Value:       a + b,
// 					}, nil
// 				},
// 			},
// 			"/": {
// 				Signature: []Parameter{
// 					{
// 						Name:     "self",
// 						TypeName: FLOAT_TYPE,
// 						Self:     true,
// 					},
// 					{
// 						Name:     "b",
// 						TypeName: FLOAT_TYPE,
// 					},
// 				},
// 				ReturnTypeName: FLOAT_TYPE,
// 				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
// 					a, _ := ((pv[0]).Value.Value).(float64)
// 					b, _ := ((pv[1]).Value.Value).(float64)
// 					t := registrar.GetOrPanic(FLOAT_TYPE)

// 					return &RuntimeValue{
// 						RuntimeType: t,
// 						Value:       a / b,
// 					}, nil
// 				},
// 			},
// 			"-": {
// 				Signature: []Parameter{
// 					{
// 						Name:     "self",
// 						TypeName: FLOAT_TYPE,
// 						Self:     true,
// 					},
// 					{
// 						Name:     "b",
// 						TypeName: FLOAT_TYPE,
// 					},
// 				},
// 				ReturnTypeName: FLOAT_TYPE,
// 				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
// 					a, _ := ((pv[0]).Value.Value).(float64)
// 					b, _ := ((pv[1]).Value.Value).(float64)
// 					t := registrar.GetOrPanic(FLOAT_TYPE)

// 					return &RuntimeValue{
// 						RuntimeType: t,
// 						Value:       a - b,
// 					}, nil
// 				},
// 			},
// 			"<->": {
// 				Signature: []Parameter{
// 					{
// 						Name:     "self",
// 						TypeName: FLOAT_TYPE,
// 						Self:     true,
// 					},
// 					{
// 						Name:     "b",
// 						TypeName: NUM_TYPE,
// 					},
// 				},
// 				ReturnTypeName: INT_TYPE,
// 				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
// 					a, _ := ((pv[0]).Value.Value).(float64)
// 					b, _ := ((pv[1]).Value.Value).(float64)
// 					var result int64 = 0
// 					if a < b {
// 						result = -1
// 					}
// 					if a > b {
// 						result = 1
// 					}
// 					t := registrar.GetOrPanic(INT_TYPE)
// 					return &RuntimeValue{
// 						RuntimeType: t,
// 						Value:       result,
// 					}, nil
// 				},
// 			},
// 			"*": {
// 				Signature: []Parameter{
// 					{
// 						Name:     "self",
// 						TypeName: FLOAT_TYPE,
// 						Self:     true,
// 					},
// 					{
// 						Name:     "b",
// 						TypeName: FLOAT_TYPE,
// 					},
// 				},
// 				ReturnTypeName: FLOAT_TYPE,
// 				Run: func(pv []ParameterValue) (*RuntimeValue, error) {
// 					a, _ := ((pv[0]).Value.Value).(float64)
// 					b, _ := ((pv[1]).Value.Value).(float64)

// 					return &RuntimeValue{
// 						RuntimeType: pv[0].Value.RuntimeType,
// 						Value:       a * b,
// 					}, nil
// 				},
// 			},
// 		},
// 	}
// }
// func MakeBoolType() *RuntimeType {
// 	return &RuntimeType{
// 		name:    BOOL_TYPE,
// 		methods: map[string]Function{}}
// }

// func MakeTypeType() *RuntimeType {
// 	return &RuntimeType{
// 		name:    TYPE_TYPE,
// 		methods: map[string]Function{}}
// }

// func (t *RuntimeType) GetDeclaredMethods() map[string]Function {
// 	return t.methods
// }

// func (t *RuntimeType) GetMethod(name string) (Function, error) {
// 	res, ok := t.methods[name]
// 	if !ok {
// 		return res, fmt.Errorf("method %s::%s not found", t.GetName(), name)
// 	}
// 	return res, nil
// }

// func MakeArrayType() *RuntimeType {
// 	return &RuntimeType{
// 		name:    ARRAY_TYPE,
// 		methods: map[string]Function{}}
// }
