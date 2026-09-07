package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v5"
)

// AnalyzeUseCase はHandlerが必要とする機能だけを定義する。
//
// Handlerは具体的なAnalyzePullRequestsUseCaseではなく、
// このinterfaceに依存する。
//
// これによりHandlerのテスト時に
// Fake UseCaseへ簡単に差し替えられる。
type AnalyzeUseCase interface {
	Execute(
		ctx context.Context,
		owner string,
		repo string,
	) (string, error)
}

// AnalyzeHandler はPR分析APIのHTTP層を担当する。
//
// GitHubやGeminiなどの具体的な処理は知らず、
// AnalyzeUseCaseを呼び出すことだけを担当する。
type AnalyzeHandler struct {
	usecase AnalyzeUseCase
}

// NewAnalyzeHandler はHandlerを生成する。
func NewAnalyzeHandler(
	uc AnalyzeUseCase,
) *AnalyzeHandler {
	return &AnalyzeHandler{
		usecase: uc,
	}
}

// AnalyzeRequest はPOST /api/analyzeの入力DTO。
type AnalyzeRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

// AnalyzeResponse は正常時のレスポンスDTO。
type AnalyzeResponse struct {
	Result string `json:"result"`
}

// ErrorResponse はエラー時のレスポンスDTO。
type ErrorResponse struct {
	Message string `json:"message"`
}

// Analyze はPR分析APIを処理する。
func (h *AnalyzeHandler) Analyze(c *echo.Context) error {
	var request AnalyzeRequest

	// JSON BodyをリクエストDTOへ変換する。
	if err := c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Message: "リクエスト形式が不正です。",
			},
		)
	}

	// MVPではownerとrepoを必須とする。
	if request.Owner == "" || request.Repo == "" {
		return c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Message: "owner と repo は必須です。",
			},
		)
	}

	// HTTP RequestのContextをそのままUseCaseへ渡す。
	result, err := h.usecase.Execute(
		c.Request().Context(),
		request.Owner,
		request.Repo,
	)
	if err != nil {
		// 内部エラーの詳細はクライアントへ露出させない。
		return c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{
				Message: "Pull Requestの分析に失敗しました。",
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		AnalyzeResponse{
			Result: result,
		},
	)
}
