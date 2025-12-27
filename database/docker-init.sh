#!/bin/bash
set -e

echo "Running database migrations..."

# Run all migration files in order
for f in /docker-entrypoint-initdb.d/migrations/*.sql; do
    if [ -f "$f" ]; then
        echo "Running migration: $f"
        psql -U postgres -d taskmanagement -f "$f"
    fi
done

echo "Migrations completed successfully!"
