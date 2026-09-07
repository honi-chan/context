package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
)

// fakeAnalyzeUseCase はHandlerテスト専用のFake。
//
// GitHub APIやGemini APIは一切呼び出さず、
// あらかじめ指定した結果を返す。
type fakeAnalyzeUseCase struct {
	result string
	err    error

	// Handlerから正しいowner/repoが
	// 渡されたことを確認するために保持する。
	receivedOwner string
	receivedRepo  string
}

// Execute は本物のUseCaseの代わりになる処理。
func (f *fakeAnalyzeUseCase) Execute(
	ctx context.Context,
	owner string,
	repo string,
) (string, error) {

	// Handlerから受け取った値を記録する。
	f.receivedOwner = owner
	f.receivedRepo = repo

	return f.result, f.err
}

// 正常系:
//
// POST /api/analyze
//
// が正しいJSONを受け取った場合に
// UseCaseを呼び出して200を返すことを確認する。
func TestAnalyzeHandler_Analyze(t *testing.T) {
	e := echo.New()

	// 実際のHTTP Requestを模擬する。
	requestBody := `{
		"owner": "yourname",
		"repo": "ai-work-assistant"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/analyze",
		strings.NewReader(requestBody),
	)

	// JSONであることをEchoへ伝える。
	req.Header.Set(
		echo.HeaderContentType,
		echo.MIMEApplicationJSON,
	)

	// Handlerが書き込むHTTP Responseを記録する。
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	// Geminiの代わりに固定結果を返す。
	fakeUseCase := &fakeAnalyzeUseCase{
		result: "PR #123 を最優先で確認してください。",
	}

	handler := NewAnalyzeHandler(fakeUseCase)

	// Handlerを実行する。
	if err := handler.Analyze(c); err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}

	// HTTP 200が返ったことを確認する。
	if rec.Code != http.StatusOK {
		t.Errorf(
			"unexpected status code: got=%d want=%d",
			rec.Code,
			http.StatusOK,
		)
	}

	// JSONのownerがUseCaseまで渡されたことを確認する。
	if fakeUseCase.receivedOwner != "yourname" {
		t.Errorf(
			"unexpected owner: got=%q want=%q",
			fakeUseCase.receivedOwner,
			"yourname",
		)
	}

	// JSONのrepoがUseCaseまで渡されたことを確認する。
	if fakeUseCase.receivedRepo != "ai-work-assistant" {
		t.Errorf(
			"unexpected repo: got=%q want=%q",
			fakeUseCase.receivedRepo,
			"ai-work-assistant",
		)
	}

	// AI分析結果がHTTP Responseに含まれることを確認する。
	expectedBody := `{"result":"PR #123 を最優先で確認してください。"}` + "\n"

	if rec.Body.String() != expectedBody {
		t.Errorf(
			"unexpected response body: got=%q want=%q",
			rec.Body.String(),
			expectedBody,
		)
	}
}

// ownerが空の場合に400になることを確認する。
func TestAnalyzeHandler_Analyze_MissingOwner(t *testing.T) {
	e := echo.New()

	requestBody := `{
		"repo": "ai-work-assistant"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/analyze",
		strings.NewReader(requestBody),
	)

	req.Header.Set(
		echo.HeaderContentType,
		echo.MIMEApplicationJSON,
	)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	fakeUseCase := &fakeAnalyzeUseCase{}

	handler := NewAnalyzeHandler(fakeUseCase)

	if err := handler.Analyze(c); err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}

	// 必須パラメータ不足なので400を期待する。
	if rec.Code != http.StatusBadRequest {
		t.Errorf(
			"unexpected status code: got=%d want=%d",
			rec.Code,
			http.StatusBadRequest,
		)
	}

	// Validationで終了しているため、
	// UseCaseが呼ばれていないことも確認する。
	if fakeUseCase.receivedRepo != "" {
		t.Errorf("usecase should not have been called")
	}
}

// UseCaseでエラーが発生した場合に
// 500を返すことを確認する。
func TestAnalyzeHandler_Analyze_UseCaseError(t *testing.T) {
	e := echo.New()

	requestBody := `{
		"owner": "yourname",
		"repo": "ai-work-assistant"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/analyze",
		strings.NewReader(requestBody),
	)

	req.Header.Set(
		echo.HeaderContentType,
		echo.MIMEApplicationJSON,
	)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// UseCaseでエラーが発生した状態を再現する。
	fakeUseCase := &fakeAnalyzeUseCase{
		err: errors.New("Gemini API error"),
	}

	handler := NewAnalyzeHandler(fakeUseCase)

	if err := handler.Analyze(c); err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf(
			"unexpected status code: got=%d want=%d",
			rec.Code,
			http.StatusInternalServerError,
		)
	}

	// Geminiなどの内部エラーを
	// HTTPレスポンスへ直接露出していないことも確認する。
	expectedBody := `{"message":"Pull Requestの分析に失敗しました。"}` + "\n"

	if rec.Body.String() != expectedBody {
		t.Errorf(
			"unexpected response body: got=%q want=%q",
			rec.Body.String(),
			expectedBody,
		)
	}
}
