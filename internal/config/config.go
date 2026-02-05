package config

import "time"

// Config 应用程序配置
type Config struct {
	// GitHub Token（可选）
	GitHubToken string
	// 请求超时
	Timeout time.Duration
}
