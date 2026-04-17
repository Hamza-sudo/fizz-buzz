# FizzBuzz API

Small production-minded Go HTTP service exposing a configurable FizzBuzz API and a statistics endpoint.

## Features

- `GET /api/v1/fizzbuzz`
- `GET /api/v1/statistics`
- `GET /health`
- strict request validation
- thread-safe in-memory request statistics
- graceful shutdown and HTTP server timeouts
- test coverage for business logic and HTTP endpoints

## Run locally

```bash
go run ./cmd/server
```

Or with `make`:

```bash
make run
```

Server listens on `:8080` by default.
You can override it with `PORT`, for example `PORT=9090 go run ./cmd/server`.

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
- In-memory statistics store protected by a mutex to remain safe under concurrent access.
- HTTP server configured with timeouts and graceful shutdown to be closer to production expectations.
- Structured JSON logs through `slog` for easier observability.

## Possible next improvements

- persist statistics in Redis or a database if cross-restart durability is needed
- add request logging / metrics
- expose OpenAPI documentation
- add configuration struct and environment validation for larger deployments
