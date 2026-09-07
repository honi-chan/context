"use client";

import { FormEvent, useState } from "react";

import { analyzeRepository } from "@/lib/api/analyze";

// AnalyzeForm は、
//
// 1. GitHub owner入力
// 2. Repository入力
// 3. 分析API実行
// 4. AI結果表示
//
// を担当するClient Component。
export default function AnalyzeForm() {
  // GitHubのowner。
  const [owner, setOwner] = useState("");

  // Repository名。
  const [repo, setRepo] = useState("");

  // Geminiから返ってきた分析結果。
  const [result, setResult] = useState("");

  // エラーメッセージ。
  const [error, setError] = useState("");

  // API通信中かどうか。
  //
  // 二重送信防止やLoading表示に利用する。
  const [isLoading, setIsLoading] = useState(false);

  // フォーム送信時の処理。
  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    // HTML Formのデフォルト送信を止める。
    event.preventDefault();

    // 前回の結果・エラーをリセットする。
    setResult("");
    setError("");

    // API通信開始。
    setIsLoading(true);

    try {
      // Go Backendへ分析リクエストを送る。
      const response = await analyzeRepository({
        owner,
        repo,
      });

      // Geminiの分析結果を画面へ表示する。
      setResult(response.result);
    } catch (err) {
      // Error型ならmessageを表示する。
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("分析に失敗しました。");
      }
    } finally {
      // 成功・失敗に関係なくLoadingを解除する。
      setIsLoading(false);
    }
  }

  return (
    <div>
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="owner">
            GitHub Owner
          </label>

          <input
            id="owner"
            type="text"
            value={owner}
            onChange={(event) => setOwner(event.target.value)}
            placeholder="yourname"
            required
          />
        </div>

        <div>
          <label htmlFor="repo">
            Repository
          </label>

          <input
            id="repo"
            type="text"
            value={repo}
            onChange={(event) => setRepo(event.target.value)}
            placeholder="ai-work-assistant"
            required
          />
        </div>

        <button
          type="submit"
          disabled={isLoading}
        >
          {isLoading ? "分析中..." : "GitHubを分析"}
        </button>
      </form>

      {/* APIでエラーになった場合 */}
      {error && (
        <div>
          <p>{error}</p>
        </div>
      )}

      {/* Geminiから結果が返ってきた場合 */}
      {result && (
        <div>
          <h2>AI分析結果</h2>

          {/* 改行を保持して表示する */}
          <p style={{ whiteSpace: "pre-wrap" }}>
            {result}
          </p>
        </div>
      )}
    </div>
  );
}