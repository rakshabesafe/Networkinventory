FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

WORKDIR /app/services/reconciliation
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /reconciliation-service .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /reconciliation-service .
CMD ["./reconciliation-service"]
