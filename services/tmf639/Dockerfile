FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy the whole monorepo to resolve go.work properly
COPY . .

# Build the specific service
WORKDIR /app/services/tmf639
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /tmf639-service .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /tmf639-service .
EXPOSE 8081
CMD ["./tmf639-service"]
