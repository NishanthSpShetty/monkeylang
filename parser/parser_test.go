package parser

import (
	"testing"

	"github.com/NishanthSpShetty/monkey/ast"
	"github.com/NishanthSpShetty/monkey/lexer"
	"github.com/stretchr/testify/assert"
)

func TestLet(t *testing.T) {
	// FIXME
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	let bar = foobar;

	`

	l := lexer.New(input)

	p := New(l)

	program := p.ParseProgram()

	assert.NotNil(t, program, "program must be non nil")

	checkParseErrors(t, p)

	assert.Equal(t, 4, len(program.Statements), "must have  3 Statements")

	tests := []struct {
		expectedIdentifier string
		val                interface{}
	}{
		{"x", 5},
		{"y", 10},
		{"foobar", 838383},
		{"bar", "foobar"},
	}

	for i, tt := range tests {
		stmnt := program.Statements[i]
		assert.True(t, testLetStmnt(t, stmnt, tt.expectedIdentifier, tt.val), "let Statements must be valid")

	}
}

func testLetStmnt(t *testing.T, s ast.Statement, name string, val interface{}) bool {
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
	testLiteralExpression(t, letStmt.Value, val)
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
}

func TestReturnStatemnt(t *testing.T) {
	// missing semicolon, works for us
	input := `
	return 5;
	return 10
	return 993322;
	`

	ret := []string{"5", "10", "993322"}

	l := lexer.New(input)

	p := New(l)

	program := p.ParseProgram()

	assert.NotNil(t, program, "program must be non nil")

	checkParseErrors(t, p)

	assert.Equal(t, 3, len(program.Statements), "must have  3 Statements")

	for i, stmnt := range program.Statements {
		retStmnt, ok := stmnt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement is not a *ast.ReturnStatement. got %T", stmnt)
		}

		assert.Equal(t, "return", retStmnt.TokenLiteral(), "invalid token")
		assert.Equal(t, ret[i], retStmnt.ReturnValue.TokenLiteral(), "must return value")
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

func TestString(t *testing.T) {
	input := `
	"hello world";
	`

	l := lexer.New(input)

	p := New(l)

	program := p.ParseProgram()

	assert.NotNil(t, program, "program must be non nil")

	checkParseErrors(t, p)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	assert.True(t, ok, "program statement must be ExpressionStatement")

	str, ok := stmt.Expression.(*ast.StringLiteral)
	assert.True(t, ok, "expression must be Identifier")

	assert.Equal(t, "hello world", str.Value)
}
