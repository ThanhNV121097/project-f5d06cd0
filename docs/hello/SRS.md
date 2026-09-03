# SRS — Hello World

Module: `hello`
Design: [View the approved design](http://localhost:8080/design/f5d06cd0-9964-4a0d-bde9-665956c142f1)
Design system: `design/design-system.md`

## 1. Purpose

`hello` provides the demo's public API and browser page. It proves the app can answer health and greeting requests, store greetings in the database, and show the stored data back to a visitor. Without it, the product is only a shell and cannot prove the full stack works end to end.

## 2. Actors

| Actor | Who they are | What they may do in this module |
|---|---|---|
| Guest | Any browser visitor or HTTP client | Call the public API, view the page, submit a greeting |
| API client | Any programmatic caller of the HTTP API | Call the public API endpoints and read returned JSON |

## 3. Scope

**In scope** — the functions specified below, by their plan titles:

- Hello API
- Persist greetings
- Hello page

**Out of scope** — name what a reader would reasonably expect here and say where it lives instead.

- Accounts, login, roles, or permissions — deliberately not built for this demo.
- Payments, subscriptions, or billing — deliberately not built for this demo.
- Any module other than `hello` — belongs to another module if it appears later.
- Design work itself — lives in `design/index.html` and `design/design-system.md`, not in this SRS.

## 4. Functional requirements

### 4.1 Hello API

**Requirement HELLO-001 — Health response**

As an API client, I want to check service health, so that I can confirm backend is reachable.

Behaviour:

1. Given a request to health check route, when caller sends request, then service returns success JSON with status `ok`.
2. Given health check route is available, when caller reads response, then response body contains only the health status needed for the check.

**Acceptance criteria**

| # | Given | When | Then |
|---|---|---|---|
| AC-1 | API server is running | client calls health check route | response is 200 and JSON is `{\"status\":\"ok\"}` |
| AC-2 | No special request data | client calls health check route | response body contains health status only |

**Failure, boundary and permission behaviour**

| Case | Condition | Expected behaviour |
|---|---|---|
| Not applicable | No roles, no writes, no alternate screen state in approved design | Not applicable: this function is a single read with no permission case and no design-visible error state |

**Data touched**

| Field | Type | Required | Rule |
|---|---|---|---|
| status | text | yes | Value is `ok` |

**Requirement HELLO-002 — Greeting text**

As an API client, I want to request a hello message with an optional name, so that I can see a friendly greeting.

Behaviour:

1. Given no name is supplied, when caller requests hello, then service returns `Hello, World!`.
2. Given a name is supplied, when caller requests hello, then service returns `Hello, <name>!`.
3. Given the request uses a name value, when caller reads response, then service uses the provided name value in the message.

**Acceptance criteria**

| # | Given | When | Then |
|---|---|---|---|
| AC-1 | No name query is sent | client calls hello route | response is 200 and JSON is `{\"message\":\"Hello, World!\"}` |
| AC-2 | Name query is `Ada` | client calls hello route | response is 200 and JSON is `{\"message\":\"Hello, Ada!\"}` |

**Failure, boundary and permission behaviour**

| Case | Condition | Expected behaviour |
|---|---|---|
| Not applicable | No roles, no writes, no design-visible error state | Not applicable: this function is a single read with no permission case and no approved error screen |

**Data touched**

| Field | Type | Required | Rule |
|---|---|---|---|
| name | text | no | When present, it becomes the name in the greeting message |

### 4.2 Persist greetings

**Requirement HELLO-003 — Store greeting**

As a guest, I want to submit a greeting, so that the greeting is stored and can be shown later.

Behaviour:

1. Given name and message are both present and valid, when guest submits greeting, then one greeting row is stored.
2. Given the row is stored, when the system returns success, then response includes that row's id and created_at.
3. Given the row was stored, when guest later asks for greetings, then stored greeting can be read back from persistence.

**Acceptance criteria**

| # | Given | When | Then |
|---|---|---|---|
| AC-1 | name and message are valid and non-empty | client submits greeting | response is success and includes id, name, message, created_at |
| AC-2 | one greeting was just saved | client requests greetings | response includes saved greeting in newest-first order |

**Failure, boundary and permission behaviour**

| Case | Condition | Expected behaviour |
|---|---|---|
| Invalid input | name is empty | request is rejected and nothing is stored |
| Invalid input | message is empty | request is rejected and nothing is stored |
| Boundary | name is longer than the agreed short limit | request is rejected with a validation error naming the field |
| Boundary | message is longer than the agreed short limit | request is rejected with a validation error naming the field |
| Not applicable | No roles or ownership rules exist | Not applicable: any guest may create a greeting |

**Data touched**

| Field | Type | Required | Rule |
|---|---|---|---|
| id | integer | yes | Server assigns unique row id |
| name | text | yes | Must be non-empty and reasonably short |
| message | text | yes | Must be non-empty and reasonably short |
| created_at | timestamp | yes | Set when row is stored |

**Requirement HELLO-004 — Read greetings newest first**

As a guest, I want to read stored greetings, so that I can see what was submitted.

Behaviour:

1. Given greetings exist in storage, when guest requests list, then service returns stored rows newest first.
2. Given no greetings exist, when guest requests list, then service returns an empty list.
3. Given rows exist, when guest reads them, then each row includes id, name, message, and created_at.

**Acceptance criteria**

| # | Given | When | Then |
|---|---|---|---|
| AC-1 | two greetings exist with different created_at values | client requests greetings | response returns both rows newest first |
| AC-2 | no greetings exist | client requests greetings | response returns empty JSON list |

**Failure, boundary and permission behaviour**

| Case | Condition | Expected behaviour |
|---|---|---|
| Not applicable | No roles, no alternate design state | Not applicable: list screen has only the default state in approved design |
| Upstream failure | database is unavailable | request fails with a clear API error response and no partial list |

**Data touched**

| Field | Type | Required | Rule |
|---|---|---|---|
| id | integer | yes | Returned for each greeting |
| name | text | yes | Returned for each greeting |
| message | text | yes | Returned for each greeting |
| created_at | timestamp | yes | Returned for each greeting |

### 4.3 Hello page

**Requirement HELLO-005 — Render live hello page**

As a guest, I want one clean page that calls the API, so that I can see the live hello message and use the form.

Behaviour:

1. Given browser opens the page, when page loads, then page calls the API for hello text and renders the result.
2. Given browser opens the page, when page loads, then page shows form for name and message and a list of stored greetings.
3. Given API cannot be reached, when page tries to load data, then page shows a friendly unreachable error to the user.

**Acceptance criteria**

| # | Given | When | Then |
|---|---|---|---|
| AC-1 | page loads in browser | visitor opens page | page shows live hello text from API, not hard-coded greeting text |
| AC-2 | API is unreachable | page tries to load data | page shows friendly unreachable error |
| AC-3 | design is rendered | visitor views page | page shows hello panel, greeting form, and stored greetings list |

**Failure, boundary and permission behaviour**

| Case | Condition | Expected behaviour |
|---|---|---|
| Upstream failure | API is unreachable | page shows friendly error state visible to user |
| Not applicable | No roles, no permissions, no alternate states beyond approved design | Not applicable: approved design shows one default page state with visible unreachable-error copy in the page content |

**Data touched**

| Field | Type | Required | Rule |
|---|---|---|---|
| hello message | text | yes | Shown from API response |
| greeting list | list | yes | Refreshed after submit |
| unreachable error copy | text | yes | Visible when API cannot be reached |

## 5. Screens

The design is the source of truth for appearance; this section maps functions onto it so nothing in the design is unaccounted for and nothing specified here is missing from the design.

| Screen | Section in the design | Functions it serves | States that must exist |
|---|---|---|---|
| Hello World demo page | Overview, Live hello, Add greeting, Stored greetings | HELLO-001, HELLO-002, HELLO-003, HELLO-004, HELLO-005 | default |

## 6. Non-functional requirements

| Area | Requirement |
|---|---|
| Performance | Page hello refresh and greetings refresh respond within 2s at 1 Mbps with cold cache |
| Accessibility | Keyboard reachable, visible focus, labelled inputs, contrast at least 4.5:1 |
| Responsive | Works at 320px and up; no horizontal page scroll |
| Privacy | Name, message, and created_at are stored only to show submitted greetings back to visitors |

## 7. Dependencies and assumptions

- **Depends on:** backend HTTP API, for hello, greeting storage, and greeting list data.
- **Depends on:** database, for persistent greeting rows.
- **Assumption:** API failure state uses the friendly unreachable copy already shown in the approved design.
- **Assumption:** short means name and message must stay within the visible max-length limits in the page form.

| Open question | Proposed default | Who decides |
|---|---|---|
| None for this module | N/A | N/A |

## 8. Traceability

| Plan item | Requirement ids | Test cases |
|---|---|---|
| Hello API | HELLO-001, HELLO-002 | `test-cases/hello-api.md` |
| Persist greetings | HELLO-003, HELLO-004 | `test-cases/persist-greetings.md` |
| Hello page | HELLO-005 | `test-cases/hello-page.md` |
