package types

const (
	INT_TYPE    = "int"
	BOOL_TYPE   = "bool"
	TYPE_TYPE   = "type"
	FLOAT_TYPE  = "float"
	NUM_TYPE    = "num"
	STRING_TYPE = "string"
	ARRAY_TYPE  = "array"
)

type RuntimeTypeType = int

type RuntimeType interface {
	GetName() string
	Match(t2 RuntimeType) error
}
