package parser

import "errors"

// ResourceType 表示 GitHub 资源类型
type ResourceType string

const (
	ResourceTypeIssue       ResourceType = "issue"
	ResourceTypePullRequest ResourceType = "pull_request"
	ResourceTypeDiscussion  ResourceType = "discussion"
)

// ResourceID 唯一标识一个 GitHub 资源
type ResourceID struct {
	Owner      string       // 仓库所有者
	Repository string       // 仓库名称
	Number     int64        // Issue/PR/Discussion 编号
	Type       ResourceType // 资源类型
}

var (
	ErrInvalidURLFormat = errors.New("invalid url format")
	ErrUnsupportedType  = errors.New("unsupported resource type")
)
