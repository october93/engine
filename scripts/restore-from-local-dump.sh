#!/usr/bin/env bash
set -e

HOST=54.201.138.198

echo "Dropping engine_local"
dropdb -h localhost -U postgres engine_local
echo "Creating engine_local"
createdb -h localhost -U postgres engine_local
echo "Importing data into database from local dump (at ~/october_db_dumps)"
pg_restore -h localhost -U postgres -d engine_local --no-privileges --no-owner ~/october_db_dumps/dump.sql
echo "Emptying device IDs"
psql -h localhost -U postgres -d engine_local -c "UPDATE users SET devices = '{}'"
