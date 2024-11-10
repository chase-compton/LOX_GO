package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chase-compton/LOX_GO/errors"
	"github.com/chase-compton/LOX_GO/interpreter"
	"github.com/chase-compton/LOX_GO/parser"
	"github.com/chase-compton/LOX_GO/resolver"
	"github.com/chase-compton/LOX_GO/scanner"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	source := string(bytes)
	interp := interpreter.NewInterpreter()
	runWithInterpreter(source, interp)

	if errors.HadError {
		os.Exit(65)
	}
	if errors.HadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	interp := interpreter.NewInterpreter() // Create a single interpreter instance
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// User signaled end of input
				fmt.Println()
				break
			} else {
				fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
				continue
			}
		}
		// Remove the trailing newline character
		line = strings.TrimRight(line, "\r\n")
		runWithInterpreter(line, interp) // Use the interpreter instance
		errors.HadError = false
		errors.HadRuntimeError = false
	}
}

func runWithInterpreter(source string, interp *interpreter.Interpreter) {
	errors.HadError = false
	errors.HadRuntimeError = false

	scanner := scanner.NewScanner(source)
	tokens := scanner.ScanTokens()

	if errors.HadError {
		// Scanning errors have occurred; do not proceed to parsing.
		return
	}

	p := parser.NewParser(tokens)
	statements, _ := p.Parse()
	if errors.HadError {
		// Parsing errors have occurred; do not proceed to interpretation.
		return
	}

	res := resolver.NewResolver(interp)
	_ = res.Resolve(statements)
	if errors.HadError {
		// Resolution errors have occurred; do not proceed to interpretation.
		return
	}

	interpretErr := interp.Interpret(statements)
	if interpretErr != nil {
		// Runtime errors are already reported by the interpreter.
		return
	}
}
