package interpreter

import "fmt"

type LoxClass struct {
    Name       string
    Methods    map[string]*LoxFunction
    Superclass *LoxClass
}

func (c *LoxClass) String() string {
    return fmt.Sprintf("<class %s>", c.Name)
}

func (c *LoxClass) Arity() int {
    initializer := c.findMethod("init")
    if initializer != nil {
        return initializer.Arity()
    }
    return 0
}

func (c *LoxClass) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
    instance := NewLoxInstance(c)
    initializer := c.findMethod("init")
    if initializer != nil {
        initializer.bind(instance).Call(interpreter, arguments)
    }
    return instance, nil
}

func (c *LoxClass) findMethod(name string) *LoxFunction {
    if method, ok := c.Methods[name]; ok {
        return method
    }
    if c.Superclass != nil {
        return c.Superclass.findMethod(name)
    }
    return nil
}
