FROM alpine:latest as certs
RUN apk add --update --no-cache ca-certificates

FROM alpine

WORKDIR /

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ADD ./out/clymene-gateway /clymene-gateway

EXPOSE 15610
EXPOSE 15690

ENTRYPOINT ["/clymene-gateway"]
