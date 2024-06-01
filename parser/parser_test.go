package parser

import (
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
