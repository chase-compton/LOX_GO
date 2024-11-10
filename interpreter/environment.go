package interpreter

import (
	"fmt"

	"github.com/chase-compton/LOX_GO/scanner"
)

type Environment struct {
	Enclosing *Environment
	values    map[string]interface{}
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Enclosing: enclosing,
		values:    make(map[string]interface{}),
	}
}

func (env *Environment) Define(name string, value interface{}) {
	env.values[name] = value
}

func (env *Environment) Get(name scanner.Token) (interface{}, error) {
	if value, ok := env.values[name.Lexeme]; ok {
		return value, nil
	}

	if env.Enclosing != nil {
		return env.Enclosing.Get(name)
	}

	return nil, &RuntimeError{
		Token:   name,
		Message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
	}
}

func (e *Environment) Assign(name scanner.Token, value interface{}) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	} else if e.Enclosing != nil {
		return e.Enclosing.Assign(name, value)
	} else {
		return &RuntimeError{
			Token:   name,
			Message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
		}
	}
}

func (env *Environment) GetAt(distance int, name string) (interface{}, error) {
	ancestor := env.ancestor(distance)
	if value, ok := ancestor.values[name]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("Undefined variable '%s'.", name)
}

func (env *Environment) AssignAt(distance int, name scanner.Token, value interface{}) error {
	ancestor := env.ancestor(distance)
	if _, ok := ancestor.values[name.Lexeme]; ok {
		ancestor.values[name.Lexeme] = value
		return nil
	}
	return fmt.Errorf("Undefined variable '%s'.", name.Lexeme)
}

func (env *Environment) ancestor(distance int) *Environment {
	environment := env
	for i := 0; i < distance; i++ {
		environment = environment.Enclosing
	}
	return environment
}
