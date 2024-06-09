package evaluator

import (
	"github.com/NishanthSpShetty/monkey/ast"
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
		return evalStatements(node.Statements)

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
		return evalPrefixExp(node.Operator, right)
		// end

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixOperator(node.Operator, left, right)
	case *ast.IfExpression:
		return evaluateIfExpression(node)

	case *ast.BlockStatement:
		return evalStatements(node.Statements)
	}

	return nil
}

func nativeBool(b bool) *object.Boolean {
	if b {
		return True
	}

	return False
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmnt := range stmts {
		result = Eval(stmnt)
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
		return Null
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
		return Null
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

	case op == "==":
		return nativeBool(left == right)

	case op == "!=":
		return nativeBool(left != right)
	}
	return Null
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
		return Null
	}

	return &object.Integer{
		Value: res,
	}
}

func evaluateIfExpression(ie *ast.IfExpression) object.Object {
	cond := Eval(ie.Condition)

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
