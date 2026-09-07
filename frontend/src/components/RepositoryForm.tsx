"use client";

import {
  FormEvent,
  useState,
} from "react";

type RepositoryFormProps = {
  loading: boolean;

  // 親側へowner/repoを渡すだけ。
  // API処理はここでは持たない。
  onAnalyze: (
    owner: string,
    repo: string,
  ) => Promise<void>;
};

export default function RepositoryForm({
  loading,
  onAnalyze,
}: RepositoryFormProps) {
  const [owner, setOwner] = useState("");
  const [repo, setRepo] = useState("");

  const handleSubmit = async (
    event: FormEvent<HTMLFormElement>,
  ) => {
    event.preventDefault();

    await onAnalyze(
      owner,
      repo,
    );
  };

  return (
    <aside className="repositoryPanel">
      <div className="sectionHeader">
        <div>
          <span className="sectionEyebrow">
            Repository
          </span>

          <h2>分析対象</h2>
        </div>

        <span className="sectionNumber">
          01
        </span>
      </div>

      <p className="sectionDescription">
        GitHub Owner と Repository を指定してください。
      </p>

      <form
        className="repositoryForm"
        onSubmit={handleSubmit}
      >
        <label className="field">
          <span className="fieldLabel">
            GitHub Owner
          </span>

          <input
            className="textInput"
            value={owner}
            onChange={(event) =>
              setOwner(event.target.value)
            }
            placeholder="honi-chan"
            autoComplete="off"
            required
          />
        </label>

        <label className="field">
          <span className="fieldLabel">
            Repository
          </span>

          <input
            className="textInput"
            value={repo}
            onChange={(event) =>
              setRepo(event.target.value)
            }
            placeholder="context"
            autoComplete="off"
            required
          />
        </label>

        <button
          className="primaryButton"
          type="submit"
          disabled={loading}
        >
          {loading ? (
            <>
              <span className="spinner" />
              分析中...
            </>
          ) : (
            "GitHubを分析"
          )}
        </button>
      </form>

      <div className="integrationStatus">
        <div className="integrationRow">
          <span className="integrationName">
            GitHub
          </span>

          <span className="integrationValue">
            <span className="integrationDot" />
            Connected
          </span>
        </div>

        <div className="integrationRow">
          <span className="integrationName">
            Gemini
          </span>

          <span className="integrationValue">
            <span className="integrationDot" />
            Ready
          </span>
        </div>
      </div>
    </aside>
  );
}