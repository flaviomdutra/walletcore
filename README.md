# WalletCore

WalletCore is a Go project that implements core wallet operations: creating clients, creating accounts, and transferring funds between accounts. It separates domain rules from application workflows and persistence interfaces.

## Features

- Create clients with required name and email validation.
- Create accounts associated with clients.
- Credit and debit account balances.
- Transfer funds with positive-amount and sufficient-balance checks.
- Coordinate persistence through injectable gateway interfaces.
- Persist clients, accounts, and transactions through SQL database adapters.
- Register event handlers and dispatch events to handlers by event name.

## Project Structure

```text
internal/
├── entity/                     # Domain models, validation, and unit tests
├── gateway/                    # Client, account, and transaction interfaces
├── database/                   # SQL persistence adapters
└── usecase/
    ├── create_client/          # Client creation workflow
    ├── create_account/         # Account creation workflow
    └── create_transaction/     # Transfer workflow
```

pkg/
└── events/                     # Event interfaces and in-memory dispatcher

Use cases accept input DTOs, operate on entities, and return output DTOs. Database adapters provide SQL-backed gateway implementations. The event dispatcher supports handler registration, duplicate-registration checks, membership checks, clearing registrations, and dispatch by event name. Tests live alongside their source files and use Go's testing package with Testify assertions and mocks.

## Getting Started

Install the Go version specified in `go.mod` (currently `1.27.3`) and ensure `go` is available on your `PATH`. From the repository root, download dependencies and build the packages:

```sh
go mod download
go build ./...
```

The HTTP API is available from `cmd/walletcore`.

Start MySQL and the API from the repository root:

```sh
docker compose up -d mysql
go run ./cmd/walletcore
```

The server connects to MySQL at `127.0.0.1:3306` by default. Set `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, or `DB_NAME` to override the connection settings. The Compose initialization script creates the database tables when MySQL initializes an empty data directory.

## Development

Run these commands from the repository root:

```sh
# Run all tests
go test ./...

# Report test coverage
go test -cover ./...

# Run static checks
go vet ./...

# Format Go sources
go fmt ./...
```

## Current Limitations

Balances use `float64`, so monetary precision needs consideration before production use. Creating a transaction immediately updates both account balances in memory; the use case then saves the transaction through its gateway. Durable balance updates, atomic persistence, and rollback on persistence failure are not implemented.

## Contributing

See [AGENTS.md](AGENTS.md) for repository conventions, testing guidance, and pull request expectations. Keep local credentials and configuration out of version control; `.gitignore` includes common environment and credential files.
