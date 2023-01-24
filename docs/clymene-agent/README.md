# Clymene-agent Getting Start
The Clymene-agent is service that collects time series data(does not use disks)  

## How to create a scrape target setting yaml
1. Config file Option  
```bash
clymene-agent --config.file=/etc/clymene/clymene.yml
```
2. How to write yaml - [See Clymene Configuration for more information](https://clymene-project.github.io/docs/service-discovery/)  
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

3. config file reload
```yaml
curl -XPOST http://clymene-agent:15692/api/reload

# check clymene-agent logs
  {"level":"info","ts":1643299385.1000407,"caller":"config/config.go:451","msg":"Loading configuration file","filename":"clymene_agent.yml"}
  {"level":"info","ts":1643299385.1012235,"caller":"config/config.go:468","msg":"Completed loading of configuration file","filename":"clymene_agent.yml"}
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
  clymene-agent:
    image: bourbonkk/clymene-agent:latest
    ports:
      - "15691:15691"
      - "15692:15692"
    environment:
      - STORAGE_TYPE=kafka
    volumes:
      - ./config/clymene_scrape_config.yml:/etc/clymene/clymene.yml
    command:
      - --log-level=debug
      - --kafka.producer.brokers=[KAFKA-IP]:9092
     # If you want to use druid, use the option below
     # - --kafka.producer.encoding=json    
     # - --kafka.producer.flatten-for-druid
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
          image: bourbonkk/clymene-agent:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 15691
            - containerPort: 15692
          args:
            - --config.file=/etc/clymene/clymene.yml
            - --kafka.producer.brokers=clymene-kafka-broker:9092
            # If you want to use druid, use the option below
            # - --kafka.producer.encoding=json    
            # - --kafka.producer.flatten-for-druid
            - --log-level=info
          env:
            - name: STORAGE_TYPE
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