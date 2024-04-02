package vm_test

import (
	"testing"

	"github.com/dani-gouken/nomad/runtime/types"
	"github.com/dani-gouken/nomad/vm"
	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	stack := vm.Stack{}
	reg := types.NewRegistrar()
	err := stack.PushInt(reg, 4)
	assert.NoError(t, err)

	v, err := stack.Current()
	assert.NoError(t, err)
	assert.Equal(t, 4, v)
}
