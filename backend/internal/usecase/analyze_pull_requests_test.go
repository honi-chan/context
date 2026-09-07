package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/yourname/ai-work-assistant/internal/domain"
)

// fakeGitHubRepository はテスト用のGitHubRepository実装。
//
// 実際のGitHub APIには接続せず、
// テスト側であらかじめ設定したPull Requestを返す。
//
// UseCaseがGitHubの具体的な実装に依存していないため、
// このように簡単なFakeへ差し替えられる。
type fakeGitHubRepository struct {
	pullRequests []domain.PullRequest
	err          error
}

// ListPullRequests はテスト用のPR一覧を返す。
//
// owner / repo は今回のテストでは利用しない。
// GitHub APIへのHTTP通信も発生しない。
func (f *fakeGitHubRepository) ListPullRequests(
	ctx context.Context,
	owner string,
	repo string,
) ([]domain.PullRequest, error) {
	return f.pullRequests, f.err
}

// fakeAnalyzer はテスト用のAI Analyzer実装。
//
// Gemini APIには接続せず、
// あらかじめ設定した分析結果を返す。
type fakeAnalyzer struct {
	result string
	err    error

	// receivedPullRequests は、
	// UseCaseからAnalyzerへ正しくPRが渡されたか確認するために保持する。
	receivedPullRequests []domain.PullRequest
}

// AnalyzePullRequests はAI分析を模擬する。
func (f *fakeAnalyzer) AnalyzePullRequests(
	ctx context.Context,
	pullRequests []domain.PullRequest,
) (string, error) {

	// UseCaseから受け取った値を保存しておく。
	// 後でテストから内容を検証する。
	f.receivedPullRequests = pullRequests

	return f.result, f.err
}

func TestAnalyzePullRequestsUseCase_Execute(t *testing.T) {
	// GitHubから取得できたことにするPRデータ。
	pullRequests := []domain.PullRequest{
		{
			ID:        123,
			Title:     "Fix login bug",
			Author:    "haruki",
			State:     "open",
			URL:       "https://github.com/yourname/example/pull/123",
			CreatedAt: time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC),
		},
	}

	// GitHubRepositoryのFakeを作る。
	//
	// Execute()からListPullRequests()が呼ばれると、
	// 上で定義したPRが返る。
	githubRepository := &fakeGitHubRepository{
		pullRequests: pullRequests,
	}

	// Geminiの代わりになるFake Analyzerを作る。
	//
	// AI分析されたことにして固定文字列を返す。
	analyzer := &fakeAnalyzer{
		result: "PR #123 を最優先で確認してください。",
	}

	// テスト対象のUseCaseを生成する。
	//
	// 本物のGitHub ClientやGemini Clientではなく、
	// interfaceを満たしているFakeを注入する。
	uc := NewAnalyzePullRequestsUseCase(
		githubRepository,
		analyzer,
	)

	// UseCaseを実行する。
	result, err := uc.Execute(
		context.Background(),
		"yourname",
		"example",
	)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	// AIの分析結果がそのまま返ってくることを確認する。
	expectedResult := "PR #123 を最優先で確認してください。"

	if result != expectedResult {
		t.Errorf(
			"unexpected result: got=%q want=%q",
			result,
			expectedResult,
		)
	}

	// GitHubから取得したPRが、
	// Analyzerへ正しく渡されたことを確認する。
	if len(analyzer.receivedPullRequests) != 1 {
		t.Fatalf(
			"unexpected received pull request count: got=%d want=1",
			len(analyzer.receivedPullRequests),
		)
	}

	receivedPR := analyzer.receivedPullRequests[0]

	if receivedPR.ID != 123 {
		t.Errorf(
			"unexpected pull request ID: got=%d want=123",
			receivedPR.ID,
		)
	}

	if receivedPR.Title != "Fix login bug" {
		t.Errorf(
			"unexpected pull request title: got=%q want=%q",
			receivedPR.Title,
			"Fix login bug",
		)
	}
}

func TestAnalyzePullRequestsUseCase_Execute_NoPullRequests(t *testing.T) {
	// PRが1件も存在しない状態を作る。
	githubRepository := &fakeGitHubRepository{
		pullRequests: []domain.PullRequest{},
	}

	analyzer := &fakeAnalyzer{
		result: "この値は呼ばれてはいけない",
	}

	uc := NewAnalyzePullRequestsUseCase(
		githubRepository,
		analyzer,
	)

	result, err := uc.Execute(
		context.Background(),
		"yourname",
		"example",
	)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	// PRが存在しない場合は、
	// Geminiを呼ばずにUseCase自身がこのメッセージを返す。
	expectedResult := "現在対応が必要なPull Requestはありません。"

	if result != expectedResult {
		t.Errorf(
			"unexpected result: got=%q want=%q",
			result,
			expectedResult,
		)
	}

	// Analyzerが呼ばれていないことを確認する。
	//
	// PRがないのにGemini APIを呼ぶと、
	// 無駄なAPIリクエストやコストが発生するため。
	if analyzer.receivedPullRequests != nil {
		t.Errorf("analyzer should not have been called")
	}
}
