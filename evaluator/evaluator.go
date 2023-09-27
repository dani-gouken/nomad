package evaluator

import "github.com/dani-gouken/nomad/parser"

type Evaluator struct {
}

func (*Evaluator) eval(program *parser.Program) {
	return
}

func NewEvaluator() Evaluator {
	return Evaluator{}
}
