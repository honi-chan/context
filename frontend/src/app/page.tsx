import AnalyzeForm from "@/components/AnalyzeForm";

// トップページ。
//
// Server ComponentのままでOK。
// ユーザー操作が必要な部分だけ
// AnalyzeFormをClient Componentとして分離している。
export default function Home() {
  return (
    <main>
      <h1>
        AI Work Assistant
      </h1>

      <p>
        GitHub RepositoryをAIが分析して、
        今優先して対応すべきことを提案します。
      </p>

      <AnalyzeForm />
    </main>
  );
}