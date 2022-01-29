## clymene-agent

clymene agent is a local daemon program which scrapes metric data.

### Synopsis

clymene agent is a daemon program that runs on every cluster and scrapes metric data.

```
clymene-agent [flags]
```

### Options

```
--opentsdb.dry-run                Don't actually send anything to the TSD, just print the datapoints.
--opentsdb.host string            Hostname to use to connect to the TSD. (default "localhost")
--opentsdb.hosts string           List of host:port to connect to tsd's (comma separated)
--opentsdb.http                   Send the data via the http interface (default 'false')
--opentsdb.http-api-path string   URL path to use for HTTP requests to TSD. (default "api/put")
--opentsdb.http-password string   Password to use for HTTP Basic Auth when sending the data via HTTP
--opentsdb.http-username string   Username to use for HTTP Basic Auth when sending the data via HTTP
--opentsdb.max-chunk int          The maximum request body size to support for incoming HTTP requests when chunking is enabled (default 512)
--opentsdb.max-tags int           The maximum number of tags to send to our TSD Instances (default 8)
--opentsdb.port int               Port to connect to the TSD instance on (default 4242)
--opentsdb.ssl                    Enable SSL - used in conjunction with http (default 'false')
--opentsdb.timeout duration       Time out when doing http insert(sec, default 10 sec) (default 10s)
```

