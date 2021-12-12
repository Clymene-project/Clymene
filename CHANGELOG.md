Changes by Version
==================

1.0.0 (2021-12-12)
------------------
Use kafka with clymene-agent and clymene-ingester.  
With only clymene-agent, you can collect timeseries of various environments and store them in DB.
Composite writer implementation allows simultaneous storage in prometheus and elasticsearch.
- Support kafka writer
- Support prometheus/cortex writer
- Support elasticsearch writer
- Support composite writer
- Support Service Discovery(https://prometheus.io/docs/prometheus/latest/configuration/configuration/)
