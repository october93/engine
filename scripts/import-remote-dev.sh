#!/usr/bin/env bash
set -e

HOST=54.191.235.10

# Ensure previously failed tunnels are closed
# ssh -4 -S remote-dump -O exit ubuntu@$HOST

echo "Setting up tunnel to ${HOST}"
ssh -4 -M -S remote-dump -fnNT -L5433:engine-db-development.cyviylciswvb.us-west-2.rds.amazonaws.com:5432 ubuntu@$HOST
echo "Invoking pg_dump and writing dump to /tmp/dump.sql"
PGPASSWORD=zvsxsnNmd93fHqYmeYRqaFz2J0bmn83D pg_dump -h localhost -U engine -p 5433 --no-owner --no-privileges engine > /tmp/dump.sql
echo "Closing tunnel"
ssh -4 -S remote-dump -O exit ubuntu@$HOST
echo "Dropping engine_local"
dropdb -h localhost -U postgres engine_local
echo "Creating engine_local"
createdb -h localhost -U postgres engine_local
echo "Importing data into database"
psql -h localhost -U postgres -q -f /tmp/dump.sql engine_local
echo "Cleaning up"
rm /tmp/dump.sql
