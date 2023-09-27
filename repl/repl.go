package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/dani-gouken/nomad/interpreter"
	"github.com/dani-gouken/nomad/vm"
)

func Start() {
	instance := vm.New()
	interpreter := interpreter.NewInterpreter()

	for {
		fmt.Print("(nomad) > ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		cmd := input.Text()
		if cmd == "exit" {
			break
		}
		err := interpreter.Interpret(cmd, instance)
		if err != nil {
			println(err.Error())
		}
	}
}
