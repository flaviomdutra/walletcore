# Repository Guidelines

## Project Structure & Module Organization

This Go module implements wallet domain logic under `internal/`:

- `internal/entity/`: clients, accounts, transactions, and domain validation.
- `internal/gateway/`: persistence interfaces used by application logic.
- `internal/usecase/create_client/`, `create_account/`, and `create_transaction/`: application workflows and input/output DTOs.
- Tests live beside the code in `*_test.go` files.

Keep domain rules in entities and coordinate persistence through gateway interfaces in use cases. No executable entry point, assets directory, or concrete persistence adapter currently exists.

## Build, Test, and Development Commands

Use mise for development tools and match the Go version declared in `go.mod` (currently `1.27.3`). No repository mise configuration or tasks are checked in; configure the matching runtime before running:

- `mise exec -- go build ./...`: compile all packages.
- `mise exec -- go test ./...`: run all tests.
- `mise exec -- go test ./internal/entity -run TestCreditAccount`: run one focused test.
- `mise exec -- go test -cover ./...`: report package coverage.
- `mise exec -- go vet ./...`: run Go static checks.
- `mise exec -- go fmt ./...`: format Go sources.

There is currently no application command for running a server locally.

## Coding Style & Naming Conventions

Use standard Go formatting with tabs and `gofmt`. Follow existing exported names such as `NewClient`, `CreateClientUseCase`, and `CreateClientInputDTO`. Use lowercase package names; existing use-case directories use underscores while package names are joined, such as `create_client` and `createclient`. Return errors explicitly and keep dependencies behind gateway interfaces.

## Testing Guidelines

Use Go's `testing` package with Testify assertions and mocks. Name tests `Test<Behavior>` or `Test<UseCase>_Execute`. Cover validation failures, balance changes, and gateway errors when changing behavior. Verify mock expectations for persistence calls. No coverage threshold is configured.

## Commit & Pull Request Guidelines

Git history is unavailable in this checkout, so no established commit convention can be verified. Use concise, imperative commit subjects. Pull requests should describe the behavior changed, link relevant issues, and list validation results and any checks that could not run. Keep changes focused and preserve unrelated work.
