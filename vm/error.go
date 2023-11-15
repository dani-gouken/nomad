package vm

import (
	"fmt"

	"github.com/dani-gouken/nomad/tokenizer"
)

func RuntimeError(message string, debugToken tokenizer.Token) error {
	return fmt.Errorf("%s: runtime error. %s", DebugToken(debugToken), message)
}
func RuntimeErrorUnsupportedOperand(operand string, typeName string, debugToken tokenizer.Token) error {
	return RuntimeError(fmt.Sprintf("unsupported operand %s on type %s", operand, typeName), debugToken)
}

func ParseError(message string, debugToken tokenizer.Token) error {
	return fmt.Errorf("parse error: %s. %s", message, DebugToken(debugToken))
}

func DebugToken(token tokenizer.Token) string {
	return fmt.Sprintf("/%d:%d:%d", token.Loc.Line, token.Loc.Start, token.Loc.End)
}
