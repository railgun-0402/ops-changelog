# ops-changelog

## 目的

障害時に「直近でこのサービス（リポジトリ）にどんな変更が入ったか」を素早く把握できるCLIツール。

## コンセプト

- シンプルに始める：まずは `repo / service / since` の3つの入力で動く
- GitHub PRのラベルでサービスを判定する（例: `service:apply-api`）
- 出力は人間が読みやすい形式とJSON形式の両対応を目指す

## 現在の要件（v0.1）

### 入力

| フラグ | 説明 | 例 |
|--------|------|-----|
| `--repo` | GitHubリポジトリ（owner/repo形式） | `--repo myorg/myrepo` |
| `--service` | サービス名（ラベル `service:<name>` で判定） | `--service apply-api` |
| `--since` | いつ以降のPRを取得するか（duration形式） | `--since 24h`, `--since 7d` |

### 出力

- merged PR一覧
- 各PRに: `merged_at`, `title`, `url`, `author`, `labels`

### ラベル判定ルール

- ラベル名が `service:<service名>` にマッチするPRを対象とする
- 例: `--service apply-api` → ラベル `service:apply-api` を持つPRのみ表示

## 将来の要件（検討中）

- [ ] 複数サービスの同時指定
- [ ] 複数リポジトリの同時指定
- [ ] JSON出力オプション（`--format json`）
- [ ] Slack通知連携
- [ ] .ops-changelog.yaml 設定ファイルでリポジトリ・サービスのマッピングを管理
- [ ] コミット単位の変更も見れるようにする
- [ ] GitHub Actions での定期実行サポート

## アーキテクチャ

```
ops-changelog/
  cmd/
    root.go       # CLIエントリポイント（cobra）
    list.go       # list サブコマンド
  internal/
    github/
      client.go   # GitHub API クライアント
      pr.go       # PR取得・フィルタリングロジック
    formatter/
      text.go     # テキスト出力フォーマッター
  main.go
  go.mod
  CLAUDE.md
  .claude/
    skills/
```

## 認証

- 環境変数 `GITHUB_TOKEN` を使用
- なければ `gh auth token` コマンドにフォールバック

## 開発方針

- Go標準 + 最小限の外部依存（cobra, google/go-github）
- エラーメッセージは障害対応者向けに明確に
- テストは `internal/` 配下のロジック層を優先的にカバー

## コマンド例

```bash
# 直近24時間にapply-apiサービスへマージされたPRを確認
ops-changelog list --repo myorg/myrepo --service apply-api --since 24h

# 直近7日間
ops-changelog list --repo myorg/myrepo --service apply-api --since 7d
```