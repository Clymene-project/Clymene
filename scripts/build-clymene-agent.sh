#!/bin/bash

BRANCH=${BRANCH:?'missing BRANCH env var'}
GOARCH=${GOARCH:-$(go env GOARCH)}
GIT_SHA=$(shell git rev-parse HEAD)
DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

go build -ldflags "-X 'main.branch=${BRANCH}(${GIT_SHA})' -X 'main.BuildTime=${DATE}" -o ./agent ./cmd/agent/main.go