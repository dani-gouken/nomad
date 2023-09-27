package vm

import (
	"fmt"

	"github.com/dani-gouken/nomad/parser"
	"github.com/dani-gouken/nomad/tokenizer"
)

func CompileExpr(expr parser.Expr) ([]Instruction, error) {
	instructions := []Instruction{}
	switch expr.Kind {
	case parser.EXPR_KIND_NOT:
		compiled, err := CompileExpr(expr.Exprs[0])
		if err != nil {
			return instructions, err
		}
		instructions = append(instructions, compiled...)
		return append(instructions, Instruction{
			Code: OP_NOT,
		}), nil
	case parser.EXPR_KIND_NEGATIVE:
		compiled, err := CompileExpr(expr.Exprs[0])
		if err != nil {
			return instructions, err
		}
		instructions = append(instructions, compiled...)
		return append(instructions, Instruction{
			Code: OP_NEGATIVE,
		}), nil
	case parser.EXPR_KIND_CONSTANT:
		t := expr.Token
		switch t.Kind {
		case tokenizer.TOKEN_KIND_TRUE:
			return []Instruction{
				{
					Code: OP_STORE_CONST,
					Arg1: VM_RUNTIME_BOOL,
					Arg2: OP_CONST_TRUE,
				},
			}, nil
		case tokenizer.TOKEN_KIND_NUM_LIT:
			return []Instruction{
				{
					Code: OP_STORE_CONST,
					Arg1: VM_RUNTIME_INT,
					Arg2: t.Content,
				},
			}, nil
		case tokenizer.TOKEN_KIND_FALSE:
			return []Instruction{
				{
					Code: OP_STORE_CONST,
					Arg1: VM_RUNTIME_BOOL,
					Arg2: OP_CONST_FALSE,
				},
			}, nil

		}
	case parser.EXPR_KIND_ADDITION:
		instructions, err := CompileBinaryExpr(expr)
		if err != nil {
			return instructions, err
		}
		return append(instructions, Instruction{
			Code: OP_ADD,
		}), nil
	case parser.EXPR_KIND_SUBSTRACTION:
		instructions, err := CompileBinaryExpr(expr)
		if err != nil {
			return instructions, err
		}
		return append(instructions, Instruction{
			Code: OP_SUB,
		}), nil
	case parser.EXPR_KIND_MULTIPLICATION:
		instructions, err := CompileBinaryExpr(expr)
		if err != nil {
			return instructions, err
		}
		return append(instructions, Instruction{
			Code: OP_MULT,
		}), nil
	case parser.EXPR_KIND_EQ:
		instructions, err := CompileBinaryExpr(expr)
		if err != nil {
			return instructions, err
		}
		return append(instructions, Instruction{
			Code: OP_EQ,
		}), nil
	}
	return instructions, fmt.Errorf("could not compile expression [%s]", expr.Kind)
}

func CompileBinaryExpr(expr parser.Expr) ([]Instruction, error) {
	instructions := []Instruction{}
	compiled, err := CompileExpr(expr.Exprs[0])
	if err != nil {
		return instructions, err
	}
	instructions = append(instructions, compiled...)
	compiled, err = CompileExpr(expr.Exprs[1])
	if err != nil {
		return instructions, err
	}
	instructions = append(instructions, compiled...)
	return instructions, nil
}

func CompileStmt(stmt parser.Stmt) ([]Instruction, error) {
	switch stmt.Kind {
	case parser.STMT_KIND_IMPLICIT_RETURN:
		instructions, err := CompileExpr(stmt.Expr)
		return instructions, err
	default:
		return []Instruction{}, fmt.Errorf("unable to compile statement [%s]", stmt.Kind)
	}
}

func Compile(program *parser.Program) ([]Instruction, error) {
	instructions := []Instruction{}
	for i := 0; i < len(program.Stmts); i++ {
		newInstructions, err := CompileStmt(program.Stmts[i])
		if err != nil {
			return instructions, err
		}
		instructions = append(instructions, newInstructions...)
	}
	instructions = append(instructions, Instruction{
		Code: OP_DEBUG_PRINT,
	})

	return instructions, nil
}
