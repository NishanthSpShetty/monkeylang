package parser

import (
	"testing"

	"github.com/NishanthSpShetty/monkey/ast"
	"github.com/NishanthSpShetty/monkey/lexer"
	"github.com/stretchr/testify/assert"
)

func TestFuncExpression(t *testing.T) {
	input := "fn(x,y){x+y;}"
	l := lexer.New(input)

	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	assert.Equal(t, 1, len(program.Statements), "must return one statements")
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok, "statement must be ExpressionStatement")

	function, ok := stmt.Expression.(*ast.FunctionLiteral)

	assert.True(t, ok, "Expression must be FunctionLiteral")
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	assert.Equal(t, 1, len(function.Body.Statements), "body must have 1 statement")

	body, ok := function.Body.Statements[0].(*ast.ExpressionStatement)

	assert.True(t, ok, "body must be ExpressionStatement")

	testInfixExpression(t, body.Expression, "x", "+", "y")
}

func TestCallExpression(t *testing.T) {
	input := "add( x, 2+4)"
	l := lexer.New(input)

	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	assert.Equal(t, 1, len(program.Statements), "must return one statements")
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok, "statement must be ExpressionStatement")

	exp, ok := stmt.Expression.(*ast.CallExpression)
	assert.True(t, ok, "Expression must be CallExpression")
	testIdentifierExpression(t, exp.Function, "add")
	assert.Equal(t, 2, len(exp.Arguments), "must have arguments")
}
