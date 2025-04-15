FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o server ./server

FROM nginx:alpine

WORKDIR /app

RUN apk add --no-cache bash gettext supervisor

COPY --from=builder /app/server .

COPY nginx.conf.template /etc/nginx/nginx.conf.template
COPY supervisord.conf /etc/supervisord.conf
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENV GIN_MODE=release

RUN mkdir -p /var/log/supervisor /var/run/supervisor /var/log/nginx /var/run/nginx /tmp/nginx

RUN chmod 777 /var/log/supervisor /var/run/supervisor /var/log/nginx /var/run/nginx /tmp/nginx

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
