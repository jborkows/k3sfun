# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY . .

# Build the application
RUN go build -o auto-transition .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder and set ownership
COPY --from=builder /app/auto-transition .
RUN chmod +x auto-transition && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    chown -R appuser:appgroup /root

USER appuser

# Run the binary
CMD ["./auto-transition"]
