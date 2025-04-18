# Build
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api-gw .

# Image
FROM debian:latest

WORKDIR /app

COPY --from=builder /app/api-gw .

RUN setcap 'cap_net_bind_service=+ep' ./api-gw

EXPOSE 8080

ENTRYPOINT ["./api-gw"]