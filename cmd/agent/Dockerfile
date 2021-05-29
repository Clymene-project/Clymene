FROM golang:1.16.2 AS build-env

ADD . /Clymene

WORKDIR /Clymene

RUN CGO_ENABLED=0 go build -ldflags "-X 'main.Version=${BRANCH_NAME}(${GIT_COMMIT})' -X 'main.BuildTime=${BUILD_ID}'" -o /agent

FROM alpine:latest as certs
RUN apk add --update --no-cache ca-certificates

FROM busybox

WORKDIR /

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=build-env /agent /

CMD ["/agent"]