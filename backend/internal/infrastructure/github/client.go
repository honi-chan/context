package github

import (
	"context"
	"fmt"

	gogithub "github.com/google/go-github/v68/github"

	"github.com/yourname/ai-work-assistant/internal/domain"
)

type Client struct {
	client *gogithub.Client
}

func NewClient(client *gogithub.Client) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) ListPullRequests(
	ctx context.Context,
	owner string,
	repo string,
) ([]domain.PullRequest, error) {

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

		// PRの変更ファイル一覧を取得する。
		files, _, err := c.client.PullRequests.ListFiles(
			ctx,
			owner,
			repo,
			pr.GetNumber(),
			&gogithub.ListOptions{
				PerPage: 100,
			},
		)
		if err != nil {
			return nil, fmt.Errorf(
				"list pull request files: pull_number=%d: %w",
				pr.GetNumber(),
				err,
			)
		}

		// PRの最新コミットSHAを取得する。
		//
		// Check Runsはコミットに紐づいているため、
		// Head側の最新SHAを指定する。
		headSHA := pr.GetHead().GetSHA()

		// GitHub Actions / CI / Testなどの
		// Check Run一覧を取得する。
		checkRunsResult, _, err := c.client.Checks.ListCheckRunsForRef(
			ctx,
			owner,
			repo,
			headSHA,
			&gogithub.ListCheckRunsOptions{
				ListOptions: gogithub.ListOptions{
					PerPage: 100,
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf(
				"list check runs: pull_number=%d sha=%s: %w",
				pr.GetNumber(),
				headSHA,
				err,
			)
		}

		// ここまで到達している時点で
		// ListCheckRunsForRefは正常終了している。
		//
		// そのためCheckRunsが0件だったとしても
		// 「取得失敗」ではなく
		// 「取得成功した結果、0件」と判断できる。
		pullRequests = append(
			pullRequests,
			mapPullRequest(
				pr,
				files,
				checkRunsResult.CheckRuns,
				true,
			),
		)
	}

	return pullRequests, nil
}
