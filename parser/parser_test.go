package parser_test

import (
	"testing"

	"github.com/dani-gouken/nomad/parser"
	"github.com/dani-gouken/nomad/tokenizer"
	"github.com/stretchr/testify/assert"
)

func TestParseNotExpr(t *testing.T) {
	actual, err := tokenizer.Tokenize("!true")
	assert.NoError(t, err)
	expected := []tokenizer.Token{
		{
			Kind: tokenizer.TOKEN_KIND_BANG,
			Loc: tokenizer.TokenLoc{
				Start: 0,
				End:   0,
				Line:  1,
			},
			Content: "!",
		},
		{
			Kind: tokenizer.TOKEN_KIND_TRUE,
			Loc: tokenizer.TokenLoc{
				Start: 1,
				End:   4,
				Line:  1,
			},
			Content: "true",
		},
	}
	assert.Equal(t, actual, expected)

	ast, err := parser.Parse(actual)
	assert.NoError(t, err)

	assert.Equal(t, ast, &parser.Program{
		Stmts: []parser.Stmt{
			{
				Kind: parser.STMT_KIND_IMPLICIT_RETURN,
				Expr: parser.Expr{
					Kind:  parser.EXPR_KIND_NOT,
					Token: actual[0],
					Exprs: []parser.Expr{
						{
							Kind:  parser.EXPR_KIND_CONSTANT,
							Token: actual[1],
						},
					},
				},
			},
		},
	})
}

func TestParseNegativeExpr(t *testing.T) {
	actual, err := tokenizer.Tokenize("-1")
	assert.NoError(t, err)
	expected := []tokenizer.Token{
		{
			Kind: tokenizer.TOKEN_KIND_MINUS,
			Loc: tokenizer.TokenLoc{
				Start: 0,
				End:   0,
				Line:  1,
			},
			Content: "-",
		},
		{
			Kind: tokenizer.TOKEN_KIND_NUM_LIT,
			Loc: tokenizer.TokenLoc{
				Start: 1,
				End:   1,
				Line:  1,
			},
			Content: "1",
		},
	}
	assert.Equal(t, actual, expected)

	ast, err := parser.Parse(actual)
	assert.NoError(t, err)

	assert.Equal(t, ast, &parser.Program{
		Stmts: []parser.Stmt{
			{
				Kind: parser.STMT_KIND_IMPLICIT_RETURN,
				Expr: parser.Expr{
					Kind:  parser.EXPR_KIND_NEGATIVE,
					Token: actual[0],
					Exprs: []parser.Expr{
						{
							Kind:  parser.EXPR_KIND_CONSTANT,
							Token: actual[1],
						},
					},
				},
			},
		},
	})
}

func TestParseIncrementExpr(t *testing.T) {
	actual, err := tokenizer.Tokenize("++a")
	assert.NoError(t, err)
	expected := []tokenizer.Token{
		{
			Kind: tokenizer.TOKEN_KIND_DB_PLUS,
			Loc: tokenizer.TokenLoc{
				Start: 0,
				End:   1,
				Line:  1,
			},
			Content: "++",
		},
		{
			Kind: tokenizer.TOKEN_KIND_ID,
			Loc: tokenizer.TokenLoc{
				Start: 2,
				End:   2,
				Line:  1,
			},
			Content: "a",
		},
	}
	assert.Equal(t, actual, expected)

	ast, err := parser.Parse(actual)
	assert.NoError(t, err)

	assert.Equal(t, &parser.Program{
		Stmts: []parser.Stmt{
			{
				Kind: parser.STMT_KIND_IMPLICIT_RETURN,
				Expr: parser.Expr{
					Kind:  parser.EXPR_KIND_LEFT_INCREMENT,
					Token: actual[0],
					Exprs: []parser.Expr{
						{
							Kind:  parser.EXPR_KIND_ID,
							Token: actual[1],
						},
					},
				},
			},
		},
	}, ast)
}

func TestParseDecrementExpr(t *testing.T) {
	actual, err := tokenizer.Tokenize("--a")
	assert.NoError(t, err)
	expected := []tokenizer.Token{
		{
			Kind: tokenizer.TOKEN_KIND_DB_MINUS,
			Loc: tokenizer.TokenLoc{
				Start: 0,
				End:   1,
				Line:  1,
			},
			Content: "--",
		},
		{
			Kind: tokenizer.TOKEN_KIND_ID,
			Loc: tokenizer.TokenLoc{
				Start: 2,
				End:   2,
				Line:  1,
			},
			Content: "a",
		},
	}
	assert.Equal(t, actual, expected)

	ast, err := parser.Parse(actual)
	assert.NoError(t, err)

	assert.Equal(t, &parser.Program{
		Stmts: []parser.Stmt{
			{
				Kind: parser.STMT_KIND_IMPLICIT_RETURN,
				Expr: parser.Expr{
					Kind:  parser.EXPR_KIND_LEFT_DECREMENT,
					Token: actual[0],
					Exprs: []parser.Expr{
						{
							Kind:  parser.EXPR_KIND_ID,
							Token: actual[1],
						},
					},
				},
			},
		},
	}, ast)
}

func TestParseAddition(t *testing.T) {
	actual, err := tokenizer.Tokenize("1+1")
	assert.NoError(t, err)
	expected := []tokenizer.Token{
		{
			Kind: tokenizer.TOKEN_KIND_NUM_LIT,
			Loc: tokenizer.TokenLoc{
				Start: 0,
				End:   0,
				Line:  1,
			},
			Content: "1",
		},
		{
			Kind: tokenizer.TOKEN_KIND_PLUS,
			Loc: tokenizer.TokenLoc{
				Start: 1,
				End:   1,
				Line:  1,
			},
			Content: "+",
		},
		{
			Kind: tokenizer.TOKEN_KIND_NUM_LIT,
			Loc: tokenizer.TokenLoc{
				Start: 2,
				End:   2,
				Line:  1,
			},
			Content: "1",
		},
	}
	assert.Equal(t, actual, expected)

	ast, err := parser.Parse(actual)
	assert.NoError(t, err)

	assert.Equal(t, ast, &parser.Program{
		Stmts: []parser.Stmt{
			{
				Kind: parser.STMT_KIND_IMPLICIT_RETURN,
				Expr: parser.Expr{
					Kind:  parser.EXPR_KIND_ADDITION,
					Token: actual[1],
					Exprs: []parser.Expr{
						{
							Kind:  parser.EXPR_KIND_CONSTANT,
							Token: actual[0],
						},
						{
							Kind:  parser.EXPR_KIND_CONSTANT,
							Token: actual[2],
						},
					},
				},
			},
		},
	})
}

func TestOperatorPrecedence(t *testing.T) {
	actual, err := tokenizer.Tokenize("1+2*3-69")
	assert.NoError(t, err)
	expected := []tokenizer.Token{
		{
			Kind: tokenizer.TOKEN_KIND_NUM_LIT,
			Loc: tokenizer.TokenLoc{
				Start: 0,
				End:   0,
				Line:  1,
			},
			Content: "1",
		},
		{
			Kind: tokenizer.TOKEN_KIND_PLUS,
			Loc: tokenizer.TokenLoc{
				Start: 1,
				End:   1,
				Line:  1,
			},
			Content: "+",
		},
		{
			Kind: tokenizer.TOKEN_KIND_NUM_LIT,
			Loc: tokenizer.TokenLoc{
				Start: 2,
				End:   2,
				Line:  1,
			},
			Content: "2",
		},
		{
			Kind: tokenizer.TOKEN_KIND_STAR,
			Loc: tokenizer.TokenLoc{
				Start: 3,
				End:   3,
				Line:  1,
			},
			Content: "*",
		},
		{
			Kind: tokenizer.TOKEN_KIND_NUM_LIT,
			Loc: tokenizer.TokenLoc{
				Start: 4,
				End:   4,
				Line:  1,
			},
			Content: "3",
		},
		{
			Kind: tokenizer.TOKEN_KIND_MINUS,
			Loc: tokenizer.TokenLoc{
				Start: 5,
				End:   5,
				Line:  1,
			},
			Content: "-",
		},
		{
			Kind: tokenizer.TOKEN_KIND_NUM_LIT,
			Loc: tokenizer.TokenLoc{
				Start: 6,
				End:   7,
				Line:  1,
			},
			Content: "69",
		},
	}
	assert.Equal(t, actual, expected)

	ast, err := parser.Parse(actual)
	assert.NoError(t, err)
	stmt := ast.Stmts[0]
	sexpr := parser.ExprToSExpr(stmt.Expr)
	assert.Equal(t, sexpr, "(- (+ 1 (* 2 3)) 69)")

}

func TestOperatorPrecedenceWithBracket(t *testing.T) {
	tokens, err := tokenizer.Tokenize("(1+2)*3-69")

	assert.NoError(t, err)
	ast, err := parser.Parse(tokens)
	assert.NoError(t, err)
	stmt := ast.Stmts[0]
	sexpr := parser.ExprToSExpr(stmt.Expr)
	assert.Equal(t, sexpr, "(- (* (+ 1 2) 3) 69)")

}
