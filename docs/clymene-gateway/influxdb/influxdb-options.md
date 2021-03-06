## Clymene-gateway

Clymene-gateway can receive data through gRPC.

```
Clymene-gateway [flags]
```

### Options

```
--influxdb.bucket string                        influx bucket, A bucket is a named location where time series data is stored
--influxdb.http.http-request-timeout duration   HTTP request timeout in sec (default 10s)
--influxdb.org string                           influx organization, An organization is a workspace for a group of users.
--influxdb.tls.ca string                        Path to a TLS CA (Certification Authority) file used to verify the remote server(s) (by default will use the system truststore)
--influxdb.tls.cert string                      Path to a TLS Certificate file, used to identify this process to the remote server(s)
--influxdb.tls.enabled                          Enable TLS when talking to the remote server(s)
--influxdb.tls.key string                       Path to a TLS Private Key file, used to identify this process to the remote server(s)
--influxdb.tls.server-name string               Override the TLS server name we expect in the certificate of the remote server(s)
--influxdb.tls.skip-host-verify                 (insecure) Skip server's certificate chain and host name verification
--influxdb.token string                         Use the Authorization header and the Token scheme
--influxdb.url string                           the influxdb url (default "http://localhost:8086")
--influxdb.write.batch-size uint                Maximum number of points sent to server in single request (default 5000)
--influxdb.write.default-tags string            Tags added to each point during writing. separated by , (TAG1=VALUE1,TAG2=VALUE2)
--influxdb.write.exponential-base uint          The base for the exponential retry delay (default 2)
--influxdb.write.flush-interval uint            Interval, in ms, in which is buffer flushed if it has not been already written (by reaching batch size) (default 1000)
--influxdb.write.max-retries uint               Maximum count of retry attempts of failed writes (default 5)
--influxdb.write.max-retry-interval uint        The maximum delay between each retry attempt in milliseconds (default 125000)
--influxdb.write.max-retry-time uint            The maximum total retry timeout in millisecond (default 180000)
--influxdb.write.precision duration             Precision to use in writes for timestamp. In unit of duration: time.Nanosecond, time.Microsecond, time.Millisecond, time.Second (default 1ns)
--influxdb.write.retry-buffer-limit uint        Maximum number of points to keep for retry. Should be multiple of BatchSize (default 50000)
--influxdb.write.retry-interval uint            Default retry interval in ms, if not sent by server (default 5000)
--influxdb.write.use-gzip                       Whether to use GZip compression in requests
```

###### Auto generated by spf13/cobra on 7-Jan-2022
