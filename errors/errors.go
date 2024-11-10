package errors

import (
	"fmt"
	"os"
)

var HadError = false
var HadRuntimeError = false

func Error(line int, message string) {
	report(line, "", message)
}

func ReportParseError(line int, lexeme string, isAtEnd bool, message string) {
	if isAtEnd {
		report(line, " at end", message)
	} else {
		report(line, fmt.Sprintf(" at '%s'", lexeme), message)
	}
}

func ReportRuntimeError(line int, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Runtime Error: %s\n", line, message)
	HadRuntimeError = true
}

func ReportResolverError(message string) {
	fmt.Fprintf(os.Stderr, "Resolver Error: %s\n", message)
	HadError = true
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	HadError = true
}
