package parser

import (
	"fmt"

	"github.com/NishanthSpShetty/monkey/ast"
	"github.com/NishanthSpShetty/monkey/token"
)

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixParser := p.prefixParserFns[p.curToken.Type]
	if prefixParser == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExpr := prefixParser()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infixParser, ok := p.infixParserFns[p.peekToken.Type]
		if !ok {
			return leftExpr
		}
		p.nextToken()

		leftExpr = infixParser(leftExpr)

	}
	return leftExpr
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()

	exp.Right = p.parseExpression(precedence)

	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}
