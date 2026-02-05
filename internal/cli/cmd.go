package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/hirolittle/issue2md/internal/config"
	"github.com/hirolittle/issue2md/internal/converter"
	"github.com/hirolittle/issue2md/internal/github"
	"github.com/hirolittle/issue2md/internal/parser"
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
