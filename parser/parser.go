package parser

import (
	"fmt"
	"strconv"

	"github.com/NishanthSpShetty/monkey/ast"
	"github.com/NishanthSpShetty/monkey/lexer"
	"github.com/NishanthSpShetty/monkey/token"
)

type (
	prefixParserFn func() ast.Expression
	infixParserFn  func(ast.Expression) ast.Expression

	Parser struct {
		l         *lexer.Lexer
		curToken  token.Token
		peekToken token.Token
		errors    []string

		prefixParserFns map[token.TokenType]prefixParserFn
		infixParserFns  map[token.TokenType]infixParserFn
	}
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:               l,
		errors:          []string{},
		prefixParserFns: make(map[token.TokenType]prefixParserFn),
		infixParserFns:  make(map[token.TokenType]infixParserFn),
	}

	p.registerPrefixParser(token.IDENT, p.parseIdentifier)
	p.registerPrefixParser(token.INT, p.parseIntegerLiteral)
	p.registerPrefixParser(token.BANG, p.parsePrefixExpression)
	p.registerPrefixParser(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixParser(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefixParser(token.TRUE, p.parseBoolean)
	p.registerPrefixParser(token.FALSE, p.parseBoolean)
	p.registerPrefixParser(token.IF, p.parseIfExpression)
	p.registerPrefixParser(token.FUNCTION, p.parseFunctionLiteral)

	// infixParserFn
	p.registerInfixParser(token.PLUS, p.parseInfixExpression)
	p.registerInfixParser(token.MINUS, p.parseInfixExpression)
	p.registerInfixParser(token.SLASH, p.parseInfixExpression)
	p.registerInfixParser(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixParser(token.EQ, p.parseInfixExpression)
	p.registerInfixParser(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixParser(token.LT, p.parseInfixExpression)
	p.registerInfixParser(token.GT, p.parseInfixExpression)
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefixParser(t token.TokenType, fn prefixParserFn) {
	p.prefixParserFns[t] = fn
}

func (p *Parser) registerInfixParser(t token.TokenType, fn infixParserFn) {
	p.infixParserFns[t] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Erors() []string {
	return p.errors
}

// ParseProgram parse whole program and turns into statements
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for p.curToken.Type != token.EOF {
		stmnt := p.parseStatement()
		if stmnt != nil {
			program.Statements = append(program.Statements, stmnt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
		// p.errors = append(p.errors, fmt.Sprintf("invalid token: %s", string(p.curToken.Type)))
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// let a = 100;
	stmnt := &ast.LetStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	// curToken would have been advanced by prev call to expectPeek
	stmnt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	// expr
	stmnt.Value = p.parseExpression(LOWEST)

	//for !p.curTokenIs(token.SEMICOLON) {
	//	p.nextToken()
	//}
	//	fmt.Printf("%+v\n", stmnt)
	return stmnt
}

func (p *Parser) peekPrecedence() int {
	if pr, ok := precedences[p.peekToken.Type]; ok {
		return pr
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if pr, ok := precedences[p.curToken.Type]; ok {
		return pr
	}
	return LOWEST
}

func (p *Parser) peekErrors(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// expectPeek checks if the next token is what we expected,
// if yes, move to that by calling nextToken
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekErrors(t)
	return false
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	st := &ast.ReturnStatement{
		Token:       p.curToken,
		ReturnValue: nil,
	}
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	//	fmt.Printf("%+v\n", stmnt)
	return st
}

// / --- expr parsers -- pratt parser

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	stmnt := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	stmnt.Value = value

	return stmnt
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{
		Token: p.curToken,
	}

	// we are at `if`, expect "("
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	// expect ")" and skip to ")"
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// skip ) and move to {
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	// } else {
	if p.peekTokenIs(token.ELSE) {
		// move to ELSE
		p.nextToken()

		// move to {
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		exp.Alternative = p.parseBlockStatement()
	}
	return exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	bs := &ast.BlockStatement{
		Token:      p.curToken,
		Statements: []ast.Statement{},
	}

	// skip the {
	p.nextToken()

	// until we hit the } or end of statements
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			bs.Statements = append(bs.Statements, stmt)
		}
		// move over
		p.nextToken()
	}

	return bs
}
