package config

import (
	"os"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
)

// Load 从环境变量加载配置
func Load() *Config {
	cfg := &Config{
		GitHubToken: os.Getenv("GITHUB_TOKEN"),
		Timeout:     defaultTimeout,
	}

	if timeoutStr := os.Getenv("ISSUE2MD_TIMEOUT"); timeoutStr != "" {
		if d, err := time.ParseDuration(timeoutStr); err == nil {
			cfg.Timeout = d
		}
		// 解析失败时保持默认值
	}

	return cfg
}
