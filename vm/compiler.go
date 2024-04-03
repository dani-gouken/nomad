package vm

import (
	"fmt"
	"strconv"
	"strings"

	nomadErrors "github.com/dani-gouken/nomad/errors"
	"github.com/dani-gouken/nomad/parser"
	"github.com/dani-gouken/nomad/runtime/types"
	"github.com/dani-gouken/nomad/tokenizer"
)

type Compiler struct {
	stmts        []*parser.Stmt
	instructions []Instruction
	cursor       int
}

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
	case parser.EXPR_KIND_LEN:
		compiled, err := CompileExpr(expr.Exprs[0])
		if err != nil {
			return instructions, err
		}
		instructions = append(instructions, compiled...)
		return append(instructions, Instruction{
			Code:       OP_LEN,
			DebugToken: expr.Token,
		}), nil
	case parser.EXPR_KIND_ANONYMOUS:
		for i := 0; i < len(expr.Exprs); i++ {
			compiled, err := CompileExpr(expr.Exprs[i])
			if err != nil {
				return instructions, err
			}
			instructions = append(instructions, compiled...)
		}
		return instructions, nil
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
					Code:       OP_PUSH_CONST,
					Arg1:       types.BOOL_TYPE,
					Arg2:       OP_CONST_TRUE,
					DebugToken: expr.Token,
				},
			}, nil
		case tokenizer.TOKEN_KIND_NUM_LIT:
			numType := types.INT_TYPE
			if strings.Contains(t.Content, ".") {
				numType = types.FLOAT_TYPE
			}
			return []Instruction{
				{
					Code:       OP_PUSH_CONST,
					Arg1:       numType,
					Arg2:       t.Content,
					DebugToken: expr.Token,
				},
			}, nil
		case tokenizer.TOKEN_KIND_STRING_LIT:
			content := expr.Token.Content
			if len(content) == 2 {
				content = ""
			} else {
				content = content[1 : len(content)-1]
			}
			return []Instruction{
				{
					Code:       OP_PUSH_CONST,
					Arg1:       types.STRING_TYPE,
					Arg2:       content,
					DebugToken: expr.Token,
				},
			}, nil
		case tokenizer.TOKEN_KIND_FALSE:
			return []Instruction{
				{
					Code:       OP_PUSH_CONST,
					Arg1:       types.BOOL_TYPE,
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
	case parser.EXPR_KIND_AND:
		instructions, err := CompileBinaryExpr(expr)
		if err != nil {
			return instructions, err
		}
		return append(instructions, Instruction{
			Code:       OP_AND,
			DebugToken: expr.Token,
		}), nil
	case parser.EXPR_KIND_OR:
		instructions, err := CompileBinaryExpr(expr)
		if err != nil {
			return instructions, err
		}
		return append(instructions, Instruction{
			Code:       OP_OR,
			DebugToken: expr.Token,
		}), nil
	case parser.EXPR_KIND_DIVISION:
		instructions, err := CompileBinaryExpr(expr)
		if err != nil {
			return instructions, err
		}
		return append(instructions, Instruction{
			Code:       OP_DIV,
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
	case parser.EXPR_KIND_LESS_THAN, parser.EXPR_KIND_MORE_THAN:
		instructions, err := CompileComp(expr)
		if err != nil {
			return instructions, err
		}
		expected := "1"
		if expr.Kind == parser.EXPR_KIND_LESS_THAN {
			expected = "-1"
		}
		instructions = append(instructions, Instruction{
			Code:       OP_PUSH_CONST,
			DebugToken: expr.Token,
			Arg1:       types.INT_TYPE,
			Arg2:       expected,
		})
		instructions = append(instructions, Instruction{
			Code:       OP_EQ,
			DebugToken: expr.Token,
		})
		return instructions, err
	case parser.EXPR_KIND_MORE_THAN_OR_EQ, parser.EXPR_LESS_THAN_OR_EQ:
		expected := "1"
		if expr.Kind == parser.EXPR_LESS_THAN_OR_EQ {
			expected = "-1"
		}
		instructions = append(instructions, Instruction{
			Code:       OP_PUSH_CONST,
			DebugToken: expr.Token,
			Arg1:       types.INT_TYPE,
			Arg2:       expected,
		})
		instructions = append(instructions, Instruction{
			Code:       OP_PUSH_CONST,
			DebugToken: expr.Token,
			Arg1:       types.INT_TYPE,
			Arg2:       "0",
		})
		exprInstructions, err := CompileComp(expr)
		if err != nil {
			return instructions, err
		}
		instructions = append(instructions, exprInstructions...)
		instructions = append(instructions, Instruction{
			Code:       OP_EQ_2,
			DebugToken: expr.Token,
		})
		return instructions, err
	case parser.EXPR_KIND_RIGHT_INCREMENT, parser.EXPR_KIND_RIGHT_DECREMENT:
		instructions, err := CompileExpr(expr.Exprs[0])
		if err != nil {
			return instructions, err
		}
		op := OP_ADD
		if expr.Kind == parser.EXPR_KIND_RIGHT_DECREMENT {
			op = OP_SUB
		}
		instructions = append(instructions, Instruction{
			Code:       OP_PUSH_CONST,
			DebugToken: expr.Token,
			Arg1:       types.INT_TYPE,
			Arg2:       "1",
		})
		instructions = append(instructions, Instruction{
			Code:       op,
			DebugToken: expr.Token,
		})
		instructions = append(instructions, Instruction{
			Code:       OP_SET_VAR,
			Arg1:       expr.Token.Content,
			DebugToken: expr.Token,
		})
		return instructions, err
	case parser.EXPR_KIND_ID:
		return append(instructions, Instruction{
			Code:       OP_LOAD_VAR,
			Arg1:       expr.Token.Content,
			DebugToken: expr.Token,
		}), nil
	case parser.EXPR_KIND_TYPE:
		return append(instructions, Instruction{
			Code:       OP_LOAD_TYPE,
			Arg1:       expr.Token.Content,
			DebugToken: expr.Token,
		}), nil
	case parser.EXPR_KIND_ARRAY:
		typeExpr := expr.Exprs[0]
		itemsExpr := expr.Exprs[1]
		typeInstr, err := CompileExpr(typeExpr)
		if err != nil {
			return instructions, err
		}
		instructions = append(instructions, typeInstr...)
		instructions = append(instructions, Instruction{
			Code:       OP_ARR_INIT,
			DebugToken: expr.Token,
		})
		for i := 0; i < len(itemsExpr.Exprs); i++ {
			exprInstructions, err := CompileExpr(itemsExpr.Exprs[i])
			if err != nil {
				return instructions, err
			}
			instructions = append(instructions, exprInstructions...)
			instructions = append(instructions, Instruction{
				Code:       OP_ARR_PUSH,
				DebugToken: itemsExpr.Exprs[i].Token,
			})
		}
		return instructions, nil
	case parser.EXPR_KIND_TYPE_ARRAY:
		typeExpr := expr.Exprs[0]
		typeInstr, err := CompileExpr(typeExpr)
		if err != nil {
			return instructions, err
		}
		instructions = append(instructions, typeInstr...)
		instructions = append(instructions, Instruction{
			Code:       OP_ARR_TYPE,
			DebugToken: expr.Token,
		})
		return instructions, nil
	case parser.EXPR_KIND_ARRAY_ACCESS:
		typeExpr := expr.Exprs[0]
		arrayInst, err := CompileExpr(typeExpr)
		if err != nil {
			return instructions, err
		}
		instructions = append(instructions, arrayInst...)
		if expr.Token.Kind == tokenizer.TOKEN_KIND_NUM_LIT {
			instructions = append(instructions, Instruction{
				Code:       OP_PUSH_CONST,
				Arg1:       types.INT_TYPE,
				Arg2:       expr.Token.Content,
				DebugToken: expr.Token,
			})
		} else {
			instructions = append(instructions, Instruction{
				Code:       OP_LOAD_VAR,
				Arg1:       expr.Token.Content,
				DebugToken: expr.Token,
			})

		}
		instructions = append(instructions, Instruction{
			Code:       OP_ARR_LOAD,
			DebugToken: expr.Token,
		})
		return instructions, nil
	}
	return instructions, fmt.Errorf("could not compile expression [%s]", expr.Kind)
}

func CompileComp(expr parser.Expr) ([]Instruction, error) {
	instructions, err := CompileBinaryExpr(expr)
	if err != nil {
		return instructions, err
	}
	return append(instructions, Instruction{
		Code:       OP_CMP,
		DebugToken: expr.Token,
	}), nil
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

func (c *Compiler) label(name string, stmt *parser.Stmt) string {
	return "__" + name + "_" + strconv.Itoa(c.cursor) + strconv.Itoa(stmt.Expr.Token.Loc.Line) + strconv.Itoa(stmt.Expr.Token.Loc.Start) + strconv.Itoa(stmt.Expr.Token.Loc.End)
}

func (c *Compiler) CompileStmt() error {
	stmt := c.peek()
	if stmt == nil {
		return nomadErrors.CompilationError("EOF")
	}
	switch stmt.Kind {
	case parser.STMT_KIND_IMPLICIT_RETURN:
		instructions, err := CompileExpr(stmt.Expr)
		c.consume()
		c.instructions = append(c.instructions, instructions...)
		return err
	case parser.STMT_KIND_IF:
		branches := []*parser.Stmt{}
		c.consume()
		branches = append(branches, stmt)
		exitIfLabel := c.label("ENDIF", stmt)
		for {
			stmt := c.peek()
			if stmt == nil {
				break
			}
			if stmt.Kind != parser.STMT_KIND_ELSE && stmt.Kind != parser.STMT_KIND_ELIF {
				break
			}

			branches = append(branches, stmt)
			c.consume()
		}
		var label string
		for i := 0; i < len(branches); i++ {
			branch := branches[i]
			if label == "" {
				label = c.label(branch.Kind, branch)
			}
			var nextLabel string = exitIfLabel
			if i+1 < len(branches) {
				nextLabel = c.label(branches[i+1].Kind, branches[i+1])
			}
			c.instructions = append(c.instructions, Instruction{
				Code: OP_LABEL,
				Arg1: label,
			})
			if branch.Kind != parser.STMT_KIND_ELSE {
				ifConditionInstructions, err := CompileExpr(branch.Expr)
				if err != nil {
					return err
				}
				c.instructions = append(c.instructions, ifConditionInstructions...)
				c.instructions = append(c.instructions, Instruction{
					Code: OP_JUMP_NOT,
					Arg1: nextLabel,
				})
			}
			branchStmts, err := CompileChunk(branch.Children)
			if err != nil {
				return err
			}
			c.instructions = append(c.instructions, branchStmts...)
			c.instructions = append(c.instructions, Instruction{
				Code: OP_JUMP,
				Arg1: exitIfLabel,
			})
			label = nextLabel
		}
		c.instructions = append(c.instructions, Instruction{
			Code: OP_LABEL,
			Arg1: exitIfLabel,
		})
		return nil
	case parser.STMT_KIND_FOR:
		c.consume()
		endForLabel := c.label("END_FOR", stmt)
		forTestLabel := c.label("FOR_TEST", stmt)
		c.instructions = append(c.instructions, Instruction{
			Code: OP_LABEL,
			Arg1: forTestLabel,
		})
		testExprInstructions, err := CompileExpr(stmt.Expr)
		if err != nil {
			return err
		}
		c.instructions = append(c.instructions, testExprInstructions...)
		c.instructions = append(c.instructions, Instruction{
			Code: OP_PUSH_CONST,
			Arg1: types.BOOL_TYPE,
			Arg2: OP_CONST_TRUE,
		})
		c.instructions = append(c.instructions, Instruction{
			Code: OP_EQ,
		})
		c.instructions = append(c.instructions, Instruction{
			Code: OP_JUMP_NOT,
			Arg1: endForLabel,
		})
		forOperationsInstructions, err := CompileChunk(stmt.Children)
		if err != nil {
			return err
		}
		c.instructions = append(c.instructions, forOperationsInstructions...)
		c.instructions = append(c.instructions, Instruction{
			Code: OP_JUMP,
			Arg1: forTestLabel,
		})
		c.instructions = append(c.instructions, Instruction{
			Code: OP_LABEL,
			Arg1: endForLabel,
		})
		return err
	case parser.STMT_KIND_ASSIGNMENT:
		varName := stmt.Data[0].Content
		compiled, err := CompileExpr(stmt.Expr)
		c.instructions = append(c.instructions, compiled...)
		c.instructions = append(c.instructions, Instruction{
			Code:       OP_SET_VAR,
			Arg1:       varName,
			DebugToken: stmt.Expr.Token,
		})
		c.instructions = append(c.instructions, Instruction{
			Code:       OP_POP_CONST,
			DebugToken: stmt.Expr.Token,
		})
		c.consume()
		return err
	case parser.STMT_KIND_DEBUG_PRINT:
		compiled, err := CompileExpr(stmt.Expr)
		c.instructions = append(c.instructions, compiled...)
		c.instructions = append(c.instructions, Instruction{
			Code: OP_DEBUG_PRINT,
		})
		c.consume()
		return err
	case parser.STMT_KIND_TYPE_DECLARATION:
		typeName := stmt.Data[0].Content
		compiled, err := CompileExpr(stmt.Expr)
		if err != nil {
			return err
		}
		c.instructions = append(c.instructions, compiled...)
		c.instructions = append(c.instructions, Instruction{
			Code:       OP_DECL_TYPE,
			Arg1:       typeName,
			DebugToken: stmt.Expr.Token,
		})
		c.consume()
		return err
	case parser.STMT_KIND_VAR_DECLARATION:
		varName := stmt.Data[0].Content
		compiled, err := CompileExpr(stmt.Expr)
		c.instructions = append(c.instructions, compiled...)
		c.instructions = append(c.instructions, Instruction{
			Code:       OP_DECL_VAR,
			Arg2:       varName,
			DebugToken: stmt.Expr.Token,
		})
		c.instructions = append(c.instructions, Instruction{
			Code:       OP_POP_CONST,
			DebugToken: stmt.Expr.Token,
		})
		c.consume()
		return err
	default:
		return fmt.Errorf("unable to compile statement [%s]", stmt.Kind)
	}
}
func (c *Compiler) GetInstructions() []Instruction {
	return c.instructions
}
func (c *Compiler) CompileChunk(stmts []*parser.Stmt) ([]Instruction, error) {
	c.stmts = stmts
	for c.peek() != nil {
		err := c.CompileStmt()
		if err != nil {
			return c.GetInstructions(), err
		}
	}
	return c.GetInstructions(), nil
}
func RemoveLabels(instructions []Instruction) ([]Instruction, error) {
	labels := map[string]int{}
	for i := 0; i < len(instructions); i++ {
		instruction := instructions[i]
		if instruction.Code == OP_LABEL {
			labels[instruction.Arg1] = i + 1
		}
	}

	for i := 0; i < len(instructions); i++ {
		instruction := &instructions[i]
		if instruction.Code == OP_JUMP || instruction.Code == OP_JUMP_NOT {
			instruction.Arg1 = strconv.Itoa(labels[instruction.Arg1])
		}
	}
	return instructions, nil
}
func (c *Compiler) Compile(stmts []*parser.Stmt) ([]Instruction, error) {
	instructions, err := c.CompileChunk(stmts)
	if err != nil {
		return instructions, err
	}
	return RemoveLabels(instructions)
}

func Compile(program []*parser.Stmt) ([]Instruction, error) {
	compiler := Compiler{}
	return compiler.Compile(program)
}
func CompileChunk(program []*parser.Stmt) ([]Instruction, error) {
	compiler := Compiler{}
	return compiler.CompileChunk(program)
}

func (c *Compiler) peek() *parser.Stmt {
	if c.cursor >= len(c.stmts) {
		return nil
	}
	return c.stmts[c.cursor]
}
func (c *Compiler) peekAt(pos int) *parser.Stmt {
	if (c.cursor+pos < 0) || (c.cursor+pos) >= len(c.stmts) {
		return nil
	}
	return c.stmts[c.cursor+pos]
}

func (c *Compiler) consume() {
	c.cursor++
}

func (c *Compiler) rollback(position int) {
	c.cursor = position
}

func DebugPrintOpCode(instructions []Instruction) {
	for i := 0; i < len(instructions); i++ {
		instruction := instructions[i]
		fmt.Printf("%s    %s %s", instruction.Code, instruction.Arg1, instruction.Arg2)
		fmt.Println()
	}
}
