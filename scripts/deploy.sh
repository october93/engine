#!/usr/bin/env bash
set -euxo pipefail

HOST=${HOST:-${ENV:-latest}.october.news}
NAME=$(git config --global user.name | awk '{print $1}')
COMMIT=$(git log -1 --pretty=%B | head -n 1)

# Build binary
go build -ldflags "-X main.environment=$ENV -X main.commit=`git rev-parse HEAD` -X main.branch=`git rev-parse --abbrev-ref HEAD` -X main.version=0.1.0" .

# # Copy configuration and assets
scp $ENV.config.toml engine@$HOST:config.toml
scp database.yml engine@$HOST:.
scp -r migrations engine@$HOST:.
ssh engine@$HOST mkdir -p ./worker/emailsender
scp -r worker/emailsender/templates engine@$HOST:./worker/emailsender

# Copy binary
scp engine engine@$HOST:engine-new

# Stop Engine and replace binary
ssh engine@$HOST sudo systemctl stop engine.service
ssh engine@$HOST mv engine-new engine

# Start Engine again and check if it is running
ssh engine@$HOST sudo systemctl start engine.service
ssh engine@$HOST pidof engine

if [ "$ENV" = "staging" ]
then
  scp scripts/update-staging.sh engine@engine.staging.october.news:.
fi

# Inform other about manual deploy
curl -X POST -H 'Content-type: application/json' --data "{'text': '$NAME has deployed $COMMIT to $HOST'}" https://hooks.slack.com/services/...
