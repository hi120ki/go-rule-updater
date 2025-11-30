# Repository Guidelines

## Project Structure & Modules

- `main.go`: entrypoint that wires environment, GitHub client, and service flow.
- `service/service.go`: orchestration for branch creation, rule append, PR handling, and optional merge.
- `env/`: configuration loading from environment variables for development vs production.
- `github/`: REST/GraphQL helpers for branches, refs, commits, files, and PR operations.
- `rule/`: YAML mutation while preserving comments; see `rule.yaml` as the default target file.
- Tests live beside code in `*_test.go`; use `rule/rule_test.go` and `github/*_test.go` for examples.
- Top-level `Makefile` contains repeatable commands; `.editorconfig` captures indentation rules.

## Build, Test, and Development Commands

- `make run`: execute the full flow locally using `GITHUB_TOKEN`/`gh auth token`.
- `make test`: run Go tests; GitHub client initialization requires `GITHUB_TOKEN`.
- `go test ./...`: standard test run (set GitHub-related env vars as needed).
- `go fmt ./... && go vet ./...`: format and static checks; keep the repo clean before pushing.
- `GITHUB_TOKEN=$(gh auth token) go run main.go`: quick ad-hoc run in development mode.

## Coding Style & Naming Conventions

- Go 1.25+, format with `gofmt` (tabs for Go per `.editorconfig`); LF line endings, final newline required.
- Keep package names short and lower-case; exported identifiers use Go’s PascalCase; tests named `TestXxx`.
- Branches typically use `add/<rule-id>` or similar imperative slugs; keep them short and descriptive.
- Avoid committing generated secrets or tokens; prefer environment variables over config files.

## Testing Guidelines

- Use Go’s standard testing package; follow table-driven patterns present in existing tests.
- Name tests after the function under test (e.g., `TestAppendRule`); co-locate fixtures near the test file.
- When adding behavior touching GitHub calls, prefer interface seams or fakes to avoid network reliance.
- Aim to keep `make test` green; add regression tests alongside fixes.

## Commit & Pull Request Guidelines

- Commit messages: short, imperative summaries (examples in history: “Add new rule”, “fix”); include scope when helpful.
- PRs: describe the change, note relevant env vars or assumptions, link issues if any, and call out test coverage.
- Screenshots/log snippets are useful when altering automation behavior or YAML mutations.
- Keep branches rebased onto `BASE_BRANCH`; clean up debugging output before opening the PR.

## Configuration & Security Notes

- Required env vars: `ENVIRONMENT`, `OWNER`, `REPOSITORY`, `RULE_PATH`, `BASE_BRANCH`, and either `GITHUB_TOKEN` (dev) or GitHub App credentials (production).
- Never commit tokens or private keys; load via shell env or secret stores.
- Default rule target is `rule.yaml`; adjust `RULE_PATH` if you add additional rule files in other repos.
