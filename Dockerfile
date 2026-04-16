# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/crochetbot cmd/server/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/crochetbot .

# Create uploads directory
RUN mkdir -p /app/uploads

# Set environment variables
ENV PORT=8080
ENV UPLOAD_DIR=/app/uploads

EXPOSE 8080

CMD ["./crochetbot"]
