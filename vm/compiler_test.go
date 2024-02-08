package vm_test

import (
	"testing"

	"github.com/dani-gouken/nomad/parser"
	"github.com/dani-gouken/nomad/tokenizer"
	"github.com/dani-gouken/nomad/vm"
	"github.com/stretchr/testify/assert"
)

func TestCompiler(t *testing.T) {
	code := `
		int a = 2; 
		int b = 0; 
		if(a == 1) {
			b  = 2;
		} elif (a == 2) {
			b = 4;
		} else {
			b = -1;
		}
	`

	tokens, err := tokenizer.Tokenize(code)
	assert.Nil(t, err)

	ast, err := parser.Parse(tokens)
	assert.Nil(t, err)

	instructions, err := vm.Compile(ast.Stmts)
	assert.Nil(t, err)
	vm.DebugPrintOpCode(instructions)

}
