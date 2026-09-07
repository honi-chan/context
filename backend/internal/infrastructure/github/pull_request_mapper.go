package github

import (
	gogithub "github.com/google/go-github/v68/github"

	"github.com/yourname/ai-work-assistant/internal/domain"
)

// mapPullRequest はGitHub APIのデータを
// アプリケーション内部のPullRequestへ変換する。
func mapPullRequest(
	pr *gogithub.PullRequest,
	files []*gogithub.CommitFile,
	githubCheckRuns []*gogithub.CheckRun,
	checkRunsFetched bool,
) domain.PullRequest {

	changedFiles := make(
		[]domain.ChangedFile,
		0,
		len(files),
	)

	for _, file := range files {
		changedFiles = append(
			changedFiles,
			domain.ChangedFile{
				Filename:  file.GetFilename(),
				Status:    file.GetStatus(),
				Additions: file.GetAdditions(),
				Deletions: file.GetDeletions(),
				Changes:   file.GetChanges(),
			},
		)
	}

	checkRuns := make(
		[]domain.CheckRun,
		0,
		len(githubCheckRuns),
	)

	for _, checkRun := range githubCheckRuns {
		checkRuns = append(
			checkRuns,
			domain.CheckRun{
				Name:       checkRun.GetName(),
				Status:     checkRun.GetStatus(),
				Conclusion: checkRun.GetConclusion(),
			},
		)
	}

	return domain.PullRequest{
		ID:     pr.GetNumber(),
		Title:  pr.GetTitle(),
		Body:   pr.GetBody(),
		Author: pr.GetUser().GetLogin(),
		State:  pr.GetState(),
		URL:    pr.GetHTMLURL(),

		HeadBranch: pr.GetHead().GetRef(),
		BaseBranch: pr.GetBase().GetRef(),
		HeadSHA:    pr.GetHead().GetSHA(),

		ChangedFiles: changedFiles,

		// Check Runs APIの取得結果を明示する。
		CheckRunsFetched: checkRunsFetched,
		CheckRuns:        checkRuns,

		CreatedAt: pr.GetCreatedAt().Time,
		UpdatedAt: pr.GetUpdatedAt().Time,
	}
}
