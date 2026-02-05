# API 草图 - Core Functionality

本文档描述 `internal/github` 和 `internal/converter` 包的主要接口签名。

---

## 1. internal/github

GitHub API 交互层，负责从 GitHub 获取数据。

### 1.1 资源类型

```go
// ResourceType 表示 GitHub 资源类型
type ResourceType string

const (
    ResourceTypeIssue       ResourceType = "issue"
    ResourceTypePullRequest ResourceType = "pull_request"
    ResourceTypeDiscussion  ResourceType = "discussion"
)
```

### 1.2 资源标识

```go
// ResourceID 唯一标识一个 GitHub 资源
type ResourceID struct {
    Owner      string        // 仓库所有者
    Repository string        // 仓库名称
    Number     int64         // Issue/PR/Discussion 编号
    Type       ResourceType  // 资源类型
}
```

### 1.3 数据模型

```go
// Resource 表示从 GitHub 获取的完整资源数据
type Resource struct {
    // 基本信息
    ID          ResourceID
    Title       string
    Body        string
    Author      string
    State       string    // open, closed, merged, etc.
    Labels      []string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    URL         string

    // 评论（扁平化，通过 InReplyTo 关联）
    Comments    []Comment
}

// Comment 表示一条评论或回复
type Comment struct {
    ID         int64
    Author     string
    Body       string
    CreatedAt  time.Time
    UpdatedAt  time.Time

    // 嵌套关系：如果是对某条评论的回复，InReplyTo 为该评论的 ID
    // 如果是对 Issue/PR 主楼的回复，InReplyTo 为 0
    InReplyTo  int64
}
```

### 1.4 Client 接口

```go
// Client 是 GitHub API 的客户端接口
type Client interface {
    // Fetch 根据 ResourceID 获取完整的资源数据
    Fetch(ctx context.Context, id ResourceID) (*Resource, error)
}

// NewClient 创建一个新的 GitHub API 客户端
// token 为空字符串时使用匿名访问（受速率限制）
func NewClient(token string) Client
```

### 1.5 错误处理

```go
var (
    ErrInvalidURL      = errors.New("invalid github url")
    ErrResourceNotFound = errors.New("resource not found")
    ErrAccessDenied     = errors.New("access denied: check token permissions")
    ErrRateLimited      = errors.New("rate limited: please provide a token")
)
```

---

## 2. internal/parser

URL 解析与类型识别层。

### 2.1 URL 解析

```go
// ParseURL 从 GitHub URL 中解析出 ResourceID
// 支持格式：
//   - https://github.com/owner/repo/issues/123
//   - https://github.com/owner/repo/pull/456
//   - https://github.com/owner/repo/discussions/789
func ParseURL(rawURL string) (*ResourceID, error)
```

### 2.2 错误处理

```go
var (
    ErrInvalidURLFormat = errors.New("invalid url format")
    ErrUnsupportedType  = errors.New("unsupported resource type")
)
```

---

## 3. internal/converter

数据转换为 Markdown 层。

### 3.1 Converter 接口

```go
// Converter 将 GitHub Resource 转换为 Markdown 格式
type Converter interface {
    // Convert 将 Resource 转换为 Markdown 字节数组
    Convert(ctx context.Context, resource *github.Resource) ([]byte, error)
}

// NewConverter 创建一个新的 Markdown 转换器
func NewConverter() Converter
```

### 3.2 输出格式

转换后的 Markdown 格式参考 `spec.md` 第 3 节。

---

## 4. internal/config

配置管理层。

### 4.1 配置结构

```go
// Config 应用程序配置
type Config struct {
    // GitHub Token（可选）
    GitHubToken string
    // 请求超时
    Timeout time.Duration
}

// Load 从环境变量加载配置
func Load() *Config
```

### 4.2 环境变量

| 变量名 | 必需 | 默认值 | 说明 |
|--------|------|--------|------|
| `GITHUB_TOKEN` | 否 | "" | GitHub Personal Access Token |
| `ISSUE2MD_TIMEOUT` | 否 | "30s" | API 请求超时时间 |

---

## 5. 交互流程示例

```
User Input URL
    │
    ▼
┌─────────────────┐
│  parser.ParseURL │ ───► ResourceID
└─────────────────┘
    │
    ▼
┌─────────────────┐
│  github.NewClient│ ───► Client
└─────────────────┘
    │
    ▼
┌─────────────────┐
│  client.Fetch   │ ───► Resource
└─────────────────┘
    │
    ▼
┌─────────────────┐
│ converter.Convert│ ───► []byte (Markdown)
└─────────────────┘
    │
    ▼
Output to stdout
```

---

## 6. 技术决策（已确认）

### 6.1 GitHub API 客户端
使用 `google/go-github` 库。

**原因**：
- 成熟稳定，社区广泛使用
- 完整覆盖 GitHub API v3
- 减少自行维护的负担

**依赖**：
```go
import "github.com/google/go-github/v69/github"
```

### 6.2 Markdown 格式
GitHub API 返回原始 Markdown 内容，直接使用。

**说明**：
- Issue/PR 的 `Body` 字段为原始 Markdown
- 评论的 `Body` 字段同样为原始 Markdown
- 无需进行 HTML → Markdown 转换

### 6.3 评论排序
按 `CreatedAt` 时间戳升序排列。

**实现**：
```go
sort.Slice(comments, func(i, j int) bool {
    return comments[i].CreatedAt.Before(comments[j].CreatedAt)
})
```
