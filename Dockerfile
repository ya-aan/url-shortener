FROM golang:1.27 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o url-shortener ./cmd/url-shortener

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/url-shortener .

CMD ["./url-shortener"]