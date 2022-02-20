FROM alpine:latest as certs
RUN apk add --update --no-cache ca-certificates

FROM alpine

WORKDIR /

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ADD ./cmd/promtail/promtail-docker-config.yaml /etc/promtail/config.yml
ADD ./out/clymene-promtail /clymene-promtail

EXPOSE 15698
EXPOSE 9080

ENTRYPOINT ["/clymene-promtail"]
CMD        [ "--config.file=/etc/promtail/config.yml"]