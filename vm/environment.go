package vm

type Scope struct {
	variables map[string]*RuntimeValue
	parent    int
}

type Environment struct {
	scope []*Scope
}

func Set()
