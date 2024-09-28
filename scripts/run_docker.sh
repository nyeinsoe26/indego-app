#!/bin/sh

# Load environment variables from the .env file
export $(grep -v '^#' .env | xargs)

# Build the Docker image
echo "Building the Docker image for the app..."
docker build -t indego-app -f docker/Dockerfile .

# Create a Docker network if it doesn't exist
docker network create nyein_indego_network || true

# Step 1: Run PostgreSQL using Docker
echo "Starting PostgreSQL container..."
docker run -d --rm --name indego_db --network nyein_indego_network -e POSTGRES_USER="$DATABASE_USER" -e POSTGRES_PASSWORD="$DATABASE_PASSWORD" \
  -e POSTGRES_DB="$DATABASE_NAME" -p "$DATABASE_PORT:$DATABASE_PORT" postgres:17.0-alpine3.20

# Step 2: Wait for PostgreSQL to be ready
echo "Waiting for Postgres to be ready..."
until docker exec indego_db pg_isready -U "$DATABASE_USER" -h "$DATABASE_HOST" -p "$DATABASE_PORT"; do
  echo "Waiting for Postgres..."
  sleep 2
done
echo "Postgres is ready!"

# Step 3: Start the app container with DATABASE_HOST set to the name of the PostgreSQL container
echo "Starting the app container..."
docker run --rm --name indego_app \
  --network nyein_indego_network \
  --env-file .env \
  -e DATABASE_HOST=indego_db \
  -p "$SERVER_PORT:$SERVER_PORT" indego-app

echo "App is running. Showing logs..."
docker logs -f indego_app
