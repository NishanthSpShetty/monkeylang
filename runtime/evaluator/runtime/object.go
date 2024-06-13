package runtime

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/NishanthSpShetty/monkey/ast"
)

type (
	ObjectType      string
	BuiltinFunction func(args ...Object) Object
)

const (
	ObjInteger  ObjectType = "Integer"
	ObjString   ObjectType = "String"
	ObjBoolean  ObjectType = "Boolean"
	ObjNull     ObjectType = "Nil"
	ObjReturn   ObjectType = "Return"
	ObjError    ObjectType = "Error"
	ObjFunction ObjectType = "Function"
	ObjBuiltin  ObjectType = "Builtin"
	ObjArray    ObjectType = "Array"
	ObjHash     ObjectType = "Hash"
)

var (
	Nil   = &NilType{}
	True  = &Boolean{Value: true}
	False = &Boolean{Value: false}
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashtable interface {
	HashKey() HashKey
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

func (i *Integer) HashKey() HashKey {
	return HashKey{
		Type:  i.Type(),
		Value: uint64(i.Value),
	}
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

func (b *Boolean) HashKey() HashKey {
	val := uint64(0)
	if b.Value {
		val = 1
	}
	return HashKey{
		Type:  b.Type(),
		Value: val,
	}
}

type NilType struct {
	Value bool
}

func (n *NilType) Inspect() string {
	return string(ObjNull)
}

func (n *NilType) Type() ObjectType {
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

type Function struct {
	Params  []*ast.Identifier
	Body    *ast.BlockStatement
	Runtime *Runtime
}

func (f *Function) Type() ObjectType { return ObjFunction }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Params {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString("\t" + f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return ObjString }
func (s *String) Inspect() string {
	return fmt.Sprintf("%s::[%s]", s.Value, s.Type())
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type Array struct {
	Elements []Object
}

func (a *Array) Len() int64 { return int64(len(a.Elements)) }

func (a *Array) Type() ObjectType { return ObjArray }
func (a *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]::[Array]")
	return out.String()
}

type HashPair struct {
	Key   Object
	Value Object
}
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return ObjHash }
func (h *Hash) Inspect() string {

	pairs := []string{}

	for _, v := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s:%s, ", v.Key.Inspect(), v.Value.Inspect()))
	}
	var out bytes.Buffer
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
