package object

import "fmt"

type ObjectType string

const (
	ObjInteger ObjectType = "Integer"
	ObjBoolean ObjectType = "Boolean"
	ObjNull    ObjectType = "Null"
	ObjReturn  ObjectType = "Return"
	ObjError   ObjectType = "Error"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d::[%s]", i.Value, i.Type())
}

func (i *Integer) Type() ObjectType {
	return ObjInteger
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t::[%s]", b.Value, b.Type())
}

func (b *Boolean) Type() ObjectType {
	return ObjBoolean
}

type Null struct {
	Value bool
}

func (n *Null) Inspect() string {
	return string(ObjNull)
}

func (n *Null) Type() ObjectType {
	return ObjNull
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

func (rv *ReturnValue) Type() ObjectType {
	return ObjReturn
}

type Error struct {
	Message string
}

func (e *Error) Inspect() string {
	return "Error: " + e.Message
}

func (e *Error) Type() ObjectType {
	return ObjError
}
