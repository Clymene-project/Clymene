FROM golang:1.17.2-bullseye AS build-env

RUN echo "deb http://deb.debian.org/debian bullseye-backports main" >> /etc/apt/sources.list
RUN apt-get update && apt-get install -t bullseye-backports -qy libsystemd-dev

ADD . /clymene
WORKDIR /clymene

#RUN mkdir ./out

RUN make build-all-platforms