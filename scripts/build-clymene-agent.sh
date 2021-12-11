#!/bin/bash

#BRANCH=${BRANCH:?'missing BRANCH env var'}
#GOARCH=${GOARCH:-$(go env GOARCH)}

GIT_SHA=${GIT_SHA:-$(shell git rev-parse HEAD)}
GIT_BRANCH=${GIT_BRANCH:-$(shell git branch)}

DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

mkdir ./out

CGO_ENABLED=0 go build -ldflags "-X 'main.Version=${GIT_BRANCH}(${GIT_SHA}))' -X 'main.BuildTime=${DATE}'" -o ./out/clymene-agent ./cmd/agent/main.go

cp ./cmd/agent/Dockerfile ./

docker build -t bourbonkk/clymene-agent:"${BRANCH}" .
docker push bourbonkk/clymene-agent:"${BRANCH}"