package domain

import "time"

// PullRequest は、このアプリケーション内で扱う
// Pull Request の共通データ形式。
//
// GitHub API固有の型を直接利用しないことで、
// 将来的にGitLabなど別サービスを追加しても
// domain層を変更せずに済むようにする。
type PullRequest struct {
	ID        int
	Title     string
	Author    string
	State     string
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
