package interpreter

import (
	"github.com/chase-compton/LOX_GO/ast"
)

type LoxFunction struct {
	Declaration   *ast.FunctionStmt
	Closure       *Environment
	IsInitializer bool
}

func NewLoxFunction(declaration *ast.FunctionStmt, closure *Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{
		Declaration:   declaration,
		Closure:       closure,
		IsInitializer: isInitializer,
	}
}

func (f *LoxFunction) Arity() int {
	return len(f.Declaration.Params)
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	environment := NewEnvironment(f.Closure)
	for i, param := range f.Declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	var returnValue interface{}
	err := interpreter.executeBlockWithReturn(f.Declaration.Body, environment, &returnValue)
	if err != nil {
		if returnErr, ok := err.(*Return); ok {
			if f.IsInitializer {
				// Return 'this' from initializer
				return f.Closure.GetAt(0, "this")
			}
			return returnErr.Value, nil
		}
		return nil, err
	}

	if f.IsInitializer {
		// Return 'this' if no explicit return in initializer
		return f.Closure.GetAt(0, "this")
	}

	return nil, nil
}

func (f *LoxFunction) bind(instance *LoxInstance) *LoxFunction {
	env := NewEnvironment(f.Closure)
	env.Define("this", instance)
	return &LoxFunction{
		Declaration:   f.Declaration,
		Closure:       env,
		IsInitializer: f.IsInitializer,
	}
}
