# Frontend (React + TypeScript + Vite)

Membership signup wizard frontend.

## Implemented in this iteration

- 3-step responsive wizard UI:
  - Step 1: title, description, member type selection
  - Step 2: name, email, phone number, birth date
  - Step 3: overview and submit
- Submission via React form action in the client app.
- If `registrationOpens` is in the future, an informational banner is shown instead of the wizard.
- On successful submission, a success banner is shown.
- Basic client-side validation aligned with backend constraints.
- Draft persistence in `sessionStorage` with 15-minute TTL:
  - Restores on refresh if still valid
  - Clears on successful submission

## API expectations

The app expects a backend at:

- `GET /api/v1/forms/public`
- `POST /api/v1/forms/public/submissions`

By default, Vite dev server proxies `/api` and `/health` to `http://localhost:8080`.

## Run locally

1. Start backend on port `8080`.
2. In another terminal:

```bash
cd frontend
npm install
npm run dev
```

Open `http://localhost:5173`.

## Build commands

Local verification build (uses `development` mode, allows relative API paths):

```bash
cd frontend
npm run build
```

Deployment build (strict production mode, requires `VITE_BACKEND_BASE_URL`):

```bash
cd frontend
npm run build:prod
npm run preview
```

For production builds, set `VITE_BACKEND_BASE_URL` in your environment or `.env.production`.
Use `.env.production.example` as the template.

## Tests

Run Playwright tests (live backend integration):

```bash
cd frontend
npm run test:e2e
```

Run Playwright in headed mode:

```bash
cd frontend
npm run test:e2e:headed
```

Playwright tests run against a real backend server (started automatically by Playwright config), using an isolated SQLite database path for deterministic test data.

## Minimal dependency choices

- React + React DOM
- TypeScript
- Vite
- Native `fetch` for API calls
- Plain CSS (no UI component library)

## Notes

- Frontend tests are Playwright-only.
- Happy-path and duplicate tests run against real backend responses.
- Registration-locked test uses a targeted route fulfill mock for form details to force `registrationOpens` into the future.
