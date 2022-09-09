# Clymene  



[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/5491/badge)](https://bestpractices.coreinfrastructure.org/projects/5491) ![CodeQL](https://github.com/clymene-project/clymene/workflows/CodeQL/badge.svg) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) ![GitHub release (latest by date)](https://img.shields.io/github/v/release/clymene-project/clymene) [![Go Reference](https://pkg.go.dev/badge/github.com/Clymene-project/Clymene.svg)](https://pkg.go.dev/github.com/Clymene-project/Clymene)  
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white) ![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white) ![ElasticSearch](https://img.shields.io/badge/-ElasticSearch-005571?style=for-the-badge&logo=elasticsearch) <img src="https://img.shields.io/badge/influxdb-%2322ADF6.svg?&style=for-the-badge&logo=influxdb&logoColor=white"/>  <img src="https://img.shields.io/badge/prometheus-%23E6522C.svg?&style=for-the-badge&logo=prometheus&logoColor=white" />   <img src="https://img.shields.io/badge/OpenTSDB-green?style=for-the-badge"> <img src="https://img.shields.io/badge/cortex-blue?style=for-the-badge">  <img src="https://img.shields.io/badge/tdengine-gray?style=for-the-badge"> <img src="https://img.shields.io/badge/Loki-yellow?style=for-the-badge">


<img align="right" width="40%" height="40%" src="https://user-images.githubusercontent.com/25188468/148681479-3ddf237c-6e5d-49a1-a517-8b3bfa92f54e.png" alt="clymene_logo">

The Clymene is a time-series data and logs collection platform for distributed systems inspired
by [Prometheus](https://prometheus.io)
and [Jaeger](https://www.jaegertracing.io). Time-series data and logs from various environments can be collected and stored in
different types of databases. It can be configured in a variety of architectures. Choose the storage that users can use easily and build a monitoring system with dashboards that users can use easily. Also, Clymene Agent's time series data collection uses fewer resources than Prometheus' remote_write. 
Various time series data analysis functions will be added.


# Architecture  

The Clymene project is a platform for collecting time series data and log data. So there are two pipelines: timeseries and logs. All agents can collect data and store it directly in the database,
Or, depending on the situation, you can use the Clymene components to configure a variety of architectures. It can be configured with an architecture that can be used with a lot of traffic. 
The Clymene Project provides various storage types, so choose a storage type you are familiar with. And build great monitoring systems using familiar dashboards.  
  
<img src="https://user-images.githubusercontent.com/25188468/165649203-7382735f-0c1f-4388-93f3-65b57b06d7f4.png" width="90%" height="90%" alt="architecture_v2.1.0">

# Components

## Pipeline 1. TimeSeries data collection   

### Clymene Agent([Getting Started](https://github.com/Clymene-project/Clymene/blob/main/docs/clymene-agent/README.md))

```dockerhub : bourbonkk/clymene-agent:v2.x.x```  
```redhatQuay: quay.io/clymene/clymene-agent:v2.x.x```  
![Docker Pulls](https://img.shields.io/docker/pulls/bourbonkk/clymene-agent.svg?maxAge=86400) [![Docker Repository on Quay](https://quay.io/repository/clymene/clymene-agent/status "Docker Repository on Quay")](https://quay.io/repository/clymene/clymene-agent)   
The Clymene-agent is service that collects time series data(does not use disks)

1. Service Discovery
    - [Prometheus's Service Discovery](https://docs.sysdig.com/en/docs/sysdig-monitor/integrations-for-sysdig-monitor/collect-prometheus-metrics/enable-prometheus-native-service-discovery/)
      feature finds Metric collection endpoints.
2. scrape time series data
3. Time-series data transfer to gateway(gRPC) (Optional)
4. Time-series data transfer to kafka (Optional)
5. Time-series data insert to Database([supported DB](https://github.com/Clymene-project/Clymene/blob/main/docs/clymene-agent/README.md#Option-description-by-storage-type)) (Optional)

### Clymene Ingester(Optional) ([Getting Started](https://github.com/Clymene-project/Clymene/blob/main/docs/clymene-ingester/README.md))

```dockerhub : bourbonkk/clymene-ingester:v2.x.x```  
```redhatQuay: quay.io/clymene/clymene-ingester:v2.x.x```  
![Docker Pulls](https://img.shields.io/docker/pulls/bourbonkk/clymene-ingester.svg?maxAge=86400) [![Docker Repository on Quay](https://quay.io/repository/clymene/clymene-ingester/status "Docker Repository on Quay")](https://quay.io/repository/clymene/clymene-ingester)  
The Clymene ingester is an optional service responsible for insert time series data loaded on kafka into the database.

1. Kafka message consume
2. Time-series data insert to Database([supported DB](https://github.com/Clymene-project/Clymene/blob/main/docs/clymene-ingester/README.md#Option-description-by-storage-type)) (Optional)

### Clymene Gateway(Optional) ([Getting Started](https://github.com/Clymene-project/Clymene/blob/main/docs/clymene-gateway/README.md))

```dockerhub : bourbonkk/clymene-gateway:v2.x.x```  
```redhatQuay: quay.io/clymene/clymene-gateway:v2.x.x```  
![Docker Pulls](https://img.shields.io/docker/pulls/bourbonkk/clymene-gateway.svg?maxAge=86400) [![Docker Repository on Quay](https://quay.io/repository/clymene/clymene-gateway/status "Docker Repository on Quay")](https://quay.io/repository/clymene/clymene-gateway)  
The Clymene Gateway is an optional service that can receive metric data from the another component through gRPC or HTTP
communication.

1. gRPC, HTTP Service
2. Time-series data insert to Database([supported DB](https://github.com/Clymene-project/Clymene/blob/main/docs/clymene-gateway/README.md#Option-description-by-storage-type)) (Optional)


## Pipeline 2. Logs collection  

### Clymene Promtail([Getting Started](https://github.com/Clymene-project/Clymene/blob/main/docs/clymene-promtail/README.md))

```dockerhub : bourbonkk/clymene-promtail:v2.x.x```  
```redhatQuay: quay.io/clymene/clymene-promtail:v2.x.x```  
![Docker Pulls](https://img.shields.io/docker/pulls/bourbonkk/clymene-promtail.svg?maxAge=86400) [![Docker Repository on Quay](https://quay.io/repository/clymene/clymene-promtail/status "Docker Repository on Quay")](https://quay.io/repository/clymene/clymene-agent)   
The Clymene-promtail customized loki's log collection agent for the Clymene project.

1. [Service Discovery](https://clymene-project.github.io/docs/service-discovery/promtail-config/)
2. log collection
3. log data transfer to gateway(gRPC or HTTP)
4. log data transfer to kafka
5. log data insert to Database([supported DB](https://github.com/Clymene-project/Clymene/blob/main/docs/clymene-promtail/README.md#Option-description-by-storage-type)) (Optional)

### Promtail Ingester(Optional) ([Getting Started](https://github.com/Clymene-project/Clymene/blob/main/docs/promtail-ingester/README.md))

```dockerhub : bourbonkk/promtail-ingester:v2.x.x```  
```redhatQuay: quay.io/clymene/promtail-ingester:v2.x.x```  
![Docker Pulls](https://img.shields.io/docker/pulls/bourbonkk/promtail-ingester.svg?maxAge=86400) [![Docker Repository on Quay](https://quay.io/repository/clymene/promtail-ingester/status "Docker Repository on Quay")](https://quay.io/repository/clymene/promtail-ingester)  
Promtail ingester is an optional service responsible for insert log data loaded on kafka into the database.

1. Kafka message consume
2. Time-series data insert to Database([supported DB](https://github.com/Clymene-project/Clymene/blob/main/docs/promtail-ingester/README.md#Option-description-by-storage-type)) (Optional)

### Promtail Gateway(Optional) ([Getting Started](https://github.com/Clymene-project/Clymene/blob/main/docs/promtail-gateway/README.md))

```dockerhub : bourbonkk/promtail-gateway:v2.x.x```  
```redhatQuay: quay.io/clymene/promtail-gateway:v2.x.x```  
![Docker Pulls](https://img.shields.io/docker/pulls/bourbonkk/promtail-gateway.svg?maxAge=86400) [![Docker Repository on Quay](https://quay.io/repository/clymene/promtail-gateway/status "Docker Repository on Quay")](https://quay.io/repository/clymene/promtail-gateway)  
The Promtail Gateway is an optional service that can receive log data from the another component through gRPC or HTTP
communication.

1. gRPC, HTTP Service
2. Time-series data insert to Database([supported DB](https://github.com/Clymene-project/Clymene/blob/main/docs/promtail-gateway/README.md#Option-description-by-storage-type)) (Optional)


## How to work it on k8s

For example, see
the [k8s example](https://clymene-project.github.io/docs/tutorials/k8s-example/)

## How to work it on docker-compose

For example, see
the [docker-compose example](https://clymene-project.github.io/docs/tutorials/docker-compose/)

## Roadmap

1. AI/ML platform for clymene(https://github.com/Clymene-project/clymene-analyzer)
2. push-type agent(node-exporter, cadvisor, process-exporter, etc)
3. Add new pipeline(trace or Something - [discussions](https://github.com/Clymene-project/Clymene/discussions))
4. eBPF trace agent(for v3.0.0 release)

## Get in Touch

Have questions, suggestions, bug reports? Reach the project community via these channels:

* GitHub [issues](https://github.com/clymene-project/clymene/issues)
  and [discussions](https://github.com/Clymene-project/Clymene/discussions)


## Further reading
To learn more about Clymene, consult the following talks and articles.

### Recent talks and articles
- Jan 2022 blog post "[Efficient time series data collection and management plan in a distributed environment](https://medium.com/@clymeneallen/efficient-time-series-data-collection-and-management-plan-in-a-distributed-environment-5c4a62d45e36)"
- Aug 2022 blog post "[Best Practices, Monitoring System for Multi-K8s Cluster Environments Using Open Source](https://medium.com/@clymeneallen/best-practices-monitoring-system-for-multi-k8s-cluster-environments-using-open-source-d85544052f37)"
- Aug 2022 blog post "[Best practice, k8s Node/POD resource usage and log monitoring system for multi-k8s cluster environment using Open source](https://medium.com/@clymeneallen/best-practice-k8s-node-pod-resource-usage-and-log-monitoring-system-for-multi-k8s-cluster-3f2e5c2c401a)"

## License

Apache License 2.0, see [LICENSE](https://github.com/Clymene-project/Clymene/blob/main/LICENSE).

# Code Style

#### 1. https://github.com/golang/go/wiki/CodeReviewComments

#### 2. https://peter.bourgon.org/go-in-production/#formatting-and-style    

# supported
The Clymene project is being developed with the support of Open Source development license of [Jetbrains'](https://www.jetbrains.com).
