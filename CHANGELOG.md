Changes by Version
==================

2.0.0 (2022-02-20)
-----------------
A log collection agent, promtail, is now available. Collected logs can be stored in loki and elasticsearch.[Getting Started](https://github.com/Clymene-project/Clymene/blob/main/docs/clymene-promtail/README.md)  

1. Add log collection component clymene-promtail issue: [#41](https://github.com/Clymene-project/Clymene/issues/41)  
2. Added http receiver function to send/receive gateway metricWriter to http issue: [#37](https://github.com/Clymene-project/Clymene/issues/37)  
3. Add elasticsearch index name change option issue: [#44](https://github.com/Clymene-project/Clymene/issues/44)  
4. Logs collected by clymene-promtail can be stored in elasticsearch issue: [#43](https://github.com/Clymene-project/Clymene/issues/43)  

1.5.0 (2022-01-30)
-----------------
Updated Clymene's service discovery feature. For details, please check the clymene official docs(https://clymene-project.github.io/docs/service-discovery/configuration/)
1. kuma(https://kuma.io/)
2. uyuni(https://www.uyuni-project.org/)
3. scaleway(https://www.scaleway.com/en/virtual-instances/)
4. puppetdb(https://puppet.com/docs/puppetdb/7/overview.html)
5. linode(https://www.linode.com/)
6. hetzner(https://www.hetzner.com/)
7. eureka(https://github.com/Netflix/eureka)
8. digitalocean(https://www.digitalocean.com/)

1.4.0 (2022-01-22)
-----------------
Now, You can use clymene using tdengine's http interface.
https://www.taosdata.com/en/
issue: #35 #31

```
--tdengine.dbname string         Destination database (default "clymene")
--tdengine.hostname string       The host to connect to TDengine server. (default "127.0.0.1")
--tdengine.max-sql-length int    Number of SQLs that can be sent at one time (default 4096)
--tdengine.password string       The password to use when connecting to the server (default "taosdata")
--tdengine.server-port int       he HTTP port number to use for the connection to TDengine server (default 6041)
--tdengine.user string           The TDengine user name to use when connecting to the server (default "root")
```

1.3.1 (2022-01-10)
-----------------
Clymene official logo added!
Creating metrics for Clymene components issue: #33  
- clymene-agent metric port = :15691/metrics  
- clymene-ingester metric port = :15694/metrics  
- clymene-gateway metric port = :15690/metrics  

1.3.0 (2022-01-08)
-----------------
influxdb is now officially supported. The options below are mandatory.
```
--influxdb.bucket string    influx bucket, A bucket is a named location where time series data is stored
--influxdb.org string       influx organization, An organization is a workspace for a group of users.
--influxdb.token string     Use the Authorization header and the Token scheme
```
TDengin support has been added to the roadmap.

- influxdb option (STORAGE_TYPE=influxdb)
- [Influxdb option](./docs/clymene-agent/influxdb/influxdb-options.md)

1.2.1 (2022-01-01)
-----------------
Separate the prometheus option and the cortex option to block confusion

- cortex option (STORAGE_TYPE=cortex)
- [Cortex option](./docs/clymene-agent/cortex/cortex-option.md)

1.2.0 (2021-12-29)
-----------------
Now, you can use the opentsdb in Clymene.

- opentsdb support (STORAGE_TYPE=opentsdb)
- opentsdb supports two methods: socket(default) and
  http([Opentsdb option](./docs/clymene-agent/opentsdb/opentsdb-option.md))

1.1.0 (2021-12-20)
-----------------
A gateway supporting gRPC communication has been added.  
Try using it in various architectures.

- gRPC support (STORAGE_TYPE=gateway)

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
