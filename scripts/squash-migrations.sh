#!/usr/bin/env bash
set -e

soda migrate -e local
pg_dump -h localhost -p 5432 -U postgres --schema-only --no-owner --no-privileges engine_local > migrations/base.sql
git rm migrations/[0-9]*
