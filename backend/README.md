# AEOlyzer backend

This module contains the first production vertical slice for guest onboarding:

```text
Layer 2 validates and authorizes
  -> Layer 3 creates the site-inspection and prompt plans
  -> Layer 6 verifies authorization and applies network policy
  -> Layer 7 reads public site metadata
  -> Layer 5 builds the dashboard presentation contract
  -> Layer 8 records sanitized outcome events
```

The API keeps onboarding state client-side for the guest session. It does not
persist raw prompts, page bodies, credentials, or project memory.

## Run

```powershell
go run ./cmd/api
```

Configuration:

```text
AEOLYZER_ADDRESS          default: 127.0.0.1:8080
AEOLYZER_FRONTEND_ORIGIN  default: http://localhost:3000
```

## Endpoints

```text
GET  /healthz
POST /v1/onboarding/inspect
POST /v1/onboarding/complete
```

The inspection endpoint accepts a guest session ID and public website URL. The
completion endpoint accepts the validated onboarding profile and returns a
versioned dashboard frame with 12 project-specific AEO prompts.

## Verify

```powershell
gofmt -w cmd internal layer_02_intake layer_03_orchestration `
  layer_05_extensions layer_06_runtime layer_07_interop `
  layer_08_observability
go test ./...
go vet ./...
go test -race ./...
```
