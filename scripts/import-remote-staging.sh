#!/usr/bin/env bash
HOST=engine.internal.staging.october.news
DBHOST=engine-db-staging.cyviylciswvb.us-west-2.rds.amazonaws.com
echo "Setting up tunnel to ${HOST}"
ssh -4 -M -S remote-dump -fnNT -L5433:$DBHOST:5432 ubuntu@$HOST
echo "Invoking pg_dump and writing dump to /tmp/dump.sql"
PGPASSWORD=2G4UnSL159W2vME92wlEDt9JRvJD0P7n pg_dump -h localhost -U engine -p 5433 --no-owner --no-privileges engine > /tmp/dump.sql
echo "Closing tunnel"
ssh -4 -S remote-dump -O exit ubuntu@$HOST
echo "Downloading snapshots archive to /tmp/snapshots.tar"
ssh engine@$HOST "tar -zcvf /tmp/snapshots.tar snapshots"
scp engine@$HOST:/tmp/snapshots.tar /tmp/snapshots.tar
echo "Dropping engine_local"
dropdb -h localhost -U postgres engine_local
echo "Creating engine_local"
createdb -h localhost -U postgres engine_local
echo "Importing data into database"
psql -h localhost -U postgres -q -f /tmp/dump.sql engine_local
rm -rf snapshots
tar -xvf /tmp/snapshots.tar
echo "Cleaning up"
ssh engine@$HOST "rm /tmp/snapshots.tar"
rm /tmp/snapshots.tar
rm /tmp/dump.sql
