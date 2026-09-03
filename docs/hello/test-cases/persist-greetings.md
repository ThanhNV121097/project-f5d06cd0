# Test cases — Persist greetings

Risk: high. Feature writes data, must prove API contract and DB persistence end to end.

## Scenario: POST valid greeting stores row and returns id plus created_at
**Given** backend running, `greetings` table empty, and client sends JSON body `{ "name": "Ada", "message": "Hello from browser" }` to `POST /api/greetings`
**When** request is submitted
**Then** response is `201`, `Location` header is `/api/greetings/{id}`, body includes non-empty string `id`, trimmed `name` `Ada`, trimmed `message` `Hello from browser`, and RFC 3339 UTC `created_at`; row exists in DB with same values
**Check:** fetch_url

## Scenario: GET greetings returns stored rows newest first
**Given** DB has at least 2 stored greetings with known creation order, older row inserted first and newer row inserted second
**When** client sends `GET /api/greetings`
**Then** response is `200`, body has `greetings` array ordered newest first, so newer row appears before older row, and response also includes `next_cursor` and `has_more`
**Check:** fetch_url

## Scenario: GET greetings on empty table returns empty collection shape
**Given** backend running and `greetings` table empty
**When** client sends `GET /api/greetings`
**Then** response is `200`, body is `{ "greetings": [], "next_cursor": null, "has_more": false }`
**Check:** fetch_url

## Scenario: POST rejects empty name or message
**Given** backend running
**When** client sends `POST /api/greetings` with body where `name` is blank after trim or `message` is blank after trim
**Then** response is `422`, error body has code `VALIDATION_FAILED`, field details name the invalid field, and no row is inserted
**Check:** fetch_url

## Scenario: POST rejects overlong name or message
**Given** backend running
**When** client sends `POST /api/greetings` with `name` longer than 80 Unicode code points after trim or `message` longer than 240 Unicode code points after trim
**Then** response is `422`, error body has code `VALIDATION_FAILED`, field details name the over-limit field, and no row is inserted
**Check:** fetch_url

## Scenario: POST ignores unknown JSON fields and stores trimmed values
**Given** backend running
**When** client sends `POST /api/greetings` with body `{ "name": "  Ada  ", "message": "  Hello  ", "extra": "ignored" }`
**Then** response is `201`, stored row uses trimmed values `Ada` and `Hello`, and extra field is not echoed back
**Check:** fetch_url

## Scenario: POST bad JSON or wrong field type returns bad request
**Given** backend running
**When** client sends `POST /api/greetings` with malformed JSON or a wrong field type such as `{"name": 1, "message": "Hi"}`
**Then** response is `400`, error body has code `BAD_REQUEST`, and no row is inserted
**Check:** fetch_url
