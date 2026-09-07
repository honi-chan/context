package usecase

import (
	"context"
	"fmt"

	"github.com/yourname/ai-work-assistant/internal/ai"
	"github.com/yourname/ai-work-assistant/internal/repository"
)

// AnalyzePullRequestsUseCase は、
//
// 1. GitHubからPRを取得
// 2. AIに分析させる
//
// という一連のユースケースを担当する。
//
// GitHub APIやGemini APIの具体的な実装については知らない。
type AnalyzePullRequestsUseCase struct {
	githubRepository repository.GitHubRepository
	analyzer         ai.Analyzer
}

// NewAnalyzePullRequestsUseCase はUseCaseを生成する。
//
// GitHubRepositoryとAnalyzerを外から渡すことで、
// 具体的な実装への依存を避ける。
func NewAnalyzePullRequestsUseCase(
	githubRepository repository.GitHubRepository,
	analyzer ai.Analyzer,
) *AnalyzePullRequestsUseCase {
	return &AnalyzePullRequestsUseCase{
		githubRepository: githubRepository,
		analyzer:         analyzer,
	}
}

// Execute はPRを取得してAIに分析させる。
func (u *AnalyzePullRequestsUseCase) Execute(
	ctx context.Context,
	owner string,
	repo string,
) (string, error) {

	// GitHubから現在OpenになっているPR一覧を取得する。
	pullRequests, err := u.githubRepository.ListPullRequests(
		ctx,
		owner,
		repo,
	)
	if err != nil {
		return "", fmt.Errorf("list pull requests: %w", err)
	}

	// PRが存在しない場合は、
	// Gemini APIを無駄に呼び出さずここで終了する。
	if len(pullRequests) == 0 {
		return "現在対応が必要なPull Requestはありません。", nil
	}

	// PR一覧をAIへ渡して分析する。
	result, err := u.analyzer.AnalyzePullRequests(
		ctx,
		pullRequests,
	)
	if err != nil {
		return "", fmt.Errorf("analyze pull requests: %w", err)
	}

	return result, nil
}
