package interpreter

import (
    "time"
)

type ClockFunction struct{}

func (c *ClockFunction) Arity() int {
    return 0
}

func (c *ClockFunction) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
    return float64(time.Now().UnixNano()) / 1e9, nil
}

func (c *ClockFunction) String() string {
    return "<native fn>"
}
