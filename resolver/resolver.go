package resolver

import (
	"fmt"

	"github.com/chase-compton/LOX_GO/ast"
	"github.com/chase-compton/LOX_GO/errors"
	"github.com/chase-compton/LOX_GO/interpreter"
	"github.com/chase-compton/LOX_GO/scanner"
)

type Resolver struct {
	interpreter     *interpreter.Interpreter
	scopes          []map[string]bool
	currentClass    ClassType
	currentFunction FunctionType
}

type ClassType int

const (
	ClassTypeNone ClassType = iota
	ClassTypeClass
	ClassTypeSubclass
)

type FunctionType int

const (
	FunctionTypeNone FunctionType = iota
	FunctionTypeFunction
	FunctionTypeInitializer
	FunctionTypeMethod
)

func (r *Resolver) VisitGetExpr(expr *ast.Get) (interface{}, error) {
	_, err := r.resolveExpr(expr.Object)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitSetExpr(expr *ast.Set) (interface{}, error) {
	_, err := r.resolveExpr(expr.Value)
	if err != nil {
		return nil, err
	}
	_, err = r.resolveExpr(expr.Object)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (r *Resolver) VisitLogicalExpr(expr *ast.Logical) (interface{}, error) {
	_, err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, err
	}
	_, err = r.resolveExpr(expr.Right)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ast.ExpressionStmt) (interface{}, error) {
	return r.resolveExpr(stmt.Expression)
}

func NewResolver(interpreter *interpreter.Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
		scopes:      make([]map[string]bool, 0),
	}
}

func (r *Resolver) Resolve(statements []ast.Stmt) error {
	err := r.resolveStatements(statements)
	if err != nil {
		errors.ReportResolverError(err.Error())
	}
	return err
}

func (r *Resolver) VisitMethodStmt(stmt *ast.FunctionStmt) (interface{}, error) {
	declaration := FunctionTypeMethod
	if stmt.Name.Lexeme == "init" {
		declaration = FunctionTypeInitializer
	}
	err := r.resolveFunction(stmt, declaration)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitReturnStmt(stmt *ast.ReturnStmt) (interface{}, error) {
	if r.currentFunction == FunctionTypeNone {
		return nil, fmt.Errorf("Can't return from top-level code.")
	}
	if r.currentFunction == FunctionTypeInitializer && stmt.Value != nil {
		return nil, fmt.Errorf("Can't return a value from an initializer.")
	}
	if stmt.Value != nil {
		_, err := r.resolveExpr(stmt.Value)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitBlockStmt(stmt *ast.BlockStmt) (interface{}, error) {
	r.beginScope()
	err := r.resolveStatements(stmt.Statements)
	if err != nil {
		return nil, err
	}
	r.endScope()
	return nil, nil
}

func (r *Resolver) VisitVarStmt(stmt *ast.VarStmt) (interface{}, error) {
	err := r.declare(stmt.Name)
	if err != nil {
		return nil, err
	}
	if stmt.Initializer != nil {
		_, err := r.resolveExpr(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}
	r.define(stmt.Name)
	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt *ast.FunctionStmt) (interface{}, error) {
	err := r.declare(stmt.Name)
	if err != nil {
		return nil, err
	}
	r.define(stmt.Name)

	err = r.resolveFunction(stmt, FunctionTypeFunction)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitIfStmt(stmt *ast.IfStmt) (interface{}, error) {
	_, err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return nil, err
	}
	_, err = r.resolveStmt(stmt.ThenBranch)
	if err != nil {
		return nil, err
	}
	if stmt.ElseBranch != nil {
		_, err = r.resolveStmt(stmt.ElseBranch)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitPrintStmt(stmt *ast.PrintStmt) (interface{}, error) {
	return r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitWhileStmt(stmt *ast.WhileStmt) (interface{}, error) {
	_, err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return nil, err
	}
	_, err = r.resolveStmt(stmt.Body)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitBinaryExpr(expr *ast.Binary) (interface{}, error) {
	_, err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, err
	}
	_, err = r.resolveExpr(expr.Right)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr *ast.Grouping) (interface{}, error) {
	return r.resolveExpr(expr.Expression)
}

func (r *Resolver) VisitLiteralExpr(expr *ast.Literal) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr *ast.Unary) (interface{}, error) {
	return r.resolveExpr(expr.Right)
}

func (r *Resolver) VisitVariableExpr(expr *ast.Variable) (interface{}, error) {
	if len(r.scopes) > 0 {
		scope := r.scopes[len(r.scopes)-1]
		if defined, ok := scope[expr.Name.Lexeme]; ok && !defined {
			// Variable is declared but not yet defined
			return nil, fmt.Errorf("Cannot read local variable '%s' in its own initializer.", expr.Name.Lexeme)
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr *ast.Assign) (interface{}, error) {
	_, err := r.resolveExpr(expr.Value)
	if err != nil {
		return nil, err
	}

	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr *ast.Call) (interface{}, error) {
	_, err := r.resolveExpr(expr.Callee)
	if err != nil {
		return nil, err
	}

	for _, arg := range expr.Arguments {
		_, err := r.resolveExpr(arg)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitClassStmt(stmt *ast.ClassStmt) (interface{}, error) {
	enclosingClass := r.currentClass
	r.currentClass = ClassTypeClass

	r.declare(stmt.Name)
	r.define(stmt.Name)

	if stmt.Superclass != nil && stmt.Name.Lexeme == stmt.Superclass.Name.Lexeme {
		return nil, fmt.Errorf("A class cannot inherit from itself.")
	}

	if stmt.Superclass != nil {
		r.currentClass = ClassTypeSubclass
		_, err := r.resolveExpr(stmt.Superclass)
		if err != nil {
			return nil, err
		}

		r.beginScope()
		r.scopes[len(r.scopes)-1]["super"] = true
	}

	r.beginScope()
	r.scopes[len(r.scopes)-1]["this"] = true

	for _, method := range stmt.Methods {
		declaration := FunctionTypeMethod
		if method.Name.Lexeme == "init" {
			declaration = FunctionTypeInitializer
		}
		err := r.resolveFunction(method, declaration)
		if err != nil {
			return nil, err
		}
	}

	r.endScope()

	if stmt.Superclass != nil {
		r.endScope()
	}

	r.currentClass = enclosingClass
	return nil, nil
}

func (r *Resolver) VisitThisExpr(expr *ast.This) (interface{}, error) {
	if r.currentClass == ClassTypeNone {
		return nil, fmt.Errorf("Can't use 'this' outside of a class.")
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) VisitSuperExpr(expr *ast.Super) (interface{}, error) {
	if r.currentClass == ClassTypeNone {
		return nil, fmt.Errorf("Can't use 'super' outside of a class.")
	} else if r.currentClass != ClassTypeSubclass {
		return nil, fmt.Errorf("Can't use 'super' in a class with no superclass.")
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name scanner.Token) error {
	if len(r.scopes) == 0 {
		return nil
	}

	scope := r.scopes[len(r.scopes)-1]
	if _, exists := scope[name.Lexeme]; exists {
		return fmt.Errorf("Already a variable with this name in this scope.")
	}

	scope[name.Lexeme] = false
	return nil
}

func (r *Resolver) define(name scanner.Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr ast.Expr, name scanner.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if isDeclared, exists := r.scopes[i][name.Lexeme]; exists {
			if !isDeclared {
				errors.Error(name.Line, "Can't read local variable in its own initializer.")
			}
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function *ast.FunctionStmt, functionType FunctionType) error {
	enclosingFunction := r.currentFunction
	r.currentFunction = functionType

	r.beginScope()
	for _, param := range function.Params {
		err := r.declare(param)
		if err != nil {
			return err
		}
		r.define(param)
	}
	err := r.resolveStatements(function.Body)
	if err != nil {
		return err
	}
	r.endScope()

	r.currentFunction = enclosingFunction
	return nil
}

func (r *Resolver) resolveStatements(statements []ast.Stmt) error {
	for _, stmt := range statements {
		_, err := r.resolveStmt(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveStmt(stmt ast.Stmt) (interface{}, error) {
	return stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr ast.Expr) (interface{}, error) {
	return expr.Accept(r)
}
