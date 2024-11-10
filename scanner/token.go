package scanner

import "fmt"

type Token struct {
    Type    TokenType
    Lexeme  string
    Literal interface{}
    Line    int
}

func (t Token) String() string {
    return fmt.Sprintf("%v %s %v", t.Type, t.Lexeme, t.Literal)
}
