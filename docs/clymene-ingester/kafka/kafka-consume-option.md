## Clymene-ingester

Clymene-ingester consumes from Kafka and send to db.

```
Clymene-ingester [flags]
```

### Options

```
      --kafka.consumer.authentication string          Authentication type used to authenticate with kafka cluster. e.g. none, kerberos, tls, plaintext (default "none")
      --kafka.consumer.brokers string                 The comma-separated list of kafka brokers. i.e. '127.0.0.1:9092,0.0.0:1234' (default "127.0.0.1:9092")
      --kafka.consumer.client-id string               The Consumer Client ID that clymene-ingester will use (default "clymene")
      --kafka.consumer.encoding string                The encoding of metrics ("json", "protobuf") consumed from kafka (default "protobuf")
      --kafka.consumer.group-id string                The Consumer Group that clymene-ingester will be consuming on behalf of (default "clymene")
      --kafka.consumer.kerberos.config-file string    Path to Kerberos configuration. i.e /etc/krb5.conf (default "/etc/krb5.conf")
      --kafka.consumer.kerberos.keytab-file string    Path to keytab file. i.e /etc/security/kafka.keytab (default "/etc/security/kafka.keytab")
      --kafka.consumer.kerberos.password string       The Kerberos password used for authenticate with KDC
      --kafka.consumer.kerberos.realm string          Kerberos realm
      --kafka.consumer.kerberos.service-name string   Kerberos service name (default "kafka")
      --kafka.consumer.kerberos.use-keytab            Use of keytab instead of password, if this is true, keytab file will be used instead of password
      --kafka.consumer.kerberos.username string       The Kerberos username used for authenticate with KDC
      --kafka.consumer.plaintext.mechanism string     The plaintext Mechanism for SASL/PLAIN authentication, e.g. 'SCRAM-SHA-256' or 'SCRAM-SHA-512' or 'PLAIN' (default "PLAIN")
      --kafka.consumer.plaintext.password string      The plaintext Password for SASL/PLAIN authentication
      --kafka.consumer.plaintext.username string      The plaintext Username for SASL/PLAIN authentication
      --kafka.consumer.protocol-version string        Kafka protocol version - must be supported by kafka server
      --kafka.consumer.tls.ca string                  Path to a TLS CA (Certification Authority) file used to verify the remote server(s) (by default will use the system truststore)
      --kafka.consumer.tls.cert string                Path to a TLS Certificate file, used to identify this process to the remote server(s)
      --kafka.consumer.tls.enabled                    Enable TLS when talking to the remote server(s)
      --kafka.consumer.tls.key string                 Path to a TLS Private Key file, used to identify this process to the remote server(s)
      --kafka.consumer.tls.server-name string         Override the TLS server name we expect in the certificate of the remote server(s)
      --kafka.consumer.tls.skip-host-verify           (insecure) Skip server's certificate chain and host name verification
      --kafka.consumer.topic string                   The name of the kafka topic to consume from (default "clymene")
```