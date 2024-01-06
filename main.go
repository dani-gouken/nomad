package main

import (
	"github.com/dani-gouken/nomad/interpreter"
	"github.com/dani-gouken/nomad/repl"
	"github.com/dani-gouken/nomad/vm"
)

func main() {

	// i,🚀 = 60,9;
	// printf "The result is: \"%s\"", i+🚀
	repl.Start()
	instance := vm.New()
	source := `(1/1)+1`
	interpreter := interpreter.NewInterpreter()
	interpreter.Interpret(source, instance)
}
