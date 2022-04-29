# Clymene-promtail Getting Start
The Clymene-promtail customized loki's log collection agent for the Clymene project. 

## How to create a target setting yaml
1. Config file Option  
```bash
clymene-promtail --config.file=/etc/promtail/config.yml
```
2. How to write yaml - [See Clymene-promtail Configuration for more information](https://clymene-project.github.io/docs/service-discovery/promtail-config/)  
```yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

scrape_configs:
  - job_name: system
    static_configs:
      - targets:
          - localhost
        labels:
          job: varlogs
          __path__: /var/log/*log

  - job_name: kafka-sasl-plain
    kafka:
      use_incoming_timestamp: false
      brokers:
        - localhost:29092
      authentication:
        type: sasl
        sasl_config:
          mechanism: PLAIN
          user: kafkaadmin
          password: kafkaadmin-pass
          use_tls: true
          ca_file: ../../../tools/kafka/secrets/promtail-kafka-ca.pem
          insecure_skip_verify: true
      group_id: kafka_group
      topics:
        - foo
        - ^promtail.*
      labels:
        job: kafka-sasl-plain
```

3. config file reload(TODO)  
will be supported

4. other options
```
  --admin.http.host-ports string                The host:ports (e.g. 127.0.0.1:15698 or :15698) for the admin server, including health check, /metrics, etc. (default ":15698")
  --clymene-promtail.batch-size-bytes int       Maximum batch size to accrue before sending. (default 1048576)
  --clymene-promtail.batch-wait duration        Maximum wait period before sending batch. (default 1s)
  --clymene-promtail.dry-run                    Start Promtail but print entries instead of sending them to Loki.
  --clymene-promtail.external-labels string     list of external labels to add to each log (e.g: --loki.client.external-labels=lb1=v1,lb2=v2)
  --clymene-promtail.inspect                    Allows for detailed inspection of pipeline stages
  --clymene-promtail.max-backoff duration       Maximum backoff time between retries. (default 5m0s)
  --clymene-promtail.max-retries int            Maximum number of retires when sending batches. (default 10)
  --clymene-promtail.min-backoff duration       Initial backoff time between retries. (default 500ms)
  --clymene-promtail.print-config-stderr        Dump the entire Loki config object to stderr
  --clymene-promtail.stream-lag-labels string   Comma-separated list of labels to use when calculating stream lag (default "filename")
  --clymene-promtail.tenant-id string           Tenant ID to use when pushing logs to Loki.
  --config.file string                          yaml file to load (default "/etc/promtail/config.yml")
  --log-level string                            Minimal allowed log Level. For more levels see https://github.com/uber-go/zap (default "info")
  --metrics-backend string                      Defines which metrics backend to use for metrics reporting: expvar, prometheus, none (default "prometheus")
  --metrics-http-route string                   Defines the route of HTTP endpoint for metrics backends that support scraping (default "/metrics")
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

- [ElasticSearch option](./elasticsearch/es-option.md)
- [Loki option](./loki/loki-option.md)
- [Kafka option](./kafka/kafka-option.md)
- [Promtail-gateway](./gateway/gateway-option.md)

### Docker-compose Example
```yaml
version: '2'
services:
  clymene-promtail:
    image: bourbonkk/clymene-promtail:latest
    ports:
      - "15698:15698"
      - "9080:9080"
    environment:
      - STORAGE_TYPE=elasticsearch
    volumes:
      - ./config/promtail-local-config.yaml.yml:/etc/promtail/config.yml
    command:
      - --log-level=debug
      - --es.server-urls=[ELASTICSEARCH-IP]:9200
```

### k8s Example
```yaml
--- # Daemonset.yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: clymene-promtail
  namespace: clymene
  labels:
    app: clymene-promtail
spec:
  selector:
    matchLabels:
      app: clymene-promtail
  template:
    metadata:
      labels:
        app: clymene-promtail
    spec:
      containers:
        - name: clymene-promtail
          image: bourbonkk/clymene-promtail:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 15698
            - containerPort: 9080
          args:
            - --config.file=/etc/promtail/config.yml
            - --es.server-urls=http://elasticsearch:9200
            - --log-level=info
          env:
            - name: STORAGE_TYPE
              value: elasticsearch
          volumeMounts:
            - mountPath: /etc/promtail/
              name: config-volume
      volumes:
        - name: config-volume
          configMap:
            name: promtail-config
      securityContext:
        runAsUser: 1000
--- # Service.yaml
apiVersion: v1
kind: Service
metadata:
  name: clymene-promtail
  namespace: clymene
  labels:
    app: clymene-promtail
spec:
  ports:
    - name: admin
      port: 15698
      targetPort: 15698
    - name: server
      port: 9080
      targetPort: 9080
  selector:
    app: clymene-promtail
--- # configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: promtail-config
  namespace: clymene
data:
  config.yml: |
    server:
      http_listen_port: 9080
      grpc_listen_port: 0
    
    positions:
      filename: /tmp/positions.yaml
    
    scrape_configs:
      - job_name: system
        static_configs:
          - targets:
              - localhost
            labels:
              job: varlogs
              __path__: /var/log/*log
    
      - job_name: kafka-sasl-plain
        kafka:
          use_incoming_timestamp: false
          brokers:
            - localhost:29092
          authentication:
            type: sasl
            sasl_config:
              mechanism: PLAIN
              user: kafkaadmin
              password: kafkaadmin-pass
              use_tls: true
              ca_file: ../../../tools/kafka/secrets/promtail-kafka-ca.pem
              insecure_skip_verify: true
          group_id: kafka_group
          topics:
            - foo
            - ^promtail.*
          labels:
            job: kafka-sasl-plain

--- # Clusterrole.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: promtail-clusterrole
rules:
  - apiGroups: [""]
    resources:
      - nodes
      - services
      - pods
    verbs:
      - get
      - watch
      - list

--- #ServiceAccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: promtail-serviceaccount

--- #Rolebinding.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: promtail-clusterrolebinding
subjects:
  - kind: ServiceAccount
    name: promtail-serviceaccount
    namespace: clymene
roleRef:
  kind: ClusterRole
  name: promtail-clusterrole
  apiGroup: rbac.authorization.k8s.io
```