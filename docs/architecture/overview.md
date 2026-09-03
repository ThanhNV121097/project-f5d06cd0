# Architecture Overview — Hello World

## Stack
- Shape: fullstack — Next.js frontend, Go HTTP API, PostgreSQL database.
- Frontend: Next.js 15 App Router, TypeScript, Tailwind v3, ESLint.
- Backend: Go 1.22+ module under `code/backend`, one binary at `cmd/api`.
- Database: PostgreSQL; backend applies SQL migrations on boot before reporting healthy.
- Runtime: `docker compose up` from repository root starts Postgres, backend, frontend.

## Repository layout
- `code/backend/cmd/api/main.go` owns process startup, migrations, health probe, and route registration.
- `code/backend/migrations/*.sql` stores ordered schema migrations; `schema_migrations` tracks applied files.
- `code/backend/.env.example` lists backend env vars.
- `code/frontend/app/page.tsx` is composition root only; story components mount there later.
- `code/frontend/app/globals.css` owns design tokens and shared base styles; story CSS modules must use tokens.
- `code/frontend/.env.example` lists browser/server API origin env vars.
- `docs/architecture/overview.md` records project-wide contracts. ERD and service contracts extend in later tasks.

## Data flow
1. Browser loads frontend from Next.js.
2. Frontend calls backend using `NEXT_PUBLIC_API_URL` in browser code and `API_ORIGIN` for server-side code.
3. Backend reads `DATABASE_URL`, applies migrations, then serves HTTP on `PORT`, `APP_PORT`, or `8080`.
4. Health is split: `/healthz` proves migrations completed and DB answers `SELECT 1`; later product API exposes `/api/health`.

## Environment variables
Backend:
- `DATABASE_URL` — PostgreSQL connection string injected by runtime.
- `PORT` — preferred HTTP listen port.
- `APP_PORT` — fallback HTTP listen port when `PORT` is absent.

Frontend:
- `NEXT_PUBLIC_API_URL` — browser-visible backend origin, e.g. `http://localhost:8080`.
- `API_ORIGIN` — server-side backend origin, e.g. `http://backend:8080`.

Root compose:
- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` — local database settings.
- `BACKEND_PORT`, `FRONTEND_PORT` — optional published ports.

## Naming and conventions
- Go packages use short lowercase names; exported names only when cross-package use needs them.
- Backend JSON fields use snake_case only when API contract says so; otherwise match SRS fields exactly.
- React component files use PascalCase and `export default function ComponentName()`.
- `app/page.tsx` stays server component; interactive components start with literal first line `"use client"`.
- CSS modules may not hardcode colours or spacing; use tokens from `app/globals.css`.
- No secrets in repository; every read env var appears in nearest `.env.example` with comment.

## Key decisions
| Decision | Why | Rejected alternative | Tradeoff |
|---|---|---|---|
| Self-migrate backend on boot | Runtime creates empty DB and no separate migrator exists | Manual migration step | Startup does more work, but deployments are reproducible |
| `/healthz` checks DB | Compose must not mark broken DB app healthy | Process-only health check | Health can fail during DB outage, which is correct |
| Keep product endpoints out of scaffold | Feature stories own `/api/*` behavior | Implement full demo now | Later work has small PRs, but scaffold proves boot only |
| Use committed Docker/CI files unchanged | Repo already supplies build contracts | Rewrite Docker/compose/workflows | Less control, less drift |
| Tailwind plus CSS custom tokens | Matches design system and CI token checks | Hardcoded utility values everywhere | More upfront tokens, fewer visual regressions |

## How to run
1. Copy `.env.example` to `.env` if local overrides are needed.
2. Run `docker compose --profile local up --build` from repo root.
3. Frontend: `http://localhost:3000`; backend health: `http://localhost:8080/healthz`.

## Verification contract
- Backend CI runs `go mod download`, `go build ./...`, `go vet ./...`, `go test ./...`.
- Frontend CI runs `npm ci`, `npm run lint`, `npm run build`, `npm test --if-present`.
- Token CI rejects undefined token use, token fallbacks, and hardcoded visual values in CSS modules.
