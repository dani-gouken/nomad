package parser

import (
	"errors"
	"fmt"

	"github.com/dani-gouken/nomad/tokenizer"
)

const (
	OPERATOR_PRECEDENCE_INVALID = iota
	OPERATOR_PRECEDENCE_MINIMUM
	OPERATOR_PRECEDENCE_REGULAR
	OPERATOR_PRECEDENCE_HIGH
	OPERATOR_PRECEDENCE_HIGHEST
)

func (p *Parser) parseExpr() (Expr, error) {
	primaryExpr, err := p.parsePrimaryExpr()
	if err != nil {
		return primaryExpr, err
	}
	binaryExpr, err := p.parseBinaryOperatorExpr(primaryExpr, OPERATOR_PRECEDENCE_MINIMUM)
	if err != nil {
		return primaryExpr, nil
	}
	return binaryExpr, err
}

func (p *Parser) parseUnaryOperatorExpr() (Expr, error) {
	t, ok := p.peek()
	if !ok {
		return Expr{}, fmt.Errorf("EOF")
	}
	switch t.Kind {
	case tokenizer.TOKEN_KIND_BANG:
		p.consume()
		expr, err := p.parseExpr()
		if err != nil {
			p.spit()
			return Expr{}, err
		}
		return Expr{
			Kind:  EXPR_KIND_NOT,
			Token: t,
			Exprs: []Expr{
				expr,
			},
		}, nil
	case tokenizer.TOKEN_KIND_MINUS:
		p.consume()
		expr, err := p.parsePrimaryExpr()
		if err != nil {
			p.spit()
			return Expr{}, err
		}
		return Expr{
			Kind:  EXPR_KIND_NEGATIVE,
			Token: t,
			Exprs: []Expr{
				expr,
			},
		}, nil

	case tokenizer.TOKEN_KIND_DB_MINUS:
		p.consume()
		expr, err := p.parseIdExpr()
		if err != nil {
			p.spit()
			return expr, err
		}
		return Expr{
			Kind:  EXPR_KIND_LEFT_DECREMENT,
			Token: t,
			Exprs: []Expr{
				expr,
			},
		}, nil
	case tokenizer.TOKEN_KIND_DB_PLUS:
		p.consume()
		expr, err := p.parseIdExpr()
		if err != nil {
			p.spit()
			return expr, err
		}
		return Expr{
			Kind:  EXPR_KIND_LEFT_INCREMENT,
			Token: t,
			Exprs: []Expr{
				expr,
			},
		}, nil
	}
	return Expr{}, errors.New("Parse error")
}

func isBinaryOperatorToken(t tokenizer.Token) bool {
	return getBinaryOperatorPrecedence(t) != OPERATOR_PRECEDENCE_INVALID
}
func getBinaryOperatorPrecedence(t tokenizer.Token) uint {
	switch t.Kind {
	case tokenizer.TOKEN_KIND_PLUS, tokenizer.TOKEN_KIND_MINUS, tokenizer.TOKEN_KIND_DB_EQUAL, tokenizer.TOKEN_KIND_SLASH:
		return OPERATOR_PRECEDENCE_REGULAR
	case tokenizer.TOKEN_KIND_STAR:
		return OPERATOR_PRECEDENCE_HIGH
	default:
		return OPERATOR_PRECEDENCE_INVALID
	}
}

func buildBinaryOpExpr(op tokenizer.Token, lhs Expr, rhs Expr) (Expr, error) {
	switch op.Kind {
	case tokenizer.TOKEN_KIND_PLUS:
		return Expr{
			Kind:  EXPR_KIND_ADDITION,
			Token: op,
			Exprs: []Expr{
				lhs, rhs,
			},
		}, nil
	case tokenizer.TOKEN_KIND_MINUS:
		return Expr{
			Kind:  EXPR_KIND_SUBSTRACTION,
			Token: op,
			Exprs: []Expr{
				lhs, rhs,
			},
		}, nil
	case tokenizer.TOKEN_KIND_SLASH:
		return Expr{
			Kind:  EXPR_KIND_DIVISION,
			Token: op,
			Exprs: []Expr{
				lhs, rhs,
			},
		}, nil
	case tokenizer.TOKEN_KIND_STAR:
		return Expr{
			Kind:  EXPR_KIND_MULTIPLICATION,
			Token: op,
			Exprs: []Expr{
				lhs, rhs,
			},
		}, nil
	case tokenizer.TOKEN_KIND_DB_EQUAL:
		return Expr{
			Kind:  EXPR_KIND_EQ,
			Token: op,
			Exprs: []Expr{
				lhs, rhs,
			},
		}, nil
	}
	return Expr{}, fmt.Errorf("expected binary operator, found %s", op.Kind)
}

func (p *Parser) parseBinaryOperatorExpr(lhs Expr, minPrecedence uint) (Expr, error) {
	lookahead, ok := p.peek()
	if !ok {
		return Expr{}, fmt.Errorf("EOF")
	}
	for isBinaryOperatorToken(lookahead) && getBinaryOperatorPrecedence(lookahead) >= minPrecedence {
		op := lookahead
		p.consume()
		rhs, err := p.parsePrimaryExpr()
		if err != nil {
			return Expr{}, fmt.Errorf("failed to parse operator %s. %s at position %d:%d:%d", op.Kind, op.Content, op.Loc.Line, op.Loc.Start, op.Loc.End)
		}
		lookahead, ok = p.peek()
		if !ok {
			return buildBinaryOpExpr(op, lhs, rhs)
		}
		opPrecedence := getBinaryOperatorPrecedence(op)
		for isBinaryOperatorToken(lookahead) && getBinaryOperatorPrecedence(lookahead) > opPrecedence {
			rhs, err = p.parseBinaryOperatorExpr(rhs, opPrecedence+1)
			if err != nil {
				return Expr{}, fmt.Errorf("failed to parse operator %s. %s at position %d:%d:%d", lookahead.Kind, lookahead.Content, lookahead.Loc.Line, lookahead.Loc.Start, lookahead.Loc.End)
			}
			lookahead, ok = p.peek()
			if !ok {
				return buildBinaryOpExpr(op, lhs, rhs)
			}
		}
		lhs, err = buildBinaryOpExpr(op, lhs, rhs)
		if err != nil {
			return lhs, err
		}

	}
	return lhs, nil
}

func (p *Parser) parseConstantExpr() (Expr, error) {
	t, ok := p.peek()
	if !ok {
		return Expr{}, fmt.Errorf("EOF")
	}
	switch t.Kind {
	case tokenizer.TOKEN_KIND_NUM_LIT, tokenizer.TOKEN_KIND_TRUE, tokenizer.TOKEN_KIND_FALSE, tokenizer.TOKEN_KIND_STRING_LIT:
		p.consume()
		return Expr{
			Kind:  EXPR_KIND_CONSTANT,
			Token: t,
		}, nil
	}
	return Expr{}, fmt.Errorf("could not parse constant")
}

func (p *Parser) parseIdExpr() (Expr, error) {
	t, ok := p.peek()
	if !ok {
		return Expr{}, fmt.Errorf("EOF")
	}
	switch t.Kind {
	case tokenizer.TOKEN_KIND_ID:
		p.consume()
		return Expr{
			Kind:  EXPR_KIND_ID,
			Token: t,
		}, nil
	}
	return Expr{}, fmt.Errorf("expected token identifier, %s: %s at position %d:%d:%d", t.Kind, t.Content, t.Loc.Line, t.Loc.Start, t.Loc.End)
}

func (p *Parser) parseBracketExpr() (Expr, error) {
	t, ok := p.peek()
	if !ok {
		return Expr{}, fmt.Errorf("EOF")
	}
	if t.Kind != tokenizer.TOKEN_KIND_LEFT_BRACKET {
		return Expr{}, fmt.Errorf("expected opening bracket")
	}
	p.consume()
	expr, err := p.parseExpr()
	if err != nil {
		return expr, err
	}
	t, ok = p.peek()
	if !ok {
		return Expr{}, errors.New("expected closing bracket, got EOF")
	}

	if t.Kind != tokenizer.TOKEN_KIND_RIGHT_BRACKET {
		return Expr{}, fmt.Errorf("expected closing bracket, got %s", t.Kind)
	}
	p.consume()
	return expr, nil
}

func (p *Parser) parsePrimaryExpr() (Expr, error) {
	expr, err := p.parseConstantExpr()
	if err == nil {
		return expr, err
	}
	expr, err = p.parseIdExpr()
	if err == nil {
		return expr, err
	}
	expr, err = p.parseBracketExpr()
	if err == nil {
		return expr, err
	}
	expr, err = p.parseUnaryOperatorExpr()
	if err == nil {
		return expr, err
	}
	token, ok := p.peek()
	if !ok {
		token, _ := p.peekAt(-1)
		return Expr{}, fmt.Errorf("failed to parse expression. unexpected end of file after token %s: %s at position %d:%d:%d", token.Kind, token.Content, token.Loc.Line, token.Loc.Start, token.Loc.End)
	}
	return Expr{}, fmt.Errorf("%s. Failed to parse token %s: %s at position %d:%d:%d", err.Error(), token.Kind, token.Content, token.Loc.Line, token.Loc.Start, token.Loc.End)
}

func ExprToSExpr(expr Expr) string {
	if !isBinaryOperatorToken(expr.Token) {
		return expr.Token.Content
	}
	sexpr := ""
	sexpr += "(" + expr.Token.Content
	for i := 0; i < len(expr.Exprs); i++ {
		sexpr += " " + ExprToSExpr(expr.Exprs[i])
	}
	sexpr += ")"
	return sexpr
}
