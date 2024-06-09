package runtime

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return ObjBuiltin }
func (b *Builtin) Inspect() string {
	return "builtin function"
}

func fnLen() *Builtin {
	return &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			default:
				return NewError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
	}
}

var builtins = map[string]*Builtin{
	"len": fnLen(),
}

func GetBuiltin(name string) (*Builtin, bool) {
	b, ok := builtins[name]
	return b, ok
}
