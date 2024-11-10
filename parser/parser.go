package parser

import (
	"fmt"

	"github.com/chase-compton/LOX_GO/ast"
	"github.com/chase-compton/LOX_GO/errors"
	"github.com/chase-compton/LOX_GO/scanner"
)

type Parser struct {
	tokens  []scanner.Token
	current int
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ([]ast.Stmt, error) {
	var statements []ast.Stmt
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}
	return statements, nil
}

func (p *Parser) printStatement() (ast.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(scanner.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return &ast.PrintStmt{Expression: value}, nil
}

func (p *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(scanner.SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}
	return &ast.ExpressionStmt{Expression: expr}, nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.MINUS, scanner.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.SLASH, scanner.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ast.Unary{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.call()
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(scanner.FALSE) {
		return &ast.Literal{Value: false}, nil
	}
	if p.match(scanner.TRUE) {
		return &ast.Literal{Value: true}, nil
	}
	if p.match(scanner.NIL) {
		return &ast.Literal{Value: nil}, nil
	}

	if p.match(scanner.NUMBER, scanner.STRING) {
		return &ast.Literal{Value: p.previous().Literal}, nil
	}

	if p.match(scanner.SUPER) {
		keyword := p.previous()
		_, err := p.consume(scanner.DOT, "Expect '.' after 'super'.")
		if err != nil {
			return nil, err
		}
		method, err := p.consume(scanner.IDENTIFIER, "Expect superclass method name.")
		if err != nil {
			return nil, err
		}
		return &ast.Super{
			Keyword: keyword,
			Method:  method,
		}, nil
	}

	if p.match(scanner.THIS) {
		return &ast.This{Keyword: p.previous()}, nil
	}

	if p.match(scanner.IDENTIFIER) {
		return &ast.Variable{Name: p.previous()}, nil
	}

	if p.match(scanner.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(scanner.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &ast.Grouping{Expression: expr}, nil
	}

	p.error(p.peek(), "Expect expression.")
	return nil, nil
}

func (p *Parser) match(types ...scanner.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scanner.EOF
}

func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() scanner.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType scanner.TokenType, message string) (scanner.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return scanner.Token{}, p.error(p.peek(), message)
}

func (p *Parser) error(token scanner.Token, message string) error {
	isAtEnd := token.Type == scanner.EOF
	lexeme := token.Lexeme
	line := token.Line
	errors.ReportParseError(line, lexeme, isAtEnd, message)
	return &ParseError{}
}

type ParseError struct{}

func (e *ParseError) Error() string {
	return "Parse error"
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == scanner.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case scanner.CLASS, scanner.FUN, scanner.VAR, scanner.FOR,
			scanner.IF, scanner.WHILE, scanner.PRINT, scanner.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) declaration() (ast.Stmt, error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(*ParseError); ok {
				p.synchronize()
			} else {
				panic(r)
			}
		}
	}()

	if p.match(scanner.CLASS) {
		return p.classDeclaration()
	}
	if p.match(scanner.FUN) {
		return p.function("function")
	}
	if p.match(scanner.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.consume(scanner.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expr
	if p.match(scanner.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(scanner.SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return &ast.VarStmt{Name: name, Initializer: initializer}, nil
}

func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.logic_or()
	if err != nil {
		return nil, err
	}

	if p.match(scanner.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if getExpr, ok := expr.(*ast.Get); ok {
			return &ast.Set{
				Object: getExpr.Object,
				Name:   getExpr.Name,
				Value:  value,
			}, nil
		} else if variable, ok := expr.(*ast.Variable); ok {
			name := variable.Name
			return &ast.Assign{
				Name:  name,
				Value: value,
			}, nil
		}

		p.error(equals, "Invalid assignment target.")
	}

	return expr, nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(scanner.FOR) {
		return p.forStatement()
	}
	if p.match(scanner.IF) {
		return p.ifStatement()
	}
	if p.match(scanner.WHILE) {
		return p.whileStatement()
	}
	if p.match(scanner.PRINT) {
		return p.printStatement()
	}
	if p.match(scanner.RETURN) {
		return p.returnStatement()
	}
	if p.match(scanner.LEFT_BRACE) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return &ast.BlockStmt{Statements: statements}, nil
	}
	return p.expressionStatement()
}

func (p *Parser) block() ([]ast.Stmt, error) {
	var statements []ast.Stmt

	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			// Optionally handle error recovery here
			p.synchronize()
			continue
		}
		statements = append(statements, stmt)
	}

	_, err := p.consume(scanner.RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) ifStatement() (ast.Stmt, error) {
	_, err := p.consume(scanner.LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch ast.Stmt
	if p.match(scanner.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &ast.IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *Parser) logic_or() (ast.Expr, error) {
	expr, err := p.logic_and()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.OR) {
		operator := p.previous()
		right, err := p.logic_and()
		if err != nil {
			return nil, err
		}
		expr = &ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) logic_and() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) whileStatement() (ast.Stmt, error) {
	_, err := p.consume(scanner.LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.RIGHT_PAREN, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) forStatement() (ast.Stmt, error) {
	_, err := p.consume(scanner.LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	// Initializer
	var initializer ast.Stmt
	if p.match(scanner.SEMICOLON) {
		initializer = nil
	} else if p.match(scanner.VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	// Condition
	var condition ast.Expr
	if !p.check(scanner.SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(scanner.SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	// Increment
	var increment ast.Expr
	if !p.check(scanner.RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(scanner.RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	// Body
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	// Desugaring
	// Increment
	if increment != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Stmt{
				body,
				&ast.ExpressionStmt{Expression: increment},
			},
		}
	}

	// Condition
	if condition == nil {
		condition = &ast.Literal{Value: true}
	}
	body = &ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}

	// Initializer
	if initializer != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Stmt{
				initializer,
				body,
			},
		}
	}

	return body, nil
}

func (p *Parser) function(kind string) (*ast.FunctionStmt, error) {
	name, err := p.consume(scanner.IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name.", kind))
	if err != nil {
		return nil, err
	}

	var parameters []scanner.Token
	if !p.check(scanner.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				p.error(p.peek(), "Cannot have more than 255 parameters.")
			}

			param, err := p.consume(scanner.IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, param)

			if !p.match(scanner.COMMA) {
				break
			}
		}
	}

	_, err = p.consume(scanner.RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body.", kind))
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &ast.FunctionStmt{
		Name:   name,
		Params: parameters,
		Body:   body,
	}, nil
}

func (p *Parser) returnStatement() (ast.Stmt, error) {
	keyword := p.previous()
	var value ast.Expr
	var err error

	if !p.check(scanner.SEMICOLON) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(scanner.SEMICOLON, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}

	return &ast.ReturnStmt{
		Keyword: keyword,
		Value:   value,
	}, nil
}

func (p *Parser) finishCall(callee ast.Expr) (ast.Expr, error) {
	var arguments []ast.Expr
	if !p.check(scanner.RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				p.error(p.peek(), "Cannot have more than 255 arguments.")
			}

			argument, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, argument)

			if !p.match(scanner.COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(scanner.RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &ast.Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}, nil
}

func (p *Parser) call() (ast.Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(scanner.LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(scanner.DOT) {
			name, err := p.consume(scanner.IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = &ast.Get{
				Object: expr,
				Name:   name,
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) classDeclaration() (ast.Stmt, error) {
	name, err := p.consume(scanner.IDENTIFIER, "Expect class name.")
	if err != nil {
		return nil, err
	}

	var superclass *ast.Variable
	if p.match(scanner.LESS) {
		_, err = p.consume(scanner.IDENTIFIER, "Expect superclass name.")
		if err != nil {
			return nil, err
		}
		superclass = &ast.Variable{Name: p.previous()}
	}

	_, err = p.consume(scanner.LEFT_BRACE, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}

	var methods []*ast.FunctionStmt
	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		method, err := p.function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}

	_, err = p.consume(scanner.RIGHT_BRACE, "Expect '}' after class body.")
	if err != nil {
		return nil, err
	}

	return &ast.ClassStmt{
		Name:       name,
		Superclass: superclass,
		Methods:    methods,
	}, nil
}
