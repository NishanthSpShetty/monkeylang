package evaluator

import (
	"github.com/NishanthSpShetty/monkey/ast"
	"github.com/NishanthSpShetty/monkey/runtime/errors"
	"github.com/NishanthSpShetty/monkey/runtime/evaluator/runtime"
)

var (
	Null  = &runtime.Null{}
	True  = &runtime.Boolean{Value: true}
	False = &runtime.Boolean{Value: false}
)

func Eval(r *runtime.Runtime, node ast.Node) runtime.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(r, node)

	case *ast.LetStatement:
		val := Eval(r, node.Value)

		if errors.IsError(val) {
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
		if errors.IsError(right) {
			return right
		}
		return evalPrefixExp(node.Operator, right)
		// end

	case *ast.InfixExpression:

		left := Eval(r, node.Left)
		if errors.IsError(left) {
			return left
		}

		right := Eval(r, node.Right)
		if errors.IsError(right) {
			return right
		}
		return evalInfixOperator(node.Operator, left, right)

	case *ast.IfExpression:
		return evaluateIfExpression(r, node)

	case *ast.BlockStatement:
		return evalBlockStmnt(r, node)

	case *ast.ReturnStatement:
		val := Eval(r, node.ReturnValue)

		if errors.IsError(val) {
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

		if errors.IsError(function) {
			return function
		}

		args := evalExpression(r, node.Arguments)
		if len(args) == 1 && errors.IsError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	}

	return errors.New("unknown program statement: %T", node)
}

func nativeBool(b bool) *runtime.Boolean {
	if b {
		return True
	}

	return False
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
		return errors.New("unknown operator: %s%s", op, right.Type())
	}
}

func evalBangOperatorExp(right runtime.Object) runtime.Object {
	switch right {
	case True:
		return False

	case False:
		return True
	case Null:
		// not of null ? == True
		return True
	default:
		return False
	}
}

func evalMinusPrefixOperator(right runtime.Object) runtime.Object {
	if right.Type() != runtime.ObjInteger {
		return errors.New("unknown operator: -%s", right.Type())
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
	case left.Type() != right.Type():
		return errors.New("type mismatch: %s %s %s", left.Type(), op, right.Type())

	case op == "==":
		return nativeBool(left == right)

	case op == "!=":
		return nativeBool(left != right)
	default:
		return errors.New("unknown operator: %s %s %s", left.Type(), op, right.Type())
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
		return errors.New("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}

	return &runtime.Integer{
		Value: res,
	}
}

func evaluateIfExpression(r *runtime.Runtime, ie *ast.IfExpression) runtime.Object {
	cond := Eval(r, ie.Condition)

	if errors.IsError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(r, ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(r, ie.Alternative)
	} else {
		return Null
	}
}

func isTruthy(obj runtime.Object) bool {
	switch obj {
	case Null:
		return false
	case True:
		return true
	case False:
		return false
	}
	return true
}

func evalIdentifier(r *runtime.Runtime, node *ast.Identifier) runtime.Object {
	val, ok := r.Get(node.Value)
	if !ok {
		return errors.New("identifier not found: %s", node.Value)
	}
	return val
}

func evalExpression(r *runtime.Runtime, exps []ast.Expression) []runtime.Object {
	var res []runtime.Object

	for _, exp := range exps {
		eval := Eval(r, exp)
		if errors.IsError(eval) {
			return []runtime.Object{eval}
		}
		res = append(res, eval)
	}
	return res
}

func applyFunction(fn runtime.Object, args []runtime.Object) runtime.Object {
	function, ok := fn.(*runtime.Function)
	if !ok {
		return errors.New("not a function: %s", fn.Type())
	}
	env := extendFunctionEnv(function, args)
	eval := Eval(env, function.Body)
	if rv, ok := eval.(*runtime.ReturnValue); ok {
		return rv.Value
	}
	return eval
}

func extendFunctionEnv(fn *runtime.Function, args []runtime.Object) *runtime.Runtime {
	env := runtime.NewScope(fn.Runtime)

	for i, param := range fn.Params {
		env.Put(param.Value, args[i])
	}
	return env
}
