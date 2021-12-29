#!/bin/bash

BRANCH=${BRANCH:?'missing BRANCH env var'}
#GOARCH=${GOARCH:-$(go env GOARCH)}

GIT_SHA=${GIT_SHA:?'missing GIT_SHA env var'}
#GIT_BRANCH=${GIT_BRANCH:-$(shell git branch)}

DATE=${DATE:?'missing DATE env var'}

mkdir ./out

CGO_ENABLED=0 go build -ldflags "-X 'main.Version=${BRANCH}(${GIT_SHA}))' -X 'main.BuildTime=${DATE}'" -o ./out/clymene-ingester ./cmd/ingester/main.go

cp ./cmd/ingester/Dockerfile ./

docker build -t bourbonkk/clymene-ingester:"${BRANCH}" .

docker tag bourbonkk/clymene-ingester:"${BRANCH}" quay.io/clymene/clymene-ingester:"${BRANCH}"

docker push bourbonkk/clymene-ingester:"${BRANCH}"
docker push quay.io/clymene/clymene-ingester:"${BRANCH}"