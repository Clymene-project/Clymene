#FROM golang:1.16.3 AS build-env
#
#ADD . /Clymene
#WORKDIR /Clymene
#RUN CGO_ENABLED=0 go build -ldflags "-X 'main.Version=${BRANCH_NAME}(${GIT_COMMIT})' -X 'main.BuildTime=${BUILD_ID}'" -o /agent

FROM alpine:latest as certs
RUN apk add --update --no-cache ca-certificates

FROM busybox

WORKDIR /

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

#COPY --from=build-env /agent /
#COPY --from=build-env /Clymene/cmd/agent/yaml/clymene.yml /etc/clymene/clymene.yml
#COPY --from=build-env /Clymene/out/clymene-agent /clymene-agent

ADD ./cmd/agent/yaml/clymene.yml /etc/clymene/clymene.yml
ADD ./out/clymene-agent /clymene-agent

EXPOSE 15691
EXPOSE 15692

ENTRYPOINT ["/clymene-agent"]
CMD        [ "--config.file=/etc/clymene/clymene.yml"]