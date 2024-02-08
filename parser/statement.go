package parser

import (
	"fmt"

	nomadError "github.com/dani-gouken/nomad/errors"
	"github.com/dani-gouken/nomad/tokenizer"
)

const (
	STMT_KIND_IMPLICIT_RETURN = "IMPLICIT_RETURN"
	STMT_KIND_VAR_DECLARATION = "VARIABLE_DECLARATION"
	STMT_KIND_IF              = "IF"
	STMT_KIND_DEBUG_PRINT     = "DEBUG_PRINT"
	STMT_KIND_ELSE            = "ELSE"
	STMT_KIND_FOR             = "FOR"
	STMT_KIND_ELIF            = "ELIF"
	STMT_KIND_SCOPE           = "SCOPE"
	STMT_KIND_ASSIGNMENT      = "ASSIGNMENT"
)

func (p *Parser) parseStmts() ([]*Stmt, *nomadError.ParseError) {
	stmts := []*Stmt{}
	for {
		_, ok := p.peek()

		if !ok {
			break
		}
		newStmts, err := p.parseStmt()

		if err != nil {
			return stmts, err
		}
		if err != nil {
			return stmts, err
		}
		stmts = append(stmts, newStmts...)
	}
	return stmts, nil
}

func (p *Parser) parseImplicitReturnStmt() ([]*Stmt, *nomadError.ParseError) {
	expr, err := p.parseExpr()
	if err != nil {
		return []*Stmt{}, err
	}

	stmt := Stmt{
		Kind: STMT_KIND_IMPLICIT_RETURN,
		Expr: expr,
	}
	p.terminateStmt(stmt)
	return []*Stmt{&stmt}, nil

}

func (p *Parser) parseStmt() ([]*Stmt, *nomadError.ParseError) {
	parseFuncs := []func() ([]*Stmt, *nomadError.ParseError){
		p.parseAssignment,
		p.parsePrint,
		p.parseVariableDeclaration,
		p.parseIfStatement,
		p.parseForLoop,
	}
	initialParseCursor := p.cursor

	for i := 0; i < len(parseFuncs); i++ {
		p.rollback(initialParseCursor)
		stmt, err := parseFuncs[i]()
		if err == nil {
			return stmt, err
		}
		if err.ShouldCrash() {
			return stmt, err
		}
	}
	return p.parseImplicitReturnStmt()
}
func (p *Parser) parseBlock() ([]*Stmt, *nomadError.ParseError) {
	stmts := []*Stmt{}
	err := p.expectNF(tokenizer.TOKEN_KIND_LEFT_CURCLY, "left curly ({)")
	if err != nil {
		return stmts, err
	}
	previousToken, _ := p.peek()
	p.consume()

	for {
		token, ok := p.peek()
		if !ok {
			return nil, nomadError.FatalParseError(fmt.Sprintf("non-terminated block, %s expected", tokenizer.TOKEN_KIND_RIGHT_CURLY), previousToken)
		}

		if token.Kind == tokenizer.TOKEN_KIND_RIGHT_CURLY {
			p.consume()
			break
		}

		blockStmts, err := p.parseStmt()

		if err != nil {
			return stmts, err
		}

		stmts = append(stmts, blockStmts...)
	}
	return stmts, nil
}

func (p *Parser) parseFlowControlStatement(tokenKind string, statementKind string, hasExpr bool) ([]*Stmt, *nomadError.ParseError) {
	stmts := []*Stmt{}
	err := p.expectNF(tokenKind, fmt.Sprintf("identifier (%s)", statementKind))
	if err != nil {
		return stmts, err
	}
	token, _ := p.peek()
	p.consume()
	var stmt Stmt
	if hasExpr {
		token, _ := p.peek()
		value, err := p.parseExpr()
		if err != nil {
			return stmts, err
		}

		stmt = Stmt{
			Data: []tokenizer.Token{
				token,
			},
			Kind: statementKind,
			Expr: value,
		}
	} else {
		stmt = Stmt{
			Data: []tokenizer.Token{
				token,
			},
			Kind: statementKind,
		}

	}
	stmts = append(stmts, &stmt)
	blockStmts, err := p.parseBlock()
	stmt.Children = blockStmts
	if err != nil {
		return stmts, err
	}
	return stmts, nil
}

func (p *Parser) parseForLoop() ([]*Stmt, *nomadError.ParseError) {
	token, _ := p.peek()
	err := p.expectNF(tokenizer.TOKEN_KIND_FOR, "for (keyword)")
	if err != nil {
		return nil, err
	}
	p.consume()
	initStmt, err := p.parseVariableDeclaration()
	if err != nil {
		initStmt, err = p.parseAssignment()
		if err != nil {
			return initStmt, err
		}
	}
	testExpr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	err = p.expectF(tokenizer.TOKEN_KIND_SEMI_COLON, "separator (;)")
	if err != nil {
		return nil, err
	}
	p.consume()
	token, _ = p.peek()
	iterStmt, err := p.parseAssignment()
	if err != nil {
		iterStmt, err = p.parseImplicitReturnStmt()
		if err != nil {
			return nil, err
		}
	}
	operations, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	return append(initStmt, &Stmt{
		Data: []tokenizer.Token{
			token,
		},
		Expr:     testExpr,
		Kind:     STMT_KIND_FOR,
		Children: append(operations, iterStmt...),
	}), nil
}
func (p *Parser) parseIfStatement() ([]*Stmt, *nomadError.ParseError) {
	stmts, err := p.parseFlowControlStatement(tokenizer.TOKEN_KIND_IF, STMT_KIND_IF, true)
	if err != nil {
		return stmts, err
	}
	for {
		token, ok := p.peek()
		if !ok {
			break
		}
		kind := ""
		hasExpr := false
		if token.Kind == tokenizer.TOKEN_KIND_ELSE {
			kind = STMT_KIND_ELSE
		}
		if token.Kind == tokenizer.TOKEN_KIND_ELIF {
			kind = STMT_KIND_ELIF
			hasExpr = true
		}
		if kind == "" {
			break
		}

		nextStmts, err := p.parseFlowControlStatement(token.Kind, kind, hasExpr)

		if err != nil {
			return stmts, err
		}
		stmts = append(stmts, nextStmts...)

		if token.Kind == tokenizer.TOKEN_KIND_ELSE {
			break
		}
	}
	return stmts, nil
}

func (p *Parser) parseAssignment() ([]*Stmt, *nomadError.ParseError) {
	stmts := []*Stmt{}
	err := p.expectNF(tokenizer.TOKEN_KIND_ID, "identifier (variable name)")
	if err != nil {
		return stmts, err
	}
	err = p.expectNextNF(tokenizer.TOKEN_KIND_EQUAL, 1, "equal (=)")
	if err != nil {
		return stmts, err
	}
	varName, _ := p.peek()
	p.consume()
	p.consume() // consume equal sign

	value, err := p.parseExpr()

	if err != nil {
		return stmts, err
	}
	stmt := Stmt{
		Data: []tokenizer.Token{varName},
		Kind: STMT_KIND_ASSIGNMENT,
		Expr: value,
	}
	p.terminateStmt(stmt)

	return append(stmts, &stmt), nil
}

func (p *Parser) parseVariableDeclaration() ([]*Stmt, *nomadError.ParseError) {
	err := p.expectNF(tokenizer.TOKEN_KIND_ID, "identifier (type)")
	if err != nil {
		return []*Stmt{}, err
	}
	err = p.expectNextF(tokenizer.TOKEN_KIND_ID, 1, "identifier (variable name)")
	if err != nil {
		return []*Stmt{}, err
	}
	err = p.expectNextF(tokenizer.TOKEN_KIND_EQUAL, 2, "equal (=)")
	if err != nil {
		return []*Stmt{}, err
	}
	varType, _ := p.peek()
	p.consume()
	varName, _ := p.peek()
	p.consume()
	p.consume() // consume equal sign

	value, err := p.parseExpr()

	if err != nil {
		return []*Stmt{}, err
	}
	stmt := Stmt{
		Data: []tokenizer.Token{varType, varName},
		Kind: STMT_KIND_VAR_DECLARATION,
		Expr: value,
	}

	p.terminateStmt(stmt)

	return []*Stmt{&stmt}, nil
}

func (p *Parser) parsePrint() ([]*Stmt, *nomadError.ParseError) {
	err := p.expectNF(tokenizer.TOKEN_KIND_PRINT, "print (keyword)")
	if err != nil {
		return []*Stmt{}, err
	}
	p.consume()
	value, err := p.parseExpr()
	if err != nil {
		return []*Stmt{}, err
	}
	stmt := Stmt{
		Kind: STMT_KIND_DEBUG_PRINT,
		Expr: value,
	}

	p.terminateStmt(stmt)

	return []*Stmt{&stmt}, nil
}

func (p *Parser) terminateStmt(stmt Stmt) *nomadError.ParseError {
	if stmt.Kind == STMT_KIND_IF {
		return nil
	}
	err := p.expectF(tokenizer.TOKEN_KIND_SEMI_COLON, "semi colon (;) or new Line")
	if err == nil {
		p.consume()
		return nil
	}
	err = p.expectF(tokenizer.TOKEN_KIND_NEW_LINE, "semi colon (;) or new Line")
	if err == nil {
		p.consume()
		return nil
	}
	if p.isEOF() {
		return nil
	}
	return err
}

func (p *Parser) expect(kind string, expected string, fatal bool) *nomadError.ParseError {
	return p.expectNext(kind, 0, expected, fatal)
}

func (p *Parser) expectF(kind string, expected string) *nomadError.ParseError {
	return p.expect(kind, expected, true)
}
func (p *Parser) expectNF(kind string, expected string) *nomadError.ParseError {
	return p.expect(kind, expected, false)
}

func (p *Parser) expectNextF(kind string, pos int, expected string) *nomadError.ParseError {
	return p.expectNext(kind, pos, expected, true)
}
func (p *Parser) expectNextNF(kind string, pos int, expected string) *nomadError.ParseError {
	return p.expectNext(kind, pos, expected, false)
}

func (p *Parser) expectNext(kind string, pos int, expected string, fatal bool) *nomadError.ParseError {
	token, ok := p.peekAt(pos)
	if !ok {
		p, _ := p.peekAt(pos - 1)
		return nomadError.NewParseError(fmt.Sprintf("non-terminated statement, expected %s or new line, got EOF", expected), p, fatal)
	}
	if token.Kind != kind {
		return nomadError.NewParseError(fmt.Sprintf("unexpected token. expected %s, got %s: %s", expected, token.Kind, token.Content), token, fatal)
	}
	return nil
}
