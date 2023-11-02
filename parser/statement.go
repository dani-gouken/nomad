package parser

import (
	"fmt"

	"github.com/dani-gouken/nomad/tokenizer"
)

const (
	STMT_KIND_IMPLICIT_RETURN = "IMPLICIT_RETURN"
	STMT_KIND_VAR_DECLARATION = "VARIABLE_DECLARATION"
)

func (p *Parser) parseStmts() ([]Stmt, error) {
	stmts := []Stmt{}
	for {
		_, ok := p.peek()

		if !ok {
			break
		}
		stmt, err := p.parseStmt()

		if err != nil {
			return stmts, err
		}
		err = p.terminateStmt()
		if err != nil {
			return stmts, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func (p *Parser) parseImplicitReturnStmt() (Stmt, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return Stmt{}, err
	}
	return Stmt{
		Kind: STMT_KIND_IMPLICIT_RETURN,
		Expr: expr,
	}, nil

}

func (p *Parser) parseStmt() (Stmt, error) {
	stmt, err := p.parseVariableDeclaration()
	if err == nil {
		return stmt, err
	}
	return p.parseImplicitReturnStmt()
}

func (p *Parser) parseVariableDeclaration() (Stmt, error) {
	stmt := Stmt{}
	err := p.expect(tokenizer.TOKEN_KIND_ID, "identifier (type)")
	if err != nil {
		return stmt, err
	}
	err = p.expectNext(tokenizer.TOKEN_KIND_ID, 1, "identifier (variable name)")
	if err != nil {
		return stmt, err
	}
	err = p.expectNext(tokenizer.TOKEN_KIND_EQUAL, 2, "equal (=)")
	if err != nil {
		return stmt, err
	}
	varType, _ := p.peek()
	p.consume()
	varName, _ := p.peek()
	p.consume()
	p.consume() // consume equal sign

	value, err := p.parseExpr()

	if err != nil {
		return stmt, err
	}
	return Stmt{
		Data: []tokenizer.Token{varType, varName},
		Kind: STMT_KIND_VAR_DECLARATION,
		Expr: value,
	}, nil
}

func (p *Parser) terminateStmt() error {
	err := p.expect(tokenizer.TOKEN_KIND_SEMI_COLON, "semi colon (;) OR New Line")
	if err == nil {
		p.consume()
		return nil
	}
	err = p.expect(tokenizer.TOKEN_KIND_NEW_LINE, "semi colon (;) OR New Line")
	if err == nil {
		p.consume()
		return nil
	}
	if p.isEOF() {
		return nil
	}
	return err
}

func (p *Parser) expect(kind string, expected string) error {
	token, ok := p.peek()
	if !ok {
		p, _ := p.peekAt(-1)
		return fmt.Errorf("non-terminated statement, expected %s or New Line, got EOF at line %d", expected, p.Loc.Line)
	}
	if token.Kind != kind {
		return fmt.Errorf("unexpected token. expected %s, got %s: %s at  at position %d:%d:%d", expected, token.Kind, token.Content, token.Loc.Line, token.Loc.Start, token.Loc.End)
	}
	return nil
}
func (p *Parser) expectNext(kind string, pos int, expected string) error {
	token, ok := p.peekAt(pos)
	if !ok {
		p, _ := p.peekAt(-1)
		return fmt.Errorf("non-terminated statement, expected %s or New Line, got EOF at line %d", expected, p.Loc.Line)
	}
	if token.Kind != kind {
		return fmt.Errorf("unexpected token. expected %s, got %s: %s at  at position %d:%d:%d", expected, token.Kind, token.Content, token.Loc.Line, token.Loc.Start, token.Loc.End)
	}
	return nil
}
