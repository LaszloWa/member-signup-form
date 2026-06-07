# Spond Fullstack Coding Test

This repository contains a small full-stack membership registration flow:

- A Go backend API with validation, duplicate prevention, and SQLite persistence.
- A React + TypeScript frontend with a 3-step wizard.
- Playwright end-to-end tests.

The implementation intentionally favors simple, readable code and minimal dependencies.

## Why these tech choices

- Backend: Go + net/http
  - Small standard-library stack keeps behavior explicit and easy to review.
  - Also the prescribed backend language
- Persistence: SQLite
  - Simple local persistence with no external database setup needed.
- Frontend: React + TypeScript + Vite
  - Lightweight and fast local iteration.
- Styling: plain CSS
  - Avoids framework overhead for a coding-test codebase.
- Testing: Playwright
  - Verifies real browser behavior and frontend-backend integration.

## Prerequisites

- Go 1.26+
- Node.js 20+
- npm

## Run backend

```bash
cd backend
go run ./cmd/api
```

Backend listens on `http://localhost:8080`.

## Run frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend runs on `http://localhost:5173`.

## Test commands

Frontend:

```bash
cd frontend
npm run lint
npx tsc -b
npm run build
npm run test:e2e
```

Backend:

```bash
cd backend
go test ./...
```

## SQLite persistence notes

- Backend defaults to SQLite file storage at `backend/data/backend.db`.
- Playwright uses an isolated test database under `backend/data/e2e/`.
- Frontend test script cleans e2e SQLite files before each run to keep tests deterministic.

## Suggested next steps

1. Add structured request logging and request IDs.
2. Add PostgreSQL repository implementation for production-oriented deployments.
3. Expand API contract tests for failure modes and edge cases.
4. Add CI workflow to run backend tests and frontend lint/build/e2e on each PR.

## AI usage disclosure

AI assistance (GitHub Copilot Chat) was used to help with code scaffolding, refactoring suggestions, and documentation drafting.
All generated changes were reviewed, edited, and validated with local tests before acceptance.
