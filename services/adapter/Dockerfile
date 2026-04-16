FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

WORKDIR /app/services/adapter
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /adapter-service .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /adapter-service .
CMD ["./adapter-service"]
