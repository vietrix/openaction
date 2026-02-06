# OpenAction

OpenAction is a next‑gen self‑hosted CI/CD platform designed to replace Jenkins with a single‑binary, low‑footprint, Go‑first stack.

## Highlights
- Single binary distribution.
- SQLite (WAL) metadata + compressed blob storage.
- WebSocket live logs.
- Release management with environment matrix.
- RBAC, audit trail, secrets management.

## Repo Layout
- `frontend/` Svelte UI
- `backend/` Go control plane (HTTP + gRPC)
- `pool/` Runner daemon (gRPC)

## Quickstart (Dev)
```bash
set OA_PORT=8080
set OA_SECRET_KEY=change-me
cd backend
go run ./cmd/openaction
```

```bash
cd frontend
npm install
npm run dev
```

## Quickstart (Binary)
```bash
set OA_PORT=8080
set OA_SECRET_KEY=change-me
cd backend
go build -o openaction ./cmd/openaction
.\openaction.exe
```

## vbuild (Fast Startup)
This repo supports [vbuild](https://github.com/vietrix/vbuild) for fast dev startup.

```bash
vbuild dev
```

```bash
vbuild build
```

```bash
vbuild test
```

## API
- Base path: `/actions`
- Public endpoints: `/public/...`

## Auth
- Browser: Cookie session (`oa_session`)
- CLI: `Authorization: Bearer <token>`

## Config
- `OA_PORT` (required)
- `OA_SECRET_KEY` (required)
- `OA_DATA_DIR` (default `../backend/data`)
- `OA_DB_PATH` (default `../backend/data/openaction.db`)
- `OA_SERVE_UI` (default `true`)
- `OA_UI_DIST` (default `../backend/web/dist`)
- `OA_TLS_CERT` / `OA_TLS_KEY` / `OA_CA_CERT` (mTLS for gRPC)
- `OA_ADMIN_EMAIL` / `OA_ADMIN_PASSWORD`
- `OA_POOL_GRPC_ADDR` (default `:7443`)

## License
Apache-2.0
