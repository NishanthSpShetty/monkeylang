package evaluator

import (
	"testing"

	"github.com/NishanthSpShetty/monkey/runtime/evaluator/runtime"
	"github.com/stretchr/testify/assert"
)

func TestFunctionLiteral(t *testing.T) {
	input := "fn(x) { x+2; };"

	eval := testEval(input)

	fn, ok := eval.(*runtime.Function)

	assert.True(t, ok, "onject must be of type Function, got %T", eval)
	assert.Equal(t, 1, len(fn.Params), "must have 1 parameter list")
	assert.Equal(t, "x", fn.Params[0].Value, "expected \"x\"")
	assert.Equal(t, "(x + 2)", fn.Body.String(), "body expected (1 + 2) got %s", fn.Body.String())
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosure(t *testing.T) {
	input := `let newAdder = fn(x){
		fn(y) { x + y; }
	};

	let addTwo = newAdder(2);

	addTwo(3)`
	testIntegerObject(t, testEval(input), 5)
}
