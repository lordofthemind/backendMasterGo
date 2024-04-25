#!/bin/sh

set -e

echo "Runing DB Migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Starting Server"

exec "$@"