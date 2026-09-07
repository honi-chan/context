package github

import (
	"context"
	"fmt"

	gogithub "github.com/google/go-github/v68/github"

	"github.com/yourname/ai-work-assistant/internal/domain"
)

// Client はGitHub Repositoryへのアクセスを担当するAdapter。
//
// GitHub SDKをrepository/usecase層から隠し、
// アプリケーション内部ではdomain.PullRequestとして扱えるようにする。
type Client struct {
	client  *gogithub.Client
	fetcher *pullRequestFetcher
}

// NewClient はGitHub Clientを生成する。
func NewClient(
	client *gogithub.Client,
) *Client {
	return &Client{
		client: client,

		// PR詳細取得はFetcherへ委譲する。
		fetcher: newPullRequestFetcher(client),
	}
}

// ListPullRequests はOpen状態のPull Request一覧を取得し、
// AI分析に必要な詳細情報を付与してdomain型へ変換する。
func (c *Client) ListPullRequests(
	ctx context.Context,
	owner string,
	repo string,
) ([]domain.PullRequest, error) {

	// まずOpen状態のPull Request一覧だけ取得する。
	githubPRs, _, err := c.client.PullRequests.List(
		ctx,
		owner,
		repo,
		&gogithub.PullRequestListOptions{
			State: "open",
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list pull requests: %w",
			err,
		)
	}

	pullRequests := make(
		[]domain.PullRequest,
		0,
		len(githubPRs),
	)

	for _, pr := range githubPRs {

		// Files / CheckRunsなど、
		// PRに紐づく詳細取得はFetcherへ任せる。
		details, err := c.fetcher.fetch(
			ctx,
			owner,
			repo,
			pr,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"fetch pull request details: %w",
				err,
			)
		}

		// GitHub固有の型からdomain型への変換は
		// Mapperへ任せる。
		pullRequest := mapPullRequest(
			pr,
			details.files,
			details.checkRuns,

			// fetch()が正常終了したため、
			// Check Runs APIの取得自体は成功している。
			true,
		)

		pullRequests = append(
			pullRequests,
			pullRequest,
		)
	}

	return pullRequests, nil
}
