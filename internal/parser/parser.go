package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ParseURL 从 GitHub URL 中解析出 ResourceID
// 支持格式：
//   - https://github.com/owner/repo/issues/123
//   - https://github.com/owner/repo/pull/456
//   - https://github.com/owner/repo/discussions/789
func ParseURL(rawURL string) (*ResourceID, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", ErrInvalidURLFormat)
	}

	// 验证域名
	if u.Host != "github.com" {
		return nil, fmt.Errorf("%w: expected github.com, got %s", ErrInvalidURLFormat, u.Host)
	}

	// 解析路径: /owner/repo/type/number
	path := strings.Trim(u.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return nil, fmt.Errorf("%w: invalid path structure", ErrInvalidURLFormat)
	}

	owner := parts[0]
	repo := parts[1]
	resourceType := parts[2]
	numberStr := parts[3]

	number, err := strconv.ParseInt(numberStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse number: %w", ErrInvalidURLFormat)
	}

	var rt ResourceType
	switch resourceType {
	case "issues":
		rt = ResourceTypeIssue
	case "pull":
		rt = ResourceTypePullRequest
	case "discussions":
		rt = ResourceTypeDiscussion
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedType, resourceType)
	}

	return &ResourceID{
		Owner:      owner,
		Repository: repo,
		Number:     number,
		Type:       rt,
	}, nil
}
