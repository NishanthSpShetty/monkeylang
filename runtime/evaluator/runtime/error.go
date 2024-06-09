package runtime

import (
	"fmt"
)

func NewError(format string, a ...interface{}) *Error {
	return &Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func IsError(obj Object) bool {
	return obj != nil && obj.Type() == ObjError
}
