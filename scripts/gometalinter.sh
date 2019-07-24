#!/usr/bin/env bash
if [ "$1" = "update" ]
then
  go get -u github.com/alecthomas/gometalinter
  gometalinter --install --update --no-vendored-linters
fi

gometalinter ./... --exclude=.*pb.go --exclude=.*github.com.* --exclude=gql/generated.go --exclude=api/generated.go --exclude=.*_gen.go --deadline 3m --vendor --disable=gotype --disable=maligned --disable=gocyclo --disable=golint --disable=gotypex --disable=aligncheck --disable=gosec
