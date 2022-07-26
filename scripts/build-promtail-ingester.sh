#!/bin/bash

BRANCH=${BRANCH:?'missing BRANCH env var'}
#GOARCH=${GOARCH:-$(go env GOARCH)}

GIT_SHA=${GIT_SHA:?'missing GIT_SHA env var'}
#GIT_BRANCH=${GIT_BRANCH:-$(shell git branch)}

DATE=${DATE:?'missing DATE env var'}

mkdir ./out

CGO_ENABLED=0 go build -ldflags "-X 'main.Version=${BRANCH}(${GIT_SHA}))' -X 'main.BuildTime=${DATE}'" -o ./out/promtail-ingester ./cmd/promtail-ingester/main.go

cp ./cmd/promtail-ingester/Dockerfile ./


if [[ ${BRANCH} == "master" ]];
then
  docker build -t bourbonkk/promtail-ingester:latest .
  docker tag bourbonkk/promtail-ingester:latest quay.io/clymene/promtail-ingester:latest
  docker push bourbonkk/promtail-ingester:latest
  docker push quay.io/clymene/promtail-ingester:latest
else
  REPLACE_BRANCH=${BRANCH//"release/"/""}
  docker build -t bourbonkk/promtail-ingester:"${REPLACE_BRANCH}" .
  docker tag bourbonkk/promtail-ingester:"${REPLACE_BRANCH}" quay.io/clymene/promtail-ingester:"${REPLACE_BRANCH}"
  docker push bourbonkk/promtail-ingester:"${REPLACE_BRANCH}"
  docker push quay.io/clymene/promtail-ingester:"${REPLACE_BRANCH}"
fi