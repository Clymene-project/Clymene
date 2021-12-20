FROM alpine:latest as certs
RUN apk add --update --no-cache ca-certificates

FROM alpine

WORKDIR /

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ADD ./out/clymene-ingester /clymene-ingester

EXPOSE 15694

ENTRYPOINT ["/clymene-ingester"]
