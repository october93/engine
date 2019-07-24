#!/usr/bin/env bash

# stop execution if any command fails
# treat unset variables as an error
# print every command before executing it
set -euxo pipefail

if [ "$TRAVIS_PULL_REQUEST" != "false" ]; then
  BRANCH=${TRAVIS_PULL_REQUEST_BRANCH:-TRAVIS_BRANCH}
else
  BRANCH=$TRAVIS_BRANCH
fi

cd $GOPATH/src/github.com/october93/engine

# Run go test and gometalinter
go test -cover -race -bench=. -benchmem ./... 2>&1 | tee new.log

# Build engine
go build

# Validate config files
./engine validate --config production.config.toml
./engine validate --config development.config.toml
./engine validate --config staging.config.toml
./engine validate --config benchmark.config.toml

# Install benchstat to summarize and compare benchmark results
go get golang.org/x/perf/cmd/benchstat

# Upload output for subsequent benchmark comparisons
aws s3 cp new.log s3://benchmarks.october.news/$BRANCH.log

# Test benchmark delta against base branch if this is a pull request and post
# it as a comment to the pull request if it has been requested with:
# - [x] Report Benchmarks
if [ "$TRAVIS_PULL_REQUEST" != "false" ]; then
  BODY=$(curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/repos/$TRAVIS_PULL_REQUEST_SLUG/issues/$TRAVIS_PULL_REQUEST| jq .body)
  if [[ $BODY == *"- [x] Report Benchmarks"* ]]; then
    # Download existing benchmark of the branch this pull request is targeting
    aws s3 cp s3://benchmarks.october.news/$TRAVIS_BRANCH.log old.log
    benchstat -alpha 10 old.log new.log > benchstat.log
    # Use Python to escape text file for JSON and sed to truncate first and last character
    curl -H "Authorization: token $GITHUB_TOKEN" -XPOST -d "{\"body\": \"\`\`\`\n$(cat benchstat.log | python -c 'import json,sys; print(json.dumps(sys.stdin.read()))' | sed 's/^.\(.*\).$/\1/')\n\`\`\`\"}" https://api.github.com/repos/$TRAVIS_PULL_REQUEST_SLUG/issues/$TRAVIS_PULL_REQUEST/comments
  fi
fi
