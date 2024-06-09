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

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.curToken}

	// we are at fn, move to (
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fn.Parameters = p.parseFunctionParameters()

	// we are at ), move to {
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fn.Body = p.parseBlockStatement()

	return fn
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	ids := []*ast.Identifier{}

	// if no args present
	if p.peekTokenIs(token.RPAREN) {
		// we are in (, move to )
		p.nextToken()
		return ids
	}
	// (a,b,c)
	// we are in (, move to first arg token
	p.nextToken()
	ids = append(ids, &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	})

	// loop till we have comms in next token
	// a,b,c)
	for p.peekTokenIs(token.COMMA) {
		// skip arg ie. a on first iteration
		p.nextToken()
		// skip ,
		p.nextToken()
		ids = append(ids, &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		})
	}
	// we should see )
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return ids
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	args := []ast.Expression{}

	// add ()
	if p.peekTokenIs(end) {
		p.nextToken()
		return args
	}

	// move to (
	// add (a, x+y, ...)
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	// A,B
	for p.peekTokenIs(token.COMMA) {

		// move to B
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	// we should see )
	if !p.expectPeek(end) {
		return nil
	}
	return args
}
