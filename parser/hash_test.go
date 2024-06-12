package parser

import (
	"testing"

	"github.com/NishanthSpShetty/monkey/ast"
	"github.com/NishanthSpShetty/monkey/lexer"
	"github.com/stretchr/testify/assert"
)

func TestHashLiteral(t *testing.T) {

	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	assert.NotNil(t, program, "program must be non nil")

	stmnt := program.Statements[0].(*ast.ExpressionStatement)

	m, ok := stmnt.Expression.(*ast.HashLiteral)
	assert.Truef(t, ok, "must be HashLiteral, got %T", stmnt)

	assert.Equalf(t, 3, len(m.Pairs), "must have 3 pairs")

	exp := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for k, v := range m.Pairs {

		key, ok := k.(*ast.StringLiteral)
		assert.Truef(t, ok, "Key must be StringLiteral, got %T", k)
		val := v.(*ast.IntegerLiteral)
		assert.Truef(t, ok, "Value must be IntegerLiteral, got %T", v)

		assert.Equalf(t, exp[key.Value], val.Value, "must match for key %s", key.Value)
		testIntegerExpression(t, val, exp[key.Value])
	}

}

func TestEmptyHashLiteral(t *testing.T) {

	input := `{}`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	assert.NotNil(t, program, "program must be non nil")

	stmnt := program.Statements[0].(*ast.ExpressionStatement)

	m, ok := stmnt.Expression.(*ast.HashLiteral)
	assert.Truef(t, ok, "must be HashLiteral, got %T", stmnt)

	assert.Equalf(t, 0, len(m.Pairs), "must have 3 pairs")
}

func TestHasExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	assert.NotNil(t, program, "program must be non nil")

	stmnt := program.Statements[0].(*ast.ExpressionStatement)

	m, ok := stmnt.Expression.(*ast.HashLiteral)
	assert.Truef(t, ok, "must be HashLiteral, got %T", stmnt)

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}
	for key, value := range m.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		assert.Truef(t, ok, "key must be StingLiteral, got %T", key)

		testFunc, ok := tests[literal.String()]
		assert.True(t, ok, "no expected test function, got")

		testFunc(value)
	}
}
