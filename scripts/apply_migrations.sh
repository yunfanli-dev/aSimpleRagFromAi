#!/usr/bin/env bash

set -euo pipefail

if [[ -z "${POSTGRES_DSN:-}" ]]; then
  echo "POSTGRES_DSN is required"
  exit 1
fi

psql "$POSTGRES_DSN" -f migrations/0001_init_extensions.sql
psql "$POSTGRES_DSN" -f migrations/0002_init_schema.sql
psql "$POSTGRES_DSN" -f migrations/0003_init_indexes.sql

echo "migrations applied"
