#!/usr/bin/env bash

# stop execution if any command fails
# treat unset variables as an error
# print every command before executing it
set -euxo pipefail

if [ "$TRAVIS_PULL_REQUEST" != "false" ] ; then
  git checkout $TRAVIS_PULL_REQUEST_SHA
else
  git checkout $TRAVIS_COMMIT
fi

# Install AWS CLI
pip install --user awscli
export PATH=$PATH:$HOME/.local/bin

