#!/bin/bash

# Load environment variables from .env file if it exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

PGPASSWORD=${DB_PASS} docker compose exec -e PGPASSWORD db psql -U ${DB_USER} -d ${DB_NAME}
