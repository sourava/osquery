FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o service ./cmd/server/main.go

FROM alpine:3.17

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/service .

EXPOSE 8080

CMD ["./service"]