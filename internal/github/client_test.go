package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hirolittle/issue2md/internal/parser"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "empty token uses default client",
			token: "",
		},
		{
			name:  "non-empty token uses oauth2 client",
			token: "test-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.token)
			if client == nil {
				t.Errorf("NewClient() returned nil")
			}
		})
	}
}

func TestClient_Fetch_Issue(t *testing.T) {
	// 创建 Mock Server
	requestCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		// 第一个请求：获取 Issue
		if requestCount == 1 && r.URL.Path == "/repos/owner/repo/issues/123" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": 1,
				"number": 123,
				"title": "Test Issue",
				"body": "This is a test issue",
				"user": {
					"login": "testuser"
				},
				"state": "open",
				"labels": [
					{"name": "bug"},
					{"name": "enhancement"}
				],
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-02T00:00:00Z",
				"html_url": "https://github.com/owner/repo/issues/123",
				"pull_request": null,
				"comments": 1
			}`))
			return
		}

		// 第二个请求：获取评论列表
		if requestCount == 2 && r.URL.Path == "/repos/owner/repo/issues/123/comments" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"id": 456,
					"user": {"login": "commenter"},
					"body": "This is a comment",
					"created_at": "2024-01-01T01:00:00Z",
					"updated_at": "2024-01-01T01:00:00Z",
					"in_reply_to": null
				}
			]`))
			return
		}

		// 未知请求
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer ts.Close()

	// 使用 Mock Server 的 HTTP 客户端
	client := newClientWithHTTP(ts.Client(), ts.URL)

	ctx := context.Background()
	id := parser.ResourceID{
		Owner:      "owner",
		Repository: "repo",
		Number:     123,
		Type:       parser.ResourceTypeIssue,
	}

	got, err := client.Fetch(ctx, id)
	if err != nil {
		t.Fatalf("Fetch() unexpected error: %v", err)
	}

	if got == nil {
		t.Fatalf("Fetch() returned nil")
	}

	// 验证返回的数据
	want := newTestResource()
	if got.Title != want.Title {
		t.Errorf("Fetch().Title = %q, want %q", got.Title, want.Title)
	}
	if got.Body != want.Body {
		t.Errorf("Fetch().Body = %q, want %q", got.Body, want.Body)
	}
	if got.Author != want.Author {
		t.Errorf("Fetch().Author = %q, want %q", got.Author, want.Author)
	}
	if got.State != want.State {
		t.Errorf("Fetch().State = %q, want %q", got.State, want.State)
	}
	if len(got.Labels) != len(want.Labels) {
		t.Errorf("Fetch().Labels length = %d, want %d", len(got.Labels), len(want.Labels))
	} else {
		for i, label := range got.Labels {
			if label != want.Labels[i] {
				t.Errorf("Fetch().Labels[%d] = %q, want %q", i, label, want.Labels[i])
			}
		}
	}

	// 验证评论
	if len(got.Comments) != len(want.Comments) {
		t.Errorf("Fetch().Comments length = %d, want %d", len(got.Comments), len(want.Comments))
	} else if len(got.Comments) > 0 {
		if got.Comments[0].Author != want.Comments[0].Author {
			t.Errorf("Fetch().Comments[0].Author = %q, want %q", got.Comments[0].Author, want.Comments[0].Author)
		}
		if got.Comments[0].Body != want.Comments[0].Body {
			t.Errorf("Fetch().Comments[0].Body = %q, want %q", got.Comments[0].Body, want.Comments[0].Body)
		}
	}
}

func TestClient_Fetch_ErrorCases(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
	}{
		{
			name:       "resource not found - 404",
			statusCode: http.StatusNotFound,
			response:   `{"message": "Not Found"}`,
		},
		{
			name:       "access denied - 403",
			statusCode: http.StatusForbidden,
			response:   `{"message": "Forbidden"}`,
		},
		{
			name:       "unauthorized - 401",
			statusCode: http.StatusUnauthorized,
			response:   `{"message": "Unauthorized"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer ts.Close()

			client := newClientWithHTTP(ts.Client(), ts.URL)
			ctx := context.Background()
			id := parser.ResourceID{
				Owner:      "owner",
				Repository: "repo",
				Number:     123,
				Type:       parser.ResourceTypeIssue,
			}

			_, err := client.Fetch(ctx, id)
			if err == nil {
				t.Errorf("Fetch() expected error, got nil")
			}
		})
	}
}

// 辅助函数：创建用于测试的 Resource
func newTestResource() *Resource {
	return &Resource{
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
		Comments: []Comment{
			{
				ID:        456,
				Author:    "commenter",
				Body:      "This is a comment",
				CreatedAt: time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
				InReplyTo: 0,
			},
		},
	}
}
