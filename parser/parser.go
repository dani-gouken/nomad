package parser

import (
	"github.com/dani-gouken/nomad/tokenizer"
)

const (
	EXPR_KIND_CONSTANT       = "CONSTANT"
	EXPR_KIND_NOT            = "NOT"
	EXPR_KIND_NEGATIVE       = "NEGATIVE"
	EXPR_KIND_LEFT_INCREMENT = "LEFT_INCREMENT"
	EXPR_KIND_LEFT_DECREMENT = "LEFT_DECREMENT"
	EXPR_KIND_ADDITION       = "ADDITION"
	EXPR_KIND_SUBSTRACTION   = "SUBSTRACTION"
	EXPR_KIND_MULTIPLICATION = "MULTIPLICATION"
	EXPR_KIND_EQ             = "EQUAL"
	EXPR_KIND_ID             = "IDENTIFIER"
)

type Parser struct {
	cursor int
	tokens []tokenizer.Token
}

type Program struct {
	Stmts []Stmt
}

type Stmt struct {
	Kind string
	Expr Expr
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
