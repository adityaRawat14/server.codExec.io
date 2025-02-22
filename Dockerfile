

FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/app/main.go
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
RUN chmod +x main
ENTRYPOINT ["/app/main"]