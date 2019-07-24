#!/usr/bin/env bash
set -e

HOST=54.201.138.198

echo "Invoking pg_dump and writing dump to /tmp/dump.sql"
ssh ubuntu@$HOST 'export PGPASSWORD="2G4UnSL159W2vME92wlEDt9JRvJD0P7n"; pg_dump -Fc -h engine-db-production.cyviylciswvb.us-west-2.rds.amazonaws.com -U engine -p 5432 --no-owner --no-privileges engine > /tmp/dump.sql'
echo "Copying dump.sql"
mkdir -p ~/october_db_dumps
rsync --progress ubuntu@$HOST:/tmp/dump.sql ~/october_db_dumps/dump.sql
echo "Dropping engine_local"
dropdb -h localhost -U postgres engine_local
echo "Creating engine_local"
createdb -h localhost -U postgres engine_local
echo "Importing data into database"
pg_restore -h localhost -U postgres -d engine_local --no-privileges --no-owner ~/october_db_dumps/dump.sql
echo "Emptying device IDs"
psql -h localhost -U postgres -d engine_local -c "UPDATE users SET devices = '{}'"
echo "Cleaning up (remote only)"
ssh ubuntu@$HOST 'rm /tmp/dump.sql'
