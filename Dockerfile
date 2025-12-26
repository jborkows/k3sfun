# syntax=docker/dockerfile:1

FROM golang:1.25.5 AS builder

WORKDIR /src

ARG BUILD_VERSION=dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.buildVersion=${BUILD_VERSION}" -o /out/shopping ./cmd/shopping

FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/shopping /app/shopping
COPY --from=builder /src/web /app/web

EXPOSE 8080

ENV ADDR=:8080

ENTRYPOINT ["/app/shopping"]
