package parser

import (
	"testing"

	"github.com/NishanthSpShetty/monkey/ast"
	"github.com/NishanthSpShetty/monkey/lexer"
	"github.com/stretchr/testify/assert"
)

func TestIfExpression(t *testing.T) {
	tests := []struct {
		input   string
		onTrue  string
		onFalse string
	}{
		{`if (x < y) { x }`, "x", ""},
		{`if (x < y) { x } else { y}`, "x", "y"},
		//{`let f = if (x < y) { x } `, "", ""},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)

		p := New(l)

		program := p.ParseProgram()
		checkParseErrors(t, p)

		assert.Equal(t, 1, len(program.Statements), "must return one statements")
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok, "statement must be ExpressionStatement")

		exp, ok := stmt.Expression.(*ast.IfExpression)

		assert.True(t, ok, "statement must be IfExpression")

		testInfixExpression(t, exp.Condition, "x", "<", "y")

		assert.Equal(t, 1, len(exp.Consequence.Statements), "Consequence must have 1 statement")
		cons, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok, "Consequence must be ExpressionStatement")
		ifTrue, ok := cons.Expression.(*ast.Identifier)
		assert.True(t, ok, "Consequence expression must be Identifier")
		assert.Equal(t, tt.onTrue, ifTrue.Value, "truthy value")
		if exp.Alternative != nil {
			f := exp.Alternative.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.Identifier)
			assert.Equal(t, tt.onFalse, f.Value, "falsy value")
		}

	}
}
