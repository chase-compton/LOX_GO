package printer

import (
	"fmt"
	"github.com/chase-compton/LOX_GO/ast"
	"strings"
)

type AstPrinter struct{}

func (p *AstPrinter) Print(expr ast.Expr) (string, error) {
	result, err := expr.Accept(p)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func (p *AstPrinter) VisitBinaryExpr(expr *ast.Binary) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitGroupingExpr(expr *ast.Grouping) (interface{}, error) {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr *ast.Literal) (interface{}, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.Value), nil
}

func (p *AstPrinter) VisitUnaryExpr(expr *ast.Unary) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *AstPrinter) VisitVariableExpr(expr *ast.Variable) (interface{}, error) {
    return expr.Name.Lexeme, nil
}

func (p *AstPrinter) VisitAssignExpr(expr *ast.Assign) (interface{}, error) {
    return p.parenthesize("assign "+expr.Name.Lexeme, expr.Value)
}

func (p *AstPrinter) VisitLogicalExpr(expr *ast.Logical) (interface{}, error) {
    return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitCallExpr(expr *ast.Call) (interface{}, error) {
	callee, err := expr.Callee.Accept(p)
	if err != nil {
		return "", err
	}
	return p.parenthesize("call "+callee.(string), expr.Arguments...)
}

func (p *AstPrinter) VisitGetExpr(expr *ast.Get) (interface{}, error) {
    return p.parenthesize("get "+expr.Name.Lexeme, expr.Object)
}

func (p *AstPrinter) VisitSetExpr(expr *ast.Set) (interface{}, error) {
    return p.parenthesize("set "+expr.Name.Lexeme, expr.Object, expr.Value)
}

func (p *AstPrinter) VisitThisExpr(expr *ast.This) (interface{}, error) {
    return "this", nil
}

func (p *AstPrinter) VisitSuperExpr(expr *ast.Super) (interface{}, error) {
    return "super", nil
}

func (p *AstPrinter) parenthesize(name string, exprs ...ast.Expr) (string, error) {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		s, err := expr.Accept(p)
		if err != nil {
			return "", err
		}
		builder.WriteString(s.(string))
	}
	builder.WriteString(")")

	return builder.String(), nil
}
