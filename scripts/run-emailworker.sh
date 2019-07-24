#!/usr/bin/env bash
go build -ldflags "-X main.commit=`git rev-parse HEAD` -X main.branch=`git rev-parse --abbrev-ref HEAD`" ./cmd/emailworker && ./emailworker --config local.config.toml
