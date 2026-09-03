# Test Cases — Hello page

Risk level: P1. Page is only UI for demo stack, and it must prove API + DB path end to end.

## Case 1
**Scenario**: Page renders live hello text from API
**Given**: App is running and backend `GET /api/hello` returns `{"message":"Hello, World!"}`
**When**: Browser opens Hello page with no name entered
**Then**: Page shows hello text from API, not hard-coded copy, and visible text is `Hello, World!`
**Check**: render_url

## Case 2
**Scenario**: Page submits new greeting and refreshes list
**Given**: App is running and greeting list is initially empty
**When**: User types valid name and message, submits form, and page reloads or refreshes list after success
**Then**: New greeting appears in stored greetings list with same name and message, newest first
**Check**: interact_page

## Case 3
**Scenario**: Stored greetings come from backend and database
**Given**: Database already contains one greeting row created outside browser flow
**When**: Browser opens Hello page
**Then**: Page shows stored greeting in list after API fetch, proving data came from backend and persisted storage
**Check**: render_url

## Case 4
**Scenario**: API unreachable shows friendly error
**Given**: Frontend is running and API service is stopped or unreachable
**When**: Browser opens Hello page or submits form
**Then**: Page shows friendly unreachable-error message to user, not blank screen or raw stack trace
**Check**: render_url

## Case 5
**Scenario**: Form does not use hard-coded greeting data
**Given**: Backend hello and greetings endpoints return changed values
**When**: Browser opens Hello page
**Then**: Page displays current API response values from network, not fixed sample text in UI
**Check**: render_url

## Case 6
**Scenario**: Form input survives API failure
**Given**: User has typed name and message into form and API becomes unreachable before submit
**When**: User submits form
**Then**: Friendly error appears and typed name/message remain in form fields
**Check**: interact_page

## Case 7
**Scenario**: Hello panel handles optional name from API
**Given**: Browser can call `GET /api/hello?name=Ada`
**When**: Page loads or refreshes hello panel with name present
**Then**: Visible hello text is `Hello, Ada!`
**Check**: render_url

## Case 8
**Scenario**: Page loads cleanly before any user action
**Given**: App starts with no prior browser state
**When**: Browser opens Hello page
**Then**: Page shows hello panel, greeting form, and greetings list without requiring login or hard-coded greeting content
**Check**: render_url

## Case 9
**Scenario**: Visual design matches approved calm demo page
**Given**: Page is open in browser
**When**: No interaction is performed
**Then**: Layout is clean, centered, near-white background with one calm accent, and no distracting animation is visible
**Check**: measure_styles
