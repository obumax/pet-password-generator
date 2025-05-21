FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bot cmd/bot/main.go

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /bot /bot
ENTRYPOINT ["/bot"]