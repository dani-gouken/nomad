package main

import (
	"os"

	"github.com/dani-gouken/nomad/interpreter"
	"github.com/dani-gouken/nomad/repl"
	"github.com/dani-gouken/nomad/vm"
)

func main() {
	repl.Start()
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
