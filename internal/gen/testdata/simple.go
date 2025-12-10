package testdata

import (
	"fmt"
	"time"
)

//go:generate optimus -type=SimpleConfig
type SimpleConfig struct {
	// Host 监听地址
	Host string `default:"0.0.0.0"`

	// Port 监听端口
	Port int `default:"8080"`

	// Timeout 请求超时
	Timeout time.Duration `default:"10 * time.Second"`

	// DBList 数据库列表 (测试切片)
	DBList []string

	// MaxConn 最大连接数 (测试指针优化)
	// nil 表示不限制
	MaxConn *int

	// Tags 元数据 (测试 Map 支持)
	Tags map[string]string

	// UseTLS 是否开启TLS (测试改名)
	UseTLS bool `opt:"WithSecure"`
}

func (c *SimpleConfig) Validate() error {
	if c.Port < 0 {
		return fmt.Errorf("port cannot be negative")
	}
	return nil
}
