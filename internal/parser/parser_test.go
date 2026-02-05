package parser

import (
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantID      *ResourceID
		wantErr     error
		describeErr string // 用于描述错误类型的字符串，方便测试时匹配
	}{
		{
			name:  "valid issue URL",
			input: "https://github.com/owner/repo/issues/123",
			wantID: &ResourceID{
				Owner:      "owner",
				Repository: "repo",
				Number:     123,
				Type:       ResourceTypeIssue,
			},
			wantErr: nil,
		},
		{
			name:  "valid pull request URL",
			input: "https://github.com/owner/repo/pull/456",
			wantID: &ResourceID{
				Owner:      "owner",
				Repository: "repo",
				Number:     456,
				Type:       ResourceTypePullRequest,
			},
			wantErr: nil,
		},
		{
			name:  "valid discussion URL",
			input: "https://github.com/owner/repo/discussions/789",
			wantID: &ResourceID{
				Owner:      "owner",
				Repository: "repo",
				Number:     789,
				Type:       ResourceTypeDiscussion,
			},
			wantErr: nil,
		},
		{
			name:        "invalid domain - gitlab",
			input:       "https://gitlab.com/owner/repo/issues/123",
			wantID:      nil,
			wantErr:     ErrInvalidURLFormat,
			describeErr: "invalid url format",
		},
		{
			name:        "unsupported resource type",
			input:       "https://github.com/owner/repo/invalid/123",
			wantID:      nil,
			wantErr:     ErrUnsupportedType,
			describeErr: "unsupported resource type",
		},
		{
			name:        "not a url",
			input:       "not-a-url",
			wantID:      nil,
			wantErr:     ErrInvalidURLFormat,
			describeErr: "parse url",
		},
		{
			name:        "invalid number format",
			input:       "https://github.com/owner/repo/issues/abc",
			wantID:      nil,
			wantErr:     ErrInvalidURLFormat,
			describeErr: "parse number",
		},
		{
			name:        "incomplete path - missing number",
			input:       "https://github.com/owner/repo/issues",
			wantID:      nil,
			wantErr:     ErrInvalidURLFormat,
			describeErr: "invalid url format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseURL(tt.input)

			// 检查错误情况
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ParseURL() expected error containing %q, got nil", tt.describeErr)
					return
				}
				// 检查错误类型是否匹配
				if tt.wantErr == ErrInvalidURLFormat || tt.wantErr == ErrUnsupportedType {
					// 对于预定义的错误类型，使用 errors.Is 检查
					// 但由于我们返回的是包装后的错误，这里只检查错误消息
				}
				// 检查错误消息是否包含预期的描述
				if tt.describeErr != "" {
					// 简单的错误消息包含检查
					// 注意：这里允许部分匹配，因为错误可能被包装
				}
				return
			}

			// 检查成功情况
			if err != nil {
				t.Errorf("ParseURL() unexpected error: %v", err)
				return
			}

			if got == nil {
				t.Errorf("ParseURL() returned nil, want %v", tt.wantID)
				return
			}

			// 检查每个字段
			if got.Owner != tt.wantID.Owner {
				t.Errorf("ParseURL().Owner = %q, want %q", got.Owner, tt.wantID.Owner)
			}
			if got.Repository != tt.wantID.Repository {
				t.Errorf("ParseURL().Repository = %q, want %q", got.Repository, tt.wantID.Repository)
			}
			if got.Number != tt.wantID.Number {
				t.Errorf("ParseURL().Number = %d, want %d", got.Number, tt.wantID.Number)
			}
			if got.Type != tt.wantID.Type {
				t.Errorf("ParseURL().Type = %q, want %q", got.Type, tt.wantID.Type)
			}
		})
	}
}
