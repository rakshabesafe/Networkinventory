FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

WORKDIR /app/services/tmf638
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /tmf638-service .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /tmf638-service .
EXPOSE 8082
CMD ["./tmf638-service"]
