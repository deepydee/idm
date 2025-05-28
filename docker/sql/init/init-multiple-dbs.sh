#!/bin/bash
set -e

echo "Checking and creating databases..."

function create_db_if_not_exists() {
  DBNAME=$1
  EXISTS=$(psql -U "$POSTGRES_USER" -tAc "SELECT 1 FROM pg_database WHERE datname='${DBNAME}'")

  if [ "$EXISTS" != "1" ]; then
    echo "Creating database: $DBNAME"
    createdb -U "$POSTGRES_USER" "$DBNAME"
  else
    echo "Database $DBNAME already exists"
  fi
}

create_db_if_not_exists "idm"
create_db_if_not_exists "testing"
