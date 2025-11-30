# go-rule-updater

go-rule-updater is a small automation tool that appends a new rule to `rule.yaml`, opens a GitHub branch, creates a pull request with a comment, and merges it automatically. It is distributed under the MIT License.

## How it works

1. Load configuration from environment variables (dotenv supported; see below) and build the GitHub client (prefers GitHub App via GCP KMS or PEM key, falls back to `GITHUB_TOKEN`).
2. Create a new branch (`add/<uuid>`) from `BASE_BRANCH`.
3. Read the target `rule.yaml` (or any path set in `RULE_PATH`).
4. Append a UUID-based rule entry while preserving existing YAML comments.
5. Commit using the GraphQL `createCommitOnBranch` mutation, then open a PR with a comment.
6. Merge the PR automatically, retrying until the PR becomes mergeable or `MERGE_MAX_RETRIES` is hit.
7. Detect any open `add/*` PRs that are conflicting and **rebuild** them on top of the latest base branch by replaying their added rules and force-updating their branches.

## Requirements

- Go 1.25 or later
- GitHub token for development (`gh auth token` is convenient)
- GitHub App credentials for production runs

## Installation

```bash
git clone https://github.com/hi120ki/go-rule-updater.git
cd go-rule-updater
go mod download
```

## Configuration

| Variable                     | Default           | Purpose                                                                                 |
| ---------------------------- | ----------------- | --------------------------------------------------------------------------------------- |
| `ENVIRONMENT`                | `development`     | `development` uses `GITHUB_TOKEN`; `production` uses GitHub App creds.                  |
| `OWNER`                      | `hi120ki`         | Repository owner.                                                                       |
| `REPOSITORY`                 | `go-rule-updater` | Target repository name.                                                                 |
| `RULE_PATH`                  | `rule.yaml`       | File path to update.                                                                    |
| `BASE_BRANCH`                | `main`            | Branch used as the base for new branches.                                               |
| `MERGE_MAX_RETRIES`          | `6`               | Number of attempts while waiting for a PR to become mergeable.                          |
| `MERGE_RETRY_DELAY_SECONDS`  | `10`              | Wait time (seconds) between merge attempts.                                             |
| `GITHUB_TOKEN`               | none              | Required in development mode (PAT with repo scope).                                     |
| `GITHUB_APP_ID`              | none              | GitHub App ID (production).                                                             |
| `GITHUB_APP_INSTALLATION_ID` | none              | Installation ID (production).                                                           |
| `GITHUB_APP_KMS_KEY_PATH`    | none              | GCP KMS key path for signing GitHub App tokens (production).                            |
| `GITHUB_APP_PRIVATE_KEY`     | none              | Path to the PEM-encoded private key file for the GitHub App (production, non-KMS flow). |

Example development run:

```bash
GITHUB_TOKEN=$(gh auth token) go run main.go
```

## Commands

- `make run`: Run the full automation flow locally (uses `gh auth token`).
- `make test`: Run all tests (requires `GITHUB_TOKEN` for GitHub client initialization).
- `go fmt ./... && go vet ./...`: Formatting and basic static checks.

## Rule file format

`rule.yaml` holds a list under `config`. Example:

```yaml
config:
  - name: first-rule
  - name: another-rule
```

Each execution appends a new `name` entry (UUID string by default) while preserving comments and indentation.

## Development notes

- Key packages: `env/` (env loading), `github/` (REST + GraphQL helpers for branches/PRs/commits), `rule/` (YAML mutation with comment preservation).
- Tests sit next to code in `*_test.go`; many are table-driven. Run `go test ./...` after changes.
- Keep code `gofmt`/`goimports` clean and avoid committing secrets.

## License

MIT License. See `LICENSE` for details.
