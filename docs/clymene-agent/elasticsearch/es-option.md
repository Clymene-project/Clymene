## clymene-agent

clymene agent is a local daemon program which scrapes metric data.

### Synopsis

clymene agent is a daemon program that runs on every cluster and scrapes metric data.

```
clymene-agent [flags]
```

### Options

```
--es-archive.bulk.actions int               The number of requests that can be enqueued before the bulk processor decides to commit (default 1000)
--es-archive.bulk.flush-interval duration   A time.Duration after which bulk requests are committed, regardless of other thresholds. Set to zero to disable. By default, this is disabled. (default 200ms)
--es-archive.bulk.size int                  The number of bytes that the bulk requests can take up before the bulk processor decides to commit (default 5000000)
--es-archive.bulk.workers int               The number of workers that are able to receive bulk requests and eventually commit them to Elasticsearch (default 1)
--es-archive.enabled                        Enable extra storage
--es-archive.log-level string               The Elasticsearch client log-level. Valid levels: [debug, info, error] (default "error")
--es-archive.max-doc-count int              The maximum document count to return from an Elasticsearch query. This will also apply to aggregations. (default 10000)
--es-archive.password string                The password required by Elasticsearch
--es-archive.remote-read-clusters string    Comma-separated list of Elasticsearch remote cluster names for cross-cluster querying.See Elasticsearch remote clusters and cross-cluster query api.
--es-archive.server-urls string             The comma-separated list of Elasticsearch servers, must be full url i.e. http://localhost:9200 (default "http://127.0.0.1:9200")
--es-archive.sniffer                        The sniffer config for Elasticsearch; client uses sniffing process to find all nodes automatically, disable if not required
--es-archive.sniffer-tls-enabled            Option to enable TLS when sniffing an Elasticsearch Cluster ; client uses sniffing process to find all nodes automatically, disabled by default
--es-archive.timeout duration               Timeout used for queries. A Timeout of zero means no timeout (default 0s)
--es-archive.tls.ca string                  Path to a TLS CA (Certification Authority) file used to verify the remote server(s) (by default will use the system truststore)
--es-archive.tls.cert string                Path to a TLS Certificate file, used to identify this process to the remote server(s)
--es-archive.tls.enabled                    Enable TLS when talking to the remote server(s)
--es-archive.tls.key string                 Path to a TLS Private Key file, used to identify this process to the remote server(s)
--es-archive.tls.server-name string         Override the TLS server name we expect in the certificate of the remote server(s)
--es-archive.tls.skip-host-verify           (insecure) Skip server's certificate chain and host name verification
--es-archive.token-file string              Path to a file containing bearer token. This flag also loads CA if it is specified.
--es-archive.username string                The username required by Elasticsearch. The basic authentication also loads CA if it is specified.
--es-archive.version uint                   The major Elasticsearch version. If not specified, the value will be auto-detected from Elasticsearch.
--es.bulk.actions int                       The number of requests that can be enqueued before the bulk processor decides to commit (default 1000)
--es.bulk.flush-interval duration           A time.Duration after which bulk requests are committed, regardless of other thresholds. Set to zero to disable. By default, this is disabled. (default 200ms)
--es.bulk.size int                          The number of bytes that the bulk requests can take up before the bulk processor decides to commit (default 5000000)
--es.bulk.workers int                       The number of workers that are able to receive bulk requests and eventually commit them to Elasticsearch (default 1)
--es.log-level string                       The Elasticsearch client log-level. Valid levels: [debug, info, error] (default "error")
--es.max-doc-count int                      The maximum document count to return from an Elasticsearch query. This will also apply to aggregations. (default 10000)
--es.password string                        The password required by Elasticsearch
--es.remote-read-clusters string            Comma-separated list of Elasticsearch remote cluster names for cross-cluster querying.See Elasticsearch remote clusters and cross-cluster query api.
--es.server-urls string                     The comma-separated list of Elasticsearch servers, must be full url i.e. http://localhost:9200 (default "http://127.0.0.1:9200")
--es.sniffer                                The sniffer config for Elasticsearch; client uses sniffing process to find all nodes automatically, disable if not required
--es.sniffer-tls-enabled                    Option to enable TLS when sniffing an Elasticsearch Cluster ; client uses sniffing process to find all nodes automatically, disabled by default
--es.timeout duration                       Timeout used for queries. A Timeout of zero means no timeout (default 0s)
--es.tls.ca string                          Path to a TLS CA (Certification Authority) file used to verify the remote server(s) (by default will use the system truststore)
--es.tls.cert string                        Path to a TLS Certificate file, used to identify this process to the remote server(s)
--es.tls.enabled                            Enable TLS when talking to the remote server(s)
--es.tls.key string                         Path to a TLS Private Key file, used to identify this process to the remote server(s)
--es.tls.server-name string                 Override the TLS server name we expect in the certificate of the remote server(s)
--es.tls.skip-host-verify                   (insecure) Skip server's certificate chain and host name verification
--es.token-file string                      Path to a file containing bearer token. This flag also loads CA if it is specified.
--es.username string                        The username required by Elasticsearch. The basic authentication also loads CA if it is specified.
--es.version uint                           The major Elasticsearch version. If not specified, the value will be auto-detected from Elasticsearch.
```

###### Auto generated by spf13/cobra on 12-Dec-2021
