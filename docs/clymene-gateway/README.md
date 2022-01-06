# Clymene Gateway Getting Start

The Clymene Gateway is an optional service that can receive metric data from the agent through gRPC communication.  
Clymene Gateway는 Clymene Agent로 부터 gRPC 통신으로 metric 데이터를 전달 받을 수 있는 선택적 서비스입니다.

1. gRPC Service
2. Time-series data insert to Database(ElasticSearch, Prometheus, ETC) (Optional) 

## How to setting gRPC server  
```
--gateway.grpc-server.host-port string      The host:port (e.g. 127.0.0.1:15610 or :15610) of the gateway's GRPC server (default ":15610")
--gateway.grpc.tls.cert string              Path to a TLS Certificate file, used to identify this server to clients
--gateway.grpc.tls.client-ca string         Path to a TLS CA (Certification Authority) file used to verify certificates presented by clients (if unset, all clients are permitted)
--gateway.grpc.tls.enabled                  Enable TLS on the server
--gateway.grpc.tls.key string               Path to a TLS Private Key file, used to identify this server to clients
```

## How to set up the Storage Type
1. Setting environmental variables

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
Several
```
TS_STORAGE_TYPE=elasticsearch,prometheus  # composite write - 여러 DB에 동시 write
```

2. Option description by storage type
- [Kafka option](./kafka/kafka-option.md)
- [ElasticSearch option](./elasticsearch/es-option.md)
- [Prometheus option](prometheus/prometheus-option.md)
- [cortex option](./cortex/cortex-option.md)
- [Opentsdb option](./opentsdb/opentsdb-option.md)
- [influxdb option](./influxdb/influxdb-options.md)

### Docker-compose Example
```yaml
version: '2'
services:
  clymene-ingester:
    image: bourbonkk/clymene-gateway:v1.2.0
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
          image: bourbonkk/clymene-gateway:v1.2.0
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