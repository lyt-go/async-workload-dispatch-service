// Package config 负责从环境变量加载服务配置。
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config 服务配置。
type Config struct {
	Addr              string
	MaxPageSize       int
	DefaultMaxRetry   int // 队列默认最大重试次数
	BackoffBaseMs     int // 重试退避基数（毫秒）
	VisibilityTimeout int // 任务可见性超时（秒）
}

// Load 从环境变量加载配置。
func Load() *Config {
	cfg := &Config{
		Addr:              ":" + getenv("PORT", "8080"),
		MaxPageSize:       getenvInt("MAX_PAGE_SIZE", 100),
		DefaultMaxRetry:   getenvInt("DEFAULT_MAX_RETRY", 3),
		BackoffBaseMs:     getenvInt("BACKOFF_BASE_MS", 1000),
		VisibilityTimeout: getenvInt("VISIBILITY_TIMEOUT", 30),
	}
	if v := os.Getenv("ADDR"); v != "" {
		cfg.Addr = v
	}
	return cfg
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return def
	}
	return n
}

func (c *Config) String() string {
	return fmt.Sprintf("addr=%s max_page_size=%d default_max_retry=%d", c.Addr, c.MaxPageSize, c.DefaultMaxRetry)
}
