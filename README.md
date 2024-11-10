# LOX GO - A LOX Interpreter in Golang by Chase Compton & Samy Paul

An implementation of the LOX interpreter in Golang. This project is based on the book "Crafting Interpreters" by Bob Nystrom. The LOX language is a dynamically typed, object-oriented programming language.

https://craftinginterpreters.com

## Table of Contents

1. [Usage](#usage)
   - [Running the REPL](#running-the-repl)
   - [Running LOX Files](#running-lox-files)
2. [Getting Started](#getting-started)
   - [Using The Provided Executable (Recommended)](#option-a-using-the-provided-executable-recommended)
        - Step 1: Locate the Executable File
        - Step 2: Running the REPL
        - Step 3: Running LOX Files
   - [Building From Source](#option-b-building-from-source)
     - Step 1: Install Go
     - Step 2: Set Up the LOX_GO Project
     - Step 3: Running the REPL
     - Step 4: Running LOX Files
3. [Testing](#testing)
   - [Overview of Tests Implemented](#overview-of-tests-implemented)
   - [Running Tests](#running-tests)
   - [Viewing Test Reports](#viewing-test-reports)

## Usage

### Running the REPL

To start an interactive session with LOX, run the interpreter in REPL mode. This allows you to type expressions, evaluate them, and see results immediately. To run the REPL, follow the steps outlined in the [Getting Started](#getting-started) section.

### Running LOX Files

You can also execute a LOX script from a file. To run a script, use the following command from the directory:

```bash
# Using Option A: Using The Provided Executable (MacOS)
./lox path/to/yourfile.lox
# Using Option B: Building From Source
go run main.go path/to/yourfile.lox
```

This will evaluate each expression in the file sequentially, displaying results for each. Example script files can be found in the `test/benchmark` directory.

```bash
# Using Option A: Running the binary_trees.lox file using the provided executable (MacOS)
./lox test/benchmark/binary_trees.lox
# Using Option B: Running the binary_trees.lox file after building from source
go run main.go test/benchmark/binary_trees.lox
```
## Getting Started

### Option A: Using The Provided Executable (Recommended)

The easiest way to get started with LOX GO is to use the provided executable file. This option requires no additional setup or installation of Go on your system.

### Step 1: Locate the Executable File Called `lox` (Mac), `lox.exe` (Windows) or `lox_linux` (Linux)

Open a terminal or command prompt and navigate to the LOX project directory: (Assuming you have source files from the submission or elsewhere)

```bash
cd path/to/your/LOX_project_directory
```

### Step 2: Running the REPL

To start the LOX REPL, run the following command:

```bash
# For MacOS
./lox
# For Windows
./lox.exe
# For Linux
./lox_linux
```

You should see the prompt (`>`) in the terminal, indicating that the REPL is ready to accept expressions.

Example:

```bash
> var a = 10;
> var b = 20;
> print a + b;
30
```

### Step 3: Running LOX Files

To run a LOX script file, use the following command:

```bash
# For MacOS
./lox path/to/yourfile.lox
# For Windows
./lox.exe path/to/yourfile.lox
# For Linux
./lox_linux path/to/yourfile.lox
```

This will evaluate each expression in the file sequentially and display the results in the terminal. See the [Running LOX Files](#running-lox-files) section for more details.

### Option B: Building From Source

This section will guide you through setting up Go, running the LOX interpreter in REPL mode or with a LOX file, and executing test cases.

LOX_GO is built with Go, so you’ll need to install Go on your system to run the interpreter.

### Step 1: Install Go

1. Download and install Go from the official website: [https://golang.org/dl/](https://golang.org/dl/).
2. Follow the instructions for your operating system to complete the installation.
3. Verify the installation by opening a terminal or command prompt and running:

   ```bash
   go version
   ```

   You should see output indicating the version of Go installed.

### Step 2: Set Up the LOX Project

1. Open a terminal or command prompt and navigate to the LOX project directory: (Assuming you have source files from the submission or elsewhere)

   ```bash
   cd path/to/your/LOX_project_directory
   ```

2. Initialize the Go module (if not already initialized):

   ```bash
   go mod init lox
   // The above will force you to change imports in various files
   go mod init github.com/chase-compton/LOX/GO
   // This should prevent needing to change import paths
   ```

3. Might have to change import paths in the source files to match the module name depending on your `go mod init`.

   Example: Change `import "github.com/chase-compton/LOX_GO/parser"` to `import "lox/parser"` in the source files.

### Step 3: Running the REPL

To start GO_LOX in interactive mode (REPL), run the following command from the project directory:

```bash
go run main.go
```

You should see the prompt (`>`) in the terminal, indicating that the REPL is ready to accept expressions.

Example:

```bash
> var a = 10;
> var b = 20;
> print a + b;
30
```

### Step 4: Running LOX Files

To run a LOX script file, use the following command:

```bash
go run main.go path/to/yourfile.lox
```

This will evaluate each expression in the file sequentially and display the results in the terminal. See the [Running LOX Files](#running-lox-files) section for more details.

## Testing

### Overview of Tests Implemented

The LOX interpreter includes a variety of test cases to ensure that core functionalities are implemented correctly.

The test cases are from author's repository, https://github.com/munificent/craftinginterpreters/tree/master/test, and are written using the Go testing framework. The tests cover various aspects of the interpreter. The test cases are located in the `test` directory. The test files themselves are located in the `test/test_files` directory.

### Running Tests

Note: The tests are written using the Go testing framework and therefore require Go to be installed on your system. See the [Getting Started](#getting-started) section for instructions on installing Go.

IMPORTANT: If you are not on MacOS, you will need to change:
    
```bash
// Run the interpreter as a subprocess
cmd := exec.Command("../lox", tmpFile.Name())
```
to

```bash
// Run the interpreter as a subprocess
cmd := exec.Command("../lox_linux", tmpFile.Name())
// or  
cmd := exec.Command("../lox.exe", tmpFile.Name())
```
in the `runLoxSource` function in `lox_test.go` file.

To run all test cases, navigate to the root directory of the project and use the following command:

```bash
go test ./tests
```

For a more detailed test report, you can use the `-v` flag:

```bash
go test -v ./tests
```

And then if you would like to receive a simple test report, you can output the results of the test to a file:

```bash
go test -v ./tests > test_report_verbose.txt # for verbose output
go test ./tests > test_report.txt # for simple output
```

### Viewing Test Reports

After running the tests, you can view the test report in the terminal or open the `test_report.txt` file depending on which method you chose.

There will be pre-ran test reports in the main directory for reference. See `test_report.txt` for a simple test report and `test_report_verbose.txt` for a verbose test report.
