#!/bin/bash

BRANCH=${BRANCH:?'missing BRANCH env var'}

GIT_SHA=${GIT_SHA:?'missing GIT_SHA env var'}

DATE=${DATE:?'missing DATE env var'}

mkdir ./out

CGO_ENABLED=0 go build -ldflags "-X 'main.Version=${BRANCH}(${GIT_SHA}))' -X 'main.BuildTime=${DATE}'" -o ./out/clymene-gateway ./cmd/gateway/main.go

cp ./cmd/gateway/Dockerfile ./




if [[ ${BRANCH} == "master" ]];
then
  docker build -t bourbonkk/clymene-gateway:latest .
  docker tag bourbonkk/clymene-gateway:latest quay.io/clymene/clymene-gateway:latest
  docker push bourbonkk/clymene-gateway:latest
  docker push quay.io/clymene/clymene-gateway:latest
else
  REPLACE_BRANCH=${BRANCH//"release/"/""}
  docker build -t bourbonkk/clymene-gateway:"${REPLACE_BRANCH}" .
  docker tag bourbonkk/clymene-gateway:"${REPLACE_BRANCH}" quay.io/clymene/clymene-gateway:"${REPLACE_BRANCH}"
  docker push bourbonkk/clymene-gateway:"${REPLACE_BRANCH}"
  docker push quay.io/clymene/clymene-gateway:"${REPLACE_BRANCH}"
fi