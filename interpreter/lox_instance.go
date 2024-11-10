package interpreter

import (
	"fmt"

	"github.com/chase-compton/LOX_GO/scanner"
)

type LoxInstance struct {
    Class  *LoxClass
    Fields map[string]interface{}
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
    return &LoxInstance{
        Class:  class,
        Fields: make(map[string]interface{}),
    }
}

func (li *LoxInstance) String() string {
    return fmt.Sprintf("<%s instance>", li.Class.Name)
}

func (li *LoxInstance) Get(name scanner.Token) (interface{}, error) {
    if value, ok := li.Fields[name.Lexeme]; ok {
        return value, nil
    }

    method := li.Class.findMethod(name.Lexeme)
    if method != nil {
        return method.bind(li), nil
    }

    return nil, &RuntimeError{
        Token:   name,
        Message: fmt.Sprintf("Undefined property '%s'.", name.Lexeme),
    }
}

func (li *LoxInstance) Set(name scanner.Token, value interface{}) {
    li.Fields[name.Lexeme] = value
}
