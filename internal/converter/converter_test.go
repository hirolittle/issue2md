package converter

import (
	"context"
	"testing"
	"time"

	"github.com/hirolittle/issue2md/internal/github"
	"github.com/hirolittle/issue2md/internal/parser"
)

func TestNewConverter(t *testing.T) {
	c := NewConverter()
	if c == nil {
		t.Errorf("NewConverter() returned nil")
	}

	// 验证返回值实现了 Converter 接口
	var _ Converter = c
}

func TestConverter_Convert_Basic(t *testing.T) {
	c := NewConverter()

	resource := &github.Resource{
		ID: parser.ResourceID{
			Owner:      "owner",
			Repository: "repo",
			Number:     123,
			Type:       parser.ResourceTypeIssue,
		},
		Title:     "Test Issue",
		Body:      "This is a test issue",
		Author:    "testuser",
		State:     "open",
		Labels:    []string{"bug", "enhancement"},
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		URL:       "https://github.com/owner/repo/issues/123",
		Comments:  []github.Comment{},
	}

	ctx := context.Background()
	got, err := c.Convert(ctx, resource)

	if err != nil {
		t.Fatalf("Convert() unexpected error: %v", err)
	}

	if got == nil {
		t.Fatalf("Convert() returned nil")
	}

	// 验证输出包含预期的内容
	output := string(got)

	// 应该包含 YAML front matter
	if !contains(output, "---") {
		t.Errorf("Convert() output missing YAML front matter delimiter")
	}

	// 应该包含标题
	if !contains(output, "# Test Issue") {
		t.Errorf("Convert() output missing title")
	}

	// 应该包含正文
	if !contains(output, "This is a test issue") {
		t.Errorf("Convert() output missing body")
	}
}

// 辅助函数：检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
