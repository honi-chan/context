package github

import (
	"testing"

	gogithub "github.com/google/go-github/v68/github"
)

func TestMapPullRequest(t *testing.T) {
	pr := &gogithub.PullRequest{
		Number: gogithub.Int(123),
		Title:  gogithub.String("Fix login bug"),
		Body:   gogithub.String("ログインエラーを修正します。"),
		State:  gogithub.String("open"),

		HTMLURL: gogithub.String(
			"https://github.com/openai/example/pull/123",
		),

		User: &gogithub.User{
			Login: gogithub.String("haruki"),
		},

		Head: &gogithub.PullRequestBranch{
			Ref: gogithub.String("feat/login-fix"),
			SHA: gogithub.String("abc123"),
		},

		Base: &gogithub.PullRequestBranch{
			Ref: gogithub.String("main"),
		},
	}

	files := []*gogithub.CommitFile{
		{
			Filename:  gogithub.String("backend/cmd/server/main.go"),
			Status:    gogithub.String("modified"),
			Additions: gogithub.Int(10),
			Deletions: gogithub.Int(2),
			Changes:   gogithub.Int(12),
		},
		{
			Filename:  gogithub.String("frontend/src/app/page.tsx"),
			Status:    gogithub.String("modified"),
			Additions: gogithub.Int(20),
			Deletions: gogithub.Int(5),
			Changes:   gogithub.Int(25),
		},
	}

	checkRuns := []*gogithub.CheckRun{
		{
			Name:       gogithub.String("test"),
			Status:     gogithub.String("completed"),
			Conclusion: gogithub.String("success"),
		},
		{
			Name:       gogithub.String("build"),
			Status:     gogithub.String("completed"),
			Conclusion: gogithub.String("failure"),
		},
	}

	result := mapPullRequest(
		pr,
		files,
		checkRuns,

		// GitHub Check Runs APIの取得に
		// 成功した状態を再現する。
		true,
	)

	if result.ID != 123 {
		t.Errorf(
			"unexpected ID: got=%d want=123",
			result.ID,
		)
	}

	if result.HeadBranch != "feat/login-fix" {
		t.Errorf(
			"unexpected HeadBranch: got=%q want=%q",
			result.HeadBranch,
			"feat/login-fix",
		)
	}

	if result.BaseBranch != "main" {
		t.Errorf(
			"unexpected BaseBranch: got=%q want=%q",
			result.BaseBranch,
			"main",
		)
	}

	if len(result.ChangedFiles) != 2 {
		t.Fatalf(
			"unexpected ChangedFiles count: got=%d want=2",
			len(result.ChangedFiles),
		)
	}

	firstFile := result.ChangedFiles[0]

	if firstFile.Filename != "backend/cmd/server/main.go" {
		t.Errorf(
			"unexpected Filename: got=%q",
			firstFile.Filename,
		)
	}

	if firstFile.Status != "modified" {
		t.Errorf(
			"unexpected Status: got=%q",
			firstFile.Status,
		)
	}

	if firstFile.Additions != 10 {
		t.Errorf(
			"unexpected Additions: got=%d want=10",
			firstFile.Additions,
		)
	}

	if firstFile.Deletions != 2 {
		t.Errorf(
			"unexpected Deletions: got=%d want=2",
			firstFile.Deletions,
		)
	}

	if firstFile.Changes != 12 {
		t.Errorf(
			"unexpected Changes: got=%d want=12",
			firstFile.Changes,
		)
	}

	if len(result.CheckRuns) != 2 {
		t.Fatalf(
			"unexpected CheckRuns count: got=%d want=2",
			len(result.CheckRuns),
		)
	}

	if result.CheckRuns[0].Name != "test" {
		t.Errorf(
			"unexpected check name: got=%q",
			result.CheckRuns[0].Name,
		)
	}

	if result.CheckRuns[0].Conclusion != "success" {
		t.Errorf(
			"unexpected check conclusion: got=%q",
			result.CheckRuns[0].Conclusion,
		)
	}

	if result.CheckRuns[1].Conclusion != "failure" {
		t.Errorf(
			"unexpected check conclusion: got=%q want=failure",
			result.CheckRuns[1].Conclusion,
		)
	}

	// Check Runs APIの取得成功状態が
	// domainへ正しく変換されていることを確認する。
	if !result.CheckRunsFetched {
		t.Errorf(
			"unexpected CheckRunsFetched: got=false want=true",
		)
	}
}
