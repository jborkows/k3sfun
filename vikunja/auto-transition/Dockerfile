# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY main.go go.mod .

# Build the application
RUN go build -o auto-transition main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/auto-transition .

# Run the binary
CMD ["./auto-transition"]
