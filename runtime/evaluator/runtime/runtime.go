package runtime

import (
	"fmt"
)

type Runtime struct {
	store map[string]Object
	outer *Runtime
}

func New() *Runtime {
	return &Runtime{
		store: map[string]Object{},
	}
}

func (r *Runtime) Put(name string, obj Object) {
	r.store[name] = obj
}

func (r *Runtime) Get(name string) (Object, bool) {
	v, ok := r.store[name]
	// if we dont find it in current scop, we will check outer scope
	if !ok && r.outer != nil {
		// could be recursive with multiple nested scope
		v, ok = r.outer.Get(name)
	}
	return v, ok
}

func (r *Runtime) PrintVars() {
	for k, v := range r.store {
		fmt.Printf(">%s = %s \n", k, v.Inspect())
	}
}

func NewScope(outer *Runtime) *Runtime {
	rt := New()
	rt.outer = outer
	return rt
}
