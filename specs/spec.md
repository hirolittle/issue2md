# issue2md 功能规格说明书

## 1. 产品概述

**issue2md** 是一个简洁的命令行工具，用于将 GitHub Issue、Pull Request 或 Discussion 转换为 Markdown 格式。

**核心使用场景**：个人归档 - 保存有价值的 GitHub 内容供个人参考。

---

## 2. 核心功能

### 2.1 基本用法

```bash
GITHUB_TOKEN=ghp_xxx issue2md https://github.com/owner/repo/issues/123 > issue-123.md
```

### 2.2 输入

- **格式**：仅支持完整的 GitHub URL
- **支持的资源类型**：
  - Issue: `https://github.com/owner/repo/issues/123`
  - Pull Request: `https://github.com/owner/repo/pull/456`
  - Discussion: `https://github.com/owner/repo/discussions/789`

### 2.3 输出

- **方式**：输出到 stdout（支持管道和重定向）
- **文件命名**（建议）：`issue-123.md` / `pr-456.md` / `discussion-789.md`

---

## 3. 输出格式

### 3.1 结构

```markdown
---
title: <标题>
author: <username>
status: <open|closed|merged|etc.>
labels: [<label1>, <label2>]
created_at: <RFC3339 timestamp>
updated_at: <RFC3339 timestamp>
url: <original_url>
type: <issue|pull_request|discussion>
---

# <标题>

<正文内容>

## Comments

### @<author1> - <timestamp>

<评论内容>

### @<author2> - <timestamp>

<评论内容>

> **Reply to @<author1>**
>
> <回复内容>
```

### 3.2 元数据字段（YAML Front Matter）

| 字段 | Issue | PR | Discussion | 说明 |
|------|-------|----|------------|------|
| title | ✓ | ✓ | ✓ | 标题 |
| author | ✓ | ✓ | ✓ | 作者用户名 |
| status | ✓ | ✓ | ✓ | 状态 |
| labels | ✓ | ✓ | ✓ | 标签列表 |
| created_at | ✓ | ✓ | ✓ | 创建时间（RFC3339） |
| updated_at | ✓ | ✓ | ✓ | 更新时间（RFC3339） |
| url | ✓ | ✓ | ✓ | 原始 URL |
| type | ✓ | ✓ | ✓ | 资源类型 |
| number | ✓ | ✓ | ✓ | Issue/PR/Discussion 编号 |
| repository | ✓ | ✓ | ✓ | owner/repo |

---

## 4. 认证

### 4.1 GitHub Token

- **传递方式**：通过环境变量 `GITHUB_TOKEN`
- **必要性**：可选
  - 无 Token：受 GitHub 匿名速率限制（60 次/小时）
  - 有 Token：提升至 5000 次/小时

### 4.2 Token 权限要求

最小权限：`public_repo`（仅访问公开仓库）

---

## 5. 技术约束

### 5.1 开发原则

- 遵循 Go 语言"少即是多"哲学
- 优先使用标准库
- 不进行不必要的抽象

### 5.2 依赖

- Go >= 1.24
- 仅使用必要的第三方依赖（如 GitHub API 客户端）

### 5.3 测试

- 所有功能必须先编写测试（TDD）
- 优先使用表格驱动测试
- 优先集成测试，避免 Mock

---

## 6. 边缘情况处理

### 6.1 输入验证

- 无效 URL → 返回错误
- 不存在的资源 → 返回错误
- 私有仓库且无 Token → 返回错误
- 网络/超时错误 → 返回错误

### 6.2 内容处理

- 空正文 → 仍生成有效 Markdown
- 无评论 → 省略 `## Comments` 章节
- 删除的内容 → 显示 `[deleted]` 标记
- Markdown 格式 → 保留原始格式（不转义）

### 6.3 特殊字符

- Issue 标题中的特殊字符 → 保留原样，在 YAML 中用引号包裹
- 用户名包含特殊字符 → 保留原样

---

## 7. 非功能需求

### 7.1 性能

- 单次请求超时：30 秒

### 7.2 兼容性

- 支持 macOS、Linux
- Go 1.24+

### 7.3 用户体验

- 错误信息清晰、可操作
- 进度反馈（可选，如获取大 Discussion 时）

---

## 8. 未来版本（Out of Scope）

以下功能**明确不在 v1 范围内**：

- 批量转换（URL 列表文件、GitHub Search query）
- 自定义输出路径/文件名
- PR diff 包含
- Discussion 特殊类型处理
- OAuth 认证流程
- 文件覆盖策略（工具只输出 stdout，不直接写文件）
- 短 URL 格式支持（如 `owner/repo#123`）
- 交互式确认

---

## 9. 验收标准

功能完成的标准：

1. 能够成功解析并转换 Issue、PR、Discussion 三种类型的 URL
2. 输出的 Markdown 包含完整的 YAML front matter 和正文内容
3. 评论及嵌套回复正确渲染
4. 有 Token 时能正常认证，无 Token 时在限额内正常工作
5. 所有边缘情况有合理的错误处理
6. 代码覆盖率 >= 80%
