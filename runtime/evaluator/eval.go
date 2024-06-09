package evaluator

import (
	"github.com/NishanthSpShetty/monkey/ast"
	"github.com/NishanthSpShetty/monkey/runtime/evaluator/runtime"
)

func Eval(r *runtime.Runtime, node ast.Node) runtime.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(r, node)

	case *ast.LetStatement:
		val := Eval(r, node.Value)

		if runtime.IsError(val) {
			return val
		}
		r.Put(node.Name.Value, val)
		return nil

	case *ast.Identifier:

		return evalIdentifier(r, node)

	case *ast.ExpressionStatement:
		return Eval(r, node.Expression)

	// expressions
	case *ast.IntegerLiteral:
		return &runtime.Integer{
			Value: node.Value,
		}

	case *ast.Boolean:
		return nativeBool(node.Value)

	case *ast.PrefixExpression:
		right := Eval(r, node.Right)
		if runtime.IsError(right) {
			return right
		}
		return evalPrefixExp(node.Operator, right)
		// end

	case *ast.InfixExpression:

		left := Eval(r, node.Left)
		if runtime.IsError(left) {
			return left
		}

		right := Eval(r, node.Right)
		if runtime.IsError(right) {
			return right
		}
		return evalInfixOperator(node.Operator, left, right)

	case *ast.IfExpression:
		return evaluateIfExpression(r, node)

	case *ast.BlockStatement:
		return evalBlockStmnt(r, node)

	case *ast.ReturnStatement:
		val := Eval(r, node.ReturnValue)

		if runtime.IsError(val) {
			return val
		}
		return &runtime.ReturnValue{
			Value: val,
		}

	case *ast.FunctionLiteral:
		return &runtime.Function{
			Params:  node.Parameters,
			Body:    node.Body,
			Runtime: r,
		}

	case *ast.CallExpression:
		function := Eval(r, node.Function)

		if runtime.IsError(function) {
			return function
		}

		args := evalExpression(r, node.Arguments)
		if len(args) == 1 && runtime.IsError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	case *ast.StringLiteral:
		return &runtime.String{Value: node.Value}

	}

	return runtime.NewError("unknown program statement: %T", node)
}

func nativeBool(b bool) *runtime.Boolean {
	if b {
		return runtime.True
	}

	return runtime.False
}

func evalProgram(r *runtime.Runtime, program *ast.Program) runtime.Object {
	var result runtime.Object
	for _, stmnt := range program.Statements {
		result = Eval(r, stmnt)
		switch result := result.(type) {
		case *runtime.ReturnValue:
			return result.Value
		case *runtime.Error:
			return result
		}
	}
	return result
}

func evalBlockStmnt(r *runtime.Runtime, block *ast.BlockStatement) runtime.Object {
	var result runtime.Object
	for _, stmnt := range block.Statements {
		result = Eval(r, stmnt)

		if result != nil {
			rt := result.Type()
			if rt == runtime.ObjReturn || rt == runtime.ObjError {
				return result
			}
		}
	}
	return result
}

func evalPrefixExp(op string, right runtime.Object) runtime.Object {
	switch op {
	case "!":
		return evalBangOperatorExp(right)
	case "-":
		return evalMinusPrefixOperator(right)
	default:
		return runtime.NewError("unknown operator: %s%s", op, right.Type())
	}
}

func evalBangOperatorExp(right runtime.Object) runtime.Object {
	switch right {
	case runtime.True:
		return runtime.False

	case runtime.False:
		return runtime.True
	case runtime.Nil:
		// not of null ? == True
		return runtime.True
	default:
		return runtime.False
	}
}

func evalMinusPrefixOperator(right runtime.Object) runtime.Object {
	if right.Type() != runtime.ObjInteger {
		return runtime.NewError("unknown operator: -%s", right.Type())
	}
	value := right.(*runtime.Integer).Value
	return &runtime.Integer{
		Value: -value,
	}
}

func evalInfixOperator(op string, left, right runtime.Object) runtime.Object {
	switch {
	case left.Type() == runtime.ObjInteger && right.Type() == runtime.ObjInteger:
		return evalIntegerInfixExpression(op, left, right)
	case left.Type() == runtime.ObjString && right.Type() == runtime.ObjString:
		return evalStringInfixExpression(op, left, right)

	case left.Type() != right.Type():
		return runtime.NewError("type mismatch: %s %s %s", left.Type(), op, right.Type())
	case op == "==":
		return nativeBool(left == right)

	case op == "!=":
		return nativeBool(left != right)
	default:
		return runtime.NewError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalStringInfixExpression(op string, left, right runtime.Object) runtime.Object {
	if op != "+" {
		return runtime.NewError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
	l := left.(*runtime.String)
	r := right.(*runtime.String)
	return &runtime.String{
		Value: l.Value + r.Value,
	}
}

func evalIntegerInfixExpression(op string, left, right runtime.Object) runtime.Object {
	lval := left.(*runtime.Integer).Value
	rval := right.(*runtime.Integer).Value

	res := int64(0)
	switch op {
	case "+":
		res = lval + rval

	case "-":
		res = lval - rval

	case "*":
		res = lval * rval

	case "/":
		res = lval / rval
	case "<":
		return nativeBool(lval < rval)

	case ">":
		return nativeBool(lval > rval)

	case "==":
		return nativeBool(lval == rval)
	case "!=":

		return nativeBool(lval != rval)

	default:
		return runtime.NewError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}

	return &runtime.Integer{
		Value: res,
	}
}

func evaluateIfExpression(r *runtime.Runtime, ie *ast.IfExpression) runtime.Object {
	cond := Eval(r, ie.Condition)

	if runtime.IsError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(r, ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(r, ie.Alternative)
	} else {
		return runtime.Nil
	}
}

func isTruthy(obj runtime.Object) bool {
	switch obj {
	case runtime.Nil:
		return false
	case runtime.True:
		return true
	case runtime.False:
		return false
	}
	return true
}

func evalIdentifier(r *runtime.Runtime, node *ast.Identifier) runtime.Object {
	val, ok := r.Get(node.Value)
	if ok {
		return val
	}
	// lets look at built in too
	if bf, ok := runtime.GetBuiltin(node.Value); ok {
		return bf
	}
	return runtime.NewError("identifier not found: %s", node.Value)
}

func evalExpression(r *runtime.Runtime, exps []ast.Expression) []runtime.Object {
	var res []runtime.Object

	for _, exp := range exps {
		eval := Eval(r, exp)
		if runtime.IsError(eval) {
			return []runtime.Object{eval}
		}
		res = append(res, eval)
	}
	return res
}

func applyFunction(fn runtime.Object, args []runtime.Object) runtime.Object {
	switch fn := fn.(type) {
	case *runtime.Function:
		env := extendFunctionEnv(fn, args)
		eval := Eval(env, fn.Body)
		if rv, ok := eval.(*runtime.ReturnValue); ok {
			return rv.Value
		}
		return eval
	case *runtime.Builtin:
		return fn.Fn(args...)
	default:
		return runtime.NewError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *runtime.Function, args []runtime.Object) *runtime.Runtime {
	env := runtime.NewScope(fn.Runtime)

	for i, param := range fn.Params {
		env.Put(param.Value, args[i])
	}
	return env
}
