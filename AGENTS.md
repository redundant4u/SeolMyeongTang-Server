# SeolMyeongTang-Server AGENTS

## Scope

- This repository contains two active services: `go/` and `socket/`.
- `legacy/` is backup-only code. Do not read it, import it, or modify it for live behavior.

## Service Context

- `go/` is the REST API for VNC session management and post reads.
- `socket/` is the WebSocket service for the terminal intro flow and CRDT canvas sync.
- Treat them as separate runtime surfaces that share product context but not the same code path.

## Structure Rules

- Keep `go/main.go` limited to wiring and startup.
- Keep shared Go helpers in `go/internal/pkg`.
- Add HTTP domains under `go/internal/api/<domain>/`.
- Keep request parsing and response mapping in handlers; keep persistence or Kubernetes logic outside handlers.
- Keep Kubernetes session lifecycle logic inside `go/internal/api/session/`.
- Keep post persistence access in a repository layer.
- Add socket features under `socket/src/<feature>/` and register them in `app.module.ts`.
- Keep gateway files focused on events and connection lifecycle; extract larger logic into services.
- Do not create new top-level directories outside `go/` and `socket/`.

## Current Behavior Constraints

- Session endpoints are `GET /session`, `POST /session`, `POST /session/client-id`, and `DELETE /session`.
- Session APIs require `X-Client-Id`.
- Session capacity is currently capped at 4 concurrent sessions per client.
- Session pods expire after 10 minutes and are cleaned by the background GC loop.
- Session pod creation currently also provisions Cloudflare tunnel config.
- Post APIs are read-only and blocked when `APP_ENV=production`.
- The terminal gateway spawns `ssh terminal` per socket and tears it down on disconnect.
- The CRDT gateway keeps a 64x64 in-memory canvas state and enforces a simple per-instance connection limit.

## Environment and Runtime Assumptions

- Required Go env values include `KUBE_SESSION_NAMESPACE`, `CF_TUNNEL_ID`, `AWS_ACCESS_KEY`, `AWS_SECRET_KEY`, and `DYNAMODB_TABLE`.
- `KUBE_CONFIG` switches between external kubeconfig and in-cluster auth.
- The Go server listens on `8090`.
- The socket server listens on `APP_PORT`, defaulting to `3000`.
- Do not change `.env` or secret values unless the task explicitly requires it.

## Working Rules

- Treat session limits, TTL, required headers, tunnel wiring, and client-visible socket behavior as high-impact changes.
- Do not copy patterns from `legacy/` into active services.
- If API behavior, socket behavior, or data contracts change, update the relevant specs first.
- The repo has very little automated coverage, so assess test additions whenever core behavior changes.
