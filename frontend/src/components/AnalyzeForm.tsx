"use client";

import { useState } from "react";

import AnalysisResult from "@/components/AnalysisResult";
import RepositoryForm from "@/components/RepositoryForm";
import { analyzeRepository } from "@/lib/api/analyze";

export default function AnalyzeForm() {
  const [result, setResult] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleAnalyze = async (
    owner: string,
    repo: string,
  ) => {
    setLoading(true);
    setResult("");
    setError("");

    try {
      // API呼び出しは親コンポーネントに残す。
      // RepositoryFormは入力UIだけに責務を限定する。
      const response = await analyzeRepository({
        owner,
        repo,
      });

      setResult(response.result);
    } catch (error) {
      console.error(
        "failed to analyze repository",
        error,
      );

      setError(
        "Repositoryの分析に失敗しました。時間を置いて再度実行してください。",
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <section className="workspace">
      <RepositoryForm
        loading={loading}
        onAnalyze={handleAnalyze}
      />

      <AnalysisResult
        result={result}
        loading={loading}
        error={error}
      />
    </section>
  );
}