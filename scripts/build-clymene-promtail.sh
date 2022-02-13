#!/bin/bash

BRANCH=${BRANCH:?'missing BRANCH env var'}
#GOARCH=${GOARCH:-$(go env GOARCH)}

GIT_SHA=${GIT_SHA:?'missing GIT_SHA env var'}
#GIT_BRANCH=${GIT_BRANCH:-$(shell git branch)}

DATE=${DATE:?'missing DATE env var'}

# install journalctl
echo "deb http://deb.debian.org/debian bullseye-backports main" >> /etc/apt/sources.list

apt-get update && apt-get install -qy tzdata ca-certificates
apt-get install -t bullseye-backports -qy libsystemd-dev

mkdir ./out

CGO_ENABLED=0 go build -ldflags "-X 'main.Version=${BRANCH}(${GIT_SHA}))' -X 'main.BuildTime=${DATE}'" -o ./out/clymene-promtail ./cmd/promtail/main.go

cp ./cmd/promtail/Dockerfile ./

if [[ ${BRANCH} == "master" ]];
then
  docker build -t bourbonkk/clymene-promtail:latest .
  docker tag bourbonkk/clymene-promtail:latest quay.io/clymene/clymene-promtail:latest
  docker push bourbonkk/clymene-promtail:latest
  docker push quay.io/clymene/clymene-promtail:latest
else
  docker build -t bourbonkk/clymene-promtail:"${BRANCH}" .
  docker tag bourbonkk/clymene-promtail:"${BRANCH}" quay.io/clymene/clymene-promtail:"${BRANCH}"
  docker push bourbonkk/clymene-promtail:"${BRANCH}"
  docker push quay.io/clymene/clymene-promtail:"${BRANCH}"
fi
