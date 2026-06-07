# Backend (Go)

Barebones backend for the Spond membership wizard task.

## Implemented

- `GET /api/v1/forms/public`
  - Returns the form configuration contract from the assignment appendix.
- `POST /api/v1/forms/public/submissions`
  - Accepts form submissions and validates/sanitizes input.
  - Returns `400` with `{ "errors": { "field": "message" } }` on validation failures.
  - Returns `409` with `{ "errors": { "submission": "..." } }` for duplicate submissions.
  - Stores accepted submissions in a persistent SQLite database by default.
- `GET /health`
  - Liveness endpoint returning plain text `ok`.
- CORS support with origin allowlist and preflight (`OPTIONS`) handling for cross-origin frontend calls.
- Consistent JSON error envelope for route misses (`404`) and wrong methods (`405`):
  - `{ "errors": { "...": "..." } }`
- Form metadata currently uses assignment defaults from the backend service.

## Input Contract (POST)

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "phoneNumber": "+47 12345678",
  "birthDate": "1990-04-21",
  "memberTypeId": "8FE4113D4E4020E0DCF887803A886981",
  "clubId": "britsport",
  "formId": "B171388180BC457D9887AD92B6CCFC86"
}
```

## Validation and Security Notes

- Strict JSON decode (`DisallowUnknownFields`).
- Request body size limit (1 MB).
- Content-Type check (`application/json`) for JSON endpoint.
  - If `Content-Type` header is present, it must be `application/json`.
  - Missing `Content-Type` is accepted.
- Input normalization (`trimSpace`, lowercase email).
- Email format validation via standard library.
- Phone format validation via conservative regex.
- Birth date format validation (`YYYY-MM-DD`) and bounds checks.
- `clubId`, `formId`, and `memberTypeId` checked against current form config.
- Submission is rejected until `registrationOpens` is reached.
- Duplicate submission rule is scoped to matching `formId`, `clubId`, `email`, `phoneNumber`, and `birthDate`.
- Security response headers applied to all routes.

## Run Locally

Requires Go `1.26+`.

```bash
cd backend
go run ./cmd/api
```

Server listens on `:8080`.
Default storage is SQLite file `./data/backend.db`.

### Optional environment configuration

- `BACKEND_ALLOWED_ORIGINS`
  - Comma-separated list of allowed origins for CORS.
  - Default: `http://localhost:5173,http://127.0.0.1:5173`
  - Example: `BACKEND_ALLOWED_ORIGINS=https://app.example.com`
- `BACKEND_REPOSITORY`
  - Backend persistence implementation.
  - Supported values: `sqlite` (default), `memory`.
- `BACKEND_SQLITE_PATH`
  - Path to SQLite database file when `BACKEND_REPOSITORY=sqlite`.
  - Default: `./data/backend.db`

## Run Tests

```bash
cd backend
go test ./...
```

## Why this structure

- `internal/http`: transport concerns (routing, handlers, response helpers)
- `internal/service`: orchestration/business rules
- `internal/validate`: request validation and normalization
- `internal/repository`: persistence abstraction
- `internal/repository/memory`: in-memory persistence (useful for tests/fallback)
- `internal/repository/sqlite`: persistent SQLite implementation

This keeps DB integration straightforward by swapping repository implementations without changing handlers.

## Local vs Production Database Choice

- Local development defaults to SQLite because it is file-based, easy to run, and requires no separate database service.
- Production should use PostgreSQL because it handles concurrency and multi-instance deployments better, and has stronger operational support for backups, monitoring, and failover.
- The repository abstraction lets us keep the same API and service logic while switching the storage backend by configuration.

## Suggested Next Steps

1. Add PostgreSQL repository implementation for production deployments.
2. Add request logging and correlation IDs.
3. Add integration tests against PostgreSQL (Docker compose).
4. Externalize form metadata configuration (for example env/DB-backed `registrationOpens`, `memberTypes`, `clubId`, and `formId`) when production flexibility is needed.
5. Keep DB unique index parity across SQLite/PostgreSQL for duplicate rule consistency.

## Migration Reference

- Initial schema is documented in `migrations/0001_submissions.sql`.
