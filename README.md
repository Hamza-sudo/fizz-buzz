# FizzBuzz API

Small production-minded Go HTTP service exposing a configurable FizzBuzz API and a statistics endpoint.

## Features

- `GET /api/v1/fizzbuzz`
- `GET /api/v1/statistics`
- `GET /health`
- strict request validation
- persistent SQLite request statistics
- graceful shutdown and HTTP server timeouts
- test coverage for business logic and HTTP endpoints

## Run locally

Prerequisite: Go `1.26.x` (or newer in the `1.26` line).

```bash
go run ./cmd/server
```

Or with `make`:

```bash
make run
```

Server listens on `:8080` by default.
You can override it with `PORT`, for example `PORT=9090 go run ./cmd/server`.
You can also configure `MAX_LIMIT` (default: `100000`), for example `MAX_LIMIT=50000 go run ./cmd/server`.
By default statistics are stored in `file:fizzbuzz_stats.db`; override with `STATS_DB_DSN` if needed.
You can tune DB operation timeout with `STATS_DB_TIMEOUT_MS` (default: `200`), for example `STATS_DB_TIMEOUT_MS=500`.

## Run with Docker

```bash
docker build -t fizz-buzz .
docker run --rm -p 8080:8080 fizz-buzz
```

## API

### Generate a FizzBuzz sequence

```bash
curl "http://localhost:8080/api/v1/fizzbuzz?int1=3&int2=5&limit=15&str1=fizz&str2=buzz"
```

Response:

```json
["1","2","fizz","4","buzz","fizz","7","8","fizz","buzz","11","fizz","13","14","fizzbuzz"]
```

### Get most frequent request

```bash
curl "http://localhost:8080/api/v1/statistics"
```

Response:

```json
{
  "params": {
    "int1": 3,
    "int2": 5,
    "limit": 15,
    "str1": "fizz",
    "str2": "buzz"
  },
  "hits": 2
}
```

When no request has been recorded yet:

```json
{
  "params": null,
  "hits": 0
}
```

### Health check

```bash
curl "http://localhost:8080/health"
```

Response:

```json
{
  "status": "ok"
}
```

### Validation error format

Example (`400 Bad Request`):

```json
{
  "error": {
    "code": "INVALID_PARAMETER",
    "message": "int1 must be greater than 0"
  }
}
```

## Validation rules

- `int1 > 0`
- `int2 > 0`
- `limit > 0`
- `limit <= MAX_LIMIT` (default `100000`)
- `str1` and `str2` must not be empty

## Tests

```bash
go test ./...
```

Or:

```bash
make test
```

## Design choices

- Standard library only for the HTTP layer: fewer dependencies, lower maintenance cost, easier long-term ownership.
- Clear separation of concerns:
  - `internal/service`: business logic
  - `internal/httpapi`: transport and request validation
  - `internal/stats`: request counting
  - `cmd/server`: application bootstrap
- SQLite-backed statistics store for durability across restarts.
- HTTP server configured with timeouts and graceful shutdown to be closer to production expectations.
- Structured JSON logs through `slog` for easier observability.

## Possible next improvements

- add request logging / metrics to improve troubleshooting and SLO tracking
- expose OpenAPI documentation to make integrations and contract validation easier
- add configuration struct and environment validation for safer deployments
- add automated DB backups / retention policies depending on deployment constraints
