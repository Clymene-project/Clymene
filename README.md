# Clymene  



[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/5491/badge)](https://bestpractices.coreinfrastructure.org/projects/5491) ![CodeQL](https://github.com/clymene-project/clymene/workflows/CodeQL/badge.svg) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) ![GitHub release (latest by date)](https://img.shields.io/github/v/release/clymene-project/clymene) [![Go Reference](https://pkg.go.dev/badge/github.com/Clymene-project/Clymene.svg)](https://pkg.go.dev/github.com/Clymene-project/Clymene)  
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white) ![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white) ![ElasticSearch](https://img.shields.io/badge/-ElasticSearch-005571?style=for-the-badge&logo=elasticsearch) <img src="https://img.shields.io/badge/influxdb-%2322ADF6.svg?&style=for-the-badge&logo=influxdb&logoColor=white"/>  <img src="https://img.shields.io/badge/prometheus-%23E6522C.svg?&style=for-the-badge&logo=prometheus&logoColor=white" />   <img src="https://img.shields.io/badge/OpenTSDB-green?style=for-the-badge"> <img src="https://img.shields.io/badge/cortex-blue?style=for-the-badge">    


<img align="right" width="400" height="400" src="https://user-images.githubusercontent.com/25188468/148681479-3ddf237c-6e5d-49a1-a517-8b3bfa92f54e.png" alt="clymene_logo">

The Clymene is a time series data collection platform for distributed systems inspired
by [Prometheus](https://prometheus.io)
and [Jaeger](https://www.jaegertracing.io). Time series data from various environments can be collected and stored in
different types of databases. It can be configured in a variety of architectures. Various time series data analysis
functions will be added.

Clymene는 [Prometheus](https://prometheus.io)와 [Jaeger](https://www.jaegertracing.io)에서 영감을 받은 분산 시스템을 위한 시계열 데이터 수집
플랫폼입니다. 다양한 환경의 시계열 데이터를 수집하고 여러 종류의 데이터베이스에 저장할 수 있습니다. 여러 컴포넌트를 이용해 다양한 아키텍처로 구성할 수 있습니다. 다양한 시계열 데이터 분석 기능이 추가될
예정입니다.

# Architecture

### Use only agent

Agent can store time series data in a variety of DBs and support Composite Writers that can be stored in multiple DBs at
once.

<img src="https://user-images.githubusercontent.com/25188468/148680397-aa9ca8a8-810c-4141-aefb-7e5d8ed87d13.png" width="70%" height="70%" alt="architecture_v1.3.0">

### Including kafka and ingester

Ingester consumes time series data loaded on kafka, can store time series data in various DBs, and supports Composite
Writers that can be stored in multiple DBs at once.

<img src="https://user-images.githubusercontent.com/25188468/148680487-de5b083d-024e-4aaa-aeac-d362bbcec8de.png" width="70%" height="70%" alt="architecture_v1.3.0_ingester">  

### Use gateway

Gateway collects time series data through gRPC communication, stores time series data in various DBs, and supports
Composite Writers that can be stored in multiple DBs at once.

<img src="https://user-images.githubusercontent.com/25188468/148680473-290f098e-fab6-4e87-a958-cdd96ff12a15.png" width="70%" height="70%" alt="architecture_v1.3.0_gateway">

# components

### Clymene Agent([Getting Started](https://github.com/Clymene-project/Clymene/blob/main/docs/clymene-agent/README.md))

```dockerhub : bourbonkk/clymene-agent:v1.x.x```  
```redhatQuay: quay.io/clymene/clymene-agent:v1.x.x```  
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

```dockerhub : bourbonkk/clymene-ingester:v1.x.x```  
```redhatQuay: quay.io/clymene/clymene-ingester:v1.x.x```  
![Docker Pulls](https://img.shields.io/docker/pulls/bourbonkk/clymene-ingester.svg?maxAge=86400) [![Docker Repository on Quay](https://quay.io/repository/clymene/clymene-ingester/status "Docker Repository on Quay")](https://quay.io/repository/clymene/clymene-ingester)  
The Clymene ingester is an optional service responsible for insert time series data loaded on kafka into the database.

1. Kafka message consume
2. Time-series data insert to Database([supported DB](https://github.com/Clymene-project/Clymene/blob/main/docs/clymene-ingester/README.md#Option-description-by-storage-type)) (Optional)

### Clymene Gateway(Optional) ([Getting Started](https://github.com/Clymene-project/Clymene/blob/main/docs/clymene-gateway/README.md))

```dockerhub : bourbonkk/clymene-gateway:v1.x.x```  
```redhatQuay: quay.io/clymene/clymene-gateway:v1.x.x```  
![Docker Pulls](https://img.shields.io/docker/pulls/bourbonkk/clymene-gateway.svg?maxAge=86400) [![Docker Repository on Quay](https://quay.io/repository/clymene/clymene-gateway/status "Docker Repository on Quay")](https://quay.io/repository/clymene/clymene-gateway)  
The Clymene Gateway is an optional service that can receive metric data from the another component through gRPC
communication.

1. gRPC Service
2. Time-series data insert to Database([supported DB](https://github.com/Clymene-project/Clymene/blob/main/docs/clymene-gateway/README.md#Option-description-by-storage-type)) (Optional)

## How to work it on k8s

For example, see
the [clymene_using_gateway.yaml](https://github.com/Clymene-project/Clymene/blob/main/k8s/Clymene_using_gateway.yaml)

## How to work it on docker-compose

For example, see
the [docker-compose-clymene_using_gateway.yml](https://github.com/Clymene-project/Clymene/blob/main/docker-compose/docker-compose-clymene_using_gateway.yml)

## Roadmap

1. Various DB support(ElasticSearch, Prometheus, Cortex, OpenTSDB, influxdb, TDengine, kdb, druid)
2. AI/ML platform for clymene(https://github.com/Clymene-project/clymene-analyzer)
3. Website development for clymene(https://github.com/Clymene-project/documentation)

## Get in Touch

Have questions, suggestions, bug reports? Reach the project community via these channels:

* GitHub [issues](https://github.com/clymene-project/clymene/issues)
  and [discussions](https://github.com/Clymene-project/Clymene/discussions)

[//]: # (## Contributing)

[//]: # ()

[//]: # ()

[//]: # (Refer to [CONTRIBUTING.md]&#40;https://github.com/clymene-project/clymene/blob/main/CONTRIBUTING.md&#41;)

## License

Apache License 2.0, see [LICENSE](https://github.com/Clymene-project/Clymene/blob/main/LICENSE).

# Code Style

#### 1. https://github.com/golang/go/wiki/CodeReviewComments

#### 2. https://peter.bourgon.org/go-in-production/#formatting-and-style    

