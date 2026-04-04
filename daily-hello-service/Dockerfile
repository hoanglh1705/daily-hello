FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git if needed for some modules
RUN apk add --no-cache git tzdata

COPY go.mod go.sum ./
COPY vendor/ vendor/
COPY . .

# Build from vendored dependencies, no network download needed
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /app/server ./cmd/server/main.go

FROM alpine:3.23

WORKDIR /app

# Install tzdata and ca-certificates
RUN apk add --no-cache tzdata ca-certificates

COPY --from=builder /app/server /app/server
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy config or any other static assets if they exist (optional, handle carefully if mapped via volumes)
COPY --from=builder /app/config /app/config

EXPOSE 8080

CMD ["./server"]
