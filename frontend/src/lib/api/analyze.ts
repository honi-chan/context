// Go Backendへ送信するRequest型。
export type AnalyzeRequest = {
  owner: string;
  repo: string;
};

// Go Backendから返ってくるResponse型。
export type AnalyzeResponse = {
  result: string;
};

// Go Backendでエラーになった場合のResponse型。
type ErrorResponse = {
  message: string;
};

// GitHub Repositoryを分析するAPI。
//
// UI Componentから直接fetchを書かず、
// API通信処理をこのファイルへまとめる。
//
// 将来的にAPI URLや認証方式が変わっても、
// UI側への影響を小さくする。
export async function analyzeRepository(
  request: AnalyzeRequest,
): Promise<AnalyzeResponse> {
  // 環境変数からGo BackendのURLを取得する。
  //
  // NEXT_PUBLIC_ を付けることで、
  // Browser側のJavaScriptから参照できる。
  const apiUrl =
    process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

  // Go Backendへ分析リクエストを送信する。
  const response = await fetch(`${apiUrl}/api/analyze`, {
    method: "POST",

    headers: {
      "Content-Type": "application/json",
    },

    body: JSON.stringify(request),
  });

  // HTTP 4xx / 5xxの場合。
  if (!response.ok) {
    // Backendが返したエラーメッセージを取得する。
    //
    // JSONとして取得できなかった場合も考慮して
    // catchでデフォルトメッセージへフォールバックする。
    const errorResponse: ErrorResponse = await response
      .json()
      .catch(() => ({
        message: "分析に失敗しました。",
      }));

    throw new Error(errorResponse.message);
  }

  // 正常時のJSONレスポンスを返す。
  return response.json() as Promise<AnalyzeResponse>;
}