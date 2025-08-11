# Configuration Management Service (Golang)

## Overview

This is a small RESTful configuration management service implemented in Go with:

- SQLite persistence for storing configurations and their versions.
- JSON Schema validation to enforce structure and types of stored configs.
- Versioning support (create, update, rollback to previous versions).
- Dynamic schema loading from JSON files in the schemas folder at startup.

The service is designed for scenarios where configuration data needs to be version-controlled, validated, and retrieved efficiently.

## Setup and Running the Application

Prerequisites:
- Go 1.20+
- Make
- Docker (optional for containerized deployment)

Local Development

### Build

```bash
make build
```

### Run locally

```bash
make run
```

### Run tests

```bash
make test
```

Or directly with Go:

```bash
go run .
```

### Running with Docker

#### Build the Docker image

```bash
make docker-build
```

#### Run in background

```bash
make docker-run
```

#### Stop running container
```bash
make docker-stop
```

Once running, the API will be available at:
http://localhost:3000/swagger/index.html


## API Documentation

*Swagger/OpenAPI*

The API is documented in OpenAPI 3.0 format and served by Swagger UI at runtime.

- Swagger UI: http://localhost:3000/swagger/index.html
- Raw OpenAPI file: docs/swagger.yaml

Endpoints include:

- POST `/configs/{schema}/{name}` – Create a config
- PUT `/configs/{schema}/{name}` – Update a config
- POST `/configs/{schema}/{name}/rollback` – Rollback to previous version
- GET `/configs/{schema}/{name}` – Fetch latest or specific config version
- GET `/configs/{schema}/{name}/versions` – List all versions

- GET `/schemas` - List of stored schema
- GET `/schemas/{schema}` - Display individual schema

## Schema Explanation

Schemas define the structure, allowed types, and constraints for each configuration type. Schemas are stored as `.json` files under the schemas directory. At service startup, all schema files are loaded into memory and used for validating incoming requests.

- Example: `schemas/payment_config.json`

```json
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "properties": {
        "max_limit": { "type": "integer" },
        "enabled": { "type": "boolean" }
    },
    "required": ["max_limit", "enabled"],
    "additionalProperties": false
}
```

## Design Decisions and Trade-offs

1. SQLite for Persistence
- Chosen for simplicity and zero external dependencies.
- Trade-off: Not ideal for high-concurrency, large-scale workloads.
	
2.	JSON Schema Validation
- Ensures configs follow strict structure before persistence.
- Trade-off: Requires upfront schema design and may reject valid but unstructured data.

3. Versioning
- Every change creates a new version for auditability and rollback.
- Trade-off: Requires more storage and careful version handling logic.

4. Dynamic Schema Loading
- Allows adding new config types without code changes.
- Trade-off: Missing or invalid schema files will prevent related config operations.

5. Containerization
- Uses multi-stage Docker build to avoid runtime library mismatches.
- Trade-off: Larger image than pure static Go binary if CGO is enabled.

## Testing API

Use the included `test.http` file for API testing via VSCode REST Client or similar tools.

## Ideas for Improvement

1. **Implement Authorization and Authentication**

Implement auth so that only authorized and authenticated user can create or update config.

2. **Per-Environment Configs**

Have separate configs for dev, staging, and prod with fallbacks.

3. **Dynamic Schema Reloading**

Watch the schema directory (e.g., via `fsnotify`) and reload schemas automatically without restarting the service.

4. **Preview Changes**

Allow a “dry-run” update that validates the new config but doesn’t save it yet.

5. **Distributed Cache**

Use Redis or Memcached to cache configs and schemas across multiple instances.