#!/usr/bin/env bash
set -Eeuo pipefail

export APP_ENV="${APP_ENV:-production}"
export SERVER_PORT="${SERVER_PORT:-8080}"
export POSTGRES_DB="${POSTGRES_DB:-vrtraining}"
export POSTGRES_USER="${POSTGRES_USER:-vrtraining}"
export POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-vrtraining_local_password}"
export DATABASE_URL="${DATABASE_URL:-postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable}"
export JWT_SECRET="${JWT_SECRET:-local-dev-secret-change-before-production}"
export CORS_ALLOWED_ORIGINS="${CORS_ALLOWED_ORIGINS:-http://localhost:3000}"
export REPORT_STORAGE_PATH="${REPORT_STORAGE_PATH:-/app/reports}"
export MIGRATIONS_PATH="${MIGRATIONS_PATH:-/app/migrations}"
export LOG_LEVEL="${LOG_LEVEL:-info}"
export COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-vrtraining}"

HEALTH_URL="http://localhost:${SERVER_PORT}/health"

if [[ "${JWT_SECRET}" == "local-dev-secret-change-before-production" ]]; then
  echo "Warning: using fallback JWT_SECRET. Set VRTRAINING_JWT_SECRET for production." >&2
fi

command -v docker >/dev/null 2>&1 || {
  echo "Docker is not installed or not available in PATH." >&2
  exit 1
}

if ! docker compose version >/dev/null 2>&1; then
  echo "Docker Compose v2 is not available. Install the Docker Compose plugin." >&2
  exit 1
fi

cleanup_on_error() {
  echo "Deployment failed. Dumping container status and recent logs..." >&2
  docker compose ps || true
  docker compose logs --no-color --tail=240 api || true
  docker compose logs --no-color --tail=160 postgres || true
}
trap cleanup_on_error ERR

echo "Stopping old VRTrainingServer containers if they exist..."
docker compose down --remove-orphans || true

echo "Building and starting VRTrainingServer containers..."
docker compose up -d --build --remove-orphans

echo "Waiting for API health endpoint: ${HEALTH_URL}"
for attempt in $(seq 1 90); do
  if curl -fsS "${HEALTH_URL}" >/tmp/vrtraining-health.json; then
    echo "API health check passed."
    cat /tmp/vrtraining-health.json
    echo
    break
  fi

  if [[ "${attempt}" == "90" ]]; then
    echo "API health check failed after 90 attempts." >&2
    exit 1
  fi

  sleep 2
done

echo "Container status:"
docker compose ps

echo "Recent API logs:"
docker compose logs --no-color --tail=200 api

echo "Recent Postgres logs:"
docker compose logs --no-color --tail=120 postgres

echo "VRTrainingServer Docker deployment completed successfully."
