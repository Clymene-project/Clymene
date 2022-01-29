# for test
FROM golang:1.16.2 AS build-env

ADD . /clymene
WORKDIR /clymene

RUN CGO_ENABLED=0 go build -ldflags "-X 'main.Version=0.1'" -o ./out/clymene-agent ./cmd/agent/main.go


FROM alpine:latest as certs
RUN apk add --update --no-cache ca-certificates

FROM alpine

WORKDIR /

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /clymene/out/clymene-agent /clymene-agent

COPY --from=build-env /clymene/cmd/agent/yaml/clymene.yml /etc/clymene/clymene.yml

EXPOSE 15691
EXPOSE 15692

ENTRYPOINT ["/clymene-agent"]
CMD        [ "--config.file=/etc/clymene/clymene.yml"]

