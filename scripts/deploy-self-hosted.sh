#!/usr/bin/env bash
set -Eeuo pipefail

SERVER_PORT="${SERVER_PORT:-8080}"
HEALTH_URL="http://localhost:${SERVER_PORT}/health"
COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-vrtraining}"
export COMPOSE_PROJECT_NAME

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
  docker compose logs --no-color --tail=200 api || true
  docker compose logs --no-color --tail=120 postgres || true
}
trap cleanup_on_error ERR

echo "Building and starting VRTrainingServer containers..."
docker compose up -d --build --remove-orphans

echo "Waiting for API health endpoint: ${HEALTH_URL}"
for attempt in $(seq 1 60); do
  if curl -fsS "${HEALTH_URL}" >/tmp/vrtraining-health.json; then
    echo "API health check passed."
    cat /tmp/vrtraining-health.json
    echo
    break
  fi

  if [[ "${attempt}" == "60" ]]; then
    echo "API health check failed after 60 attempts." >&2
    exit 1
  fi

  sleep 2
done

echo "Container status:"
docker compose ps

echo "Recent API logs:"
docker compose logs --no-color --tail=160 api

echo "Recent Postgres logs:"
docker compose logs --no-color --tail=80 postgres

echo "VRTrainingServer Docker deployment completed successfully."
