package vm_test

import (
	"testing"

	"github.com/dani-gouken/nomad/vm"
	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	stack := vm.Stack{}
	err := stack.PushInt(4)
	assert.NoError(t, err)

	v, err := stack.Current()
	assert.NoError(t, err)
	assert.Equal(t, 4, v)
}
