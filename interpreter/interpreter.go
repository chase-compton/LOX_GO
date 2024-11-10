package interpreter

import (
	"fmt"
	"github.com/chase-compton/LOX_GO/ast"
	"github.com/chase-compton/LOX_GO/errors"
	"github.com/chase-compton/LOX_GO/scanner"
)

type Interpreter struct {
	environment *Environment
	globals     *Environment
	locals      map[ast.Expr]int
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment(nil)
	interpreter := &Interpreter{
		globals:     globals,
		environment: globals,
		locals:      make(map[ast.Expr]int),
	}

	// Define native functions
	interpreter.globals.Define("clock", &ClockFunction{})

	return interpreter
}

var _ ast.ExprVisitor = &Interpreter{}
var _ ast.StmtVisitor = &Interpreter{}

func (i *Interpreter) Interpret(statements []ast.Stmt) error {
	defer func() {
		if r := recover(); r != nil {
			// Recover from panic if any
		}
	}()

	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			if runtimeErr, ok := err.(*RuntimeError); ok {
				errors.ReportRuntimeError(runtimeErr.Token.Line, runtimeErr.Message)
				return err
			}
		}
	}
	return nil
}

func (i *Interpreter) Resolve(expr ast.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) VisitExpressionStmt(stmt *ast.ExpressionStmt) (interface{}, error) {
	_, err := i.evaluate(stmt.Expression)
	return nil, err
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) (interface{}, error) {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(stringify(value))
	return nil, nil
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) (interface{}, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case scanner.MINUS:
		number, ok := right.(float64)
		if !ok {
			return nil, i.newRuntimeError(expr.Operator, "Operand must be a number.")
		}
		return -number, nil
	case scanner.BANG:
		return !isTruthy(right), nil
	}

	// Unreachable
	return nil, nil
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) (interface{}, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case scanner.MINUS:
		leftNum, ok1 := left.(float64)
		rightNum, ok2 := right.(float64)
		if !ok1 || !ok2 {
			return nil, i.newRuntimeError(expr.Operator, "Operands must be numbers.")
		}
		return leftNum - rightNum, nil

	case scanner.SLASH:
		leftNum, ok1 := left.(float64)
		rightNum, ok2 := right.(float64)
		if !ok1 || !ok2 {
			return nil, i.newRuntimeError(expr.Operator, "Operands must be numbers.")
		}
		if rightNum == 0 {
			return nil, i.newRuntimeError(expr.Operator, "Division by zero.")
		}
		return leftNum / rightNum, nil
	case scanner.STAR:
		l, ok1 := left.(float64)
		r, ok2 := right.(float64)
		if !ok1 || !ok2 {
			return nil, i.newRuntimeError(expr.Operator, "Operands must be numbers.")
		}
		return l * r, nil
	case scanner.PLUS:
		switch left.(type) {
		case float64:
			if r, ok := right.(float64); ok {
				return left.(float64) + r, nil
			}
			// Left is number, right is not
			return nil, i.newRuntimeError(expr.Operator, "Operands must be two numbers or two strings.")
		case string:
			if r, ok := right.(string); ok {
				return left.(string) + r, nil
			}
			// Left is string, right is not
			return nil, i.newRuntimeError(expr.Operator, "Operands must be two numbers or two strings.")
		default:
			// Left is neither number nor string
			return nil, i.newRuntimeError(expr.Operator, "Operands must be two numbers or two strings.")
		}
	case scanner.GREATER:
		leftNum, ok1 := left.(float64)
		rightNum, ok2 := right.(float64)
		if !ok1 || !ok2 {
			return nil, i.newRuntimeError(expr.Operator, "Operands must be numbers.")
		}
		return leftNum > rightNum, nil
	case scanner.GREATER_EQUAL:
		leftNum, ok1 := left.(float64)
		rightNum, ok2 := right.(float64)
		if !ok1 || !ok2 {
			return nil, i.newRuntimeError(expr.Operator, "Operands must be numbers.")
		}
		return leftNum >= rightNum, nil
	case scanner.LESS:
		leftNum, ok1 := left.(float64)
		rightNum, ok2 := right.(float64)
		if !ok1 || !ok2 {
			return nil, i.newRuntimeError(expr.Operator, "Operands must be numbers.")
		}
		return leftNum < rightNum, nil
	case scanner.LESS_EQUAL:
		leftNum, ok1 := left.(float64)
		rightNum, ok2 := right.(float64)
		if !ok1 || !ok2 {
			return nil, i.newRuntimeError(expr.Operator, "Operands must be numbers.")
		}
		return leftNum <= rightNum, nil
	case scanner.BANG_EQUAL:
		return !isEqual(left, right), nil
	case scanner.EQUAL_EQUAL:
		return isEqual(left, right), nil
	}

	// Unreachable
	return nil, nil
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping) (interface{}, error) {
	return i.evaluate(expr.Expression)
}
func (i *Interpreter) VisitVariableExpr(expr *ast.Variable) (interface{}, error) {
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) lookUpVariable(name scanner.Token, expr ast.Expr) (interface{}, error) {
	if distance, ok := i.locals[expr]; ok {
		return i.environment.GetAt(distance, name.Lexeme)
	} else {
		return i.globals.Get(name)
	}
}

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) (interface{}, error) {
	var value interface{}
	if stmt.Initializer != nil {
		var err error
		value, err = i.evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}

	i.environment.Define(stmt.Name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) VisitAssignExpr(expr *ast.Assign) (interface{}, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	distance, ok := i.locals[expr]
	if ok {
		err = i.environment.AssignAt(distance, expr.Name, value)
	} else {
		err = i.globals.Assign(expr.Name, value)
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) (interface{}, error) {
	return i.executeBlock(stmt.Statements, NewEnvironment(i.environment))
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) (interface{}, error) {
	conditionValue, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}

	if isTruthy(conditionValue) {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
	}
	return nil, nil
}

func (i *Interpreter) VisitLogicalExpr(expr *ast.Logical) (interface{}, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == scanner.OR {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.WhileStmt) (interface{}, error) {
	for {
		conditionValue, err := i.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}

		if !isTruthy(conditionValue) {
			break
		}

		_, err = i.execute(stmt.Body)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *ast.FunctionStmt) (interface{}, error) {
	function := NewLoxFunction(stmt, i.environment, false)
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil, nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) (interface{}, error) {
	var value interface{}
	var err error
	if stmt.Value != nil {
		value, err = i.evaluate(stmt.Value)
		if err != nil {
			return nil, err
		}
	}

	return nil, &Return{Value: value}
}

func (i *Interpreter) VisitCallExpr(expr *ast.Call) (interface{}, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	var arguments []interface{}
	for _, argument := range expr.Arguments {
		argValue, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argValue)
	}

	function, ok := callee.(Callable)
	if !ok {
		return nil, &RuntimeError{
			Token:   expr.Paren,
			Message: "Can only call functions and classes.",
		}
	}

	if len(arguments) != function.Arity() {
		return nil, &RuntimeError{
			Token:   expr.Paren,
			Message: fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments)),
		}
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) VisitClassStmt(stmt *ast.ClassStmt) (interface{}, error) {
	var superclass *LoxClass
	if stmt.Superclass != nil {
		sc, err := i.evaluate(stmt.Superclass)
		if err != nil {
			return nil, err
		}
		var ok bool
		superclass, ok = sc.(*LoxClass)
		if !ok {
			return nil, &RuntimeError{
				Token:   stmt.Superclass.Name,
				Message: "Superclass must be a class.",
			}
		}
	}

	i.environment.Define(stmt.Name.Lexeme, nil)

	if stmt.Superclass != nil {
		// Begin a new scope for 'super'
		i.environment = NewEnvironment(i.environment)
		i.environment.Define("super", superclass)
	}

	methods := make(map[string]*LoxFunction)
	for _, method := range stmt.Methods {
		isInitializer := method.Name.Lexeme == "init"
		function := NewLoxFunction(method, i.environment, isInitializer)
		methods[method.Name.Lexeme] = function
	}

	class := &LoxClass{
		Name:       stmt.Name.Lexeme,
		Methods:    methods,
		Superclass: superclass,
	}

	if stmt.Superclass != nil {
		// Restore the previous environment
		i.environment = i.environment.Enclosing
	}

	err := i.environment.Assign(stmt.Name, class)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *Interpreter) VisitGetExpr(expr *ast.Get) (interface{}, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	if instance, ok := object.(*LoxInstance); ok {
		return instance.Get(expr.Name)
	}

	return nil, &RuntimeError{
		Token:   expr.Name,
		Message: "Only instances have properties.",
	}
}

func (i *Interpreter) VisitSetExpr(expr *ast.Set) (interface{}, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	if instance, ok := object.(*LoxInstance); ok {
		value, err := i.evaluate(expr.Value)
		if err != nil {
			return nil, err
		}
		instance.Set(expr.Name, value)
		return value, nil
	}

	return nil, &RuntimeError{
		Token:   expr.Name,
		Message: "Only instances have fields.",
	}
}

func (i *Interpreter) VisitThisExpr(expr *ast.This) (interface{}, error) {
	return i.lookUpVariable(expr.Keyword, expr)
}

func (i *Interpreter) VisitSuperExpr(expr *ast.Super) (interface{}, error) {
	distance, ok := i.locals[expr]
	if !ok {
		return nil, fmt.Errorf("Undefined 'super' expression.")
	}

	// Get the superclass
	superclassInterface, err := i.environment.GetAt(distance, "super")
	if err != nil {
		return nil, err
	}
	superclass := superclassInterface.(*LoxClass)

	// 'this' is always one level nearer than 'super's environment
	objectInterface, err := i.environment.GetAt(distance-1, "this")
	if err != nil {
		return nil, err
	}
	object := objectInterface.(*LoxInstance)

	// Look up the method in the superclass
	method := superclass.findMethod(expr.Method.Lexeme)
	if method == nil {
		return nil, &RuntimeError{
			Token:   expr.Method,
			Message: fmt.Sprintf("Undefined property '%s'.", expr.Method.Lexeme),
		}
	}

	// Bind the method to 'this' instance
	return method.bind(object), nil
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, environment *Environment) (interface{}, error) {
	previous := i.environment
	i.environment = environment
	defer func() {
		i.environment = previous
	}()

	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) execute(stmt ast.Stmt) (interface{}, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) executeBlockWithReturn(statements []ast.Stmt, environment *Environment, returnValue *interface{}) error {
	previous := i.environment
	i.environment = environment
	defer func() {
		i.environment = previous
	}()

	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			if _, ok := err.(*Return); ok {
				*returnValue = err.(*Return).Value
				return err
			}
			return err
		}
	}
	return nil
}

func stringify(value interface{}) string {
	if value == nil {
		return "nil"
	}

	switch v := value.(type) {
	case float64:
		// Format numbers to avoid trailing .0 for integers
		text := fmt.Sprintf("%g", v)
		return text
	case bool:
		return fmt.Sprintf("%t", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (i *Interpreter) evaluate(expr ast.Expr) (interface{}, error) {
	return expr.Accept(i)
}

func isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}
	if b, ok := value.(bool); ok {
		return b
	}
	return true
}

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func (i *Interpreter) newRuntimeError(token scanner.Token, message string) error {
	return &RuntimeError{
		Token:   token,
		Message: message,
	}
}
