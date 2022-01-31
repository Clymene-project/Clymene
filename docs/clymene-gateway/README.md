# Clymene Gateway Getting Start

The Clymene Gateway is an optional service that can receive metric data from the agent through gRPC, HTTP communication.  

1. gRPC, HTTP Service
2. Time-series data insert to Database(ElasticSearch, Prometheus, ETC) (Optional) 

## How to setting gateway
```
--admin.http.host-ports string              The host:ports (e.g. 127.0.0.1:15690 or :15690) for the admin server, including health check, /metrics, etc. (default ":15690")
--gateway.grpc-server.host-port string      The host:port (e.g. 127.0.0.1:15610 or :15610) of the gateway's GRPC server (default ":15610")
--gateway.grpc.tls.cert string              Path to a TLS Certificate file, used to identify this server to clients
--gateway.grpc.tls.client-ca string         Path to a TLS CA (Certification Authority) file used to verify certificates presented by clients (if unset, all clients are permitted)
--gateway.grpc.tls.enabled                  Enable TLS on the server
--gateway.grpc.tls.key string               Path to a TLS Private Key file, used to identify this server to clients
--gateway.http-server.host-port string      The host:port (e.g. 127.0.0.1:15610 or :15611) of the gateway's HTTP server (default ":15611")
--gateway.http.tls.cert string              Path to a TLS Certificate file, used to identify this server to clients
--gateway.http.tls.client-ca string         Path to a TLS CA (Certification Authority) file used to verify certificates presented by clients (if unset, all clients are permitted)
--gateway.http.tls.enabled                  Enable TLS on the server
--gateway.http.tls.key string               Path to a TLS Private Key file, used to identify this server to clients
--log-level string                          Minimal allowed log Level. For more levels see https://github.com/uber-go/zap (default "info")
--metrics-backend string                    Defines which metrics backend to use for metrics reporting: expvar, prometheus, none (default "prometheus")
--metrics-http-route string                 Defines the route of HTTP endpoint for metrics backends that support scraping (default "/metrics")
```

## How to set up the Storage Type
#### 1. Setting environmental variables

ElasticSearch
```
TS_STORAGE_TYPE=elasticsearch
```
Kafka
```
TS_STORAGE_TYPE=kafka
```
prometheus
```
TS_STORAGE_TYPE=prometheus
```
cortex
```
TS_STORAGE_TYPE=cortex
```
opentsdb
```
TS_STORAGE_TYPE=opentsdb
```
influxdb
```
TS_STORAGE_TYPE=influxdb
```
tdengine
```
TS_STORAGE_TYPE=tdengine
```
Several
```
TS_STORAGE_TYPE=elasticsearch,prometheus  # composite write
```

#### 2. Option description by storage type
- [Kafka option](./kafka/kafka-option.md)
- [ElasticSearch option](./elasticsearch/es-option.md)
- [Prometheus option](prometheus/prometheus-option.md)
- [cortex option](./cortex/cortex-option.md)
- [Opentsdb option](./opentsdb/opentsdb-option.md)
- [influxdb option](./influxdb/influxdb-options.md)
- [TDengine options](./tdengine/tdengine-options.md)

### Docker-compose Example
```yaml
version: '2'
services:
  clymene-ingester:
    image: bourbonkk/clymene-gateway:latest
    ports:
      - "15610:15610"
    environment:
      - TS_STORAGE_TYPE=elasticsearch
    command:
      - --log-level=debug
      - --es.server-urls=http://[ELASTICSEARCH-IP]:9200
```

### k8s Example
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: clymene-gateway
  namespace: clymene
  labels:
    app: clymene-gateway
spec:
  selector:
    matchLabels:
      app: clymene-gateway
  replicas: 1
  template:
    metadata:
      labels:
        app: clymene-gateway
    spec:
      containers:
        - name: clymene-gateway
          image: bourbonkk/clymene-gateway:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 15610
          args:
            - --es.server-urls=http://[ELASTICSEARCH-IP]:9200
            - --log-level=info
          env:
            - name: TS_STORAGE_TYPE
              value: elasticsearch
      securityContext:
        runAsUser: 1000
```