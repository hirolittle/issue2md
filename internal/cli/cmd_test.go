package cli

import (
	"context"
	"strings"
	"testing"
)

func TestRun_URLParsingErrors(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr string
	}{
		{
			name:    "invalid URL format",
			url:     "not-a-url",
			wantErr: "parse url",
		},
		{
			name:    "unsupported URL type",
			url:     "https://github.com/owner/repo/tree/main",
			wantErr: "parse number", // tree 不能解析为数字
		},
		{
			name:    "invalid domain",
			url:     "https://gitlab.com/owner/repo/issues/123",
			wantErr: "invalid url format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := Run(ctx, tt.url)

			if err == nil {
				t.Errorf("Run() expected error, got nil")
				return
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Run() error = %q, want containing %q", err.Error(), tt.wantErr)
			}
		})
	}
}
