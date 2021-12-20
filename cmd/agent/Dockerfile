FROM alpine:latest as certs
RUN apk add --update --no-cache ca-certificates

FROM alpine

WORKDIR /

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ADD ./cmd/agent/yaml/clymene.yml /etc/clymene/clymene.yml
ADD ./out/clymene-agent /clymene-agent

EXPOSE 15691
EXPOSE 15692

ENTRYPOINT ["/clymene-agent"]
CMD        [ "--config.file=/etc/clymene/clymene.yml"]