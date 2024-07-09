# Build stage
FROM golang:1.22.4 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chatapp .

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/chatapp .

EXPOSE 8080
CMD ["./chatapp"]