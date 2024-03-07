package parser

import (
	"fmt"
	"strings"

	"github.com/dani-gouken/nomad/tokenizer"
)

const (
	EXPR_KIND_CONSTANT        = "CONSTANT"
	EXPR_KIND_NOT             = "NOT"
	EXPR_KIND_NEGATIVE        = "NEGATIVE"
	EXPR_KIND_LEFT_INCREMENT  = "LEFT_INCREMENT"
	EXPR_KIND_RIGHT_INCREMENT = "RIGHT_INCREMENT"
	EXPR_KIND_LEFT_DECREMENT  = "LEFT_DECREMENT"
	EXPR_KIND_RIGHT_DECREMENT = "RIGHT_DECREMENT"
	EXPR_KIND_ADDITION        = "ADDITION"
	EXPR_KIND_DIVISION        = "DIVISION"
	EXPR_KIND_SUBSTRACTION    = "SUBSTRACTION"
	EXPR_KIND_MULTIPLICATION  = "MULTIPLICATION"
	EXPR_KIND_LESS_THAN       = "LESS_THAN"
	EXPR_LESS_THAN_OR_EQ      = "LESS_THAN_OR_EQ"
	EXPR_KIND_MORE_THAN       = "MORE_THAN"
	EXPR_KIND_MORE_THAN_OR_EQ = "MORE_THAN_OR_EQ"
	EXPR_KIND_EQ              = "EQUAL"
	EXPR_KIND_OR              = "OR"
	EXPR_KIND_AND             = "AND"
	EXPR_KIND_ID              = "IDENTIFIER"
	EXPR_KIND_LOOP            = "LOOP"
)

type Parser struct {
	cursor int
	tokens []tokenizer.Token
}

type Program struct {
	Stmts []*Stmt
}

type Stmt struct {
	Data     []tokenizer.Token
	Kind     string
	Expr     Expr
	Children []*Stmt
}

type Expr struct {
	Kind  string
	Token tokenizer.Token
	Exprs []Expr
}

func (p *Parser) parse() (*Program, error) {

	return p.parseProgram()
}

func (p *Parser) parseProgram() (*Program, error) {
	stmts, err := p.parseStmts()
	if err != nil {
		return &Program{
			Stmts: stmts,
		}, err
	}
	program := &Program{
		Stmts: stmts,
	}
	return program, nil
}

func NewParser(tokens []tokenizer.Token) Parser {
	return Parser{
		tokens: tokens,
	}
}

func (p *Parser) peek() (tokenizer.Token, bool) {
	if p.cursor >= len(p.tokens) {
		return tokenizer.Token{}, false
	}
	return p.tokens[p.cursor], true
}
func (p *Parser) peekAt(pos int) (tokenizer.Token, bool) {
	if (p.cursor+pos < 0) || (p.cursor+pos) >= len(p.tokens) {
		return tokenizer.Token{}, false
	}
	return p.tokens[p.cursor+pos], true
}

func (p *Parser) consume() {
	p.cursor++
}

func (p *Parser) rollback(position int) {
	p.cursor = position
}

func (p *Parser) spit() {
	p.cursor--
}

func (p *Parser) isEOF() bool {
	_, ok := p.peek()
	return !ok
}

func Parse(tokens []tokenizer.Token) (*Program, error) {
	p := NewParser(tokens)
	return p.parse()
}

func DebugPrintParseTree(stmts []*Stmt, indentLevel int) {
	for i := 0; i < len(stmts); i++ {
		stmt := stmts[i]
		fmt.Print(strings.Repeat(" ", indentLevel) + stmt.Kind)
		fmt.Print(" ")
		if stmt.Expr.Kind != "" {
			fmt.Print(ExprToSExpr(stmt.Expr))
			fmt.Print(" ")
		}
		if len(stmt.Data) > 0 {
			fmt.Print("(")
			for k := 0; k < len(stmt.Data); k++ {
				fmt.Print(stmt.Data[k].Content)
				if k < len(stmt.Data)-1 {
					fmt.Print(" ")
				}
			}
			fmt.Print(")")
		}
		fmt.Println()
		if len(stmt.Children) > 0 {
			DebugPrintParseTree(stmt.Children, indentLevel+1)
		}
	}
}
