FROM golang:1.16.2 AS build-env

ADD . /datagate

WORKDIR /datagate

WORKDIR /Clymene/cmd/agent
RUN CGO_ENABLED=0 go build -ldflags "-X 'main.Version=${bamboo.inject.tag}(${bamboo.inject.buildVersion})' -X 'main.BuildTime=${bamboo.inject.date}'" -o /agent

FROM alpine:latest as certs
RUN apk add --update --no-cache ca-certificates

FROM busybox

WORKDIR /

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=build-env /agent /

CMD ["/agent"]