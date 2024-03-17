package parser

import (
	nomadError "github.com/dani-gouken/nomad/errors"
	"github.com/dani-gouken/nomad/tokenizer"
)

func (p *Parser) parseTypeExpr() (Expr, *nomadError.ParseError) {
	return p.parseTypeAliasExpr()
}

func (p *Parser) parseTypeAliasExpr() (Expr, *nomadError.ParseError) {
	err := p.expectF(tokenizer.TOKEN_KIND_ID, "identifier (type name)")

	if err != nil {
		return Expr{}, err
	}
	typeName, _ := p.peek()
	p.consume()
	return Expr{
		Kind:  EXPR_KIND_TYPE,
		Token: typeName,
	}, nil
}
