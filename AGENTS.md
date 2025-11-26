# Repository Guidelines

## Project Structure & Module Organization

- `main.go`: CLI entrypoint that loads env vars, opens a GitHub client, creates a branch, updates the rule file, and opens/merges a PR.
- `env/`: Environment loading via `envconfig`; defaults exist for `OWNER`, `REPOSITORY`, `RULE_PATH` (`rule.yaml`), and `BASE_BRANCH`.
- `github/`: REST + GraphQL helpers for branches, commits, files, refs, and PR flows; paired `*_test.go` files show expected API interactions.
- `rule/`: YAML manipulation that appends config entries while preserving comments using `goccy/go-yaml`; validated by table-driven tests.
- `rule.yaml`: Default rule store; structure is a `config` array of `name` values and can be relocated by setting `RULE_PATH`.

## Build, Test, and Development Commands

- `make run` (or `GITHUB_TOKEN=$(gh auth token) go run main.go`): Executes the end-to-end update flow against the configured repo.
- `make test` (or `GITHUB_TOKEN=$(gh auth token) go test ./...`): Runs all unit tests; token is required because some tests hit GitHub APIs.
- `go fmt ./... && go vet ./...`: Format and basic static analysis before opening a PR.
- Production mode: set `ENVIRONMENT=production` and provide `GITHUB_APP_ID`, `GITHUB_APP_INSTALLATION_ID`, and `GITHUB_APP_PRIVATE_KEY`.

## Coding Style & Naming Conventions

- Go 1.25; keep files `gofmt`-clean (tabs, trailing newline) and run `goimports` if you add imports.
- Exported identifiers use CamelCase; prefer clear, imperative function names (e.g., `CreatePullRequest`, `Add`).
- Keep package-scoped variables minimal; favor passing context explicitly as in existing functions.
- YAML additions should use the helpers in `rule/` to preserve ordering and comments.

## Testing Guidelines

- Tests live alongside code in `*_test.go`; name tests `TestXxx_Subject` for clarity.
- Prefer table-driven tests for permutations (see `rule/rule_test.go`); mock GitHub interactions via existing patterns in `github/` tests.
- Run `go test ./...` before commits; add coverage when modifying request/response handling or YAML mutation logic.

## Commit & Pull Request Guidelines

- Use concise, imperative commit subjects (history shows short action verbs like “Add new rule”; add context such as scope when possible).
- PRs should state what changed, why, and how it was validated (commands run, env vars used). Link related issues and include screenshots only if UI output is affected (rare).
- Ensure CI-like checks locally (`go fmt`, `go vet`, `go test ./...`) and note any required secrets or GitHub App settings in the description.

## Security & Configuration Tips

- Never commit tokens or private keys; use `gh auth token` for local runs and environment variables for automation.
- Limit repo/branch scope via `OWNER`, `REPOSITORY`, and `BASE_BRANCH` when testing; use disposable branches for experiments.
