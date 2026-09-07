type AnalysisResultProps = {
  result: string;
  loading: boolean;
  error: string;
};

export default function AnalysisResult({
  result,
  loading,
  error,
}: AnalysisResultProps) {
  return (
    <section className="analysisPanel">
      <div className="sectionHeader">
        <div>
          <span className="sectionEyebrow">
            AI Analysis
          </span>

          <h2>最優先アクション</h2>
        </div>

        <span className="sectionNumber">
          02
        </span>
      </div>

      <p className="sectionDescription">
        Repositoryの状態から、今もっとも優先すべき対応を提案します。
      </p>

      {!loading &&
        !result &&
        !error && (
          <EmptyState />
        )}

      {loading && (
        <LoadingState />
      )}

      {error && (
        <div className="errorCard">
          <span className="errorLabel">
            ANALYSIS FAILED
          </span>

          <h3>
            分析できませんでした
          </h3>

          <p>{error}</p>
        </div>
      )}

      {result &&
        !loading && (
          <div className="resultArea">
            <div className="priorityHeader">
              <div>
                <span className="priorityLabel">
                  Highest Priority
                </span>

                <h3>
                  AIが今もっとも重要と判断した内容
                </h3>
              </div>

              <div className="priorityState">
                Priority
              </div>
            </div>

            <div className="resultCard">
              {/*
                現時点ではGemini結果を文字列表示。
                後で構造化レスポンスに変える予定。
              */}
              <div className="resultText">
                {result}
              </div>
            </div>
          </div>
        )}
    </section>
  );
}

function EmptyState() {
  return (
    <div className="emptyState">
      <div className="emptyIcon">
        <span />
        <span />
        <span />
      </div>

      <h3>
        Repositoryを分析してください
      </h3>

      <p>
        Pull Request、コード差分、CI、
        ブランチ構成をAIが確認します。
      </p>
    </div>
  );
}

function LoadingState() {
  return (
    <div className="analysisLoading">
      <div className="loadingHeader">
        <span className="spinnerDark" />

        <div>
          <h3>
            Repositoryを分析中
          </h3>

          <p>
            開発上の重要なシグナルを確認しています。
          </p>
        </div>
      </div>

      <div className="loadingChecks">
        <LoadingRow label="Pull Requests" />
        <LoadingRow label="Code changes" />
        <LoadingRow label="Branch dependencies" />
        <LoadingRow label="CI / Check Runs" />
      </div>
    </div>
  );
}

type LoadingRowProps = {
  label: string;
};

function LoadingRow({
  label,
}: LoadingRowProps) {
  return (
    <div className="loadingRow">
      <span className="loadingRowDot" />
      <span>{label}</span>
    </div>
  );
}