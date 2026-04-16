FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN go build -o /server ./cmd/server

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /server /server

USER app

EXPOSE 8080

ENV PORT=8080

ENTRYPOINT ["/server"]
