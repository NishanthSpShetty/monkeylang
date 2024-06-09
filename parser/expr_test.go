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
		integerValue interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
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

		testLiteralExpression(t, expr.Right, tt.integerValue)

	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"6 + 5;", 6, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
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

		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func testInfixExpression(t *testing.T, stmt ast.Expression, left interface{},
	operator string, right interface{},
) {
	exp, ok := stmt.(*ast.InfixExpression)
	assert.True(t, ok, "expression must be PrefixExpression")

	testLiteralExpression(t, exp.Left, left)
	assert.Equal(t, operator, exp.Operator, "operator must macth")
	testLiteralExpression(t, exp.Right, right)
}

func testIntegerExpression(t *testing.T, ie ast.Expression, val int64) {
	literal, ok := ie.(*ast.IntegerLiteral)
	assert.True(t, ok, "expression must be IntegerLiteral")

	assert.Equal(t, val, literal.Value)
	assert.Equal(t, fmt.Sprintf("%d", val), literal.TokenLiteral())
}

func testIdentifierExpression(t *testing.T, ie ast.Expression, val string) {
	ident, ok := ie.(*ast.Identifier)
	assert.True(t, ok, "expression must be Identifier")

	assert.Equal(t, val, ident.Value)
	assert.Equal(t, val, ident.TokenLiteral())
}

func testLiteralExpression(t *testing.T, ie ast.Expression, val interface{}) {
	switch v := val.(type) {
	case int:
		testIntegerExpression(t, ie, int64(v))
	case int64:
		testIntegerExpression(t, ie, v)
	case string:
		testIdentifierExpression(t, ie, v)
	case bool:
		testBooleanLiteral(t, ie, v)
	}
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) {
	be, ok := exp.(*ast.Boolean)
	assert.True(t, ok, "must be boolean literal")
	assert.Equal(t, value, be.Value, "expected boolean value matches")

	assert.Equal(t, fmt.Sprintf("%t", value), be.TokenLiteral(), "token value must match")
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},

		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c),)) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a,b,1,(2 * 3),(4 + 5),add(6,(7 * 8),),)",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g),)",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1,2,3,4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])),(b[1]),(2 * ([1,2][1])),)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParseErrors(t, p)

		assert.Equal(t, tt.expected, program.String(), "parsed expression")

	}
}
