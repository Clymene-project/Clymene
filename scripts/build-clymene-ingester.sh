#!/bin/bash

BRANCH=${BRANCH:?'missing BRANCH env var'}
GOARCH=${GOARCH:-$(go env GOARCH)}

mkdir ./out

CGO_ENABLED=0 go build -ldflags "-X 'main.Version=${BRANCH}'" -o ./out/clymene-ingester ./cmd/ingester/main.go

cp ./cmd/ingester/Dockerfile ./

docker build -t bourbonkk/clymene-ingester:"${BRANCH}" .
docker push bourbonkk/clymene-ingester:"${BRANCH}"