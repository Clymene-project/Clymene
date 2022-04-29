# Promtail Ingester Getting Start

The Promtail ingester is an optional service responsible for insert logs data loaded on kafka into the database  

1. Kafka message consume
2. logs data insert to Database(ElasticSearch, Loki, ETC) (Optional)

## How to setting kafka consumer

```
--admin.http.host-ports string                  The host:ports (e.g. 127.0.0.1:15694 or :15694) for the admin server, including health check, /metrics, etc. (default ":15694")
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
--kafka.consumer.promtail.topic string          The name of the promtail kafka topic to consume from (default "clymene-logs")
--kafka.consumer.protocol-version string        Kafka protocol version - must be supported by kafka server
--kafka.consumer.tls.ca string                  Path to a TLS CA (Certification Authority) file used to verify the remote server(s) (by default will use the system truststore)
--kafka.consumer.tls.cert string                Path to a TLS Certificate file, used to identify this process to the remote server(s)
--kafka.consumer.tls.enabled                    Enable TLS when talking to the remote server(s)
--kafka.consumer.tls.key string                 Path to a TLS Private Key file, used to identify this process to the remote server(s)
--kafka.consumer.tls.server-name string         Override the TLS server name we expect in the certificate of the remote server(s)
--kafka.consumer.tls.skip-host-verify           (insecure) Skip server's certificate chain and host name verification
--log-level string                              Minimal allowed log Level. For more levels see https://github.com/uber-go/zap (default "info")
--metrics-backend string                        Defines which metrics backend to use for metrics reporting: expvar, prometheus, none (default "prometheus")
--metrics-http-route string                     Defines the route of HTTP endpoint for metrics backends that support scraping (default "/metrics")
--promtail-ingester.deadlockInterval duration   Interval to check for deadlocks. If no messages gets processed in given time, clymene-ingester app will exit. Value of 0 disables deadlock check. (default 0s)
--promtail-ingester.parallelism string          The number of messages to process in parallel (default "1000")
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

- [ElasticSearch option](../clymene-promtail/elasticsearch/es-option.md)
- [Loki option](../clymene-promtail/loki/loki-option.md)
- [Kafka option](../clymene-promtail/kafka/kafka-option.md)
- [Promtail-gateway](../clymene-promtail/gateway/gateway-option.md)


### Docker-compose Example

```yaml
version: '2'
services:
  promtail-ingester:
    image: bourbonkk/promtail-ingester:latest
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
  name: promtail-ingester
  namespace: clymene
  labels:
    app: promtail-ingester
spec:
  selector:
    matchLabels:
      app: promtail-ingester
  replicas: 1
  template:
    metadata:
      labels:
        app: promtail-ingester
    spec:
      containers:
        - name: promtail-ingester
          image: bourbonkk/promtail-ingester:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 15694
          args:
            - --es.server-urls=http://[ELASTICSEARCH-IP]:9200
            - --log-level=info
            - --kafka.consumer.brokers=clymene-kafka-broker:9092
          env:
            - name: STORAGE_TYPE
              value: elasticsearch
      securityContext:
        runAsUser: 1000
```