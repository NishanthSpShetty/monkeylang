package parser

import (
	"fmt"
	"testing"

	"github.com/NishanthSpShetty/monkey/ast"
	"github.com/NishanthSpShetty/monkey/lexer"
	"github.com/stretchr/testify/assert"
)

func TestLiteralExpression(t *testing.T) {
	input := `
	50;
	`

	l := lexer.New(input)

	p := New(l)

	program := p.ParseProgram()

	assert.NotNil(t, program, "program must be non nil")

	checkParseErrors(t, p)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	assert.True(t, ok, "program statement must be ExpressionStatement")

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	assert.True(t, ok, "expression must be IntegerLiteral")

	assert.Equal(t, int64(50), literal.Value)
	assert.Equal(t, "50", literal.TokenLiteral())
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {

		l := lexer.New(tt.input)

		p := New(l)

		program := p.ParseProgram()

		assert.NotNil(t, program, "program must be non nil")

		checkParseErrors(t, p)

		assert.Equal(t, 1, len(program.Statements), "must have 1 program statement")

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok, "program statement must be ExpressionStatement")

		expr, ok := stmt.Expression.(*ast.PrefixExpression)
		assert.True(t, ok, "expression must be PrefixExpression")

		assert.Equal(t, tt.operator, expr.Operator, "operator must macth")

		testIntegerLiteral(t, expr.Right, tt.integerValue)

	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"6 + 5;", 6, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		//	{"5 != 5 + 10;", 5, "!=", 5},
	}

	for _, tt := range infixTests {

		l := lexer.New(tt.input)

		p := New(l)

		program := p.ParseProgram()

		assert.NotNil(t, program, "program must be non nil")

		checkParseErrors(t, p)

		assert.Equal(t, 1, len(program.Statements), "must have 1 program statement")

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		assert.True(t, ok, "program statement must be ExpressionStatement")

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		assert.True(t, ok, "expression must be PrefixExpression")

		testIntegerLiteral(t, exp.Left, tt.leftValue)
		assert.Equal(t, tt.operator, exp.Operator, "operator must macth")
		testIntegerLiteral(t, exp.Right, tt.rightValue)

	}
}

func testIntegerLiteral(t *testing.T, ie ast.Expression, val int64) {
	literal, ok := ie.(*ast.IntegerLiteral)
	assert.True(t, ok, "expression must be IntegerLiteral")

	assert.Equal(t, val, literal.Value)
	assert.Equal(t, fmt.Sprintf("%d", val), literal.TokenLiteral())
}
