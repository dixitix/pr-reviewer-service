# Стейдж сборки
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o pr-reviewer-service ./cmd/pr-reviewer-service

# Стейдж рантайма
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/pr-reviewer-service /usr/local/bin/pr-reviewer-service

CMD ["pr-reviewer-service"]
