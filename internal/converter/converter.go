package converter

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hirolittle/issue2md/internal/github"
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

// Convert 实现 Converter 接口
func (c *converter) Convert(ctx context.Context, resource *github.Resource) ([]byte, error) {
	var buf strings.Builder

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

	return []byte(buf.String()), nil
}

// writeFrontMatter 写入 YAML Front Matter
func (c *converter) writeFrontMatter(buf *strings.Builder, r *github.Resource) {
	buf.WriteString("---\n")
	fmt.Fprintf(buf, "title: %q\n", r.Title)
	fmt.Fprintf(buf, "author: %s\n", r.Author)
	fmt.Fprintf(buf, "state: %s\n", r.State)
	fmt.Fprintf(buf, "url: %s\n", r.URL)
	fmt.Fprintf(buf, "created_at: %s\n", r.CreatedAt.Format(time.RFC3339))
	fmt.Fprintf(buf, "updated_at: %s\n", r.UpdatedAt.Format(time.RFC3339))

	// 标签列表
	if len(r.Labels) > 0 {
		buf.WriteString("labels:\n")
		for _, label := range r.Labels {
			fmt.Fprintf(buf, "  - %s\n", label)
		}
	} else {
		buf.WriteString("labels: []\n")
	}

	buf.WriteString("---\n\n")
}

// writeComments 写入评论
func (c *converter) writeComments(buf *strings.Builder, comments []github.Comment) {
	// 按创建时间排序
	sorted := make([]github.Comment, len(comments))
	copy(sorted, comments)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CreatedAt.Before(sorted[j].CreatedAt)
	})

	buf.WriteString("## Comments\n\n")

	// 构建评论 ID 到作者的映射
	authorMap := make(map[int64]string)
	for _, comment := range sorted {
		authorMap[comment.ID] = comment.Author
	}

	for _, comment := range sorted {
		if comment.InReplyTo > 0 {
			// 嵌套回复格式
			targetAuthor := findAuthor(comment.InReplyTo, authorMap)
			buf.WriteString(fmt.Sprintf("> **Reply to @%s**\n>\n", targetAuthor))
			lines := strings.Split(comment.Body, "\n")
			for _, line := range lines {
				buf.WriteString(fmt.Sprintf("> %s\n", line))
			}
			buf.WriteString("\n")
		} else {
			// 顶级评论格式
			fmt.Fprintf(buf, "### @%s - %s\n\n", comment.Author, comment.CreatedAt.Format(time.RFC3339))
			fmt.Fprintf(buf, "%s\n\n", comment.Body)
		}
	}
}

// findAuthor 查找评论作者
func findAuthor(id int64, authorMap map[int64]string) string {
	if author, ok := authorMap[id]; ok {
		return author
	}
	return "unknown"
}
