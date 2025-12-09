package config

import "time"

//go:generate go run ../cmd/option-gen/main.go -type=ServerConfig
type ServerConfig struct {
	// Host 默认为 localhost
	// 注意：string 类型，生成器会自动给它加引号
	Host string `default:"localhost"`

	// Port 默认为 8080
	// 注意：int 类型，生成器直接使用原值
	Port int `default:"8080"`

	// Timeout 默认为 10秒
	// 注意：这里我们直接写 Go 代码片段，生成器会把它复制过去
	Timeout time.Duration `default:"10 * time.Second"`

	// EnableLog 默认为 true
	EnableLog bool `default:"true"`

	// Whitelist 切片暂不支持 default（因为切片初始化比较复杂），保持为空
	Whitelist []string
}
