package evaluator

import (
	"github.com/NishanthSpShetty/monkey/ast"
	"github.com/NishanthSpShetty/monkey/runtime/errors"
	"github.com/NishanthSpShetty/monkey/runtime/object"
)

var (
	Null  = &object.Null{}
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// expressions
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}

	case *ast.Boolean:
		return nativeBool(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if errors.IsError(right) {
			return right
		}
		return evalPrefixExp(node.Operator, right)
		// end

	case *ast.InfixExpression:

		left := Eval(node.Left)
		if errors.IsError(left) {
			return left
		}

		right := Eval(node.Right)
		if errors.IsError(right) {
			return right
		}
		return evalInfixOperator(node.Operator, left, right)

	case *ast.IfExpression:
		return evaluateIfExpression(node)

	case *ast.BlockStatement:
		return evalBlockStmnt(node)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)

		if errors.IsError(val) {
			return val
		}
		return &object.ReturnValue{
			Value: val,
		}

	}

	return nil
}

func nativeBool(b bool) *object.Boolean {
	if b {
		return True
	}

	return False
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, stmnt := range program.Statements {
		result = Eval(stmnt)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStmnt(block *ast.BlockStatement) object.Object {
	var result object.Object
	for _, stmnt := range block.Statements {
		result = Eval(stmnt)

		if result != nil {
			rt := result.Type()
			if rt == object.ObjReturn || rt == object.ObjError {
				return result
			}
		}
	}
	return result
}

func evalPrefixExp(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExp(right)
	case "-":
		return evalMinusPrefixOperator(right)
	default:
		return errors.New("unknown operator: %s%s", op, right.Type())
	}
}

func evalBangOperatorExp(right object.Object) object.Object {
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

func evalMinusPrefixOperator(right object.Object) object.Object {
	if right.Type() != object.ObjInteger {
		return errors.New("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{
		Value: -value,
	}
}

func evalInfixOperator(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.ObjInteger && right.Type() == object.ObjInteger:
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

func evalIntegerInfixExpression(op string, left, right object.Object) object.Object {
	lval := left.(*object.Integer).Value
	rval := right.(*object.Integer).Value

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

	return &object.Integer{
		Value: res,
	}
}

func evaluateIfExpression(ie *ast.IfExpression) object.Object {
	cond := Eval(ie.Condition)

	if errors.IsError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return Null
	}
}

func isTruthy(obj object.Object) bool {
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
