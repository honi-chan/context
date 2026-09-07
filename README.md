# AI Work Assistant

> **Know what matters before you even ask.**

AI Work Assistant は、開発者の作業コンテキストを理解し、
**「今、何に対応すべきか」をAIが先回りして提案する開発支援ツール**です。

## Concept

現在のAIツールでは、多くの場合、人間が情報を集めてAIに質問する必要があります。

```text
問題が発生
    ↓
人間が気づく
    ↓
GitHub / CI / Logs などを確認
    ↓
必要な情報を集める
    ↓
AIに質問する
    ↓
AIが回答する
```

AI Work Assistant が目指しているのは、その逆です。

```text
GitHub / CI / Logs / Monitoring
              ↓
        AI Work Assistant
              ↓
      状況を継続的に理解
              ↓
「今これを確認した方がいい」
              ↓
            Human
```

**人間がAIのところへ情報を持っていくのではなく、AIが人間の仕事の場所へ来る。**

これがこのプロジェクトの中心となる考え方です。

---

## Current MVP

最初のMVPでは **GitHub Repository Analysis** に対象を絞っています。

Repositoryを指定すると、GitHubからPull Requestに関する情報を取得し、Geminiが分析します。

現在利用している主なシグナル：

* Open Pull Requests
* PR title / description
* Head / Base branch
* Changed files
* Additions / deletions
* Code diff (`patch`)
* GitHub Check Runs / CI status

これらをもとにAIが、

> **今もっとも優先して確認・対応すべきこと**

を1つだけ選びます。

分析結果には、

* 対象Pull Request
* 優先すべき理由
* 次に取るべき行動

を含めます。

---

## Architecture

```text
┌──────────────────────┐
│       Next.js        │
│      Frontend        │
└──────────┬───────────┘
           │
           │ HTTP
           ▼
┌──────────────────────┐
│       Go / Echo      │
│       Backend        │
└──────────┬───────────┘
           │
           ├─────────────────┐
           │                 │
           ▼                 ▼
┌─────────────────┐  ┌─────────────────┐
│   GitHub API    │  │   Gemini API    │
│                 │  │                 │
│ Repository Data │  │   AI Analysis   │
└─────────────────┘  └─────────────────┘
```

Backendでは外部APIのモデルを直接UseCaseへ持ち込まず、

```text
GitHub API
    ↓
GitHub Adapter
    ↓
Mapper
    ↓
Domain Model
    ↓
UseCase
    ↓
AI Analyzer
```

という形で責務を分離しています。

---

## Tech Stack

### Frontend

* TypeScript
* Next.js
* React

### Backend

* Go
* Echo

### AI

* Gemini API

### External Services

* GitHub REST API

### Development

* Docker
* Docker Compose

---

## Design Principles

### Architecture

設計・リファクタリングでは Refactoring.Guru の考え方を参考にしています。

特定のDesign Patternを使うこと自体を目的にはせず、
**変更しやすさ・責務の明確さ・外部サービスとの疎結合**を優先します。

Examples:

* Adapter
* Mapper
* Dependency Inversion
* Separation of Concerns

### UI

UIのデザインシステムは Google Labs の `DESIGN.md` specification を基準に管理します。

`DESIGN.md` に、

* Colors
* Typography
* Spacing
* Rounded corners
* Components
* Design rationale

を定義し、Frontend実装の基準とします。

---

## Project Structure

```text
.
├── backend
│   ├── cmd
│   │   └── server
│   └── internal
│       ├── ai
│       ├── domain
│       ├── handler
│       ├── infrastructure
│       │   ├── gemini
│       │   └── github
│       ├── repository
│       └── usecase
│
├── frontend
│   ├── DESIGN.md
│   └── src
│       ├── app
│       ├── components
│       └── lib
│
└── docker-compose.yml
```

---

## Getting Started

必要な環境変数を設定します。

```env
GITHUB_TOKEN=your_github_token
GEMINI_API_KEY=your_gemini_api_key
GEMINI_MODEL=your_gemini_model
```

Frontend:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

Docker Composeで起動します。

```bash
docker compose up --build
```

Frontend:

```text
http://localhost:3000
```

---

## Roadmap

### Phase 1 — GitHub Intelligence

GitHub Repositoryから開発状況を理解する。

* Pull Request analysis
* Code diff analysis
* CI / Check Runs
* Branch dependency detection
* Review status

### Phase 2 — Work Context

GitHub以外の開発コンテキストを統合する。

```text
GitHub
Slack
AWS
New Relic
CI/CD
   ↓
AI Work Assistant
```

複数サービスの情報を組み合わせ、問題の関連性まで判断できるようにします。

### Phase 3 — Proactive Assistant

最終的には、人間が質問する前にAIが重要な変化を検出します。

例えば、

> CIが失敗しています。
> 30分前にマージされたPRの変更と関連している可能性があります。

のように、AI側から開発者へ知らせることを目指します。

---

## Vision

AIの能力が高くなっても、
人間が毎回コンテキストを集めて質問しなければならないのであれば、仕事の流れそのものは大きく変わりません。

AI Work Assistant が目指すのは、

**AIを「質問すると答えるツール」から「仕事を理解している存在」へ変えること。**

開発者がツールを巡回して問題を探すのではなく、
重要なことだけが開発者のところへ届く世界を目指します。
