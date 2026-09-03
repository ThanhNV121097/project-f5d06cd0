# Hello API

As a guest, I want backend health and hello endpoints, so that browser and API clients can confirm service is up and get a live greeting.

## In scope
- GET `/api/health` returns `{"status":"ok"}` with 200.
- GET `/api/hello` returns `{"message":"Hello, World!"}` when no `name` query is sent.
- GET `/api/hello?name=<name>` returns `{"message":"Hello, <name>!"}` when `name` is present.
- Responses use exact JSON field names and no extra data.

## Out of scope
- Greeting persistence endpoints.
- Frontend page, form, list rendering, or API unreachable error state.
- Accounts, login, payments, or permissions.
- Database schema changes beyond what this story needs for read-only API behavior.

## UI scope
- No UI in this story.
- Approved design sections touched later by other stories only.

## Acceptance criteria
- Health route responds 200 with exact JSON `{"status":"ok"}`.
- Hello route responds 200 with exact JSON `{"message":"Hello, World!"}` when `name` is absent.
- Hello route responds 200 with exact JSON `{"message":"Hello, Ada!"}` when `name=Ada`.
- Response body contains only required field for each route.
- Empty or missing `name` uses default greeting, not error.

## Dependencies
- Backend HTTP API scaffold exists.
- Project routing and JSON response helpers already in place.
- No external accounts or secrets required.
