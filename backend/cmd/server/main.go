package main

import (
	"context"
	"log"
	"net/http"
	"os"

	gogithub "github.com/google/go-github/v68/github"
	"github.com/labstack/echo/v5"
	"google.golang.org/genai"

	"github.com/yourname/ai-work-assistant/internal/handler"
	geminiinfra "github.com/yourname/ai-work-assistant/internal/infrastructure/gemini"
	githubinfra "github.com/yourname/ai-work-assistant/internal/infrastructure/github"
	"github.com/yourname/ai-work-assistant/internal/usecase"
)

func main() {
	ctx := context.Background()

	// ========================================
	// 環境変数
	// ========================================

	// GitHub APIアクセス用Token。
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("GITHUB_TOKEN is required")
	}

	// Gemini APIアクセス用Key。
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Fatal("GEMINI_API_KEY is required")
	}

	// 使用するGeminiモデル。
	//
	// モデル名は環境変数で変更できるようにしておく。
	// 未設定の場合だけデフォルト値を利用する。
	geminiModel := os.Getenv("GEMINI_MODEL")
	if geminiModel == "" {
		geminiModel = "gemini-3.7-flash"
	}

	// ========================================
	// GitHub Client
	// ========================================

	// go-githubのClientを生成する。
	//
	// WithAuthTokenによって、
	// GitHub APIへのリクエストにTokenを付与する。
	githubClient := gogithub.NewClient(nil).
		WithAuthToken(githubToken)

	// GitHub APIの形式を、
	// アプリ内部のdomain.PullRequestへ変換するAdapter。
	githubAdapter := githubinfra.NewClient(
		githubClient,
	)

	// ========================================
	// Gemini Client
	// ========================================

	// Google公式Gemini SDKのClientを生成する。
	//
	// APIキーの読み込みはInfrastructure側ではなく
	// Composition Rootであるmain.goが担当する。
	geminiClient, err := genai.NewClient(
		ctx,
		&genai.ClientConfig{
			APIKey:  geminiAPIKey,
			Backend: genai.BackendGeminiAPI,
		},
	)
	if err != nil {
		log.Fatalf("failed to create Gemini client: %v", err)
	}

	// Gemini SDKを、
	// ai.Analyzerとして利用できるAdapterへ変換する。
	geminiAdapter := geminiinfra.NewClient(
		geminiClient,
		geminiModel,
	)

	// ========================================
	// UseCase
	// ========================================

	// GitHub AdapterとGemini AdapterをUseCaseへDIする。
	//
	// UseCase自身は、
	// GitHubやGeminiという具体的なサービスを知らない。
	analyzeUseCase := usecase.NewAnalyzePullRequestsUseCase(
		githubAdapter,
		geminiAdapter,
	)

	// ========================================
	// Handler
	// ========================================

	analyzeHandler := handler.NewAnalyzeHandler(
		analyzeUseCase,
	)

	// ========================================
	// Echo
	// ========================================

	e := echo.New()

	// 動作確認用Health Check。
	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(
			http.StatusOK,
			map[string]string{
				"status": "ok",
			},
		)
	})

	// GitHub PRをGeminiに分析させるAPI。
	e.POST(
		"/api/analyze",
		analyzeHandler.Analyze,
	)

	// ========================================
	// Server Start
	// ========================================

	log.Println("server started on :8080")

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
