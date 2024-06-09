package object

import "fmt"

type ObjectType string

const (
	ObjInteger ObjectType = "INTEGER"
	ObjBoolean ObjectType = "BOOLEAN"
	ObjNull    ObjectType = "NULL"
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
