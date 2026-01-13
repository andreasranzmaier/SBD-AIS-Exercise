#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="${SCRIPT_DIR}/../debug.env"

# Load and export debug.env 
echo "Exporting variables from $ENV_FILE to current shell..."
set -a
source "$ENV_FILE"
set +a

CONTAINER_NAME="${CONTAINER_NAME:-postgres18}"
VOLUME_NAME="${VOLUME_NAME:-pg18_data}"

# Enforce storage layout:
PG_PARENT="/var/lib/postgresql/18"
PGDATA_DIR="${PGDATA_DIR:-/var/lib/postgresql/18/docker}"

RESET="${1:-}"  # use --reset to wipe data

# Check if a container exists
exists_container() { docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}\$"; }

remove_container_if_exists() {
  param="$1"
  if exists_container; then
    echo "Removing existing container: $param"
    docker rm -f "$param" >/dev/null
  fi
}

create_volume_if_missing() {
  if ! docker volume inspect "$VOLUME_NAME" >/dev/null 2>&1; then
    echo "Creating volume: $VOLUME_NAME"
    docker volume create "$VOLUME_NAME" >/dev/null
  fi
}

wipe_volume() {
  echo "Deleting volume (data loss): $VOLUME_NAME"
  docker volume rm "$VOLUME_NAME" >/dev/null || true
  echo "Recreating volume: $VOLUME_NAME"
  docker volume create "$VOLUME_NAME" >/dev/null
}

prepare_volume_permissions() {
  echo "Preparing volume permissions for postgres (uid 999) under ${PG_PARENT}..."
  docker run --rm -v "${VOLUME_NAME}:${PG_PARENT}" alpine sh -lc "
    mkdir -p '${PGDATA_DIR}' &&
    chown -R 999:999 '${PG_PARENT}'
  "
}

#  MAIN 
remove_container_if_exists "$CONTAINER_NAME"
wipe_volume
create_volume_if_missing
prepare_volume_permissions

# -d deamonized runs in background
# script continues after this
echo "Starting container: $CONTAINER_NAME (postgres:18)"
docker run -d \
  --name "$CONTAINER_NAME" \
  --restart unless-stopped \
  --env-file $ENV_FILE \
  -p "${POSTGRES_TCP_PORT}:5432" \
  -v "${VOLUME_NAME}:${PG_PARENT}" \
  postgres:18 >/dev/null

echo "Container '$CONTAINER_NAME' is starting."
echo "Port mapping: localhost:${POSTGRES_TCP_PORT} -> container:5432"
echo "Data volume:  ${VOLUME_NAME} mounted at ${PGDATA_DIR}"
echo "Waiting for Postgres to become ready..."

# Wait until ready to accept connections
until docker exec "$CONTAINER_NAME" pg_isready -U "${POSTGRES_USER}" >/dev/null 2>&1; do
  printf "."
  sleep 1
done

# Choose admin user: prefer POSTGRES_USER, fallback to 'postgres'
ADMIN_USER="$POSTGRES_USER"
if ! docker exec "$CONTAINER_NAME" psql -U "$ADMIN_USER" -d postgres -tAc "SELECT 1" >/dev/null 2>&1; then
  if docker exec "$CONTAINER_NAME" psql -U postgres -d postgres -tAc "SELECT 1" >/dev/null 2>&1; then
    ADMIN_USER="postgres"
  fi
fi

# sql to check if the docker user and DB exist, create them if not to stay consistent with debug.env
# grant all privileges on the DB to the user  
echo "Postgres is ready as '$ADMIN_USER'"
docker exec -i "$CONTAINER_NAME" psql -U "$ADMIN_USER" -d postgres -v ON_ERROR_STOP=1 <<SQL
DO \$\$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '${POSTGRES_USER}') THEN
    CREATE ROLE ${POSTGRES_USER} LOGIN PASSWORD '${POSTGRES_PASSWORD}';
  END IF;
END\$\$ ;

DO \$\$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_database WHERE datname = '${POSTGRES_DB}') THEN
    EXECUTE format('CREATE DATABASE %I OWNER %I TEMPLATE template1', '${POSTGRES_DB}', '${POSTGRES_USER}');
  END IF;
END\$\$;

GRANT ALL PRIVILEGES ON DATABASE "${POSTGRES_DB}" TO "${POSTGRES_USER}";
SQL

echo
echo "Ready! Connect with:"
echo "PGPASSWORD='${POSTGRES_PASSWORD}' psql -h ${DB_HOST:-127.0.0.1} -p ${POSTGRES_TCP_PORT} -U ${POSTGRES_USER} -d ${POSTGRES_DB}"

remove_container_if_exists "ordersystem-container"
# docker build 
# Build
docker build -t ordersystem-image -f ../Dockerfile ..

# Run (note the position of --,env-file)
docker run -d \
  --name ordersystem-container \
  --network=host \
  --env-file "$ENV_FILE" \
  ordersystem-image
