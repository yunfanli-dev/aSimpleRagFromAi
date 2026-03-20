#!/usr/bin/env bash

set -euo pipefail

docker compose up -d
echo "dependencies started: postgres, redis"
