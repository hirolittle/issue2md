package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// 保存原始环境变量
	originalToken := os.Getenv("GITHUB_TOKEN")
	originalTimeout := os.Getenv("ISSUE2MD_TIMEOUT")

	// 测试完成后恢复环境变量
	defer func() {
		if originalToken != "" {
			os.Setenv("GITHUB_TOKEN", originalToken)
		} else {
			os.Unsetenv("GITHUB_TOKEN")
		}
		if originalTimeout != "" {
			os.Setenv("ISSUE2MD_TIMEOUT", originalTimeout)
		} else {
			os.Unsetenv("ISSUE2MD_TIMEOUT")
		}
	}()

	tests := []struct {
		name           string
		setToken       string
		setTimeout     string
		wantToken      string
		wantTimeout    time.Duration
	}{
		{
			name:        "no environment variables",
			setToken:    "",
			setTimeout:  "",
			wantToken:   "",
			wantTimeout: 30 * time.Second,
		},
		{
			name:        "GITHUB_TOKEN set",
			setToken:    "test-token",
			setTimeout:  "",
			wantToken:   "test-token",
			wantTimeout: 30 * time.Second,
		},
		{
			name:        "ISSUE2MD_TIMEOUT set",
			setToken:    "",
			setTimeout:  "60s",
			wantToken:   "",
			wantTimeout: 60 * time.Second,
		},
		{
			name:        "both environment variables set",
			setToken:    "xxx",
			setTimeout:  "60s",
			wantToken:   "xxx",
			wantTimeout: 60 * time.Second,
		},
		{
			name:        "invalid ISSUE2MD_TIMEOUT uses default",
			setToken:    "",
			setTimeout:  "invalid",
			wantToken:   "",
			wantTimeout: 30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清除环境变量
			os.Unsetenv("GITHUB_TOKEN")
			os.Unsetenv("ISSUE2MD_TIMEOUT")

			// 设置测试环境变量
			if tt.setToken != "" {
				os.Setenv("GITHUB_TOKEN", tt.setToken)
			}
			if tt.setTimeout != "" {
				os.Setenv("ISSUE2MD_TIMEOUT", tt.setTimeout)
			}

			// 调用 Load
			got := Load()

			// 验证结果
			if got == nil {
				t.Fatalf("Load() returned nil, want non-nil Config")
			}

			if got.GitHubToken != tt.wantToken {
				t.Errorf("Load().GitHubToken = %q, want %q", got.GitHubToken, tt.wantToken)
			}

			if got.Timeout != tt.wantTimeout {
				t.Errorf("Load().Timeout = %v, want %v", got.Timeout, tt.wantTimeout)
			}
		})
	}
}
