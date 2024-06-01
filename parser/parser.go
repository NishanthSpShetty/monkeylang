package parser

import (
	"fmt"

	"github.com/NishanthSpShetty/monkey/ast"
	"github.com/NishanthSpShetty/monkey/lexer"
	"github.com/NishanthSpShetty/monkey/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Erors() []string {
	return p.errors
}

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
		p.errors = append(p.errors, fmt.Sprintf("invalid token: %s", string(p.curToken.Type)))
	}

	return nil
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

	// expr

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	//	fmt.Printf("%+v\n", stmnt)
	return stmnt
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
