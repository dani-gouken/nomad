package main

import (
	"github.com/dani-gouken/nomad/interpreter"
	"github.com/dani-gouken/nomad/vm"
)

func main() {

	// i,🚀 = 60,9;
	// printf "The result is: \"%s\"", i+🚀
	source := `
	for int i = 0; i < 10; i = i+1 {
		print i
	}
	`
	instance := vm.New()

	interpreter := interpreter.NewInterpreter()
	err := interpreter.Interpret(source, instance)
	if err != nil {
		println(err.Error())
	}

}
