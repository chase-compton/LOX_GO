package ast

import "github.com/chase-compton/LOX_GO/scanner"

type Stmt interface {
	Accept(visitor StmtVisitor) (interface{}, error)
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt *ExpressionStmt) (interface{}, error)
	VisitPrintStmt(stmt *PrintStmt) (interface{}, error)
	VisitVarStmt(stmt *VarStmt) (interface{}, error)
	VisitBlockStmt(stmt *BlockStmt) (interface{}, error)
	VisitIfStmt(stmt *IfStmt) (interface{}, error)
	VisitWhileStmt(stmt *WhileStmt) (interface{}, error)
	VisitFunctionStmt(stmt *FunctionStmt) (interface{}, error)
	VisitReturnStmt(stmt *ReturnStmt) (interface{}, error)
	VisitClassStmt(stmt *ClassStmt) (interface{}, error)
}

type VarStmt struct {
	Name        scanner.Token
	Initializer Expr
}

func (s *VarStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitVarStmt(s)
}

type ExpressionStmt struct {
	Expression Expr
}

func (s *ExpressionStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitExpressionStmt(s)
}

type PrintStmt struct {
	Expression Expr
}

func (s *PrintStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitPrintStmt(s)
}

type BlockStmt struct {
    Statements []Stmt
}

func (s *BlockStmt) Accept(visitor StmtVisitor) (interface{}, error) {
    return visitor.VisitBlockStmt(s)
}

type IfStmt struct {
    Condition  Expr
    ThenBranch Stmt
    ElseBranch Stmt
}

func (s *IfStmt) Accept(visitor StmtVisitor) (interface{}, error) {
    return visitor.VisitIfStmt(s)
}

type WhileStmt struct {
    Condition Expr
    Body      Stmt
}

func (s *WhileStmt) Accept(visitor StmtVisitor) (interface{}, error) {
    return visitor.VisitWhileStmt(s)
}

type FunctionStmt struct {
    Name   scanner.Token
    Params []scanner.Token
    Body   []Stmt
}

func (s *FunctionStmt) Accept(visitor StmtVisitor) (interface{}, error) {
    return visitor.VisitFunctionStmt(s)
}

type ReturnStmt struct {
    Keyword scanner.Token
    Value   Expr
}

func (s *ReturnStmt) Accept(visitor StmtVisitor) (interface{}, error) {
    return visitor.VisitReturnStmt(s)
}

type ClassStmt struct {
    Name       scanner.Token
    Superclass *Variable // For inheritance
    Methods    []*FunctionStmt
}

func (s *ClassStmt) Accept(visitor StmtVisitor) (interface{}, error) {
    return visitor.VisitClassStmt(s)
}
