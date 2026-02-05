# issue2md 开发任务列表

> **版本**: v1.0
> **生成日期**: 2025-02-05
> **基于**: `plan.md` v1.0

---

## 任务说明

### 符号约定
- `[P]` - 可并行执行的任务（无相互依赖）
- `T-xxx` - 测试任务
- `I-xxx` - 实现任务
- `ID: xxx` - 任务唯一标识符

### 执行原则
1. **TDD 强制**：每个功能模块必须先完成测试任务（T），再完成实现任务（I）
2. **依赖关系**：任务必须按依赖顺序执行
3. **原子化**：每个任务只涉及一个主要文件的修改或创建

---

## Phase 1: Foundation (数据结构定义)

> 本阶段定义跨包共享的核心数据类型，无外部依赖，可快速完成。

### 1.1 Parser 类型定义

**ID**: `PH1-001-T`
**任务**: 编写 `internal/parser/types_test.go`
**描述**: 定义 ResourceType 和 ResourceID 的单元测试
**测试用例**:
- ResourceType 常量值验证
- ResourceID 结构体字段验证
**依赖**: 无
**并行**: [P]

---

**ID**: `PH1-001-I`
**任务**: 创建 `internal/parser/types.go`
**描述**: 实现 ResourceType 和 ResourceID 类型定义
**内容**:
```go
package parser

type ResourceType string

const (
    ResourceTypeIssue       ResourceType = "issue"
    ResourceTypePullRequest ResourceType = "pull_request"
    ResourceTypeDiscussion  ResourceType = "discussion"
)

type ResourceID struct {
    Owner      string
    Repository string
    Number     int64
    Type       ResourceType
}
```
**依赖**: `PH1-001-T`
**并行**: -

### 1.2 GitHub 类型定义

**ID**: `PH1-002-T`
**任务**: 编写 `internal/github/types_test.go`
**描述**: 定义 Resource 和 Comment 结构体的单元测试
**测试用例**:
- Resource 结构体字段验证
- Comment 结构体字段验证
- InReplyTo 字段语义验证（0 表示顶级评论）
**依赖**: `PH1-001-I`
**并行**: [P]

---

**ID**: `PH1-002-I`
**任务**: 创建 `internal/github/types.go`
**描述**: 实现 Resource 和 Comment 数据模型
**内容**:
```go
package github

import "time"

type Resource struct {
    ID          parser.ResourceID
    Title       string
    Body        string
    Author      string
    State       string
    Labels      []string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    URL         string
    Comments    []Comment
}

type Comment struct {
    ID        int64
    Author    string
    Body      string
    CreatedAt time.Time
    UpdatedAt time.Time
    InReplyTo int64
}
```
**依赖**: `PH1-002-T`
**并行**: -

### 1.3 GitHub 错误定义

**ID**: `PH1-003-T`
**任务**: 编写 `internal/github/errors_test.go`
**描述**: 验证错误变量正确定义
**测试用例**:
- 错误变量非空验证
- 错误信息格式验证
**依赖**: 无
**并行**: [P]

---

**ID**: `PH1-003-I`
**任务**: 创建 `internal/github/errors.go`
**描述**: 定义 GitHub 包的错误变量
**内容**:
```go
package github

import "errors"

var (
    ErrResourceNotFound = errors.New("resource not found")
    ErrAccessDenied     = errors.New("access denied: check token permissions")
    ErrRateLimited      = errors.New("rate limited: please provide a token")
)
```
**依赖**: `PH1-003-T`
**并行**: -

### 1.4 Config 结构定义

**ID**: `PH1-004-T`
**任务**: 编写 `internal/config/config_test.go`
**描述**: 定义 Config 结构体的单元测试
**测试用例**:
- Config 结构体字段验证
- 默认值验证
**依赖**: 无
**并行**: [P]

---

**ID**: `PH1-004-I`
**任务**: 创建 `internal/config/config.go`
**描述**: 实现 Config 配置结构
**内容**:
```go
package config

import "time"

type Config struct {
    GitHubToken string
    Timeout     time.Duration
}
```
**依赖**: `PH1-004-T`
**并行**: -

### 1.5 Parser 错误定义

**ID**: `PH1-005-T`
**任务**: 编写 `internal/parser/errors_test.go`
**描述**: 验证 Parser 错误变量正确定义
**测试用例**:
- ErrInvalidURLFormat 验证
- ErrUnsupportedType 验证
**依赖**: 无
**并行**: [P]

---

**ID**: `PH1-005-I`
**任务**: 更新 `internal/parser/types.go`
**描述**: 添加 Parser 包的错误变量
**内容**: 在 types.go 中添加
```go
import "errors"

var (
    ErrInvalidURLFormat = errors.New("invalid url format")
    ErrUnsupportedType  = errors.New("unsupported resource type")
)
```
**依赖**: `PH1-005-T`, `PH1-001-I`
**并行**: -

---

## Phase 2: GitHub Fetcher (API 交互逻辑)

> 本阶段实现 GitHub API 交互，使用 `google/go-github` 库。

### 2.1 URL 解析器

**ID**: `PH2-001-T`
**任务**: 编写 `internal/parser/parser_test.go`
**描述**: ParseURL 函数的表格驱动测试
**测试用例**:
| 输入 | 预期 ResourceID | 预期错误 |
|------|----------------|----------|
| `https://github.com/owner/repo/issues/123` | owner/repo/123/issue | nil |
| `https://github.com/owner/repo/pull/456` | owner/repo/456/pull_request | nil |
| `https://github.com/owner/repo/discussions/789` | owner/repo/789/discussion | nil |
| `https://gitlab.com/owner/repo/issues/123` | - | ErrInvalidURLFormat |
| `https://github.com/owner/repo/invalid/123` | - | ErrUnsupportedType |
| `not-a-url` | - | parse error |
| `https://github.com/owner/repo/issues/abc` | - | parse error |
| `https://github.com/owner/repo/issues` | - | ErrInvalidURLFormat |
**依赖**: `PH1-005-I`
**并行**: -

---

**ID**: `PH2-001-I`
**任务**: 创建 `internal/parser/parser.go`
**描述**: 实现 ParseURL 函数
**功能**:
- 使用 net/url 解析 URL
- 验证域名为 github.com
- 解析路径: /owner/repo/type/number
- 根据 type 映射到 ResourceType
- 返回 ResourceID 或错误
**错误处理**: 使用 `fmt.Errorf("...: %w", err)` 包装错误
**依赖**: `PH2-001-T`
**并行**: -

### 2.2 GitHub Client 初始化

**ID**: `PH2-002-T`
**任务**: 编写 `internal/github/client_test.go` - NewClient 测试
**描述**: 测试 GitHub 客户端初始化
**测试用例**:
| Token 输入 | 预期行为 |
|-----------|----------|
| 空字符串 | 使用默认 http.Client |
| 非空字符串 | 使用 oauth2 配置的 http.Client |
**依赖**: `PH1-003-I`, `PH2-001-I`
**并行**: -

---

**ID**: `PH2-002-I`
**任务**: 创建 `internal/github/client.go`
**描述**: 实现 Client 接口和 NewClient 函数
**内容**:
```go
type Client interface {
    Fetch(ctx context.Context, id parser.ResourceID) (*Resource, error)
}

type client struct {
    gh *github.Client
}

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
```
**依赖**: `PH2-002-T`
**并行**: -

### 2.3 Issue/PR 获取逻辑

**ID**: `PH2-003-T`
**任务**: 编写 `internal/github/fetch_test.go` - fetchIssueOrPR 测试
**描述**: 测试 Issue 和 PR 的数据获取
**策略**:
- 使用 httptest 模拟 GitHub API 响应
- 测试 Issue 响应解析
- 测试 PR 响应解析
- 测试评论获取（包括嵌套关系）
- 测试错误场景（404, 403, 401）
**依赖**: `PH2-002-I`
**并行**: -

---

**ID**: `PH2-003-I`
**任务**: 创建 `internal/github/fetch.go`
**描述**: 实现 fetchIssueOrPR 函数
**功能**:
- 根据 ResourceID 调用 GitHub Issues API
- 解析 Issue/PR 基本信息到 Resource
- 获取评论列表并解析为 Comment
- 映射 InReplyTo 关系
- 处理状态: open/closed/merged
**依赖**: `PH2-003-T`
**并行**: -

### 2.4 Discussion 获取逻辑

**ID**: `PH2-004-T`
**任务**: 编写 `internal/github/fetch_discussion_test.go`
**描述**: 测试 Discussion 的数据获取（GraphQL）
**策略**:
- 使用 httptest 模拟 GitHub GraphQL API 响应
- 测试 Discussion 响应解析
- 测试评论和嵌套回复解析
- 测试错误场景
**依赖**: `PH2-002-I`
**并行**: [P] 与 PH2-003 并行

---

**ID**: `PH2-004-I`
**任务**: 创建 `internal/github/fetch_discussion.go`
**描述**: 实现 fetchDiscussion 函数（GraphQL）
**功能**:
- 构造 GraphQL 查询
- 调用 GitHub GraphQL API
- 解析 Discussion 基本信息
- 解析评论和回复（replyTo 关系）
- 映射到 Resource 和 Comment
**依赖**: `PH2-004-T`
**并行**: -

### 2.5 Client Fetch 实现

**ID**: `PH2-005-T`
**任务**: 更新 `internal/github/client_test.go` - Fetch 测试
**描述**: 测试 Fetch 方法的路由逻辑
**测试用例**:
- Issue 类型 → 调用 fetchIssueOrPR
- PR 类型 → 调用 fetchIssueOrPR
- Discussion 类型 → 调用 fetchDiscussion
- 未知类型 → 返回错误
**依赖**: `PH2-003-I`, `PH2-004-I`
**并行**: -

---

**ID**: `PH2-005-I`
**任务**: 更新 `internal/github/client.go`
**描述**: 实现 Fetch 方法
**内容**:
```go
func (c *client) Fetch(ctx context.Context, id parser.ResourceID) (*Resource, error) {
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
**依赖**: `PH2-005-T`
**并行**: -

---

## Phase 3: Markdown Converter (转换逻辑)

> 本阶段实现 Resource 到 Markdown 的转换，无外部依赖。

### 3.1 Converter 接口定义

**ID**: `PH3-001-T`
**任务**: 编写 `internal/converter/converter_test.go` - 接口测试
**描述**: 定义 Converter 接口的测试框架
**测试用例**:
- NewConverter 返回非 nil
- 返回值实现 Converter 接口
**依赖**: `PH2-005-I`
**并行**: -

---

**ID**: `PH3-001-I`
**任务**: 创建 `internal/converter/converter.go`
**描述**: 定义 Converter 接口和基本结构
**内容**:
```go
type Converter interface {
    Convert(ctx context.Context, resource *github.Resource) ([]byte, error)
}

type converter struct{}

func NewConverter() Converter {
    return &converter{}
}
```
**依赖**: `PH3-001-T`
**并行**: -

### 3.2 YAML Front Matter 生成

**ID**: `PH3-002-T`
**任务**: 编写 `internal/converter/format_test.go` - writeFrontMatter 测试
**描述**: 测试 YAML Front Matter 生成
**测试用例**:
- 完整字段生成
- 空标签处理
- 特殊字符转义（标题含引号）
**依赖**: `PH3-001-I`
**并行**: -

---

**ID**: `PH3-002-I`
**任务**: 创建 `internal/converter/format.go`
**描述**: 实现 writeFrontMatter 函数
**功能**:
- 生成 YAML front matter 格式
- 包含所有 spec.md 要求的字段
- 正确处理字符串引号
- 标签列表格式化
**依赖**: `PH3-002-T`
**并行**: -

### 3.3 评论格式化

**ID**: `PH3-003-T`
**任务**: 编写 `internal/converter/format_test.go` - writeComments 测试
**描述**: 测试评论格式化
**测试用例**:
| 场景 | 预期输出 |
|------|----------|
| 无评论 | 不输出 Comments 章节 |
| 单条顶级评论 | ### @author - timestamp |
| 嵌套回复 | > **Reply to @author** |
| 多条评论按时间排序 | 时间顺序正确 |
**依赖**: `PH3-002-I`
**并行**: -

---

**ID**: `PH3-003-I`
**任务**: 更新 `internal/converter/format.go`
**描述**: 实现 writeComments 函数
**功能**:
- 按 CreatedAt 排序评论
- 生成顶级评论格式
- 生成嵌套回复引用格式
- 查找回复目标作者
**依赖**: `PH3-003-T`
**并行**: -

### 3.4 Convert 方法实现

**ID**: `PH3-004-T`
**任务**: 更新 `internal/converter/converter_test.go` - Convert 集成测试
**描述**: 测试完整的 Convert 流程
**测试用例**:
| 场景 | 验证点 |
|------|--------|
| 完整 Resource（含评论） | 输出包含 front matter + 正文 + 评论 |
| 空 Body | 仍生成有效 Markdown |
| 无标签 | labels 为空列表 |
| 多个标签 | YAML 中数组格式 |
| 嵌套评论 | 引用格式正确 |
**依赖**: `PH3-003-I`
**并行**: -

---

**ID**: `PH3-004-I`
**任务**: 更新 `internal/converter/converter.go`
**描述**: 实现 Convert 方法
**功能**:
- 调用 writeFrontMatter
- 写入标题 (# title)
- 写入正文（如果有）
- 调用 writeComments（如果有评论）
- 返回字节数组
**依赖**: `PH3-004-T`
**并行**: -

---

## Phase 4: CLI Assembly (命令行入口集成)

> 本阶段组装各模块，实现完整的 CLI 工具。

### 4.1 Config 加载逻辑

**ID**: `PH4-001-T`
**任务**: 编写 `internal/config/load_test.go`
**描述**: 测试配置加载逻辑
**测试用例**:
| 环境变量设置 | 预期 Config |
|-------------|-------------|
| 无环境变量 | Token="", Timeout=30s |
| GITHUB_TOKEN=xxx | Token=xxx, Timeout=30s |
| ISSUE2MD_TIMEOUT=60s | Token="", Timeout=60s |
| 两者都设置 | Token=xxx, Timeout=60s |
| ISSUE2MD_TIMEOUT=invalid | Token="", Timeout=30s（默认值）|
**依赖**: `PH1-004-I`
**并行**: [P] 与 Phase 3 任务并行

---

**ID**: `PH4-001-I`
**任务**: 创建 `internal/config/load.go`
**描述**: 实现 Load 函数
**功能**:
- 读取 GITHUB_TOKEN 环境变量
- 读取 ISSUE2MD_TIMEOUT 环境变量
- 解析超时值，失败时使用默认值 30s
**依赖**: `PH4-001-T`
**并行**: -

### 4.2 CLI Run 函数

**ID**: `PH4-002-T`
**任务**: 编写 `internal/cli/cmd_test.go`
**描述**: 测试 CLI Run 函数（端到端测试）
**策略**:
- 使用 httptest 模拟 GitHub API
- 测试完整流程: ParseURL → Fetch → Convert → Output
- 测试错误场景（URL 解析失败、API 错误、转换错误）
**依赖**: `PH3-004-I`, `PH4-001-I`
**并行**: -

---

**ID**: `PH4-002-I`
**任务**: 创建 `internal/cli/cmd.go`
**描述**: 实现 Run 函数
**功能**:
1. 加载配置 (config.Load)
2. 创建带超时的上下文
3. 解析 URL (parser.ParseURL)
4. 创建 GitHub 客户端 (github.NewClient)
5. 获取数据 (client.Fetch)
6. 转换为 Markdown (converter.Convert)
7. 输出到 stdout
**错误处理**: 每一步都显式处理并包装错误
**依赖**: `PH4-002-T`
**并行**: -

### 4.3 主程序入口

**ID**: `PH4-003-T`
**任务**: 编写 `cmd/issue2md/main_test.go`
**描述**: 主程序基础测试
**测试用例**:
- main 函数可正常退出
- 参数解析验证
**依赖**: `PH4-002-I`
**并行**: -

---

**ID**: `PH4-003-I`
**任务**: 创建 `cmd/issue2md/main.go`
**描述**: 实现主程序入口
**功能**:
- 解析命令行参数（URL）
- 调用 cli.Run
- 处理错误并输出到 stderr
- 设置正确的退出码
**依赖**: `PH4-003-T`
**并行**: -

---

## 任务依赖图

```
Phase 1: Foundation
├── PH1-001-T ──────────────────────┐
├── PH1-001-I (types.go) ───────────┤
│                                   ▼
├── PH1-002-T ────────────────► PH1-002-I (github/types.go) ──────┐
│                                                                   │
├── PH1-003-T ────────────────► PH1-003-I (errors.go)              │
│                                                                   │
├── PH1-004-T ────────────────► PH1-004-I (config.go)              │
│                                                                   │
├── PH1-005-T ────────────────► PH1-005-I (parser/errors.go) ──────┤
│                                                                   │
└───────────────────────────────────────────────────────────────────┤
                                                                    │
Phase 2: GitHub Fetcher                                             │
├── PH2-001-T ────────────────► PH2-001-I (parser.go)               │
│                                                                   │
├── PH2-002-T ────────────────► PH2-002-I (client.go)               │
│                                                                   │
├── PH2-003-T ────────────────► PH2-003-I (fetch.go) ───────────────┤
│                                                                   │
├── PH2-004-T ────────────────► PH2-004-I (fetch_discussion.go) ────┤
│                                                                   │
├── PH2-005-T ────────────────► PH2-005-I (client.Fetch()) ─────────┤
│                                                                   │
└───────────────────────────────────────────────────────────────────┤
                                                                    │
Phase 3: Markdown Converter                                         │
├── PH3-001-T ────────────────► PH3-001-I (converter.go)            │
│                                                                   │
├── PH3-002-T ────────────────► PH3-002-I (format.go - front matter)│
│                                                                   │
├── PH3-003-T ────────────────► PH3-003-I (format.go - comments)    │
│                                                                   │
├── PH3-004-T ────────────────► PH3-004-I (converter.Convert())     │
│                                                                   │
└───────────────────────────────────────────────────────────────────┤
                                                                    │
Phase 4: CLI Assembly                                               │
├── PH4-001-T ────────────────► PH4-001-I (config/load.go)          │
│                                                                   │
├── PH4-002-T ────────────────► PH4-002-I (cli/cmd.go) ─────────────┤
│                                                                   │
└── PH4-003-T ────────────────► PH4-003-I (cmd/issue2md/main.go) ───┘
```

---

## 并行执行建议

### 可并行的任务组
1. **Foundation 阶段**: 所有 `PH1-xxx-T` 测试任务可并行
2. **GitHub Fetcher 阶段**: `PH2-003` 和 `PH2-004` 可并行
3. **Config 加载**: `PH4-001` 可与 Phase 3 任务并行
4. **跨阶段**: `PH2-003` (fetch.go) 可与 `PH3-001` (converter 初始化) 并行

---

## 验收检查清单

完成所有任务后，验证以下功能：

- [ ] 能正确解析 Issue URL (`.../issues/123`)
- [ ] 能正确解析 PR URL (`.../pull/456`)
- [ ] 能正确解析 Discussion URL (`.../discussions/789`)
- [ ] 无效 URL 返回清晰错误
- [ ] 无 Token 时能在限额内正常工作
- [ ] 有 Token 时能正确认证
- [ ] Issue 数据正确转换为 Markdown
- [ ] PR 数据正确转换为 Markdown
- [ ] Discussion 数据正确转换为 Markdown
- [ ] 嵌套评论正确渲染
- [ ] YAML front matter 包含所有必需字段
- [ ] 输出到 stdout 正常工作
- [ ] 超时控制有效（默认 30s，可通过环境变量覆盖）
- [ ] 错误信息清晰、可操作
- [ ] 代码覆盖率 >= 80%

---

## 附录：go.mod 依赖

```go
module github.com/yourusername/issue2md

go 1.24

require (
    github.com/google/go-github/v69 v69.0.0
    golang.org/x/oauth2 v0.24.0
)
```

执行 `go mod tidy` 前请确保已创建上述 `go.mod` 文件。
