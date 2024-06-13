package evaluator

import (
	"testing"

	"github.com/NishanthSpShetty/monkey/lexer"
	"github.com/NishanthSpShetty/monkey/parser"
	"github.com/NishanthSpShetty/monkey/runtime/evaluator/runtime"
	"github.com/stretchr/testify/assert"
)

func TestEvalIntegerExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"85", 85},
		{"-85", -85},
		{"5+5", 10},
		{"5 + 5 -16", -6},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) runtime.Object {
	return Eval(runtime.New(),
		parser.New(lexer.New(input)).ParseProgram())
}

func testIntegerObject(t *testing.T, obj runtime.Object, exp int64) {
	er, ok := obj.(*runtime.Error)
	if ok {
		assert.False(t, ok, er.Inspect())
		return
	}
	io, ok := obj.(*runtime.Integer)
	assert.Truef(t, ok, "runtime must be Integer object, got %T", obj)
	if io == nil {
		return
	}
	assert.Equal(t, exp, io.Value)
	assert.Equal(t, runtime.ObjInteger, io.Type())
}

func TestEvalBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"1+1 != 2", false},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		testBoolObject(t, evaluated, tt.expected)
	}
}

func testBoolObject(t *testing.T, obj runtime.Object, exp bool) {
	io, ok := obj.(*runtime.Boolean)
	assert.True(t, ok, "runtime must be Boolean object")
	assert.Equal(t, exp, io.Value)
	assert.Equal(t, runtime.ObjBoolean, io.Type())
}

func TestBangExp(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoolObject(t, evaluated, tt.expected)
	}
}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		eval := testEval(tt.input)
		i, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, eval, int64(i))
		} else {
			testNilObject(t, eval)
		}
	}
}

func testNilObject(t *testing.T, obj runtime.Object) {
	assert.Equal(t, runtime.Nil, obj, "runtime must be of Null")
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
if (10 > 1) {
	if (10 > 1) {
		return 10;
	}
	return 1;
}`, 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: Integer + Boolean",
		},
		{
			"5 + true; 5;",
			"type mismatch: Integer + Boolean",
		},
		{
			"-true",
			"unknown operator: -Boolean",
		},
		{
			"true + false;",
			"unknown operator: Boolean + Boolean",
		},
		{
			"5; true + false; 5",
			"unknown operator: Boolean + Boolean",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: Boolean + Boolean",
		},
		{
			`if (10 > 1) {
if (10 > 1) {
return true + false; 10;
}
return 1;
}
`,
			"unknown operator: Boolean + Boolean",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*runtime.Error)
		assert.Truef(t, ok, "result must be Error Object, got %T", evaluated)
		assert.Equal(t, tt.expectedMessage, errObj.Message)
	}
}

func TestString(t *testing.T) {
	input := `"hello world"`

	eval := testEval(input)
	str, ok := eval.(*runtime.String)
	assert.True(t, ok, "expected string literal")
	assert.Equal(t, "hello world", str.Value, "string dint match")
}

func TestStringConcat(t *testing.T) {
	input := `"hello" + " " + "nishanth!"`
	eval := testEval(input)

	str, ok := eval.(*runtime.String)
	assert.True(t, ok, "expected string literal")
	assert.Equal(t, "hello nishanth!", str.Value, "string dint match")
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
	result, ok := evaluated.(*runtime.Array)
	assert.True(t, ok, "expected array literal")
	assert.Equal(t, 3, len(result.Elements), "must have 3 elements")
	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNilObject(t, evaluated)
		}
	}
}

func TestHashLiteral(t *testing.T) {
	input := `let two = "two";
{
"one": 10 - 9,
two: 1 + 1,
"thr" + "ee": 6 / 2,
4: 4,
true: 5,
false: 6
}`
	evaluated := testEval(input)
	result, ok := evaluated.(*runtime.Hash)
	assert.True(t, ok, "evaluated result is a hash")

	expected := map[runtime.HashKey]int64{
		(&runtime.String{Value: "one"}).HashKey():   1,
		(&runtime.String{Value: "two"}).HashKey():   2,
		(&runtime.String{Value: "three"}).HashKey(): 3,
		(&runtime.Integer{Value: 4}).HashKey():      4,
		runtime.True.HashKey():                      5,
		runtime.False.HashKey():                     6,
	}

	assert.Equal(t, len(result.Pairs), len(expected), "must return expected number of pairs")

	for ek, ev := range expected {
		val, ok := result.Pairs[ek]
		assert.True(t, ok, "must have a value in hash")

		assert.Equal(t, ev, val, "value must match")
	}

}
