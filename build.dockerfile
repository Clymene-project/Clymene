FROM golang:1.16.2 AS build-env

ADD . /clymene
WORKDIR /clymene

RUN mkdir ./out

RUN make build-all-platforms