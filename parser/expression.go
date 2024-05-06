package parser

import (
	"fmt"

	nomadError "github.com/dani-gouken/nomad/errors"
	"github.com/dani-gouken/nomad/tokenizer"
)

const (
	OPERATOR_PRECEDENCE_INVALID = iota
	OPERATOR_PRECEDENCE_MINIMUM
	OPERATOR_PRECEDENCE_LOW
	OPERATOR_PRECEDENCE_REGULAR
	OPERATOR_PRECEDENCE_HIGH
	OPERATOR_PRECEDENCE_HIGHEST
)

func (p *Parser) parseExpr() (Expr, *nomadError.ParseError) {
	expr, err := p.parseFuncExpr()
	if err == nil || err.ShouldCrash() {
		return expr, err
	}
	return p.parseBaseExpr()
}

func (p *Parser) parseFuncExpr() (Expr, *nomadError.ParseError) {
	beginning := p.cursor
	err := p.expectNF(tokenizer.TOKEN_KIND_LEFT_BRACKET, "opening bracket")
	if err != nil {
		return Expr{}, err
	}
	p.consume()

	paramListExpr, err := p.parseFuncParamListExpr()

	if err != nil {
		p.rollback(beginning)
		return Expr{}, err
	}

	err = p.expectF(tokenizer.TOKEN_KIND_RIGHT_BRACKET, "closing bracket")
	if err != nil {
		return Expr{}, err
	}
	p.consume()

	retTypeExpr, err := p.parseTypeExpr()
	if err != nil {
		p.rollback(beginning)
		return Expr{}, err
	}

	block, err := p.parseBlock()
	if err != nil {
		return Expr{}, nomadError.NewParseErrorFromMessage(err.Error(), true)
	}
	return Expr{
		Kind: EXPR_KIND_FUNC,
		Children: []Expr{
			paramListExpr,
			retTypeExpr,
		},
		Block: block,
	}, nil

}
func (p *Parser) parseFuncParamListExpr() (Expr, *nomadError.ParseError) {
	expr := Expr{
		Kind:     EXPR_KIND_FUNC_PARAM_LIST,
		Children: []Expr{},
	}

	for {
		t, _ := p.peek()
		if t.Kind == tokenizer.TOKEN_KIND_RIGHT_BRACKET {
			return expr, nil
		}
		param, err := p.parseFuncParamExpr()
		if err != nil {
			return expr, err
		}

		expr.Children = append(expr.Children, param)

		t, _ = p.peek()
		if t.Kind != tokenizer.TOKEN_KIND_COMMA && t.Kind != tokenizer.TOKEN_KIND_RIGHT_BRACKET {
			return expr, nomadError.FatalParseError("expected end of parameter list or closing bracket", t)
		}
		if t.Kind == tokenizer.TOKEN_KIND_COMMA {
			p.consume()
		}
	}
}

func (p *Parser) parseArgumentList() (Expr, *nomadError.ParseError) {
	args := []Expr{}
	var hasNamedArgument bool = false
	for {
		t, _ := p.peek()
		argExpr, err := p.parseArgument()
		if hasNamedArgument && argExpr.Kind == EXPR_KIND_FUNC_ARG {
			return Expr{}, nomadError.FatalParseError("positional arguments are not allowed after named argument", argExpr.Token)
		}
		if argExpr.Kind == EXPR_KIND_FUNC_NAMED_ARG && !hasNamedArgument {
			hasNamedArgument = true
		}
		if err != nil {
			return Expr{
				Kind:     EXPR_KIND_FUNC_ARG_LIST,
				Children: args,
			}, err
		}
		args = append(args, argExpr)
		t, _ = p.peek()

		if t.Kind != tokenizer.TOKEN_KIND_COMMA && t.Kind != tokenizer.TOKEN_KIND_RIGHT_BRACKET {
			return Expr{}, nomadError.FatalParseError(fmt.Sprintf("expected end of argument list or closing bracket, got %s", t.Kind), t)
		}

		if t.Kind == tokenizer.TOKEN_KIND_COMMA {
			p.consume()
		}

		if t.Kind == tokenizer.TOKEN_KIND_RIGHT_BRACKET {
			return Expr{
				Kind:     EXPR_KIND_FUNC_ARG_LIST,
				Children: args,
			}, nil
		}
	}
}

func (p *Parser) parseDirectArgument() (Expr, *nomadError.ParseError) {
	t, _ := p.peek()
	expr, err := p.parseExpr()
	if err != nil {
		return Expr{}, err
	}
	return Expr{
		Kind: EXPR_KIND_FUNC_ARG,
		Children: []Expr{
			expr,
		},
		Token: t,
	}, nil
}

func (p *Parser) parseNamedArgument() (Expr, *nomadError.ParseError) {
	err := p.expectNF(tokenizer.TOKEN_KIND_ID, "parameter name")
	if err != nil {
		return Expr{}, err
	}
	err = p.expectNextNF(tokenizer.TOKEN_KIND_COLON, 1, "colon")
	if err != nil {
		return Expr{}, err
	}

	name, _ := p.peek()
	p.consume()
	p.consume()

	valueExpr, err := p.parseExpr()
	if err != nil {
		return Expr{}, err
	}
	return Expr{
		Kind: EXPR_KIND_FUNC_NAMED_ARG,
		Children: []Expr{
			valueExpr,
		},
		Token: name,
	}, nil
}

func (p *Parser) parseArgument() (Expr, *nomadError.ParseError) {
	argExpr, err := p.parseNamedArgument()
	if err != nil {
		argExpr, err = p.parseDirectArgument()
	}
	return argExpr, err
}

func (p *Parser) parseFuncParamExpr() (Expr, *nomadError.ParseError) {
	typeExpr, err := p.parseTypeExpr()
	if err != nil {
		return Expr{}, err
	}

	err = p.expectNF(tokenizer.TOKEN_KIND_ID, "parameter name")
	if err != nil {
		return Expr{}, err
	}

	name, _ := p.peek()
	p.consume()

	t, _ := p.peek()
	expr := Expr{
		Kind:  EXPR_KIND_FUNC_PARAM,
		Token: name,
		Children: []Expr{
			typeExpr,
		},
	}
	if t.Kind != tokenizer.TOKEN_KIND_DB_COLON {
		return expr, nil
	}
	p.consume()
	defaultValueExpr, err := p.parseBasePrimaryExpr()
	if err != nil {
		return expr, nomadError.NewParseErrorFromMessage(err.Error(), true)
	}

	expr.Children = append(expr.Children, defaultValueExpr)
	if err != nil {
		return expr, err
	}

	return expr, nil
}

func (p *Parser) parseBaseExpr() (Expr, *nomadError.ParseError) {
	primaryExpr, err := p.parsePrimaryExpr()
	if err != nil {
		return primaryExpr, err
	}
	return p.parseBinaryOperatorExpr(primaryExpr, OPERATOR_PRECEDENCE_MINIMUM)
}

func (p *Parser) parseUnaryOperatorExpr() (Expr, *nomadError.ParseError) {
	t, ok := p.peek()
	if !ok {
		return Expr{}, nomadError.FatalParseError("EOF", tokenizer.Token{})
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
			Children: []Expr{
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
			Children: []Expr{
				expr,
			},
		}, nil
	case tokenizer.TOKEN_KIND_LEN:
		p.consume()
		expr, err := p.parsePrimaryExpr()
		if err != nil {
			p.spit()
			return Expr{}, err
		}
		return Expr{
			Kind:  EXPR_KIND_LEN,
			Token: t,
			Children: []Expr{
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
			Children: []Expr{
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
			Children: []Expr{
				expr,
			},
		}, nil
	case tokenizer.TOKEN_KIND_ID:
		op, ok := p.peekAt(1)
		if !ok {
			break
		}
		var exprKind string
		if op.Kind == tokenizer.TOKEN_KIND_DB_PLUS {
			exprKind = EXPR_KIND_RIGHT_INCREMENT
		}
		if op.Kind == tokenizer.TOKEN_KIND_DB_MINUS {
			exprKind = EXPR_KIND_RIGHT_DECREMENT
		}
		if exprKind == "" {
			break
		}
		expr, err := p.parseIdExpr()
		if err != nil {
			p.spit()
			return expr, err
		}
		p.consume()
		return Expr{
			Kind:  exprKind,
			Token: t,
			Children: []Expr{
				expr,
			},
		}, nil
	}
	return Expr{}, nomadError.FatalParseError("failed to parse unary operator", t)
}

func isBinaryOperatorToken(t tokenizer.Token) bool {
	return getBinaryOperatorPrecedence(t) != OPERATOR_PRECEDENCE_INVALID
}
func getBinaryOperatorPrecedence(t tokenizer.Token) uint {
	switch t.Kind {
	case tokenizer.TOKEN_KIND_EQUAL:
		return OPERATOR_PRECEDENCE_LOW
	case tokenizer.TOKEN_KIND_PLUS,
		tokenizer.TOKEN_KIND_MINUS,
		tokenizer.TOKEN_KIND_SLASH,
		tokenizer.TOKEN_KIND_INFERIOR_SIGN,
		tokenizer.TOKEN_KIND_INFERIOR_OR_EQ_SIGN,
		tokenizer.TOKEN_KIND_SUPERIOR_SIGN,
		tokenizer.TOKEN_KIND_SUPERIOR_OR_EQ_SIGN,
		tokenizer.TOKEN_KIND_AND,
		tokenizer.TOKEN_KIND_BAR:
		return OPERATOR_PRECEDENCE_REGULAR
	case tokenizer.TOKEN_KIND_STAR:
		return OPERATOR_PRECEDENCE_HIGH
	default:
		return OPERATOR_PRECEDENCE_INVALID
	}
}

func buildBinaryOpExpr(op tokenizer.Token, lhs Expr, rhs Expr) (Expr, *nomadError.ParseError) {
	switch op.Kind {
	case tokenizer.TOKEN_KIND_PLUS:
		return Expr{
			Kind:  EXPR_KIND_ADDITION,
			Token: op,
			Children: []Expr{
				lhs, rhs,
			},
		}, nil
	case tokenizer.TOKEN_KIND_MINUS:
		return Expr{
			Kind:  EXPR_KIND_SUBSTRACTION,
			Token: op,
			Children: []Expr{
				lhs, rhs,
			},
		}, nil
	case tokenizer.TOKEN_KIND_SLASH:
		return Expr{
			Kind:  EXPR_KIND_DIVISION,
			Token: op,
			Children: []Expr{
				lhs, rhs,
			},
		}, nil
	case tokenizer.TOKEN_KIND_STAR:
		return Expr{
			Kind:  EXPR_KIND_MULTIPLICATION,
			Token: op,
			Children: []Expr{
				lhs, rhs,
			},
		}, nil
	case tokenizer.TOKEN_KIND_INFERIOR_SIGN:
		return Expr{
			Kind:  EXPR_KIND_LESS_THAN,
			Token: op,
			Children: []Expr{
				lhs, rhs,
			},
		}, nil
	case tokenizer.TOKEN_KIND_SUPERIOR_SIGN:
		return Expr{
			Kind:  EXPR_KIND_MORE_THAN,
			Token: op,
			Children: []Expr{
				lhs, rhs,
			},
		}, nil
	case tokenizer.TOKEN_KIND_INFERIOR_OR_EQ_SIGN:
		return Expr{
			Kind:  EXPR_LESS_THAN_OR_EQ,
			Token: op,
			Children: []Expr{
				lhs, rhs,
			},
		}, nil
	case tokenizer.TOKEN_KIND_SUPERIOR_OR_EQ_SIGN:
		return Expr{
			Kind:  EXPR_KIND_MORE_THAN_OR_EQ,
			Token: op,
			Children: []Expr{
				lhs, rhs,
			},
		}, nil
	case tokenizer.TOKEN_KIND_EQUAL:
		return Expr{
			Kind:  EXPR_KIND_EQ,
			Token: op,
			Children: []Expr{
				lhs, rhs,
			},
		}, nil
	case tokenizer.TOKEN_KIND_AND:
		return Expr{
			Kind:  EXPR_KIND_AND,
			Token: op,
			Children: []Expr{
				lhs, rhs,
			},
		}, nil
	case tokenizer.TOKEN_KIND_BAR:
		return Expr{
			Kind:  EXPR_KIND_OR,
			Token: op,
			Children: []Expr{
				lhs, rhs,
			},
		}, nil
	}
	return Expr{}, nomadError.FatalParseError(fmt.Sprintf("unknown binary operator %s", op.Kind), op)
}

func (p *Parser) parseBinaryOperatorExpr(lhs Expr, minPrecedence uint) (Expr, *nomadError.ParseError) {
	lookahead, _ := p.peek()
	var ok bool
	for isBinaryOperatorToken(lookahead) && getBinaryOperatorPrecedence(lookahead) >= minPrecedence {
		op := lookahead
		p.consume()
		rhs, err := p.parsePrimaryExpr()
		if err != nil {
			return Expr{}, nomadError.FatalParseError(fmt.Sprintf("failed to parse operator %s: %s", op.Kind, err.Error()), op)
		}
		lookahead, ok = p.peek()
		if !ok {
			return buildBinaryOpExpr(op, lhs, rhs)
		}
		opPrecedence := getBinaryOperatorPrecedence(op)
		for isBinaryOperatorToken(lookahead) && getBinaryOperatorPrecedence(lookahead) > opPrecedence {
			rhs, err = p.parseBinaryOperatorExpr(rhs, opPrecedence+1)
			if err != nil {
				return Expr{}, nomadError.FatalParseError(fmt.Sprintf("failed to parse operator %s: %s", lookahead.Kind, err.Error()), lookahead)
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

func (p *Parser) parseConstantExpr() (Expr, *nomadError.ParseError) {
	t, _ := p.peek()
	switch t.Kind {
	case tokenizer.TOKEN_KIND_NUM_LIT, tokenizer.TOKEN_KIND_TRUE, tokenizer.TOKEN_KIND_FALSE, tokenizer.TOKEN_KIND_STRING_LIT:
		p.consume()
		return Expr{
			Kind:  EXPR_KIND_CONSTANT,
			Token: t,
		}, nil
	}
	return Expr{}, nomadError.FatalParseError("could not parse constant", t)
}

func (p *Parser) parseIdExpr() (Expr, *nomadError.ParseError) {
	t, ok := p.peek()
	if !ok {
		return Expr{}, nomadError.FatalParseError("EOF", tokenizer.Token{})
	}
	switch t.Kind {
	case tokenizer.TOKEN_KIND_ID:
		p.consume()
		return Expr{
			Kind:  EXPR_KIND_ID,
			Token: t,
		}, nil
	}
	return Expr{}, nomadError.FatalParseError(fmt.Sprintf("expected token identifier, %s: %s", t.Kind, t.Content), t)
}

func (p *Parser) parseBracketExpr(parseFunc func() (Expr, *nomadError.ParseError)) (Expr, *nomadError.ParseError) {
	t, _ := p.peek()
	if t.Kind != tokenizer.TOKEN_KIND_LEFT_BRACKET {
		return Expr{}, nomadError.FatalParseError("expected opening bracket", t)
	}
	p.consume()
	expr, err := parseFunc()
	if err != nil {
		return expr, err
	}
	t, ok := p.peek()
	if !ok {
		return Expr{}, nomadError.FatalParseError("expected opening bracket, got EOF", t)
	}

	if t.Kind != tokenizer.TOKEN_KIND_RIGHT_BRACKET {
		return Expr{}, nomadError.FatalParseError(fmt.Sprintf("expected closing bracket, got %s", t.Kind), t)
	}
	p.consume()
	return expr, nil
}
func (p *Parser) parseArrayExpr() (Expr, *nomadError.ParseError) {
	pos := p.cursor
	err := p.expectNF(tokenizer.TOKEN_KIND_LEFT_SQUARE_BRACKET, "opening square bracket ([)")
	if err != nil {
		return Expr{}, err
	}
	p.consume()
	arrayTypeExpr, err := p.parseTypeExpr()
	if err != nil {
		p.rollback(pos)
		return Expr{}, err
	}

	err = p.expectF(tokenizer.TOKEN_KIND_RIGHT_SQUARE_BRACKET, "closing square bracket (])")
	if err != nil {
		return Expr{}, err
	}
	p.consume()

	err = p.expectF(tokenizer.TOKEN_KIND_LEFT_CURCLY, "opening bracket ({)")
	if err != nil {
		return Expr{}, err
	}
	p.consume()

	itemsExpr, err := p.parseExprList(tokenizer.TOKEN_KIND_RIGHT_CURLY)
	if err != nil {
		return Expr{}, err
	}

	err = p.expectF(tokenizer.TOKEN_KIND_RIGHT_CURLY, "opening bracket (})")
	p.consume()
	if err != nil {
		return Expr{}, err
	}

	return Expr{
		Kind:  EXPR_KIND_ARRAY,
		Token: arrayTypeExpr.Token,
		Children: []Expr{
			arrayTypeExpr,
			itemsExpr,
		},
	}, nil

}

func (p *Parser) parseExprList(endTokenKind string) (Expr, *nomadError.ParseError) {
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
		expr, err := p.parseExpr()

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

func (p *Parser) parsePrimaryExpr() (Expr, *nomadError.ParseError) {
	primaryExpr, err := p.parseBasePrimaryExpr()
	if err != nil {
		return primaryExpr, err
	}

	primaryExpr, err = p.parseAccessExpression(primaryExpr)
	if err != nil && err.ShouldCrash() {
		return primaryExpr, err
	}

	return primaryExpr, nil
}

func (p *Parser) parseAccessExpression(baseExpr Expr) (Expr, *nomadError.ParseError) {
	baseExpr, err := p.parseArrayAccess(baseExpr)
	if err != nil && err.ShouldCrash() {
		return baseExpr, err
	}
	baseExpr, err = p.parseObjectAccess(baseExpr)
	if err != nil && err.ShouldCrash() {
		return baseExpr, err
	}
	baseExpr, err = p.parseFuncCall(baseExpr)
	if err != nil && err.ShouldCrash() {
		return baseExpr, err
	}
	return p.parseObjectDefaultAccess(&baseExpr)

}

func (p *Parser) parseBasePrimaryExpr() (Expr, *nomadError.ParseError) {
	expr, err := p.parseConstantExpr()
	if err == nil {
		return expr, err
	}
	expr, err = p.parseUnaryOperatorExpr()
	if err == nil {
		return expr, err
	}
	expr, err = p.parseArrayExpr()
	if err == nil {
		return expr, err
	}

	expr, err = p.parseIdExpr()
	if err == nil {
		return expr, err
	}

	expr, err = p.parseBracketExpr(p.parseExpr)
	if err == nil {
		return expr, err
	}

	expr, err = p.parseObjectExpr()
	if err == nil {
		return expr, err
	}

	token, ok := p.peek()
	if !ok {
		token, _ := p.peekAt(-1)
		return Expr{}, nomadError.FatalParseError(fmt.Sprintf("unexpected end of file after token %s", token.Kind), token)
	}
	return Expr{}, nomadError.FatalParseError(fmt.Sprintf("failed to parse expression: %s", token.Kind), token)
}

func ExprToSExpr(expr Expr) string {
	if !isBinaryOperatorToken(expr.Token) {
		return expr.Token.Content
	}
	sexpr := ""
	sexpr += "(" + expr.Token.Content
	for i := 0; i < len(expr.Children); i++ {
		sexpr += " " + ExprToSExpr(expr.Children[i])
	}
	sexpr += ")"
	return sexpr
}

func (p *Parser) parseObjectTypeField() (Expr, *nomadError.ParseError) {
	pos := p.cursor
	typeExpr, err := p.parseTypeExpr()
	if err != nil {
		return Expr{}, err
	}
	err = p.expectNF(tokenizer.TOKEN_KIND_ID, "identifier (variable name)")
	if err != nil {
		p.rollback(pos)
		return Expr{}, err
	}
	varName, _ := p.peek()
	p.consume()
	err = p.expectNF(tokenizer.TOKEN_KIND_DB_COLON, "double colon (::)")
	if err != nil {
		return Expr{}, err
	}
	p.consume() // consume equal sign
	value, err := p.parseExpr()
	if err != nil {
		return Expr{}, err
	}
	return Expr{
		Kind:  EXPR_KIND_TYPE_OBJ_FIELD,
		Token: varName,
		Children: []Expr{
			typeExpr,
			value,
		},
	}, nil
}

func (p *Parser) parseObjectField() (Expr, *nomadError.ParseError) {
	pos := p.cursor
	err := p.expectNF(tokenizer.TOKEN_KIND_ID, "identifier (variable name)")
	if err != nil {
		p.rollback(pos)
		return Expr{}, err
	}
	varName, _ := p.peek()
	p.consume()
	err = p.expectNF(tokenizer.TOKEN_KIND_DB_COLON, "double colon (::)")
	if err != nil {
		p.rollback(pos)
		return Expr{}, err
	}
	p.consume() // consume equal sign
	value, err := p.parseExpr()
	if err != nil {
		return Expr{}, err
	}
	return Expr{
		Kind:  EXPR_KIND_OBJ_FIELD,
		Token: varName,
		Children: []Expr{
			value,
		},
	}, nil
}

func (p *Parser) parseObjectExpr() (Expr, *nomadError.ParseError) {
	err := p.expectNF(tokenizer.TOKEN_KIND_NEW, "new")
	if err != nil {
		return Expr{}, err
	}
	p.consume()

	err = p.expectF(tokenizer.TOKEN_KIND_ID, "type")
	if err != nil {
		return Expr{}, err
	}

	err = p.expectNextF(tokenizer.TOKEN_KIND_LEFT_CURCLY, 1, "opening curly bracket ({)")
	if err != nil {
		return Expr{}, err
	}

	t, _ := p.peek()
	p.consume()
	p.consume()
	p.cleanupNewLines()

	declarations := []Expr{}

	for {
		fieldDeclr, err := p.parseObjectField()
		if err != nil {
			break
		}
		declarations = append(declarations, fieldDeclr)
		sep, _ := p.peek()
		if sep.Kind == tokenizer.TOKEN_KIND_COMMA {
			p.consume()
		}
		p.cleanupNewLines()
	}
	err = p.expectF(tokenizer.TOKEN_KIND_RIGHT_CURLY, "closing curly bracket (})")
	p.consume()
	if err != nil {
		return Expr{}, err
	}
	return Expr{
		Kind:     EXPR_KIND_OBJ,
		Children: declarations,
		Token:    t,
	}, nil
}

func (p *Parser) parseObjectAccess(baseExpr Expr) (Expr, *nomadError.ParseError) {
	err := p.expectNF(tokenizer.TOKEN_KIND_DOT, "dot")
	if err != nil {
		return baseExpr, nil
	}

	err = p.expectNextF(tokenizer.TOKEN_KIND_ID, 1, "identifier")
	if err != nil {
		return baseExpr, nil
	}

	p.consume()
	field, _ := p.peek()
	p.consume()

	return p.parseAccessExpression(Expr{
		Kind:     EXPR_KIND_OBJ_ACCESS,
		Children: []Expr{baseExpr},
		Token:    field,
	})
}

func (p *Parser) parseFuncCall(baseExpr Expr) (Expr, *nomadError.ParseError) {
	begin := p.cursor
	t, _ := p.peek()
	err := p.expectNF(tokenizer.TOKEN_KIND_LEFT_BRACKET, "function call")
	if err != nil {
		return baseExpr, nil
	}
	p.consume()
	argList, err := p.parseArgumentList()
	if err != nil {
		p.rollback(begin)
		return baseExpr, err
	}
	err = p.expectF(tokenizer.TOKEN_KIND_RIGHT_BRACKET, "function call")
	if err != nil {
		return baseExpr, err
	}
	p.consume()

	return Expr{
		Kind:     EXPR_KIND_FUNC_CALL,
		Token:    t,
		Children: []Expr{baseExpr, argList},
	}, nil

}

func (p *Parser) parseObjectDefaultAccess(baseExpr *Expr) (Expr, *nomadError.ParseError) {
	if baseExpr.Kind != EXPR_KIND_ID {
		return *baseExpr, nil
	}

	err := p.expectNF(tokenizer.TOKEN_KIND_HASH, "hash")
	if err != nil {
		return *baseExpr, nil
	}

	err = p.expectNextF(tokenizer.TOKEN_KIND_ID, 1, "identifier")
	if err != nil {
		return *baseExpr, nil
	}

	p.consume()
	field, _ := p.peek()
	p.consume()
	baseExpr.Kind = EXPR_KIND_TYPE

	return p.parseAccessExpression(Expr{
		Kind:     EXPR_KIND_OBJ_DEFAULT_ACCESS,
		Children: []Expr{*baseExpr},
		Token:    field,
	})
}

func (p *Parser) parseArrayAccess(baseExpr Expr) (Expr, *nomadError.ParseError) {
	err := p.expectNF(tokenizer.TOKEN_KIND_LEFT_SQUARE_BRACKET, "opening bracket ([)")
	if err != nil {
		return baseExpr, nil
	}

	err = p.expectNextF(tokenizer.TOKEN_KIND_NUM_LIT, 1, "index")
	if err != nil {
		err = p.expectNextF(tokenizer.TOKEN_KIND_ID, 1, "identifier")
		if err != nil {
			return baseExpr, nil
		}
	}
	err = p.expectNextNF(tokenizer.TOKEN_KIND_RIGHT_SQUARE_BRACKET, 2, "closing bracket (])")
	if err != nil {
		return Expr{}, err
	}
	p.consume()
	index, _ := p.peek()
	p.consume()
	p.consume()

	return p.parseArrayAccess(Expr{
		Kind:     EXPR_KIND_ARRAY_ACCESS,
		Children: []Expr{baseExpr},
		Token:    index,
	})
}
