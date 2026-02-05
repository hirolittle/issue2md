# issue2md 技术实现方案

> **版本**: v1.0
> **日期**: 2025-02-05
> **状态**: 待评审

---

## 1. 技术上下文

### 1.1 技术栈

| 层级 | 技术选型 | 说明 |
|------|----------|------|
| 语言 | Go >= 1.24 | 主要编程语言 |
| Web框架 | `net/http` (标准库) | 遵循简单性原则，不引入第三方框架 |
| GitHub API | `google/go-github` v69 | REST API v3 + GraphQL API v4 |
| Markdown | 标准库 `strings`/`fmt` | 不使用第三方模板库 |
| 数据存储 | N/A | 实时获取，无需持久化 |
| 测试 | `testing` (标准库) | 表格驱动测试，集成测试优先 |

### 1.2 外部依赖

```go
// go.mod
require (
    github.com/google/go-github/v69 v69.0.0
    golang.org/x/oauth2 v0.24.0 // GitHub Token 认证
)
```

### 1.3 技术约束

- **Web 服务**: 仅使用标准库 `net/http`
- **GitHub API**: 使用 `google/go-github` 库
- **Markdown 处理**: 不使用第三方库，直接拼接字符串
- **数据存储**: 无需数据库

---

## 2. "合宪性"审查

本方案逐条对照 `constitution.md` 进行审查：

### 2.1 第一条：简单性原则 ✓

| 宪法条款 | 本方案措施 | 状态 |
|----------|------------|------|
| 1.1 YAGNI | 仅实现 spec.md 明确要求的功能，Out of Scope 功能不在 v1 实现 | ✓ |
| 1.2 标准库优先 | Web 框架使用 `net/http`，Markdown 处理使用标准库 | ✓ |
| 1.3 反过度工程 | 核心数据结构使用简单 `struct`，避免复杂接口体系 | ✓ |

### 2.2 第二条：测试先行铁律 ✓

| 宪法条款 | 本方案措施 | 状态 |
|----------|------------|------|
| 2.1 TDD 循环 | 每个包从失败的测试开始开发 | ✓ |
| 2.2 表格驱动 | 所有单元测试采用表格驱动风格 | ✓ |
| 2.3 拒绝 Mock | 优先编写集成测试，使用真实 GitHub API（或测试环境） | ✓ |

### 2.3 第三条：明确性原则 ✓

| 宪法条款 | 本方案措施 | 状态 |
|----------|------------|------|
| 3.1 错误处理 | 所有错误显式处理，使用 `fmt.Errorf("...: %w", err)` 包装 | ✓ |
| 3.2 无全局变量 | 所有依赖通过构造函数注入，无全局状态 | ✓ |

**审查结论**: 本方案完全符合 `constitution.md` 所有条款。

---

## 3. 项目结构细化

### 3.1 目录树

```
issue2md/
├── cmd/
│   ├── issue2md/           # CLI 入口
│   │   └── main.go
│   └── issue2mdweb/        # Web 入口 (future)
│       └── main.go
├── internal/
│   ├── cli/                # 命令行接口层
│   │   ├── cmd.go          # 命令定义与执行
│   │   └── error.go        # CLI 错误处理
│   ├── config/             # 配置管理层
│   │   ├── config.go       # 配置结构
│   │   └── load.go         # 环境变量加载
│   ├── parser/             # URL 解析层
│   │   ├── parser.go       # URL 解析逻辑
│   │   └── types.go        # ResourceID 等类型定义
│   ├── github/             # GitHub API 交互层
│   │   ├── client.go       # Client 实现
│   │   ├── fetch.go        # 数据获取逻辑
│   │   ├── types.go        # Resource, Comment 等数据结构
│   │   └── errors.go       # 错误定义
│   └── converter/          # Markdown 转换层
│       ├── converter.go    # Converter 实现
│       └── format.go       # Markdown 格式化逻辑
├── web/
│   ├── templates/          # Web 模板 (future)
│   └── static/             # 静态资源 (future)
├── specs/                  # 规格文档
│   ├── spec.md
│   └── 001-core-functionality/
│       ├── api-sketch.md
│       └── plan.md
├── constitution.md         # 项目宪法
├── CLAUDE.md               # AI 协作指令
├── go.mod
├── go.sum
└── Makefile
```

### 3.2 包职责与依赖关系

```
┌─────────────────────────────────────────────────────────────────┐
│                         cmd/issue2md                            │
│                     (main.go - 程序入口)                         │
└────────────────────────────────┬────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                      internal/cli                               │
│                  (命令行参数处理、错误输出)                       │
└────────────┬────────────────────────────────┬───────────────────┘
             │                                │
             ▼                                ▼
┌─────────────────────────┐      ┌────────────────────────────────┐
│   internal/config       │      │      internal/parser           │
│  (环境变量配置加载)       │      │     (URL 解析与类型识别)        │
└─────────────────────────┘      └────────────┬───────────────────┘
                                              │
                                              ▼
                                   ┌────────────────────────────────┐
                                   │      internal/github            │
                                   │  (GitHub API 客户端、数据获取)   │
                                   └────────────┬───────────────────┘
                                                │
                                                ▼
                                   ┌────────────────────────────────┐
                                   │     internal/converter          │
                                   │    (Resource → Markdown 转换)   │
                                   └────────────────────────────────┘
```

### 3.3 依赖规则

- **cmd** → internal/cli
- **internal/cli** → internal/config, internal/parser, internal/github, internal/converter
- **internal/github** → (仅外部依赖: google/go-github)
- **internal/converter** → internal/github (使用 Resource 类型)
- **internal/parser** → 无依赖
- **internal/config** → 无依赖

---

## 4. 核心数据结构

### 4.1 跨包共享类型 (`internal/parser/types.go`)

```go
package parser

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
```

### 4.2 GitHub 数据模型 (`internal/github/types.go`)

```go
package github

import "time"

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
```

### 4.3 配置结构 (`internal/config/config.go`)

```go
package config

import "time"

// Config 应用程序配置
type Config struct {
    // GitHub Token（可选）
    GitHubToken string
    // 请求超时
    Timeout time.Duration
}
```

---

## 5. 接口设计

### 5.1 URL 解析器 (`internal/parser/parser.go`)

```go
package parser

import (
    "errors"
    "net/url"
    "strconv"
    "strings"
)

var (
    ErrInvalidURLFormat = errors.New("invalid url format")
    ErrUnsupportedType  = errors.New("unsupported resource type")
)

// ParseURL 从 GitHub URL 中解析出 ResourceID
// 支持格式：
//   - https://github.com/owner/repo/issues/123
//   - https://github.com/owner/repo/pull/456
//   - https://github.com/owner/repo/discussions/789
func ParseURL(rawURL string) (*ResourceID, error) {
    u, err := url.Parse(rawURL)
    if err != nil {
        return nil, fmt.Errorf("parse url: %w", err)
    }

    // 验证域名
    if u.Host != "github.com" {
        return nil, fmt.Errorf("%w: expected github.com, got %s", ErrInvalidURLFormat, u.Host)
    }

    // 解析路径: /owner/repo/type/number
    parts := strings.Split(strings.Trim(u.Path, "/"), "/")
    if len(parts) < 4 {
        return nil, fmt.Errorf("%w: invalid path structure", ErrInvalidURLFormat)
    }

    owner := parts[0]
    repo := parts[1]
    resourceType := parts[2]
    numberStr := parts[3]

    number, err := strconv.ParseInt(numberStr, 10, 64)
    if err != nil {
        return nil, fmt.Errorf("parse number: %w", err)
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
```

### 5.2 GitHub 客户端 (`internal/github/client.go`)

```go
package github

import (
    "context"
    "os"

    "github.com/google/go-github/v69/github"
    "golang.org/x/oauth2"
)

var (
    ErrResourceNotFound = errors.New("resource not found")
    ErrAccessDenied     = errors.New("access denied: check token permissions")
    ErrRateLimited      = errors.New("rate limited: please provide a token")
)

// Client 是 GitHub API 的客户端接口
type Client interface {
    // Fetch 根据 ResourceID 获取完整的资源数据
    Fetch(ctx context.Context, id ResourceID) (*Resource, error)
}

// client 实现 Client 接口
type client struct {
    gh *github.Client
}

// NewClient 创建一个新的 GitHub API 客户端
// token 为空字符串时使用匿名访问（受速率限制）
func NewClient(token string) Client {
    var httpClient *http.Client
    if token != "" {
        ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
        httpClient = oauth2.NewClient(context.Background(), ts)
    }
    return &client{
        gh: github.NewClient(httpClient),
    }
}

// Fetch 实现接口
func (c *client) Fetch(ctx context.Context, id ResourceID) (*Resource, error) {
    // 根据 ResourceID.Type 调用不同的 GitHub API
    switch id.Type {
    case parser.ResourceTypeIssue, parser.ResourceTypePullRequest:
        return c.fetchIssueOrPR(ctx, id)
    case parser.ResourceTypeDiscussion:
        return c.fetchDiscussion(ctx, id)
    default:
        return nil, fmt.Errorf("unsupported resource type: %s", id.Type)
    }
}
```

### 5.3 Markdown 转换器 (`internal/converter/converter.go`)

```go
package converter

import (
    "context"
    "fmt"

    "github.com/owner/project/internal/github"
)

// Converter 将 GitHub Resource 转换为 Markdown 格式
type Converter interface {
    // Convert 将 Resource 转换为 Markdown 字节数组
    Convert(ctx context.Context, resource *github.Resource) ([]byte, error)
}

// converter 实现 Converter 接口
type converter struct{}

// NewConverter 创建一个新的 Markdown 转换器
func NewConverter() Converter {
    return &converter{}
}

// Convert 实现接口
func (c *converter) Convert(ctx context.Context, resource *github.Resource) ([]byte, error) {
    var buf bytes.Buffer

    // 1. 写入 YAML Front Matter
    c.writeFrontMatter(&buf, resource)

    // 2. 写入标题
    fmt.Fprintf(&buf, "# %s\n\n", resource.Title)

    // 3. 写入正文
    if resource.Body != "" {
        fmt.Fprintf(&buf, "%s\n\n", resource.Body)
    }

    // 4. 写入评论
    if len(resource.Comments) > 0 {
        c.writeComments(&buf, resource.Comments)
    }

    return buf.Bytes(), nil
}

func (c *converter) writeFrontMatter(buf *bytes.Buffer, r *github.Resource) {
    // YAML front matter 格式化
    fmt.Fprintf(buf, "---\n")
    fmt.Fprintf(buf, "title: %q\n", r.Title)
    fmt.Fprintf(buf, "author: %s\n", r.Author)
    fmt.Fprintf(buf, "status: %s\n", r.State)
    // ... 其他字段
    fmt.Fprintf(buf, "---\n\n")
}

func (c *converter) writeComments(buf *bytes.Buffer, comments []github.Comment) {
    // 按时间戳排序
    sorted := make([]github.Comment, len(comments))
    copy(sorted, comments)
    sort.Slice(sorted, func(i, j int) bool {
        return sorted[i].CreatedAt.Before(sorted[j].CreatedAt)
    })

    fmt.Fprintf(buf, "## Comments\n\n")

    for _, comment := range sorted {
        if comment.InReplyTo > 0 {
            // 嵌套回复格式
            fmt.Fprintf(buf, "> **Reply to @%s**\n>\n", findAuthor(comment.InReplyTo, sorted))
            fmt.Fprintf(buf, "> %s\n\n", comment.Body)
        } else {
            // 顶级评论格式
            fmt.Fprintf(buf, "### @%s - %s\n\n", comment.Author, comment.CreatedAt.Format(time.RFC3339))
            fmt.Fprintf(buf, "%s\n\n", comment.Body)
        }
    }
}
```

### 5.4 CLI 入口 (`internal/cli/cmd.go`)

```go
package cli

import (
    "context"
    "fmt"
    "os"

    "github.com/owner/project/internal/config"
    "github.com/owner/project/internal/converter"
    "github.com/owner/project/internal/github"
    "github.com/owner/project/internal/parser"
)

// Run 执行 CLI 命令
func Run(ctx context.Context, url string) error {
    // 1. 加载配置
    cfg := config.Load()

    // 2. 创建带超时的上下文
    ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
    defer cancel()

    // 3. 解析 URL
    resourceID, err := parser.ParseURL(url)
    if err != nil {
        return fmt.Errorf("parse url: %w", err)
    }

    // 4. 创建 GitHub 客户端
    client := github.NewClient(cfg.GitHubToken)

    // 5. 获取数据
    resource, err := client.Fetch(ctx, *resourceID)
    if err != nil {
        return fmt.Errorf("fetch resource: %w", err)
    }

    // 6. 转换为 Markdown
    conv := converter.NewConverter()
    markdown, err := conv.Convert(ctx, resource)
    if err != nil {
        return fmt.Errorf("convert markdown: %w", err)
    }

    // 7. 输出到 stdout
    fmt.Fprint(os.Stdout, string(markdown))

    return nil
}
```

### 5.5 配置加载 (`internal/config/load.go`)

```go
package config

import (
    "os"
    "strconv"
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
    }

    return cfg
}
```

---

## 6. 实现计划（TDD 顺序）

遵循测试先行原则，按以下顺序实现：

### Phase 1: URL 解析 (无外部依赖)
- [ ] `internal/parser/parser_test.go` - 表格驱动测试
- [ ] `internal/parser/parser.go` - 实现

### Phase 2: GitHub API 集成
- [ ] `internal/github/client_test.go` - 集成测试（使用测试仓库）
- [ ] `internal/github/client.go` - 实现
- [ ] `internal/github/fetch.go` - Issue/PR 获取
- [ ] `internal/github/fetch_discussion.go` - Discussion 获取（GraphQL）

### Phase 3: Markdown 转换
- [ ] `internal/converter/converter_test.go` - 表格驱动测试
- [ ] `internal/converter/converter.go` - 实现
- [ ] `internal/converter/format.go` - 格式化逻辑

### Phase 4: 配置与 CLI
- [ ] `internal/config/load_test.go` - 测试
- [ ] `internal/config/load.go` - 实现
- [ ] `internal/cli/cmd_test.go` - 端到端测试
- [ ] `internal/cli/cmd.go` - 实现

### Phase 5: 主程序
- [ ] `cmd/issue2md/main.go` - 入口点

---

## 7. 测试策略

### 7.1 单元测试

**适用范围**: 纯逻辑，无外部依赖
- `internal/parser` - URL 解析逻辑
- `internal/converter` - Markdown 格式化
- `internal/config` - 环境变量解析

**风格**: 表格驱动测试

### 7.2 集成测试

**适用范围**: 涉及外部 API 调用
- `internal/github` - GitHub API 交互

**策略**:
- 使用 GitHub 的真实 API（测试仓库）
- 或使用 `httptest` 模拟 GitHub API 响应

### 7.3 端到端测试

**适用范围**: 完整流程
- `internal/cli` - 完整的 URL → Markdown 流程

**策略**:
- 使用测试 Issue/PR/Discussion
- 验证输出格式正确性

---

## 8. 错误处理策略

所有错误必须显式处理，使用 `fmt.Errorf("...: %w", err)` 包装：

```go
// 示例
if err != nil {
    return fmt.Errorf("fetch issue: %w", err)
}
```

错误类型：
- 输入验证错误 → `ErrInvalidURLFormat`
- API 错误 → `ErrResourceNotFound`, `ErrAccessDenied`, `ErrRateLimited`
- 网络错误 → 包装原始错误

---

## 9. 验收标准

| # | 标准 | 验证方式 |
|---|------|----------|
| 1 | 支持三种 URL 类型解析 | 单元测试覆盖 |
| 2 | 正确获取 Issue/PR/Discussion 数据 | 集成测试 |
| 3 | 输出符合 spec.md 格式 | 端到端测试 |
| 4 | Token 认证正常工作 | 集成测试 |
| 5 | 错误处理清晰 | 手动测试 + 错误用例测试 |
| 6 | 代码覆盖率 >= 80% | `go test -cover` |

---

## 10. 附录

### 10.1 GraphQL 查询示例 (Discussion)

```graphql
query GetDiscussion($owner: String!, $repo: String!, $number: Int!) {
    repository(owner: $owner, name: $repo) {
        discussion(number: $number) {
            title
            body
            author { login }
            createdAt
            updatedAt
            stateReason
            labels(first: 10) { nodes { name } }
            comments(first: 100) {
                nodes {
                    id
                    body
                    author { login }
                    createdAt
                    updatedAt
                    replyTo { id }
                }
            }
        }
    }
}
```

### 10.2 参考文档

- [GitHub REST API v3](https://docs.github.com/en/rest)
- [GitHub GraphQL API v4](https://docs.github.com/en/graphql)
- [google/go-github](https://github.com/google/go-github)
- [Go Package Naming Conventions](https://go.dev/doc/effective_go#names)
