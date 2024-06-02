package parser

import (
	"fmt"
	"testing"

	"github.com/NishanthSpShetty/monkey/ast"
	"github.com/NishanthSpShetty/monkey/lexer"
	"github.com/stretchr/testify/assert"
)

func TestLet(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;

	`

	l := lexer.New(input)

	p := New(l)

	program := p.ParseProgram()

	assert.NotNil(t, program, "program must be non nil")

	checkParseErrors(t, p)

	assert.Equal(t, 3, len(program.Statements), "must have  3 Statements")

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmnt := program.Statements[i]
		assert.True(t, testLetStmnt(t, stmnt, tt.expectedIdentifier), "let Statements must be valid")

	}
}

func testLetStmnt(t *testing.T, s ast.Statement, name string) bool {
	if "let" != s.TokenLiteral() {
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Fatalf("s not *ast.LetStatement. got %T", s)

		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s\n", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s\n", name, letStmt.Name)
		return false
	}
	return true
}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Erors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, m := range errors {
		t.Errorf("Parse error: %s", m)
	}
	t.FailNow()
}

func TestReturnStatemnt(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	l := lexer.New(input)

	p := New(l)

	program := p.ParseProgram()

	assert.NotNil(t, program, "program must be non nil")

	checkParseErrors(t, p)

	assert.Equal(t, 3, len(program.Statements), "must have  3 Statements")

	for _, stmnt := range program.Statements {
		retStmnt, ok := stmnt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement is not a *ast.ReturnStatement. got %T", stmnt)
		}

		assert.Equal(t, "return", retStmnt.TokenLiteral(), "invalid token")
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `
	foobar;
	`

	l := lexer.New(input)

	p := New(l)

	program := p.ParseProgram()

	assert.NotNil(t, program, "program must be non nil")

	checkParseErrors(t, p)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	assert.True(t, ok, "program statement must be ExpressionStatement")

	ident, ok := stmt.Expression.(*ast.Identifier)
	assert.True(t, ok, "expression must be Identifier")

	assert.Equal(t, "foobar", ident.Value)
	assert.Equal(t, "foobar", ident.TokenLiteral())
}

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
