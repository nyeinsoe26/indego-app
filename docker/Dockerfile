# Stage 1: Build the Go app
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum, and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the migration binary
RUN go build -o /bin/migrate migration/migrate.go

# Build the main application binary
RUN go build -o /bin/app cmd/main.go

# Stage 2: Final image with only the binaries
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy binaries from the builder stage
COPY --from=builder /bin/migrate /bin/migrate
COPY --from=builder /bin/app /bin/app

# Copy config and migration files
COPY config.yaml /app/config.yaml
COPY .env /app/.env
COPY templates /app/templates
COPY migration /app/migration

# Run the migration and start the app
CMD /bin/migrate --config /app/config.yaml --up && /bin/app --config /app/config.yaml
