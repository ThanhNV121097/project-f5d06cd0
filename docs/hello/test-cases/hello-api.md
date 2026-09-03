# Test Cases — Hello API

Risk level: medium. API has public reads only, but contract exactness matters for frontend and later backend verification.

## HELLO-001 — Health response

**Scenario**: health route returns ok JSON
**Given** API server is running
**When** client calls health check route
**Then** response is 200 and JSON is `{"status":"ok"}`
Check: fetch_url

**Scenario**: health response body contains only status
**Given** API server is running
**When** client calls health check route with no special request data
**Then** response body contains only health status needed for check and no extra JSON fields
Check: fetch_url

## HELLO-002 — Greeting text

**Scenario**: hello defaults to World when name missing
**Given** no name query is sent
**When** client calls hello route
**Then** response is 200 and JSON is `{"message":"Hello, World!"}`
Check: fetch_url

**Scenario**: hello uses supplied name
**Given** name query is `Ada`
**When** client calls hello route
**Then** response is 200 and JSON is `{"message":"Hello, Ada!"}`
Check: fetch_url

**Scenario**: hello trims supplied name before message
**Given** name query is `  Ada  `
**When** client calls hello route
**Then** response is 200 and JSON is `{"message":"Hello, Ada!"}`
Check: fetch_url

**Scenario**: hello rejects name longer than 80 code points
**Given** name query is longer than 80 Unicode code points after URL decoding
**When** client calls hello route
**Then** response is 422 and error code is `VALIDATION_FAILED`
Check: fetch_url
