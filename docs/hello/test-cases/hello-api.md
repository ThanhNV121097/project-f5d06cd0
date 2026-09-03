# Test Cases — Hello API

Risk level: low. Read-only API. Cover success shape, defaulting, no-body behavior, and documented boundary on over-limit name.

## Scenario 1: Health returns exact ok body
**Given** backend is running and reachable
**When** client sends `GET /api/health`
**Then** response status is `200` and response body is exactly `{"status":"ok"}` with no extra fields
Check: fetch_url

Trace: HELLO-001 AC-1, service contract 3.1 success shape

## Scenario 2: Hello returns default greeting with no name
**Given** backend is running and reachable
**When** client sends `GET /api/hello` with no `name` query parameter
**Then** response status is `200` and response body is exactly `{"message":"Hello, World!"}`
Check: fetch_url

Trace: HELLO-002 AC-1, service contract 3.2 default response

## Scenario 3: Hello returns supplied name in message
**Given** backend is running and reachable
**When** client sends `GET /api/hello?name=Ada`
**Then** response status is `200` and response body is exactly `{"message":"Hello, Ada!"}`
Check: fetch_url

Trace: HELLO-002 AC-2, service contract 3.2 named response

## Scenario 4: Hello ignores request body
**Given** backend is running and reachable
**When** client sends `GET /api/hello?name=Ada` with any request body
**Then** response status is `200` and response body is still exactly `{"message":"Hello, Ada!"}`
Check: fetch_url

Trace: service contract 3.2 request body ignored

## Scenario 5: Hello rejects over-limit name
**Given** backend is running and reachable
**When** client sends `GET /api/hello?name=` followed by more than 80 Unicode code points after URL decoding
**Then** response status is `422` and response body uses error envelope with `code` `VALIDATION_FAILED`
Check: fetch_url

Trace: service contract 3.2 validation boundary, error contract 2.3
