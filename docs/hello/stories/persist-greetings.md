# Persist greetings

As a guest, I want to submit greetings and read them back, so that greeting data persists in the demo and proves backend plus database work end to end.

## In scope
- `POST /api/greetings` stores one greeting row in database table `greetings`.
- `GET /api/greetings` returns stored greetings newest first.
- Validate `name` and `message` are non-empty and reasonably short.
- Return created greeting with `id`, `name`, `message`, and `created_at`.
- Support the Hello page by providing persistent data it can refresh after submit.

## Out of scope
- `GET /api/health` and `GET /api/hello`.
- Frontend form layout, live hello display, and error handling UI.
- Accounts, login, permissions, payments, or any other module.
- Deleting or editing greetings.

## UI scope
- No dedicated screen in approved design.
- This story only supplies backend data for the Hello World demo page’s stored greetings area and submit flow.
- Runtime states covered here are API success and validation failure; no design-visible error page belongs to this story.

## Acceptance criteria
1. When guest sends valid non-empty `name` and `message` to `POST /api/greetings`, service stores exactly one row and returns 201 with the stored row including `id` and `created_at`.
2. When guest sends empty `name` or empty `message`, service rejects request and stores nothing.
3. When guest sends overlong `name` or overlong `message`, service rejects request with a validation error naming the field and stores nothing.
4. When greetings exist, `GET /api/greetings` returns them newest first.
5. When no greetings exist, `GET /api/greetings` returns an empty JSON list.
6. When a greeting is stored, later reads return the persisted row with `id`, `name`, `message`, and `created_at`.

## Dependencies
- Backend HTTP API and PostgreSQL database.
- Backend schema for table `greetings(id, name, message, created_at)`.
- `Hello API` story for service boot and shared API conventions.
- `Hello page` story for the browser flow that consumes this data.
