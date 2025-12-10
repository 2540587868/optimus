# Optimus: Functional Options Generator for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/2540587868/optimus)](https://goreportcard.com/report/github.com/2540587868/optimus)
[![Go Version](https://img.shields.io/badge/go-1.18+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[ English | [中文](README_zh.md) ]

**Optimus** is a command-line tool that automates the creation of the **Functional Options Pattern** in Go. It uses static analysis (AST) to parse your struct definitions and generates type-safe, boilerplate-free configuration code.

## What Problem Does It Solve?

In Go, the Functional Options Pattern is a best practice for creating clean, extensible, and easy-to-use APIs for object configuration. However, writing `With...` option functions by hand is:

-   **Tedious**: For a struct with 10 fields, you need to write 10 repetitive functions.
-   **Error-Prone**: Manual implementation can lead to copy-paste errors.
-   **Hard to Maintain**: Adding, removing, or changing a field requires updating multiple places.

Optimus solves this by generating all the necessary code from a single source of truth: your struct definition with tags.

## Features

-   **AST-based Generation**: No reflection is used. The generated code is as fast and efficient as hand-written code.
-   **Struct Tag Driven**: Control generation behavior using simple struct tags.
    -   Define **default values** with `default:"..."`.
    -   **Ignore fields** with `opt:"-"`.
    -   **Rename option functions** with `opt:"WithName"`.
-   **Smart Type Handling**:
    -   **Slices**: Automatically generates `AddXxx(...)` for appending values.
    -   **Pointers**: Generates option functions with value parameters for a better developer experience (e.g., `WithTimeout(10)` instead of `WithTimeout(&ten)`).
-   **Automatic Formatting**: Uses `goimports` to automatically format the generated code and manage imports.
-   **Constructor Generation**: Automatically generates a `New...` constructor that applies default values.

## Installation

```bash
go install github.com/2540587868/optimus/cmd/optimus@latest
```
*(Ensure that your `$GOPATH/bin` directory is in your system's `PATH`.)*

## Usage

### 1. Define Your Struct

Create a config struct and add the `//go:generate` directive, along with `default` and `opt` tags.

**`config/server.go`**
```go
package config

import "time"

//go:generate optimus -type=ServerConfig
type ServerConfig struct {
	// Host is the listening address.
	Host string `default:"0.0.0.0"`

	// Port is the listening port.
	Port int `default:"8080"`

	// Timeout for requests.
	Timeout time.Duration `default:"10 * time.Second"`

	// Whitelist for IPs (generates AddWhitelist).
	Whitelist []string

	// MaxConnections limit (generates pointer-friendly WithMaxConnections).
	MaxConnections *int

	// InternalState should be ignored.
	InternalState int `opt:"-"`

	// UseTLS will be renamed to WithSecure.
	UseTLS bool `opt:"WithSecure"`
}
```

### 2. Generate the Code

Run `go generate` from your project's root directory:
```bash
go generate ./...
```
This will create a `config/serverconfig_options.go` file containing the `NewServerConfig` constructor and all the option functions.

### 3. Use the Generated Code

Now you can use the generated code to initialize your config cleanly.

**`main.go`**
```go
package main

import (
	"fmt"
	"github.com/your-repo/optimus/config"
)

func main() {
	// Initialize with defaults, then apply custom options.
	cfg := config.NewServerConfig(
		config.WithHost("127.0.0.1"),
		config.WithPort(9000),
		config.WithSecure(true),
		config.AddWhitelist("192.168.1.100", "10.0.0.1"),
		config.WithMaxConnections(500),
	)

	fmt.Printf("Server starting on %s:%d\n", cfg.Host, cfg.Port)
	fmt.Printf("Secure mode: %v\n", cfg.UseTLS)
	fmt.Printf("Whitelist: %v\n", cfg.Whitelist)
	fmt.Printf("Max Connections: %d\n", *cfg.MaxConnections)
}
```


## License

This project is licensed under the [MIT License](LICENSE).