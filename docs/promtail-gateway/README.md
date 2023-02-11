# Promtail Gateway Getting Start

The Promtail Gateway is an optional service that can receive logs data from the agent through gRPC, HTTP communication.  

1. gRPC, HTTP Service
2. Logs data insert to Database(ElasticSearch, Loki, ETC) (Optional) 

## How to setting gateway
```
--admin.http.host-ports string              The host:ports (e.g. 127.0.0.1:15690 or :15690) for the admin server, including health check, /metrics, etc. (default ":15690")
--gateway.grpc-server.host-port string       The host:port (e.g. 127.0.0.1:15610 or :15610) of the gateway's GRPC server (default ":15610")
--gateway.grpc-server.tls.cert string        Path to a TLS Certificate file, used to identify this server to clients
--gateway.grpc-server.tls.client-ca string   Path to a TLS CA (Certification Authority) file used to verify certificates presented by clients (if unset, all clients are permitted)
--gateway.grpc-server.tls.enabled            Enable TLS on the server
--gateway.grpc-server.tls.key string         Path to a TLS Private Key file, used to identify this server to clients
--gateway.http-server.host-port string       The host:port (e.g. 127.0.0.1:15610 or :15611) of the gateway's HTTP server (default ":15611")
--gateway.http-server.tls.cert string        Path to a TLS Certificate file, used to identify this server to clients
--gateway.http-server.tls.client-ca string   Path to a TLS CA (Certification Authority) file used to verify certificates presented by clients (if unset, all clients are permitted)
--gateway.http-server.tls.enabled            Enable TLS on the server
--gateway.http-server.tls.key string         Path to a TLS Private Key file, used to identify this server to clients--log-level string                          Minimal allowed log Level. For more levels see https://github.com/uber-go/zap (default "info")
--metrics-backend string                    Defines which metrics backend to use for metrics reporting: expvar, prometheus, none (default "prometheus")
--metrics-http-route string                 Defines the route of HTTP endpoint for metrics backends that support scraping (default "/metrics")
```

## How to set up the Storage Type
#### 1. Setting environmental variables


ElasticSearch
```
STORAGE_TYPE=elasticsearch
```
Loki
```
STORAGE_TYPE=loki
```

Promtail-gateway
```
STORAGE_TYPE=gateway
```

Kafka
```
STORAGE_TYPE=kafka
```

#### 2. Option description by storage type

- [ElasticSearch option](https://clymene-project.github.io/docs/database-options/elasticsearch/)
- [Loki option](https://clymene-project.github.io/docs/database-options/loki/)
- [Kafka option](https://clymene-project.github.io/docs/database-options/kafka/)

### Docker-compose Example
```yaml
version: '2'
services:
  clymene-ingester:
    image: bourbonkk/promtail-gateway:latest
    ports:
      - "15610:15610"
    environment:
      - STORAGE_TYPE=elasticsearch
    command:
      - --log-level=debug
      - --es.server-urls=http://[ELASTICSEARCH-IP]:9200
```

### k8s Example
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: promtail-gateway
  namespace: clymene
  labels:
    app: promtail-gateway
spec:
  selector:
    matchLabels:
      app: promtail-gateway
  replicas: 1
  template:
    metadata:
      labels:
        app: promtail-gateway
    spec:
      containers:
        - name: promtail-gateway
          image: bourbonkk/promtail-gateway:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 15610
          args:
            - --es.server-urls=http://[ELASTICSEARCH-IP]:9200
            - --log-level=info
          env:
            - name: STORAGE_TYPE
              value: elasticsearch
      securityContext:
        runAsUser: 1000
```