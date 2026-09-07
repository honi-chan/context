package repository

import (
	"context"

	"github.com/yourname/ai-work-assistant/internal/domain"
)

// GitHubRepository は、GitHubから必要な情報を取得するための
// インターフェース。
//
// 呼び出し側はGitHub APIの具体的な実装を知る必要がない。
// infrastructure層がこのinterfaceを実装する。
type GitHubRepository interface {

	// ListPullRequests は指定したリポジトリの
	// Pull Request一覧を取得する。
	ListPullRequests(
		ctx context.Context,
		owner string,
		repo string,
	) ([]domain.PullRequest, error)
}
