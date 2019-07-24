#!/usr/bin/env bash

# stop execution if any command fails
# treat unset variables as an error
# print every command before executing it
set -euxo pipefail

if [ "$TRAVIS_BRANCH" == "master" ] && [ "$TRAVIS_PULL_REQUEST" == "false" ]; then
  instances=( 34.221.108.135 54.245.187.132 )
  for i in "${instances[@]}"
  do
    HOST=$i ENV=production ./scripts/ci/deploy_to_env.sh
  done

  if git diff --name-only $TRAVIS_COMMIT_RANGE | grep -q "emailsender"; then
    TASK=emailsender ENV=production ./scripts/deploy-worker.sh
  fi
  if git diff --name-only $TRAVIS_COMMIT_RANGE | grep -q "activityrecorder"; then
    TASK=emailsender ENV=production ./scripts/deploy-worker.sh
  fi
  # JIRA webhook to move issues into current fix version
  curl -X POST -H 'Content-type: application/json' https://automation.codebarrel.io/pro/hooks/2a8b8bc2ea67c173a6fc6aeed132e01a40acafc8
fi

if [ "$TRAVIS_BRANCH" == "development" ] && [ "$TRAVIS_PULL_REQUEST" == "false" ]; then
  instances=( 54.191.235.10 34.219.213.219 )
  for i in "${instances[@]}"
  do
    HOST=$i ENV=development ./scripts/ci/deploy_to_env.sh
  done

  if git diff --name-only $TRAVIS_COMMIT_RANGE | grep -q "emailsender"; then
    TASK=emailsender ENV=development ./scripts/deploy-worker.sh
  fi
  if git diff --name-only $TRAVIS_COMMIT_RANGE | grep -q "activityrecorder"; then
    TASK=emailsender ENV=development ./scripts/deploy-worker.sh
  fi
fi
