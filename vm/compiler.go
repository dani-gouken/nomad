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
			Code:       OP_NOT,
			DebugToken: expr.Token,
		}), nil
	case parser.EXPR_KIND_NEGATIVE:
		compiled, err := CompileExpr(expr.Exprs[0])
		if err != nil {
			return instructions, err
		}
		instructions = append(instructions, compiled...)
		return append(instructions, Instruction{
			Code:       OP_NEGATIVE,
			DebugToken: expr.Token,
		}), nil
	case parser.EXPR_KIND_CONSTANT:
		t := expr.Token
		switch t.Kind {
		case tokenizer.TOKEN_KIND_TRUE:
			return []Instruction{
				{
					Code:       OP_STORE_CONST,
					Arg1:       BOOL_TYPE,
					Arg2:       OP_CONST_TRUE,
					DebugToken: expr.Token,
				},
			}, nil
		case tokenizer.TOKEN_KIND_NUM_LIT:
			return []Instruction{
				{
					Code:       OP_STORE_CONST,
					Arg1:       INT_TYPE,
					Arg2:       t.Content,
					DebugToken: expr.Token,
				},
			}, nil
		case tokenizer.TOKEN_KIND_FALSE:
			return []Instruction{
				{
					Code:       OP_STORE_CONST,
					Arg1:       BOOL_TYPE,
					Arg2:       OP_CONST_FALSE,
					DebugToken: expr.Token,
				},
			}, nil

		}
	case parser.EXPR_KIND_ADDITION:
		instructions, err := CompileBinaryExpr(expr)
		if err != nil {
			return instructions, err
		}
		return append(instructions, Instruction{
			Code:       OP_ADD,
			DebugToken: expr.Token,
		}), nil
	case parser.EXPR_KIND_SUBSTRACTION:
		instructions, err := CompileBinaryExpr(expr)
		if err != nil {
			return instructions, err
		}
		return append(instructions, Instruction{
			Code:       OP_SUB,
			DebugToken: expr.Token,
		}), nil
	case parser.EXPR_KIND_MULTIPLICATION:
		instructions, err := CompileBinaryExpr(expr)
		if err != nil {
			return instructions, err
		}
		return append(instructions, Instruction{
			Code:       OP_MULT,
			DebugToken: expr.Token,
		}), nil
	case parser.EXPR_KIND_EQ:
		instructions, err := CompileBinaryExpr(expr)
		if err != nil {
			return instructions, err
		}
		return append(instructions, Instruction{
			Code:       OP_EQ,
			DebugToken: expr.Token,
		}), nil
	case parser.EXPR_KIND_ID:
		return append(instructions, Instruction{
			Code:       OP_LOAD_VAR,
			Arg1:       expr.Token.Content,
			DebugToken: expr.Token,
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
	case parser.STMT_KIND_VAR_DECLARATION:
		instructions := []Instruction{}
		varType := stmt.Data[0].Content
		varName := stmt.Data[1].Content
		compiled, err := CompileExpr(stmt.Expr)
		instructions = append(instructions, compiled...)
		instructions = append(instructions, Instruction{
			Code:       OP_STORE_VAR,
			Arg1:       varType,
			Arg2:       varName,
			DebugToken: stmt.Expr.Token,
		})
		instructions = append(instructions, Instruction{
			Code:       OP_POP_CONST,
			Arg1:       varName,
			DebugToken: stmt.Expr.Token,
		})
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
