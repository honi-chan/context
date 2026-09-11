package github

import (
	"context"
	"fmt"

	gogithub "github.com/google/go-github/v68/github"
)

// pullRequestDetails は、1つのPull Requestについて
// GitHub APIから追加取得した情報をまとめる内部用データ。
//
// domain層には公開せず、GitHub Adapter内部だけで利用する。
type pullRequestDetails struct {
	files     []*gogithub.CommitFile
	checkRuns []*gogithub.CheckRun
}

// pullRequestFetcher は、Pull Requestに紐づく
// 詳細情報の取得を担当する。
//
// client.goからFilesやCheckRunsの取得処理を分離し、
// GitHub API呼び出しの責務をこの型に集約する。
type pullRequestFetcher struct {
	client *gogithub.Client
}

// newPullRequestFetcher はFetcherを生成する。
func newPullRequestFetcher(
	client *gogithub.Client,
) *pullRequestFetcher {
	return &pullRequestFetcher{
		client: client,
	}
}

// fetch は、1件のPull Requestについて
// AI分析に必要な詳細情報をGitHubから取得する。
func (f *pullRequestFetcher) fetch(
	ctx context.Context,
	owner string,
	repo string,
	pr *gogithub.PullRequest,
) (*pullRequestDetails, error) {

	// Pull Requestで変更されているファイル一覧を取得する。
	files, _, err := f.client.PullRequests.ListFiles(
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

	// Check RunsはPull Request番号ではなく、
	// PRの最新コミットSHAに紐づいている。
	headSHA := pr.GetHead().GetSHA()

	checkRunsResult, _, err := f.client.Checks.ListCheckRunsForRef(
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

	return &pullRequestDetails{
		files:     files,
		checkRuns: checkRunsResult.CheckRuns,
	}, nil
}
