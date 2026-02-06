# OpenAction Backend

## Quickstart

```bash
set OA_PORT=8080
go run ./cmd/openaction
```

## Config

- `OA_PORT` (required)
- `OA_DATA_DIR` (default `../backend/data`)
- `OA_DB_PATH` (default `../backend/data/openaction.db`)
- `OA_SERVE_UI` (default `true`)
- `OA_UI_DIST` (default `../backend/web/dist`)
- `OA_TLS_CERT` / `OA_TLS_KEY` / `OA_CA_CERT` (mTLS for gRPC)
- `OA_ADMIN_EMAIL` / `OA_ADMIN_PASSWORD`
- `OA_POOL_GRPC_ADDR` (default `:7443`)

## Auth

- Browser: Cookie session (`oa_session`)
- CLI: `Authorization: Bearer <token>`

Tokens are stored as SHA256 hashes only.
