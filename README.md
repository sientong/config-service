# Configuration Management Service (Golang)

## Overview
A small RESTful configuration management service implemented in Go with SQLite persistence and JSON Schema validation. It supports create, update, rollback, fetch (latest or specific version), and list versions.

Schemas are hardcoded and validated with JSON Schema. Persistence uses SQLite (file `config_database.db`).

## Run

### Running in Local Machine

Requirements: Go 1.20+, `make`.

```bash
make build
make run
make test
```

Or directly:

```bash
go run ./...
```

### Running as Docker image

Requirements: docker, `make`.

```bash
make docker-build
make docker-run
make docker-stop
```

API is available at `http://localhost:3000`.

## Makefile targets
- `build` - build executable
- `run` - run service locally
- `test` - run unit and integration tests
- `clean` - delete local build files
- `docker-build` - build docker image
- `docker-run` - run docker image as background process
- `docker-stop` - stop running docker image

## OpenAPI
See `openapi.yaml` for API contract.

## Testing API
See `test.http` for API testing.

## Notes
- Persistence: SQLite file `config_database.db` created in working directory.
- Validation: JSON Schema (hardcoded in `schema.go`).
