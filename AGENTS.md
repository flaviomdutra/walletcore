# Repository Guidelines

## Project Structure & Module Organization

This Go module implements wallet domain logic under `internal/` and event dispatching under `pkg/`:

- `internal/entity/`: clients, accounts, transactions, and domain validation.
- `internal/gateway/`: persistence interfaces used by application logic.
- `internal/database/`: SQL database adapters for clients, accounts, and transactions.
- `internal/usecase/create_client/`, `create_account/`, and `create_transaction/`: application workflows and input/output DTOs.
- `pkg/events/`: event and handler interfaces plus an in-memory event dispatcher.
- Tests live beside the code in `*_test.go` files.

Keep domain rules in entities and coordinate persistence through gateway interfaces in use cases. Database adapters implement those interfaces. The event dispatcher registers handlers by event name, prevents duplicate registrations, checks registrations with `Has`, and dispatches events to registered handlers. There is no executable entry point or server command.

## Build, Test, and Development Commands

Use the Go version declared in `go.mod` (currently `1.27.3`). Run standard Go commands directly.

- `go build ./...`: compile all packages.
- `go test ./...`: run all tests.
- `go test ./internal/entity -run TestCreditAccount`: run one focused test.
- `go test -cover ./...`: report package coverage.
- `go vet ./...`: run Go static checks.
- `go fmt ./...`: format Go sources.

There is currently no application command for running a server locally.

## Coding Style & Naming Conventions

Use standard Go formatting with tabs and `gofmt`. Follow existing exported names such as `NewClient`, `CreateClientUseCase`, and `CreateClientInputDTO`. Use lowercase package names; existing use-case directories use underscores while package names are joined, such as `create_client` and `createclient`. Return errors explicitly and keep dependencies behind gateway interfaces.

## Testing Guidelines

Use Go's `testing` package with Testify assertions and mocks. Name tests `Test<Behavior>` or `Test<UseCase>_Execute`. Cover validation failures, balance changes, and gateway errors when changing behavior. Verify mock expectations for persistence calls. No coverage threshold is configured.

## Commit & Pull Request Guidelines

Git history is unavailable in this checkout, so no established commit convention can be verified. Use concise, imperative commit subjects. Pull requests should describe the behavior changed, link relevant issues, and list validation results and any checks that could not run. Keep changes focused and preserve unrelated work.
