FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/
COPY --from=builder /app/configs /app/configs/
COPY --from=builder /app/.env /app/
COPY --from=builder /app/keys /app/keys/

EXPOSE 8080
CMD ["/app/main"]
