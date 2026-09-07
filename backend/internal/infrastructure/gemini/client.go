package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"google.golang.org/genai"

	"github.com/yourname/ai-work-assistant/internal/domain"
)

type Client struct {
	client *genai.Client
	model  string
}

func NewClient(
	client *genai.Client,
	model string,
) *Client {
	return &Client{
		client: client,
		model:  model,
	}
}

func (c *Client) AnalyzePullRequests(
	ctx context.Context,
	pullRequests []domain.PullRequest,
) (string, error) {

	prJSON, err := json.MarshalIndent(
		pullRequests,
		"",
		"  ",
	)
	if err != nil {
		return "", fmt.Errorf(
			"marshal pull requests: %w",
			err,
		)
	}

	prompt := fmt.Sprintf(`
あなたはソフトウェア開発チームを支援するAIアシスタントです。

以下は現在OpenになっているPull Request一覧です。

この中から、
「今最も優先して確認・対応すべきこと」を1つだけ選んでください。

回答には以下を含めてください。

- 対象Pull Request
- 優先すべき理由
- 次に取るべき行動

重要なルール:

- 与えられたPull Request情報だけを根拠に判断してください。
- 記載されていないプロジェクト事情を推測しないでください。
- 判断材料が不足している場合は、そのことを明示してください。
- Pull Requestのタイトルだけで変更内容を断定しないでください。
- CheckRunsFetched=true かつ CheckRuns=[] の場合、
  「CI情報を取得できなかった」と判断しないでください。
  GitHubからCheck Runsの取得には成功したものの、
  対象コミットにCheck Runが存在しなかった状態として扱ってください。
- BaseBranchや差分内容からブランチ作成経緯・コミット履歴を推測する場合は、
  事実として断定せず「可能性がある」と表現してください。
- Gitのコミット履歴が提供されていない場合、
  「設定ミス」「コミット混入」「rebaseが必要」などを確定事項として扱わないでください。

Pull Requests:

%s
`, string(prJSON))

	const maxAttempts = 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		response, err := c.client.Models.GenerateContent(
			ctx,
			c.model,
			genai.Text(prompt),
			nil,
		)

		if err == nil {
			return response.Text(), nil
		}

		// Gemini APIの失敗内容を各試行ごとに記録する。
		// 503 / 429 / その他のエラーを切り分けやすくする。
		log.Printf(
			"Gemini GenerateContent failed: attempt=%d/%d model=%s error=%v",
			attempt,
			maxAttempts,
			c.model,
			err,
		)

		if attempt == maxAttempts {
			return "", fmt.Errorf(
				"generate Gemini content after %d attempts: %w",
				maxAttempts,
				err,
			)
		}

		// 1秒 → 2秒 → 4秒 の指数バックオフ。
		waitDuration := time.Second * time.Duration(1<<(attempt-1))

		select {
		case <-ctx.Done():
			return "", ctx.Err()

		case <-time.After(waitDuration):
		}
	}

	return "", fmt.Errorf(
		"generate Gemini content failed",
	)
}
