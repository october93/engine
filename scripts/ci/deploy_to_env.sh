#!/usr/bin/env bash

# stop execution if any command fails
# treat unset variables as an error
# print every command before executing it
set -euxo pipefail

REMOTE_USER=engine
HOST=${HOST:-${ENV:-latest}.october.news}

PREFIX="$ENV."
if [ "$ENV" = "production" ]
then
  PREFIX=""
fi

# Inform everyone about deploy start
curl -X POST -H 'Content-type: application/json' --data "{'text': 'Deploying Engine to https://$HOST'}" https://hooks.slack.com/services/...

# Build binaries
go build -ldflags "-X main.environment=$ENV -X main.commit=$TRAVIS_COMMIT -X main.branch=$TRAVIS_BRANCH -X main.version=0.1.0" github.com/october93/engine
go build -ldflags "-X main.environment=$ENV -X main.commit=$TRAVIS_COMMIT -X main.branch=$TRAVIS_BRANCH -X main.version=0.1.0" github.com/october93/engine/cmd/indexer

# Validate config before continuing
./engine validate --config $ENV.config.toml

# Copy binary
rsync engine $REMOTE_USER@$HOST:./engine-new
rsync engine cron@${PREFIX}cron.october.news:./engine
rsync indexer $REMOTE_USER@$HOST:./indexer

# Copy configuration file
rsync $ENV.config.toml $REMOTE_USER@$HOST:./config.toml
rsync $ENV.config.toml cron@${PREFIX}cron.october.news:./config.toml

# Copy template files
ssh $REMOTE_USER@$HOST mkdir -p worker/emailsender
rsync -r worker/emailsender/templates $REMOTE_USER@$HOST:./worker/emailsender

# Copy assets
ssh $REMOTE_USER@$HOST mkdir -p data
rsync -r data/assets $REMOTE_USER@$HOST:./data
ssh $REMOTE_USER@$HOST mkdir -p public/images
rsync -r public/images/anonymous $REMOTE_USER@$HOST:./public/images

# Copy database schema
ssh $REMOTE_USER@$HOST mkdir -p migrations
rsync database.yml $REMOTE_USER@$HOST:.
rsync -r migrations/* $REMOTE_USER@$HOST:./migrations --delete

# Replace old binary
ssh $REMOTE_USER@$HOST "mv engine-new engine"

# Make sure engine will auto-start on boot
ssh $REMOTE_USER@$HOST "sudo systemctl enable engine.service"

# Start engine
ssh $REMOTE_USER@$HOST "sudo systemctl restart engine.service"

# Inform everyone about deploy start

curl -X POST -H 'Content-type: application/json' --data "{'text': 'Commit “$TRAVIS_COMMIT_MESSAGE” was deployed to https://$HOST (Engine)'}" https://hooks.slack.com/services/...

# Check if engine is running
ssh $REMOTE_USER@$HOST pidof engine

rc=$?
if [[ $rc != 0 ]]; then
  curl -X POST -H 'Content-type: application/json' --data "{'text': 'Engine (https://$HOST) did not start'}" https://hooks.slack.com/services/...
fi

# Copy staging reset script if this is a staging deploy
if [ "$ENV" = "staging" ]
then
  rsync scripts/update-staging.sh engine@52.24.33.78:.
fi

# Generate and upload API documentation
go generate github.com/october93/engine/rpc
rsync -r docs $REMOTE_USER@$HOST:.
ssh $REMOTE_USER@$HOST chmod +x docs

# Sync emoji avatars to S3
aws s3 sync data/assets/emojis s3://assets.${ENV}.october.news/emojis
aws s3 sync data/assets/backgrounds s3://assets.${ENV}.october.news/cards/backgrounds
aws s3 sync data/assets/system s3://assets.${ENV}.october.news/system

exit $rc
