package config

import "time"

//go:generate go run ../cmd/option-gen/main.go -type=ServerConfig
type ServerConfig struct {
	// Host is the binding address
	Host string

	// Port 监听端口
	// 默认 8080
	Port int

	// Timeout 请求超时时间
	Timeout time.Duration

	// Whitelist 这里的 []string 会自动生成 AddWhitelist 方法
	Whitelist []string

	// Internal 这里的 opt:"-" 表示不对外暴露，不生成 Option
	Internal int `opt:"-"`

	// TLS 这里的 opt:"WithSecure" 表示生成的函数改名叫 WithSecure
	TLS bool `opt:"WithSecure"`
}