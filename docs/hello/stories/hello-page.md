# Story — Hello page

## User story
As a guest, I want one clean page that calls the API, so that I can see the live hello message, submit a greeting, and view stored greetings.

## In scope
- One browser page for the `hello` module.
- Page loads live hello text from `GET /api/hello`, not hard-coded copy.
- Page shows form fields for name and message.
- Page submits new greetings to `POST /api/greetings`.
- Page refreshes stored greetings from `GET /api/greetings` after submit and shows newest-first data from backend.
- Page shows a friendly error state when the API cannot be reached.

## Out of scope
- Health endpoint behavior itself.
- Greeting validation rules beyond displaying API responses.
- Accounts, login, roles, payments, or any other module.
- Backend implementation, database schema, or migration work.
- Design changes outside the approved page structure.

## UI scope
- One default page state in the approved `Hello World demo page` design.
- Uses the live hello panel, greeting form, stored greetings list, and unreachable-error copy shown in the design system.
- Covers empty list, successful submit refresh, and API-unreachable state.

## Acceptance criteria
1. When a guest opens the page, the page requests live hello text from the API and renders the returned message.
2. When the API is unavailable, the page shows a friendly unreachable error to the guest.
3. When a guest enters a non-empty name and message and submits, the page sends the greeting to the API.
4. When the API returns a stored greeting, the page updates the stored greetings list from backend data instead of local mock data.
5. When greetings exist, the page shows them newest first after refresh.
6. The page works without any hard-coded greeting list data.

## Dependencies
- `Hello API` story must provide `GET /api/hello`.
- `Persist greetings` story must provide `POST /api/greetings` and `GET /api/greetings`.
- Approved design and design system for page structure and states.
