#!/bin/bash

# Load environment variables from .env file if it exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# Check if a script file is provided as an argument or if input is piped
if [ -z "$1" ] && [ -t 0 ]; then
  echo "Usage: $0 <script-file> or pipe input to the script"
  exit 1
fi

# If a script file is provided, use it
if [ -n "$1" ]; then
  SCRIPT_FILE=$1
  PGPASSWORD=${DB_PASS} docker compose exec -e PGPASSWORD db psql -U ${DB_USER} -d ${DB_NAME} -f ${SCRIPT_FILE}
else
  # Otherwise, read from standard input
  PGPASSWORD=${DB_PASS} docker compose exec -e PGPASSWORD -T db psql -U ${DB_USER} -d ${DB_NAME} -c "$(cat)"
fi
