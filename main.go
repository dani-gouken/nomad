package main

import (
	"os"

	"github.com/dani-gouken/nomad/interpreter"
	"github.com/dani-gouken/nomad/vm"
)

func main() {

	if len(os.Args) < 2 {
		panic("source file is needed")
	}
	sourceFile := os.Args[1]
	bytes, err := os.ReadFile(sourceFile)
	if err != nil {
		panic(err)
	}

	instance := vm.New()

	interpreter := interpreter.NewInterpreter()
	err = interpreter.Interpret(string(bytes), instance)
	if err != nil {
		println(err.Error())
	}

}
