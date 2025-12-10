# Optimus: Go 函数式选项模式生成器

[![Go Report Card](https://goreportcard.com/badge/github.com/2540587868/optimus)](https://goreportcard.com/report/github.com/2540587868/optimus)
[![Go Version](https://img.shields.io/badge/go-1.18+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[ [English](README.md) | 中文 ]

**Optimus** 是一个命令行工具，用于在 Go 中自动生成**函数式选项模式 (Functional Options Pattern)** 的代码。它通过静态分析（AST）解析你的结构体定义，并生成类型安全、无样板的配置代码。

## 解决什么问题？

在 Go 中，函数式选项模式是创建清晰、可扩展且易于使用的对象配置 API 的最佳实践。然而，手动编写 `With...` 选项函数存在以下问题：

-   **枯燥乏味**：对于一个有 10 个字段的结构体，你需要手写 10 个重复的函数。
-   **容易出错**：手动实现可能导致复制粘贴错误。
-   **难以维护**：增加、删除或修改一个字段需要同步更新多个地方。

Optimus 通过唯一的代码源——带标签的结构体定义——来自动生成所有必需的代码，从而解决了这些问题。

## 功能特性

-   **基于 AST 生成**：不使用反射。生成的代码与手写代码一样高效。
-   **结构体标签驱动**：通过简单的结构体标签来控制生成行为。
    -   使用 `default:"..."` 定义**默认值**。
    -   使用 `opt:"-"` **忽略字段**。
    -   使用 `opt:"WithName"` **重命名选项函数**。
-   **智能类型处理**：
    -   **切片 (Slice)**：自动生成用于追加元素的 `AddXxx(...)` 方法。
    -   **指针 (Pointer)**：为指针字段生成接收值类型参数的选项函数，提供更好的开发体验（例如，使用 `WithTimeout(10)` 而不是 `WithTimeout(&ten)`）。
-   **自动格式化**：使用 `goimports` 自动格式化生成的代码并管理 `import`。
-   **构造函数生成**：自动生成应用了所有默认值的 `New...` 构造函数。

## 安装

```bash
go install github.com/2540587868/optimus/cmd/optimus@latest
```
*(请确保你的 `$GOPATH/bin` 目录已添加到系统的 `PATH` 环境变量中。)*

## 用法

### 1. 定义你的结构体

创建一个配置结构体，并添加 `//go:generate` 指令以及 `default` 和 `opt` 标签。

**`config/server.go`**
```go
package config

import "time"

//go:generate optimus -type=ServerConfig
type ServerConfig struct {
	// Host 是监听地址
	Host string `default:"0.0.0.0"`

	// Port 是监听端口
	Port int `default:"8080"`

	// Timeout 是请求超时时间
	Timeout time.Duration `default:"10 * time.Second"`

	// Whitelist 是 IP 白名单 (将生成 AddWhitelist)
	Whitelist []string

	// MaxConnections 是最大连接数限制 (将生成对指针友好的 WithMaxConnections)
	MaxConnections *int

	// InternalState 是内部状态，应该被忽略
	InternalState int `opt:"-"`

	// UseTLS 将被重命名为 WithSecure
	UseTLS bool `opt:"WithSecure"`
}
```

### 2. 生成代码

在你的项目根目录下运行 `go generate`：
```bash
go generate ./...
```
这将在 `config/` 目录下创建一个 `serverconfig_options.go` 文件，其中包含了 `NewServerConfig` 构造函数和所有的选项函数。

### 3. 使用生成的代码

现在你可以用生成好的代码来清晰地初始化你的配置。

**`main.go`**
```go
package main

import (
	"fmt"
	"github.com/your-repo/optimus/config"
)

func main() {
	// 使用默认值进行初始化，然后应用自定义选项
	cfg := config.NewServerConfig(
		config.WithHost("127.0.0.1"),
		config.WithPort(9000),
		config.WithSecure(true),
		config.AddWhitelist("192.168.1.100", "10.0.0.1"),
		config.WithMaxConnections(500),
	)

	fmt.Printf("服务器启动于 %s:%d\n", cfg.Host, cfg.Port)
	fmt.Printf("安全模式: %v\n", cfg.UseTLS)
	fmt.Printf("白名单: %v\n", cfg.Whitelist)
	fmt.Printf("最大连接数: %d\n", *cfg.MaxConnections)
}
```

## 许可证

本项目基于 [MIT 许可证](LICENSE) 。