package evaluator

import (
	"testing"

	"github.com/NishanthSpShetty/monkey/runtime/evaluator/runtime"
	"github.com/stretchr/testify/assert"
)

func TestBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},

		{`len(1)`, "argument to `len` not supported, got Integer"},
	}

	for _, tt := range tests {
		eval := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, eval, int64(expected))
		case string:
			err, ok := eval.(*runtime.Error)
			if !ok {
				assert.Failf(t, "assert failed", "result must be error, got %T", eval)
				continue
			}
			assert.Equal(t, expected, err.Message, "error message dint match")
		}
	}
}
