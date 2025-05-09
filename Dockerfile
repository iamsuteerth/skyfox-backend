FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /app/skyfox ./server/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/skyfox /app/skyfox

COPY --from=builder /app/assets /app/assets

RUN chmod +x /app/skyfox

EXPOSE 8080

ENTRYPOINT ["/app/skyfox"]
