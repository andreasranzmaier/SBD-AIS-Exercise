#!/usr/bin/env bash
set -euo pipefail

ENV_FILE="./debug.env"
CONTAINER_DB="postgres18"
CONTAINER_APP="ordersystem-container"
VOLUME_DB="pg18_data"
IMAGE_APP="ordersystem-image"
PGDATA_DIR="/var/lib/postgresql/data"

set -a
source "$ENV_FILE"
set +a

docker network create ordersys-net >/dev/null 2>&1 || true

exists_container() { docker ps -a --format '{{.Names}}' | grep -q "^${1}\$"; }
rm_if_exists() { exists_container "$1" && docker rm -f "$1" >/dev/null || true; }

rm_if_exists "$CONTAINER_DB"
docker volume rm "$VOLUME_DB" >/dev/null 2>&1 || true
docker volume create "$VOLUME_DB" >/dev/null

docker run --rm -v "${VOLUME_DB}:${PGDATA_DIR}" alpine sh -lc "mkdir -p '${PGDATA_DIR}' && chown -R 999:999 '${PGDATA_DIR}'"

docker run -d \
  --name "$CONTAINER_DB" \
  --restart unless-stopped \
  --env-file "$ENV_FILE" \
  --network ordersys-net \
  -v "${VOLUME_DB}:${PGDATA_DIR}" \
  -p 5432:5432 \
  postgres:18 >/dev/null

until docker exec "$CONTAINER_DB" pg_isready -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" >/dev/null 2>&1; do sleep 1; done

rm_if_exists "$CONTAINER_APP"
docker build -t "$IMAGE_APP" -f Dockerfile . >/dev/null

docker run -d \
  --name "$CONTAINER_APP" \
  --network ordersys-net \
  --env-file "$ENV_FILE" \
  -p 3000:3000 \
  "$IMAGE_APP" >/dev/null

echo "up: psql PGPASSWORD='${POSTGRES_PASSWORD}' psql -h ${DB_HOST:-127.0.0.1} -p ${POSTGRES_TCP_PORT} -U ${POSTGRES_USER} -d ${POSTGRES_DB}"
