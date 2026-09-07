package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
)

type AnalyzeUseCase interface {
	Execute(
		ctx context.Context,
		owner string,
		repo string,
	) (string, error)
}

type AnalyzeHandler struct {
	usecase AnalyzeUseCase
}

func NewAnalyzeHandler(uc AnalyzeUseCase) *AnalyzeHandler {
	return &AnalyzeHandler{
		usecase: uc,
	}
}

type AnalyzeRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type AnalyzeResponse struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (h *AnalyzeHandler) Analyze(c *echo.Context) error {
	var request AnalyzeRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Message: "リクエスト形式が不正です。",
			},
		)
	}

	if request.Owner == "" || request.Repo == "" {
		return c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Message: "owner と repo は必須です。",
			},
		)
	}

	result, err := h.usecase.Execute(
		c.Request().Context(),
		request.Owner,
		request.Repo,
	)
	if err != nil {
		// 内部エラーの詳細はサーバーログへ出す。
		// これでGitHub API / Gemini APIのどこで
		// 失敗したのか確認できる。
		log.Printf(
			"analyze pull requests failed: owner=%s repo=%s error=%v",
			request.Owner,
			request.Repo,
			err,
		)

		// クライアントには内部情報をそのまま返さない。
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
