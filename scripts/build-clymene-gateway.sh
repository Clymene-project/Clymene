#!/bin/bash

BRANCH=${BRANCH:?'missing BRANCH env var'}

GIT_SHA=${GIT_SHA:?'missing GIT_SHA env var'}

DATE=${DATE:?'missing DATE env var'}

mkdir ./out

CGO_ENABLED=0 go build -ldflags "-X 'main.Version=${BRANCH}(${GIT_SHA}))' -X 'main.BuildTime=${DATE}'" -o ./out/clymene-gateway ./cmd/gateway/main.go

cp ./cmd/gateway/Dockerfile ./

docker build -t bourbonkk/clymene-gateway:"${BRANCH}" .

docker tag bourbonkk/clymene-gateway:"${BRANCH}" quay.io/clymene/clymene-gateway:"${BRANCH}"

docker push bourbonkk/clymene-gateway:"${BRANCH}"
docker push quay.io/clymene/clymene-gateway:"${BRANCH}"