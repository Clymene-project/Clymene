FROM alpine:latest as certs
RUN apk add --update --no-cache ca-certificates

FROM alpine

WORKDIR /

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ADD ./out/promtail-ingester /promtail-ingester

EXPOSE 15694

ENTRYPOINT ["/promtail-ingester"]
