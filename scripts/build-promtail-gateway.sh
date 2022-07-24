#!/bin/bash

BRANCH=${BRANCH:?'missing BRANCH env var'}

GIT_SHA=${GIT_SHA:?'missing GIT_SHA env var'}

DATE=${DATE:?'missing DATE env var'}

mkdir ./out

CGO_ENABLED=0 go build -ldflags "-X 'main.Version=${BRANCH}(${GIT_SHA}))' -X 'main.BuildTime=${DATE}'" -o ./out/promtail-gateway ./cmd/promtail-gateway/main.go

cp ./cmd/promtail-gateway/Dockerfile ./




if [[ ${BRANCH} == "master" ]];
then
  docker build -t bourbonkk/promtail-gateway:latest .
  docker tag bourbonkk/promtail-gateway:latest quay.io/clymene/promtail-gateway:latest
  docker push bourbonkk/promtail-gateway:latest
  docker push quay.io/clymene/promtail-gateway:latest
else
  REPLACE_BRANCH=${BRANCH//"release/"/""}
  docker build -t bourbonkk/promtail-gateway:"${REPLACE_BRANCH}" .
  docker tag bourbonkk/promtail-gateway:"${REPLACE_BRANCH}" quay.io/clymene/promtail-gateway:"${REPLACE_BRANCH}"
  docker push bourbonkk/promtail-gateway:"${REPLACE_BRANCH}"
  docker push quay.io/clymene/promtail-gateway:"${REPLACE_BRANCH}"
fi