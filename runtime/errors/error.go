package errors

import (
	"fmt"

	"github.com/NishanthSpShetty/monkey/runtime/evaluator/runtime"
)

func New(format string, a ...interface{}) *runtime.Error {
	return &runtime.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func IsError(obj runtime.Object) bool {
	return obj != nil && obj.Type() == runtime.ObjError
}
