# SeolMyeongTang-Server AGENTS

## Required Reading

- When working inside the SMT spec workspace, read `../constitution.md` and
  `../AGENTS.md` before applying these local rules. Resolve paths relative to this file.

## Scope

- `go/` serves HTTP APIs; `socket/` handles WebSocket events. Keep their runtime
  and code paths separate.
- `legacy/` is backup-only code. Do not read it, import it, or modify it for live behavior.

## Domain Structure

Keep Go HTTP domains under `go/internal/api/<domain>/`:

```text
<domain>/
├── routes.go         # Route registration and dependency wiring
├── handler.go        # HTTP parsing, boundary validation, and response mapping
├── dto.go            # Request/response types and mapping
├── model.go          # Feature data types, when needed
├── repository.go     # Persistence access, when needed
└── *_test.go         # Tests for the owning package
```

- Create only files the feature needs; this layout is not a mandatory scaffold.
- Keep persistence and external integration logic outside handlers, within the
  owning feature. Kubernetes session lifecycle code belongs in the session domain.
- Keep shared Go helpers in `go/internal/pkg` and `go/main.go` limited to startup wiring.
- Keep socket features under `socket/src/<feature>/`, with `<feature>.module.ts`
  for provider registration and `<feature>.gateway.ts` for events and connection lifecycle.
  Add `<feature>.service.ts` only when processing needs to be separated from the gateway.
- Register feature modules in `app.module.ts`; do not add new top-level source
  directories outside `go/` and `socket/`.

## Naming And Dependencies

- Match Go package names to their domain directories. Use the file responsibilities
  above rather than introducing parallel naming or layer hierarchies.
- Handlers may call feature repositories and integration helpers; those helpers
  must not depend on handlers. Shared `internal/pkg` helpers must not import feature domains.
- Use the existing Go validation and response helpers in `go/internal/pkg/validator`
  and `go/internal/pkg/response` instead of creating endpoint-specific copies.
- Shape response DTOs from the shared contracts and retain wire field names and
  casing. Do not return persistence/provider objects directly or expose undeclared fields.
- Pass the caller's context through Go external I/O operations; do not replace
  request cancellation or deadlines with `context.Background()` in request paths.
- Socket gateways may call feature services and helpers; services must not depend
  on gateways or own socket event registration.

## Boundary Validation And Resource Lifecycle

- Validate HTTP and WebSocket payload structure, bounds, and resource access scope
  before changing state or performing external writes or deletions.
- Client-supplied identifiers alone are not proof of authorization to a resource.
- Repeated initialization, disconnect, and cleanup must not leak subprocesses or
  register duplicate listeners; release only resources owned by that operation or connection.
- For multi-resource creation, define cleanup of resources created by the operation
  when a later step fails, and make that cleanup safe to repeat.

## Observability

- Use `go/internal/pkg/logger` for Go events and pass the request context when
  available to preserve correlation fields and redaction.
- Reuse `go/internal/pkg/metrics` helpers. Keep labels bounded, such as route patterns,
  outcomes, and error codes; do not label by client/session/request IDs or raw URLs.

## Working Rules

- Do not change `.env` or secret values unless the task explicitly requires it.
- When resource lifecycle changes, review both request-driven operations and
  background cleanup paths.
- Core logic includes validation and access checks, contract mapping, capacity
  limits, resource lifecycle transitions, and shared-state synchronization.

## Go Verification

- Keep Go tests beside the package as `*_test.go`. From `go/`, start with
  `go test ./internal/<affected-package>` and broaden coverage when the change
  affects shared packages or multiple domains.
