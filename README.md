# ops-changelog

A CLI tool to quickly see what changed in a service — useful during incidents.

When something breaks in production, the first question is often: _"What changed recently?"_
`ops-changelog` answers that by listing merged PRs for a specific service, filtered by GitHub labels.

## Features

- List merged PRs for a service filtered by label (`service:<name>`)
- Filter by time range (`--since 1h`, `--since 7d`, etc.)
- Show PR title, URL, author, merged date, and labels
- Omit `--service` to list all merged PRs in a repository
- Authenticates via `GITHUB_TOKEN` or `gh auth login`

## Installation

**Using Go:**

```bash
go install github.com/railgun-0402/ops-changelog@latest
```

**From GitHub Releases:**

Download the binary for your platform from the [Releases page](https://github.com/railgun-0402/ops-changelog/releases).

## Usage

```bash
# List merged PRs for a specific service in the last 24 hours
ops-changelog list --repo myorg/myrepo --service apply-api --since 24h

# List all merged PRs in the last 7 days
ops-changelog list --repo myorg/myrepo --since 7d
```

**Output example:**

```
Merged PRs for service:apply-api in myorg/myrepo (since 2026-03-14 09:00 UTC)
────────────────────────────────────────────────────────────────────────
[2026-03-14 10:30]  #123  Fix null pointer in payment handler
  Author : alice
  URL    : https://github.com/myorg/myrepo/pull/123
  Labels : service:apply-api, fix
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--repo` | GitHub repository in `owner/repo` format (required) | |
| `--service` | Service name — matches label `service:<name>` | all services |
| `--since` | How far back to look (`1h`, `24h`, `7d`, `30d`) | `24h` |

## License

MIT
