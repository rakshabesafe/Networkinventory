# Instructions for AI Agents

Welcome to the Cloud-Native Network Inventory System repository. When assisting with this codebase, please adhere to the following guidelines:

## Code Organization & Go Workspace
1. **Monorepo Structure:** This repository uses a Go Workspace (`go.work`). It is not a single Go module. The root directory contains `go.work`, while individual microservices (e.g., `services/tmf638`, `services/tmf639`) and shared libraries (e.g., `pkg`) are separate modules.
2. **Adding Dependencies:** When adding a new Go package dependency, ensure you run `go get` within the specific module directory that requires it (e.g., `cd services/tmf639 && go get ...`), not at the repository root.
3. **Shared Code:** Any structs, models, or utility functions that need to be accessed by multiple microservices MUST be placed in the `pkg/` directory.

## Database Interactions (Neo4j)
1. **Fallback Mechanism:** All HTTP Handlers currently implement a fallback mechanism. If the `neo4jDB` instance is `nil` (meaning the DB failed to connect on startup), the handlers MUST fall back to reading/writing from thread-safe mock maps (e.g., `mockResources`, `mockMutex`). Ensure any new handlers maintain this dual-path logic to support testing without a live database.
2. **Cypher Queries:** Always use parameterized Cypher queries to prevent injection attacks.
3. **Idempotency:** When creating new entities in the database, prefer the `MERGE` Cypher command over `CREATE` to maintain idempotency based on the entity's ID.

## Testing & Pre-Commit Validation
1. **Running Tests:** Always run tests from the root of the repository using `go test ./...`. This command will test all modules in the workspace.
2. **Graceful Degradation:** Tests should not fail if a local Neo4j instance is not running. The fallback mock maps in the handlers ensure that HTTP endpoint tests pass regardless of database availability. Do not remove this capability.
3. **Linting and Vetting:** Before submitting changes, it is highly recommended to run `go vet ./...` and ensure standard formatting using `gofmt`.

## Deployment Architecture
1. **Docker Compose:** The primary local deployment method is via `docker-compose.yml`.
2. **Healthchecks:** Microservices in the docker-compose file depend on the `neo4j` container being explicitly marked as `healthy` (via cypher-shell healthcheck) before starting. Do not remove this dependency condition.
3. **Go Version:** The Go version specified in the `Dockerfile`s (`golang:1.24-alpine`) MUST match or exceed the version specified in the `go.mod` files across the workspace to prevent automated toolchain download failures during container builds.