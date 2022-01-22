#!/bin/bash

BRANCH=${BRANCH:?'missing BRANCH env var'}
#GOARCH=${GOARCH:-$(go env GOARCH)}

GIT_SHA=${GIT_SHA:?'missing GIT_SHA env var'}
#GIT_BRANCH=${GIT_BRANCH:-$(shell git branch)}

DATE=${DATE:?'missing DATE env var'}

mkdir ./out

CGO_ENABLED=0 go build -ldflags "-X 'main.Version=${BRANCH}(${GIT_SHA}))' -X 'main.BuildTime=${DATE}'" -o ./out/clymene-agent ./cmd/agent/main.go

cp ./cmd/agent/Dockerfile ./

if [[ ${BRANCH} == "master" ]];
then
  docker build -t bourbonkk/clymene-agent:latest .
  docker tag bourbonkk/clymene-agent:latest quay.io/clymene/clymene-agent:latest
  docker push bourbonkk/clymene-agent:latest
  docker push quay.io/clymene/clymene-agent:latest
else
  docker build -t bourbonkk/clymene-agent:"${BRANCH}" .
  docker tag bourbonkk/clymene-agent:"${BRANCH}" quay.io/clymene/clymene-agent:"${BRANCH}"
  docker push bourbonkk/clymene-agent:"${BRANCH}"
  docker push quay.io/clymene/clymene-agent:"${BRANCH}"
fi
