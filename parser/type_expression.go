package parser

import (
	nomadError "github.com/dani-gouken/nomad/errors"
	"github.com/dani-gouken/nomad/tokenizer"
)

func (p *Parser) parseTypeExpr() (Expr, *nomadError.ParseError) {

	t, _ := p.peek()

	if t.Kind == tokenizer.TOKEN_KIND_LEFT_BRACKET {
		return p.parseBracketExpr(p.parseTypeExpr)
	}

	if t.Kind == tokenizer.TOKEN_KIND_ID {
		p.consume()
		return Expr{
			Kind:  EXPR_KIND_TYPE,
			Token: t,
		}, nil
	}
	if t.Kind == tokenizer.TOKEN_KIND_LEFT_SQUARE_BRACKET {
		expr, err := p.parseArrayTypeExpr()
		if err != nil {
			return Expr{}, err
		}
		return expr, nil
	}

	if t.Kind == tokenizer.TOKEN_KIND_LEFT_CURCLY {
		expr, err := p.parseObjectTypeExpr()
		if err != nil {
			return Expr{}, err
		}
		return expr, nil
	}

	if t.Kind == tokenizer.TOKEN_KIND_FUNC {
		expr, err := p.parseFuncTypeExpr()
		if err != nil {
			return Expr{}, err
		}
		return expr, nil
	}
	return Expr{}, nomadError.NonFatalParseError("could not parse type expression", t)
}

func (p *Parser) parseArrayTypeExpr() (Expr, *nomadError.ParseError) {
	err := p.expectNF(tokenizer.TOKEN_KIND_LEFT_SQUARE_BRACKET, "opening bracket ([)")
	if err != nil {
		return Expr{}, err
	}
	t, _ := p.peek()
	p.consume()
	typeExpr, err := p.parseTypeExpr()
	if err != nil {
		p.spit()
		return Expr{}, err
	}
	err = p.expectF(tokenizer.TOKEN_KIND_RIGHT_SQUARE_BRACKET, "closing bracket (])")
	p.consume()
	if err != nil {
		return Expr{}, err
	}
	return Expr{
		Kind:     EXPR_KIND_TYPE_ARRAY,
		Children: []Expr{typeExpr},
		Token:    t,
	}, nil
}

func (p *Parser) parseFuncTypeExpr() (Expr, *nomadError.ParseError) {
	err := p.expectNF(tokenizer.TOKEN_KIND_FUNC, "keyword (func)")
	if err != nil {
		return Expr{}, nil
	}

	t, _ := p.peek()
	p.consume()
	funcExpr := Expr{
		Kind:  EXPR_KIND_TYPE_FUNC,
		Token: t,
	}

	t, _ = p.peek()
	if t.Kind != tokenizer.TOKEN_KIND_LEFT_BRACKET {
		return funcExpr, nil
	}
	p.consume()

	paramTypeList, err := p.parseTypeExprList(tokenizer.TOKEN_KIND_RIGHT_BRACKET)
	if err != nil {
		return funcExpr, err
	}
	p.consume()

	err = p.expectF(tokenizer.TOKEN_KIND_ARROW, "Arrow (->)")
	if err != nil {
		return funcExpr, err
	}
	p.consume()

	returnType, err := p.parseTypeExpr()
	if err != nil {
		return funcExpr, err
	}
	funcExpr.Children = []Expr{
		paramTypeList,
		returnType,
	}

	return funcExpr, nil
}

func (p *Parser) parseObjectTypeExpr() (Expr, *nomadError.ParseError) {
	err := p.expectNF(tokenizer.TOKEN_KIND_LEFT_CURCLY, "opening curly bracket ({)")
	if err != nil {
		return Expr{}, err
	}

	t, _ := p.peek()
	p.consume()
	p.cleanupNewLines()

	declarations := []Expr{}

	for {
		fieldDeclr, err := p.parseObjectTypeField()
		if err != nil {
			break
		}
		declarations = append(declarations, fieldDeclr)
		p.cleanupNewLines()
	}
	err = p.expectF(tokenizer.TOKEN_KIND_RIGHT_CURLY, "closing curly bracket (})")
	p.consume()
	if err != nil {
		return Expr{}, err
	}
	return Expr{
		Kind:     EXPR_KIND_TYPE_OBJ,
		Children: declarations,
		Token:    t,
	}, nil
}

func (p *Parser) parseTypeExprList(endTokenKind string) (Expr, *nomadError.ParseError) {
	list := []Expr{}
	for {
		token, _ := p.peek()
		if token.Kind == tokenizer.TOKEN_KIND_NEW_LINE {
			p.consume()
			continue
		}
		if token.Kind == endTokenKind {
			return Expr{
				Kind:     EXPR_KIND_ANONYMOUS,
				Children: list,
			}, nil
		}
		expr, err := p.parseTypeExpr()

		if err != nil {
			return Expr{
				Kind:     EXPR_KIND_ANONYMOUS,
				Children: list,
			}, err
		}

		list = append(list, expr)
		token, _ = p.peek()

		if token.Kind == tokenizer.TOKEN_KIND_COMMA {
			p.consume()
		}
	}
}
