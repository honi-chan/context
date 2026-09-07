package domain

import "time"

// ChangedFile はPR内の1ファイル分の変更情報。
type ChangedFile struct {
	Filename  string
	Status    string
	Additions int
	Deletions int
	Changes   int

	// GitHubが返すdiff形式の変更内容。
	// AIが実際のコード変更内容を判断するために使用する。
	Patch string
}

// CheckRun はPRの最新コミットに対して実行された
// CI / Test / LintなどのCheck結果。
type CheckRun struct {
	Name       string
	Status     string
	Conclusion string
}

// PullRequest はアプリケーション内部で扱うPR情報。
type PullRequest struct {
	ID     int
	Title  string
	Body   string
	Author string
	State  string
	URL    string

	HeadBranch string
	BaseBranch string
	HeadSHA    string

	ChangedFiles []ChangedFile

	// CheckRunsFetched はGitHub APIから
	// Check Runsの取得に成功したかを表す。
	//
	// true + CheckRunsが空:
	//   API取得には成功したがCheck Runが存在しない。
	//
	// true + CheckRunsあり:
	//   CI/Test等のCheck Runが存在する。
	//
	// 現在は取得エラー時に処理全体をerrorにしているため、
	// Geminiへ渡るPullRequestでは基本的にtrueになる。
	CheckRunsFetched bool

	CheckRuns []CheckRun

	CreatedAt time.Time
	UpdatedAt time.Time
}
