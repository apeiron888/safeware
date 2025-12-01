#!/bin/bash
# Migration runner script for SIMS

set -e

echo "Running SIMS database migrations..."

# Check if migrate is installed
if ! command -v migrate &> /dev/null; then
    echo "Installing golang-migrate..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
fi

# Database URL
DB_URL="postgresql://postgres:postgres@localhost:5432/sims?sslmode=disable"

# Run migrations
migrate -path migrations -database "$DB_URL" up

echo "Migrations completed successfully!"
echo "You can check the status with: migrate -path migrations -database '$DB_URL' version"
