package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	gogithub "github.com/google/go-github/v68/github"
)

func TestClient_ListPullRequests(t *testing.T) {
	// GitHub APIの代わりになる偽HTTPサーバーを起動する。
	//
	// 実際のGitHub APIを呼ばないため、
	// ネットワークやGitHubの状態に依存しないテストになる。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Adapterが想定したGitHub APIのパスへ
		// リクエストしているか確認する。
		expectedPath := "/repos/openai/example/pulls"

		if r.URL.Path != expectedPath {
			t.Errorf(
				"unexpected request path: got=%s want=%s",
				r.URL.Path,
				expectedPath,
			)
		}

		// GitHub APIが返すPull Requestレスポンスを模擬する。
		//
		// 重要なのは、
		// ここではdomain.PullRequestではなく
		// GitHub API側のJSON形式を返すこと。
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, `[
			{
				"number": 123,
				"title": "Fix login bug",
				"state": "open",
				"html_url": "https://github.com/openai/example/pull/123",
				"user": {
					"login": "haruki"
				},
				"created_at": "2026-09-01T10:00:00Z",
				"updated_at": "2026-09-07T10:00:00Z"
			}
		]`)
	}))
	defer server.Close()

	// go-githubのClientを生成する。
	githubClient := gogithub.NewClient(nil)

	// 通常は api.github.com を参照するが、
	// テスト中だけ先ほど作った偽サーバーへ向ける。
	baseURL := server.URL + "/"

	var err error
	githubClient.BaseURL, err = githubClient.BaseURL.Parse(baseURL)
	if err != nil {
		t.Fatalf("failed to set GitHub BaseURL: %v", err)
	}

	// 今回テストするAdapterを生成する。
	client := NewClient(githubClient)

	// 実際にPR一覧取得処理を実行する。
	pullRequests, err := client.ListPullRequests(
		context.Background(),
		"openai",
		"example",
	)
	if err != nil {
		t.Fatalf("ListPullRequests returned error: %v", err)
	}

	// PRが1件取得できることを確認する。
	if len(pullRequests) != 1 {
		t.Fatalf(
			"unexpected pull request count: got=%d want=1",
			len(pullRequests),
		)
	}

	pr := pullRequests[0]

	// GitHub API形式から、
	// domain.PullRequestへ正しく変換されたか確認する。
	if pr.ID != 123 {
		t.Errorf("unexpected ID: got=%d want=123", pr.ID)
	}

	if pr.Title != "Fix login bug" {
		t.Errorf(
			"unexpected Title: got=%s want=%s",
			pr.Title,
			"Fix login bug",
		)
	}

	if pr.Author != "haruki" {
		t.Errorf(
			"unexpected Author: got=%s want=%s",
			pr.Author,
			"haruki",
		)
	}

	if pr.State != "open" {
		t.Errorf(
			"unexpected State: got=%s want=%s",
			pr.State,
			"open",
		)
	}

	if pr.URL != "https://github.com/openai/example/pull/123" {
		t.Errorf(
			"unexpected URL: got=%s",
			pr.URL,
		)
	}
}
