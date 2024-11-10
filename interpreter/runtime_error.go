package interpreter

import "github.com/chase-compton/LOX_GO/scanner"

type RuntimeError struct {
	Token   scanner.Token
	Message string
}

func (e *RuntimeError) Error() string {
	return e.Message
}
