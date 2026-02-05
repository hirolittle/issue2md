package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-github/v69/github"
	"golang.org/x/oauth2"

	"github.com/hirolittle/issue2md/internal/parser"
)

// Client 是 GitHub API 的客户端接口
type Client interface {
	// Fetch 根据 ResourceID 获取完整的资源数据
	Fetch(ctx context.Context, id parser.ResourceID) (*Resource, error)
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

// newClientWithHTTP 创建一个使用自定义 HTTP 客户端的 GitHub 客户端
// 用于测试，支持 httptest.Server
func newClientWithHTTP(httpClient *http.Client, baseURL string) Client {
	ghClient := github.NewClient(httpClient)
	if baseURL != "" {
		// BaseURL 必须以 / 结尾
		if baseURL[len(baseURL)-1] != '/' {
			baseURL += "/"
		}
		parsedURL, err := url.Parse(baseURL)
		if err == nil {
			ghClient.BaseURL = parsedURL
		}
	}
	return &client{
		gh: ghClient,
	}
}

// Fetch 实现 Client 接口
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

// fetchIssueOrPR 获取 Issue 或 PR 数据
func (c *client) fetchIssueOrPR(ctx context.Context, id parser.ResourceID) (*Resource, error) {
	issue, resp, err := c.gh.Issues.Get(ctx, id.Owner, id.Repository, int(id.Number))
	if err != nil {
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusNotFound:
				return nil, fmt.Errorf("%w: %w", ErrResourceNotFound, err)
			case http.StatusForbidden, http.StatusUnauthorized:
				return nil, fmt.Errorf("%w: %w", ErrAccessDenied, err)
			}
		}
		return nil, fmt.Errorf("get issue: %w", err)
	}

	// 获取评论
	comments := []Comment{}
	issueComments, _, err := c.gh.Issues.ListComments(ctx, id.Owner, id.Repository, int(id.Number), nil)
	if err == nil {
		for _, ic := range issueComments {
			inReplyTo := int64(0)
			// go-github v69: InReplyTo 字段可能在 IssueComment 中
			// 如果不存在，默认为 0（顶级评论）
			comments = append(comments, Comment{
				ID:        ic.GetID(),
				Author:    ic.GetUser().GetLogin(),
				Body:      ic.GetBody(),
				CreatedAt: ic.GetCreatedAt().Time,
				UpdatedAt: ic.GetUpdatedAt().Time,
				InReplyTo: inReplyTo,
			})
		}
	}

	// 解析标签
	labels := make([]string, 0, len(issue.Labels))
	for _, label := range issue.Labels {
		labels = append(labels, label.GetName())
	}

	// 确定 PR 的状态（merged/closed/open）
	state := issue.GetState()
	// 对于 PR，需要检查是否已合并
	// 由于 Issue API 不提供合并状态，这里暂时使用 issue 的状态
	// TODO: 对于 PR 类型，调用 PullRequests API 获取合并状态

	return &Resource{
		ID:        id,
		Title:     issue.GetTitle(),
		Body:      issue.GetBody(),
		Author:    issue.GetUser().GetLogin(),
		State:     state,
		Labels:    labels,
		CreatedAt: issue.GetCreatedAt().Time,
		UpdatedAt: issue.GetUpdatedAt().Time,
		URL:       issue.GetHTMLURL(),
		Comments:  comments,
	}, nil
}

// fetchDiscussion 获取 Discussion 数据（使用 GraphQL）
func (c *client) fetchDiscussion(ctx context.Context, id parser.ResourceID) (*Resource, error) {
	// TODO: 实现 GraphQL 查询
	return nil, fmt.Errorf("discussions not yet implemented")
}
