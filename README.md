# Makego
A simple build automation tool similar to Make, but with a simpler featureset.

The tool is able to execute targets with dependencies and commands, focusing on the core functionality of a build system.

## Functions:
### `Make() error`
Make is a function that parse a makefile specified by the command line flag -f or the makefile in the current directory then execute the targets specified by the command line arguments or execute the default target

## Features:
### 1. Configuration File:
   - Parse a custom configuration file `Makefile`
   - Support defining targets, their dependencies, and associated commands

### 2. Dependency Resolution:
   - Detect and report circular dependencies

### 3. Command Execution:
   - Execute shell commands associated with each target
   - Support running multiple commands for a single target

### 4. CLI Interface:
   - Accept target names as command-line arguments
   - If no target is specified, run the default target (if defined)
   - Accept a custom Makefile via optional '-f' flag

### 5. Concurrency:
   - Execute independent targets concurrently using goroutines
   - Implement proper synchronization for dependent targets

### 6. Error Handling:
   - Handle and report errors in Makefile parsing and commands execution

## How to Use:
### Step 1: Install the tool using `go get`

  ```bash
  go get github.com/codescalersinternships/Makego-Fatma-Ebrahim
  ```

This command fetches the tool and adds it to your project's `go.mod` file.

### Step 2: Import and use the tool in your code

  After running `go get`, you can import the tool into your project and use the functions as described:
```
package main

import (
	"fmt"

	makego "github.com/codescalersinternships/Makego-Fatma-Ebrahim/pkg"
)

func main() {
	err := makego.Make()
	if err != nil {
		fmt.Println(err)
	}
}
```

### Step 3: Specify the targets
You can use the default target or specify targets via command-line arguments as follows:
```
go run main.go clean
```
Also, you can specify a custom Makefile via command-line flag as follows:
```
go run main.go -f ../Makefile
```
