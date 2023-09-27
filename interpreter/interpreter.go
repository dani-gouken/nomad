package interpreter

import (
	"github.com/dani-gouken/nomad/parser"
	"github.com/dani-gouken/nomad/tokenizer"
	"github.com/dani-gouken/nomad/vm"
)

type Interpreter struct{}

func NewInterpreter() Interpreter {
	return Interpreter{}
}
func (p *Interpreter) Interpret(code string, instance *vm.Vm) error {
	tokens, err := tokenizer.Tokenize(code)
	if err != nil {
		return err
	}
	program, err := parser.Parse(tokens)
	if err != nil {
		return err
	}
	opCode, err := vm.Compile(program)
	if err != nil {
		return err
	}
	return instance.Interpret(opCode)

}
