package github

import (
	"context"

	"github.com/google/go-github/v68/github"

	"github.com/yourname/ai-work-assistant/internal/domain"
)

// Client はGitHub APIと通信するAdapter。
//
// 外部サービスであるGitHub APIのレスポンスを
// アプリ内部のdomain.PullRequestへ変換する責務を持つ。
type Client struct {
	client *github.Client
}

// NewClient はGitHub Adapterを生成する。
func NewClient(client *github.Client) *Client {
	return &Client{
		client: client,
	}
}

// ListPullRequests はGitHub APIからPR一覧を取得し、
// GitHub固有のデータ形式からdomain.PullRequestへ変換する。
func (c *Client) ListPullRequests(
	ctx context.Context,
	owner string,
	repo string,
) ([]domain.PullRequest, error) {

	// GitHub APIを呼び出す。
	githubPRs, _, err := c.client.PullRequests.List(
		ctx,
		owner,
		repo,
		&github.PullRequestListOptions{
			State: "open",
		},
	)
	if err != nil {
		return nil, err
	}

	// GitHub APIのレスポンスをそのまま上位層へ返さず、
	// アプリ内部で利用する共通形式へ変換する。
	pullRequests := make([]domain.PullRequest, 0, len(githubPRs))

	for _, pr := range githubPRs {
		pullRequests = append(pullRequests, domain.PullRequest{
			ID:        pr.GetNumber(),
			Title:     pr.GetTitle(),
			Author:    pr.GetUser().GetLogin(),
			State:     pr.GetState(),
			URL:       pr.GetHTMLURL(),
			CreatedAt: pr.GetCreatedAt().Time,
			UpdatedAt: pr.GetUpdatedAt().Time,
		})
	}

	return pullRequests, nil
}
