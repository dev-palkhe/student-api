# Stage 1: Builder stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Verify the binary exists (for debugging)
RUN ls -l /app

# Stage 2: Final stage (smaller image)
FROM scratch

WORKDIR /app

# Explicitly copy to /app/main (best practice)
COPY --from=builder /app/main /app/main

COPY migrations /app/migrations
COPY .env .

ENV PORT=8080
ENV DATABASE_URL=""

EXPOSE 8080

CMD ["/app/main"]