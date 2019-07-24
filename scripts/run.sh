#!/usr/bin/env bash
if [ "$1" = "quickstart" ]
then
  ./scripts/build.sh && ./engine --config local.config.toml
else
  go generate rpc/rpc.go && go generate gql/graphql.go && ./scripts/build.sh && ./engine --config local.config.toml
fi
