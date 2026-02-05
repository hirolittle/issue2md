package github

import (
	"time"

	"github.com/hirolittle/issue2md/internal/parser"
)

// Resource 表示从 GitHub 获取的完整资源数据
type Resource struct {
	// 基本信息
	ID          parser.ResourceID
	Title       string
	Body        string
	Author      string
	State       string // open, closed, merged, etc.
	Labels      []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	URL         string

	// 评论（扁平化，通过 InReplyTo 关联）
	Comments []Comment
}

// Comment 表示一条评论或回复
type Comment struct {
	ID        int64
	Author    string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time

	// 嵌套关系：如果是对某条评论的回复，InReplyTo 为该评论的 ID
	// 如果是对 Issue/PR 主楼的回复，InReplyTo 为 0
	InReplyTo int64
}
