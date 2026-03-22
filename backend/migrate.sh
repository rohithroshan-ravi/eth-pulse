#!/bin/bash

# Database Migration Script
# Usage: ./migrate.sh [up|down|status]

set -e

COMMAND="${1:-up}"

echo "🚀 Running migration: $COMMAND"

cd "$(dirname "$0")/.."

case "$COMMAND" in
  up)
    go run ./migrations/migrate.go up
    ;;
  down)
    go run ./migrations/migrate.go down
    ;;
  status)
    go run ./migrations/migrate.go status
    ;;
  *)
    echo "Usage: $0 [up|down|status]"
    exit 1
    ;;
esac
