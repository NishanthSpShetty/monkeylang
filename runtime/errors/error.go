package errors

import (
	"fmt"

	"github.com/NishanthSpShetty/monkey/runtime/object"
)

func New(format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func IsError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ObjError
}
