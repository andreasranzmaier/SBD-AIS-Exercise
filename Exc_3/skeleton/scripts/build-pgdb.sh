#!/usr/bin/env bash
set -euo pipefail

# Location of the environment file:
ENV_FILE="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)/debug.env"

source "$ENV_FILE"

CONTAINER_NAME="postgres18"
VOLUME_NAME="pg18_data"

# Create the persistent Docker volume if it doesn't exist already
if ! docker volume inspect "$VOLUME_NAME" >/dev/null 2>&1; then
  echo "Creating volume: $VOLUME_NAME"
  docker volume create "$VOLUME_NAME" >/dev/null
fi

# If a container with the same name exists (even if stopped), remove it
if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}\$"; then
  echo "Removing existing container: $CONTAINER_NAME"
  docker rm -f "$CONTAINER_NAME" >/dev/null
fi

# can still lead to port conflicts if another container is using the same port

# Run PostgreSQL 18 with persistent storage and healthcheck
echo "Starting container: $CONTAINER_NAME (postgres:18)"

docker run -d \
  --name "$CONTAINER_NAME" \
  --restart unless-stopped \
  --env-file "$ENV_FILE" \
  -e POSTGRES_DB="$POSTGRES_DB" \
  -e POSTGRES_USER="$POSTGRES_USER" \
  -e POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
  -p "${POSTGRES_TCP_PORT}:${POSTGRES_TCP_PORT}" \
  -v "${VOLUME_NAME}:/var/lib/postgresql/18/docker" \
  postgres:18

echo "Container '$CONTAINER_NAME' is starting."
echo "Port mapping: localhost:${POSTGRES_TCP_PORT} -> container:${POSTGRES_TCP_PORT}"
echo "Data volume:  ${VOLUME_NAME} mounted at /var/lib/postgresql/18/docker"
