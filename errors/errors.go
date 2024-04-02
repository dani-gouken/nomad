package error

import (
	"fmt"

	"github.com/dani-gouken/nomad/tokenizer"
)

type ParseError struct {
	crash   bool
	message string
}

func (e *ParseError) Error() string {
	return e.message
}

func (e *ParseError) ShouldCrash() bool {
	return e.crash
}

func RuntimeError(message string, debugToken tokenizer.Token) error {
	return fmt.Errorf("%s: runtime error. %s", DebugToken(debugToken), message)
}

func CompilationError(message string) error {
	return fmt.Errorf("compilation error. %s", message)
}

func RuntimeErrorUnsupportedOperand(operand string, typeName string, debugToken tokenizer.Token) error {
	return RuntimeError(fmt.Sprintf("unsupported operand %s on type %s", operand, typeName), debugToken)
}

func NewParseError(message string, debugToken tokenizer.Token, crash bool) *ParseError {
	return &ParseError{message: fmt.Sprintf("%s: parse error. %s", DebugToken(debugToken), message), crash: crash}
}

func NewParseErrorFromMessage(message string, crash bool) *ParseError {
	return &ParseError{message: message, crash: crash}
}

func FatalParseError(message string, debugToken tokenizer.Token) *ParseError {
	return NewParseError(message, debugToken, true)
}

func NonFatalParseError(message string, debugToken tokenizer.Token) *ParseError {
	return NewParseError(message, debugToken, false)
}

func DebugToken(token tokenizer.Token) string {
	return fmt.Sprintf("/%d:%d:%d", token.Loc.Line, token.Loc.Start, token.Loc.End)
}
