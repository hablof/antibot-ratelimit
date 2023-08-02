FROM golang:1.20-alpine AS builder

WORKDIR /home/app

COPY . .

RUN go mod download

RUN go build -o ratelimiter cmd/main.go



FROM alpine:latest AS app

WORKDIR /root/app

COPY --from=builder /home/app/ratelimiter .

COPY config.yml .

EXPOSE 8050

CMD ["./ratelimiter"]