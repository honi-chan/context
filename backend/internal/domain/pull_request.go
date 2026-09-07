package domain

import "time"

// PullRequest は、このアプリケーション内で扱う
// Pull Requestの共通データ形式。
//
// GitHub API固有の型を直接利用しないことで、
// 将来的にGitLabなど別サービスへ対応する場合でも
// domain層への影響を小さくする。
type PullRequest struct {
	ID    int
	Title string

	// Body はPull Requestの説明文。
	//
	// タイトルだけではAIが変更内容を正確に判断できないため、
	// PR作成者が記載した目的・変更内容・注意事項などを
	// AIの判断材料として利用する。
	Body string

	Author string
	State  string
	URL    string

	CreatedAt time.Time
	UpdatedAt time.Time
}
