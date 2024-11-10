package test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestLox(t *testing.T) {
	testDir := "./test_files"

	var passed, failed int

	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".lox") {
			testName := strings.TrimPrefix(path, testDir+"/")
			t.Run(testName, func(t *testing.T) {
				runTestFile(t, path)
				if t.Failed() {
					failed++
					fmt.Printf("Test %s: FAIL\n", testName)
				} else {
					passed++
					fmt.Printf("Test %s: PASS\n", testName)
				}
			})
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk through test directory: %v", err)
	}

	fmt.Printf("Total tests: %d, Passed: %d, Failed: %d\n", passed+failed, passed, failed)
}

func runTestFile(t *testing.T, path string) {
    content, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("Failed to read file %s: %v", path, err)
    }
    source := string(content)

    // Determine if an error is expected based on comments
    errorExpected := strings.Contains(source, "// Error")

    // Run the interpreter and capture output and errors
    stdout, stderr, err := runLoxSource(source)
    if err != nil {
        t.Fatalf("Failed to run interpreter: %v", err)
    }

    // Check if an error occurred
    errorOccurred := len(strings.TrimSpace(stderr)) > 0

    // Evaluate the test result
    if errorExpected {
        if errorOccurred {
            // Test passes
            return
        } else {
            t.Errorf("Expected an error but none occurred in %s", path)
            return
        }
    } else {
        if errorOccurred {
            t.Errorf("Unexpected error in %s:\n%s", path, stderr)
            return
        }

        // Extract expected output
        expectedOutput := extractExpectedOutput(source)

        // Compare outputs
        if strings.TrimSpace(stdout) != strings.TrimSpace(expectedOutput) {
            t.Errorf("Test failed for %s\nExpected:\n%s\nActual:\n%s", path, expectedOutput, stdout)
        }
    }
}


func extractExpectedOutput(source string) string {
    lines := strings.Split(source, "\n")
    var outputs []string
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if strings.Contains(line, "// expect:") {
            parts := strings.Split(line, "// expect:")
            output := strings.TrimSpace(parts[1])
            outputs = append(outputs, output)
        }
    }
    return strings.Join(outputs, "\n")
}


func runLoxSource(source string) (stdout string, stderr string, err error) {
    // Create a temporary file to hold the source code
    tmpFile, err := os.CreateTemp("", "test-*.lox")
    if err != nil {
        return "", "", err
    }
    defer os.Remove(tmpFile.Name())

    // Write the source code to the temporary file
    if _, err := tmpFile.WriteString(source); err != nil {
        return "", "", err
    }
    tmpFile.Close()

    // Run the interpreter as a subprocess
    cmd := exec.Command("../lox", tmpFile.Name())

    // Capture stdout and stderr
    var outBuf, errBuf bytes.Buffer
    cmd.Stdout = &outBuf
    cmd.Stderr = &errBuf

    // Run the command
    err = cmd.Run()

    stdout = outBuf.String()
    stderr = errBuf.String()

    return stdout, stderr, nil
}
