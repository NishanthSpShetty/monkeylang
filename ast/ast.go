package ast

import (
	"fmt"

	"github.com/NishanthSpShetty/monkey/token"
)

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

// TokenLiteral implements Node.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatement) String() string {
	return fmt.Sprintf("[Token: %s, Name: %s, Value: %s]", ls.Token, ls.Name, ls.Value)
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	return fmt.Sprintf("<Token: %v, Value: %s >", i.Token, i.Value)
}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
