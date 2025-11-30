# go-rule-updater

go-rule-updater is a small automation tool that appends a new rule to `rule.yaml`, opens a GitHub branch, creates a pull request with a comment, and merges it automatically. It is distributed under the MIT License.

## How it works

1. Load configuration from environment variables (see below).
2. Create a new branch from `BASE_BRANCH`.
3. Read the target `rule.yaml` (or any path set in `RULE_PATH`).
4. Append a UUID-based rule entry while preserving existing YAML comments.
5. Commit, open a PR with a comment.
6. **Automatically detect and update any existing conflicting PRs** by rebasing them with the latest base branch.
7. (Optional) Merge the PR automatically (currently commented out).

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

| Variable                     | Default           | Purpose                                                                |
| ---------------------------- | ----------------- | ---------------------------------------------------------------------- |
| `ENVIRONMENT`                | `development`     | `development` uses `GITHUB_TOKEN`; `production` uses GitHub App creds. |
| `OWNER`                      | `hi120ki`         | Repository owner.                                                      |
| `REPOSITORY`                 | `go-rule-updater` | Target repository name.                                                |
| `RULE_PATH`                  | `rule.yaml`       | File path to update.                                                   |
| `BASE_BRANCH`                | `main`            | Branch used as the base for new branches.                              |
| `GITHUB_TOKEN`               | none              | Required in development mode (PAT with repo scope).                    |
| `GITHUB_APP_ID`              | none              | GitHub App ID (production).                                            |
| `GITHUB_APP_INSTALLATION_ID` | none              | Installation ID (production).                                          |
| `GITHUB_APP_PRIVATE_KEY`     | none              | PEM-encoded private key content (production).                          |

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
