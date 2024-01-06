package vm

type Function struct {
	Name           string
	Signature      []Parameter
	ReturnTypeName string
	Run            func([]ParameterValue) (*RuntimeValue, error)
}

type Parameter struct {
	Name     string
	Self     bool
	TypeName string
}
type ParameterValue struct {
	Parameter
	Value *RuntimeValue
}

type RuntimeValue struct {
	TypeName      string
	Value         interface{}
	PossibleTypes []string
}
