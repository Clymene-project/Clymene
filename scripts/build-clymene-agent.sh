#!/bin/bash

BRANCH=${BRANCH:?'missing BRANCH env var'}
GOARCH=${GOARCH:-$(go env GOARCH)}
GIT_SHA=$(shell git rev-parse HEAD)
#DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

mkdir ./out

go build -ldflags "-X 'main.Version=${BRANCH}(${GIT_SHA})' -X 'main.BuildTime=${DATE}" -o ./out/clymene-agent ./cmd/agent/main.go

docker build -f ../cmd/agent/Dockerfile -t bourbonkk/clymene:"${BRANCH}"
docker push bourbonkk/clymene:"${BRANCH}"