# Clymene-agent Getting Start
The Clymene-agent is service that collects time series data(does not use disks)  
Clymene-agent는 시계열 데이터를 수집하는 서비스입니다(디스크 사용 안함)

## How to create a scrape target setting yaml
1. Config file Option  
```bash
clymene-agent --config.file=/etc/clymene/clymene.yml
```
2. How to write yaml - [See Prometheus Configuration for more information](https://prometheus.io/docs/prometheus/latest/configuration/configuration/)  
```yaml
global:
   scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.  
   evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.

scrape_configs:
- job_name: 'localhost'  
  static_configs:  
    - targets: [ 'localhost:9100' ]   

- job_name: 'kubernetes-kubelet'  
  scheme: https  
  tls_config:  
  ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt  
  bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token  
  kubernetes_sd_configs:  
    - role: node  
      relabel_configs:  
    - action: labelmap  
      regex: __meta_kubernetes_node_label_(.+)  
    - target_label: __address__  
      replacement: kubernetes.default.svc:443  
    - source_labels: [__meta_kubernetes_node_name]  
      regex: (.+)  
      target_label: __metrics_path__  
      replacement: /api/v1/nodes/${1}/proxy/metrics
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

### Docker-compose Example
```yaml
version: '2'
services:
  clymene-agent:
    image: bourbonkk/clymene-agent:main
    ports:
      - "15691:15691"
      - "15692:15692"
    environment:
      - TS_STORAGE_TYPE=kafka
    volumes:
      - ./config/clymene_scrape_config.yml:/etc/clymene/clymene.yml
    command:
      - --log-level=debug
      - --kafka.producer.brokers=[KAFKA-IP]:9092
```

### k8s Example
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: clymene-agent
  namespace: clymene
  labels:
    app: clymene-agent
spec:
  selector:
    matchLabels:
      app: clymene-agent
  replicas: 1
  template:
    metadata:
      labels:
        app: clymene-agent
    spec:
      containers:
        - name: clymene-agent
          image: bourbonkk/clymene-agent:main
          imagePullPolicy: Always
          ports:
            - containerPort: 15691
            - containerPort: 15692
          args:
            - --config.file=/etc/clymene/clymene.yml
            - --kafka.producer.brokers=clymene-kafka-broker:9092
            - --log-level=info
          env:
            - name: TS_STORAGE_TYPE
              value: kafka
#              value: prometheus
          volumeMounts:
            - mountPath: /etc/clymene/
              name: config-volume
      volumes:
        - name: config-volume
          configMap:
            name: clymene-agent-config
      securityContext:
        runAsUser: 1000
---
apiVersion: v1
kind: Service
metadata:
  name: clymene-agent
  namespace: clymene
  labels:
    app: clymene-agent
spec:
  ports:
    - name: metric
      port: 15691
      targetPort: 15691
    - name: admin
      port: 15692
      targetPort: 15692
  selector:
    app: clymene-agent
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: clymene-agent-config
  namespace: clymene
data:
  clymene.yml: |
    global:
      scrape_interval:     15s
    scrape_configs:
      - job_name: 'kubernetes-kubelet'
        scheme: https
        tls_config:
          ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
        kubernetes_sd_configs:
          - role: node
        relabel_configs:
          - action: labelmap
            regex: __meta_kubernetes_node_label_(.+)
          - target_label: __address__
            replacement: kubernetes.default.svc:443
          - source_labels: [__meta_kubernetes_node_name]
            regex: (.+)
            target_label: __metrics_path__
            replacement: /api/v1/nodes/${1}/proxy/metrics
```