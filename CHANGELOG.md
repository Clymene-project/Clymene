Changes by Version
==================

1.2.1 (2022-01-01)
-----------------
Separate the prometheus option and the cortex option to block confusion

- cortex option (TS_STORAGE_TYPE=cortex)
- [Cortex option](./docs/clymene-agent/cortex/cortex-option.md)

1.2.0 (2021-12-29)
-----------------
Now, you can use the opentsdb in Clymene.

- opentsdb support (TS_STORAGE_TYPE=opentsdb)
- opentsdb supports two methods: socket(default) and
  http([Opentsdb option](./docs/clymene-agent/opentsdb/opentsdb-option.md))

1.1.0 (2021-12-20)
-----------------
A gateway supporting gRPC communication has been added.  
Try using it in various architectures.

- gRPC support (TS_STORAGE_TYPE=gateway)

1.0.0 (2021-12-12)
------------------
Use kafka with clymene-agent and clymene-ingester.  
With only clymene-agent, you can collect timeseries of various environments and store them in DB. Composite writer
implementation allows simultaneous storage in prometheus and elasticsearch.

- Support kafka writer
- Support prometheus/cortex writer
- Support elasticsearch writer
- Support composite writer
- Support Service Discovery(https://prometheus.io/docs/prometheus/latest/configuration/configuration/)
