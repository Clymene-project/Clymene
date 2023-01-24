# Clymene Ingester Getting Start

The Clymene ingester is an optional service responsible for insert time series data loaded on kafka into the database  

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

#### 1. Setting environmental variables

ElasticSearch

```
STORAGE_TYPE=elasticsearch
```

Kafka

```
STORAGE_TYPE=kafka
```

prometheus

```
STORAGE_TYPE=prometheus
```

cortex

```
STORAGE_TYPE=cortex
```

gateway

```
STORAGE_TYPE=gateway
```

opentsdb

```
STORAGE_TYPE=opentsdb
```

influxdb

```
STORAGE_TYPE=influxdb
```
tdengine
```
STORAGE_TYPE=tdengine
```
druid
```
# env setting
STORAGE_TYPE=kafka
# arg
--kafka.producer.encoding=json
--kafka.producer.flatten-for-druid
```

Several

```
STORAGE_TYPE=elasticsearch,prometheus  # composite write
```

#### 2. Option description by storage type

- [Kafka option](https://clymene-project.github.io/docs/database-options/kafka/)
- [ElasticSearch option](https://clymene-project.github.io/docs/database-options/elasticsearch/)
- [Prometheus option](https://clymene-project.github.io/docs/database-options/prometheus/)
- [cortex option](https://clymene-project.github.io/docs/database-options/cortex/)
- [gateway option](https://clymene-project.github.io/docs/database-options/gateway/)
- [Opentsdb option](https://clymene-project.github.io/docs/database-options/opentsdb/)
- [influxdb option](https://clymene-project.github.io/docs/database-options/influxdb/)
- [TDengine options](https://clymene-project.github.io/docs/database-options/tdengine/)

### Docker-compose Example

```yaml
version: '2'
services:
  clymene-ingester:
    image: bourbonkk/clymene-ingester:latest
    ports:
      - "15694:15694"
    environment:
      #      - STORAGE_TYPE=elasticsearch,prometheus   # use composite writer
      - STORAGE_TYPE=elasticsearch
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
          image: bourbonkk/clymene-ingester:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 15694
          args:
            - --prometheus.remote.url=http://prometheus:9090/api/v1/write
            - --log-level=info
            - --kafka.consumer.brokers=clymene-kafka-broker:9092
          env:
            - name: STORAGE_TYPE
              value: prometheus
      securityContext:
        runAsUser: 1000
```