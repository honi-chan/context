import AnalyzeForm from "@/components/AnalyzeForm";

export default function Home() {
  return (
    <main className="page">
      <header className="appHeader">
        <div className="appHeaderInner">
          <div className="brand">
            <div className="brandMark">A</div>

            <div className="brandText">
              <span className="brandName">AI Work Assistant</span>
              <span className="brandDescription">
                Repository Intelligence
              </span>
            </div>
          </div>

          <div className="headerStatus">
            <span className="statusIndicator" />
            Ready
          </div>
        </div>
      </header>

      <div className="pageContainer">
        <section className="hero">
          <div className="heroLabel">
            GitHub Repository Analysis
          </div>

          <h1>
            今、何を優先すべきかを
            <br />
            AIが判断する。
          </h1>

          <p>
            Pull Request、変更内容、ブランチ依存、CI情報を分析し、
            開発者が次に対応すべきことを1つに絞って提案します。
          </p>
        </section>

        <AnalyzeForm />
      </div>
    </main>
  );
}