# Stage 1: Builder stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Stage 2: Final stage (smaller image)
FROM scratch

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy necessary files (migrations if needed)
COPY migrations /app/migrations

# Copy .env file (optional, but recommended to use environment variables)

# Set environment variables (these can be overridden at runtime)
ENV PORT=8080      
ENV DATABASE_URL="" 
# Correct format (equals sign separated)
#Important: Set this at runtime

# Expose the port
EXPOSE 8080

# Command to run the application
CMD ["/app/main"]
