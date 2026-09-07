package gemini

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/genai"

	"github.com/yourname/ai-work-assistant/internal/domain"
)

// Client はGemini APIを利用するAI Adapter。
//
// application層から渡されたdomain.PullRequestを
// Geminiが理解できるプロンプトへ変換して送信する。
type Client struct {
	client *genai.Client
	model  string
}

// NewClient はGemini Adapterを生成する。
//
// genai.Clientを外から受け取ることで、
// Client自身はAPIキーの読み込みなどを担当しない。
func NewClient(
	client *genai.Client,
	model string,
) *Client {
	return &Client{
		client: client,
		model:  model,
	}
}

// AnalyzePullRequests はPR一覧をGeminiへ渡し、
// 「今最も優先して対応すべきこと」を分析させる。
func (c *Client) AnalyzePullRequests(
	ctx context.Context,
	pullRequests []domain.PullRequest,
) (string, error) {

	// PR一覧をJSONへ変換する。
	//
	// Geminiへ構造化された情報を渡すことで、
	// PRが複数あっても内容を理解しやすくする。
	prJSON, err := json.MarshalIndent(pullRequests, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal pull requests: %w", err)
	}

	// Geminiへ渡すプロンプト。
	//
	// MVPでは「1つだけ選ぶ」という役割に限定する。
	// AIに何でも判断させず、責務を明確にしておく。
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

Pull Requests:

%s
`, string(prJSON))

	// Gemini APIを呼び出す。
	response, err := c.client.Models.GenerateContent(
		ctx,
		c.model,
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("generate Gemini content: %w", err)
	}

	// Gemini SDKのレスポンスから
	// テキスト部分だけを取り出してapplication層へ返す。
	return response.Text(), nil
}
