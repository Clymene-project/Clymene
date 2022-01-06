# Clymene Ingester Getting Start

The Clymene ingester is an optional service responsible for insert time series data loaded on kafka into the database  
Clymene ingester는 카프카에 적재된 시계열 데이터를 데이터베이스에 삽입하는 선택적 서비스이다

1. Kafka message consume
2. Time-series data insert to Database(ElasticSearch, Prometheus, ETC) (Optional)

## How to setting kafka consumer

```
--kafka.consumer.authentication string          Authentication type used to authenticate with kafka cluster. e.g. none, kerberos, tls, plaintext (default "none")
--kafka.consumer.brokers string                 The comma-separated list of kafka brokers. i.e. '127.0.0.1:9092,0.0.0:1234' (default "127.0.0.1:9092")
--kafka.consumer.client-id string               The Consumer Client ID that clymene-ingester will use (default "clymene")
--kafka.consumer.encoding string                The encoding of metrics ("json", "protobuf") consumed from kafka (default "protobuf")
--kafka.consumer.group-id string                The Consumer Group that clymene-ingester will be consuming on behalf of (default "clymene")
--kafka.consumer.kerberos.config-file string    Path to Kerberos configuration. i.e /etc/krb5.conf (default "/etc/krb5.conf")
--kafka.consumer.kerberos.keytab-file string    Path to keytab file. i.e /etc/security/kafka.keytab (default "/etc/security/kafka.keytab")
--kafka.consumer.kerberos.password string       The Kerberos password used for authenticate with KDC
--kafka.consumer.kerberos.realm string          Kerberos realm
--kafka.consumer.kerberos.service-name string   Kerberos service name (default "kafka")
--kafka.consumer.kerberos.use-keytab            Use of keytab instead of password, if this is true, keytab file will be used instead of password
--kafka.consumer.kerberos.username string       The Kerberos username used for authenticate with KDC
--kafka.consumer.plaintext.mechanism string     The plaintext Mechanism for SASL/PLAIN authentication, e.g. 'SCRAM-SHA-256' or 'SCRAM-SHA-512' or 'PLAIN' (default "PLAIN")
--kafka.consumer.plaintext.password string      The plaintext Password for SASL/PLAIN authentication
--kafka.consumer.plaintext.username string      The plaintext Username for SASL/PLAIN authentication
--kafka.consumer.protocol-version string        Kafka protocol version - must be supported by kafka server
--kafka.consumer.tls.ca string                  Path to a TLS CA (Certification Authority) file used to verify the remote server(s) (by default will use the system truststore)
--kafka.consumer.tls.cert string                Path to a TLS Certificate file, used to identify this process to the remote server(s)
--kafka.consumer.tls.enabled                    Enable TLS when talking to the remote server(s)
--kafka.consumer.tls.key string                 Path to a TLS Private Key file, used to identify this process to the remote server(s)
--kafka.consumer.tls.server-name string         Override the TLS server name we expect in the certificate of the remote server(s)
--kafka.consumer.tls.skip-host-verify           (insecure) Skip server's certificate chain and host name verification
--kafka.consumer.topic string                   The name of the kafka topic to consume from (default "clymene")
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

gateway

```
TS_STORAGE_TYPE=gateway
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
- [gateway option](./gateway/gataway-option.md)
- [Opentsdb option](./opentsdb/opentsdb-option.md)
- [influxdb option](./influxdb/influxdb-options.md)
-

### Docker-compose Example

```yaml
version: '2'
services:
  clymene-ingester:
    image: bourbonkk/clymene-ingester:main
    ports:
      - "15694:15694"
    environment:
      #      - TS_STORAGE_TYPE=elasticsearch,prometheus   # use composite writer
      - TS_STORAGE_TYPE=elasticsearch
    command:
      - --log-level=debug
      - --kafka.consumer.brokers=[KAFKA-IP]:9092
      - --es.server-urls=http://[ELASTICSEARCH-IP]:9200
#      - --prometheus.remote.url=http://prometheus:9090/api/v1/write
```

### k8s Example

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: clymene-ingester
  namespace: clymene
  labels:
    app: clymene-ingester
spec:
  selector:
    matchLabels:
      app: clymene-ingester
  replicas: 1
  template:
    metadata:
      labels:
        app: clymene-ingester
    spec:
      containers:
        - name: clymene-ingester
          image: bourbonkk/clymene-ingester:main
          imagePullPolicy: Always
          ports:
            - containerPort: 15694
          args:
            - --prometheus.remote.url=http://prometheus:9090/api/v1/write
            - --log-level=info
            - --kafka.consumer.brokers=clymene-kafka-broker:9092
          env:
            - name: TS_STORAGE_TYPE
              value: prometheus
      securityContext:
        runAsUser: 1000
```