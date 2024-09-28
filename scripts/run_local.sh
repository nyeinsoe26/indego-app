#!/bin/sh

# Load environment variables from the .env file
export $(grep -v '^#' .env | xargs)

# Step 1: Run PostgreSQL using Docker
echo "Starting PostgreSQL container..."
docker run -d --rm --name indego_db -e POSTGRES_USER="$DATABASE_USER" -e POSTGRES_PASSWORD="$DATABASE_PASSWORD" \
  -e POSTGRES_DB="$DATABASE_NAME" -p "$DATABASE_PORT:$DATABASE_PORT" postgres:17.0-alpine3.20

# Step 2: Wait for PostgreSQL to be ready
echo "Waiting for Postgres to be ready..."
until docker exec indego_db pg_isready -U "$DATABASE_USER" -h "$DATABASE_HOST" -p "$DATABASE_PORT"; do
  echo "Waiting for Postgres..."
  sleep 2
done
echo "Postgres is ready!"

# Step 3: Run migration
echo "Running database migrations..."
go run migration/migrate.go --up --config config.yaml

# Step 4: Run the Go application locally
echo "Starting the Go application..."
go run cmd/main.go --config config.yaml
