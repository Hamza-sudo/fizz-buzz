FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app
RUN chown -R app:app /app

COPY --from=builder /server /server

USER app

EXPOSE 8080

ENV PORT=8080
ENV STATS_DB_DSN=file:/app/fizzbuzz_stats.db

HEALTHCHECK --interval=30s --timeout=3s \
	CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["/server"]
