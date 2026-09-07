package ai

import (
	"context"

	"github.com/yourname/ai-work-assistant/internal/domain"
)

// Analyzer は、Pull RequestをAIに分析させるためのインターフェース。
//
// application層はGeminiなどの具体的なAIサービスを直接知らず、
// このinterfaceだけに依存する。
//
// 将来的にGeminiからOpenAIやClaudeへ変更しても、
// application層を変更しなくて済むようにする。
type Analyzer interface {
	AnalyzePullRequests(
		ctx context.Context,
		pullRequests []domain.PullRequest,
	) (string, error)
}
