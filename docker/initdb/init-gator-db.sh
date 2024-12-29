#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER gator;
	CREATE DATABASE gator;
	GRANT ALL PRIVILEGES ON DATABASE gator TO gator;
EOSQL

