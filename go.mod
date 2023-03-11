module github.com/Clymene-project/Clymene

go 1.18

require (
	cloud.google.com/go/pubsub v1.15.0
	github.com/Azure/azure-sdk-for-go v58.3.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.23
	github.com/Azure/go-autorest/autorest/adal v0.9.18
	github.com/Shopify/sarama v1.30.0
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137
	github.com/aws/aws-sdk-go v1.42.28
	github.com/bmatcuk/doublestar v1.2.2
	github.com/bsm/sarama-cluster v2.1.15+incompatible
	github.com/buger/jsonparser v1.1.1
	github.com/c2h5oh/datasize v0.0.0-20200112174442-28bbd4740fee
	github.com/cespare/xxhash/v2 v2.1.2
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/digitalocean/godo v1.73.0
	github.com/docker/docker v20.10.12+incompatible
	github.com/envoyproxy/go-control-plane v0.10.1
	github.com/envoyproxy/protoc-gen-validate v0.6.2
	github.com/fatih/color v1.13.0
	github.com/felixge/fgprof v0.9.1
	github.com/fsnotify/fsnotify v1.5.1
	github.com/go-kit/log v0.2.0
	github.com/go-logfmt/logfmt v0.5.1
	github.com/go-zookeeper/zk v1.0.2
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/golang/snappy v0.0.4
	github.com/google/go-cmp v0.5.6
	github.com/gophercloud/gophercloud v0.24.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/grafana/dskit v0.0.0-20220209070952-ea22a8f662d0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/hashicorp/consul/api v1.12.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/hetznercloud/hcloud-go v1.33.1
	github.com/hpcloud/tail v1.0.0
	github.com/imdario/mergo v0.3.12
	github.com/influxdata/go-syslog/v3 v3.0.1-0.20201128200927-a1889d947b48
	github.com/influxdata/influxdb-client-go/v2 v2.6.0
	github.com/influxdata/telegraf v1.20.0
	github.com/jaegertracing/jaeger v1.29.0
	github.com/jmespath/go-jmespath v0.4.0
	github.com/json-iterator/go v1.1.12
	github.com/klauspost/compress v1.15.0 // indirect
	github.com/kolo/xmlrpc v0.0.0-20201022064351-38db28db192b
	github.com/linode/linodego v1.2.1
	github.com/miekg/dns v1.1.45
	github.com/mitchellh/mapstructure v1.4.3
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f
	github.com/oklog/run v1.1.0
	github.com/olivere/elastic v6.2.37+incompatible
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.1
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.32.1
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.7.0.20210223165440-c65ae3540d44
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
	github.com/taosdata/taosadapter v0.0.0-20220119131555-6d401bb95e11
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible
	github.com/weaveworks/common v0.0.0-20211015155308-ebe5bdc2c89e
	github.com/wenj91/taos-driver v0.0.0-20211207130521-b61711c4da89
	github.com/xdg-go/scram v1.0.2
	go.uber.org/atomic v1.9.0
	go.uber.org/goleak v1.1.12 // indirect
	go.uber.org/zap v1.19.1
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	golang.org/x/sys v0.5.0
	golang.org/x/time v0.0.0-20211116232009-f0f3c7e86c11
	google.golang.org/api v0.63.0
	google.golang.org/genproto v0.0.0-20211223182754-3ac035c7e7cb
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.22.4
	k8s.io/apimachinery v0.22.4
	k8s.io/client-go v0.22.4
)

require (
	github.com/google/renameio/v2 v2.0.0
	gopkg.in/Graylog2/go-gelf.v2 v2.0.0-20191017102106-1550ee647df0
)

require (
	cloud.google.com/go v0.99.0 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/VividCortex/gohistogram v1.0.0 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/armon/go-metrics v0.3.9 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/census-instrumentation/opencensus-proto v0.2.1 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cncf/udpa/go v0.0.0-20210930031921-04548b0d99d4 // indirect
	github.com/cncf/xds/go v0.0.0-20211011173535-cb28da3451f1 // indirect
	github.com/containerd/containerd v1.5.18 // indirect
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deepmap/oapi-codegen v1.8.2 // indirect
	github.com/docker/distribution v2.8.0+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/go-kit/kit v0.12.0 // indirect
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/go-resty/resty/v2 v2.1.1-0.20191201195748-d7b97669fe48 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/status v1.1.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/pprof v0.0.0-20211214055906-6f57359322fd // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gax-go/v2 v2.1.1 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.0.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-msgpack v0.5.5 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/memberlist v0.3.0 // indirect
	github.com/hashicorp/serf v0.9.6 // indirect
	github.com/huandu/xstrings v1.3.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/influxdata/line-protocol v0.0.0-20200327222509-2487e7298839 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.2 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/leodido/ragel-machinery v0.0.0-20181214104525-299bdde78165 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.0 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/opentracing-contrib/go-grpc v0.0.0-20210225150812-73cb765af46e // indirect
	github.com/opentracing-contrib/go-stdlib v1.0.0 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/node_exporter v1.0.0-rc.0.0.20200428091818-01054558c289 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529 // indirect
	github.com/sercand/kuberesolver v2.4.0+incompatible // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.34.0 // indirect
	github.com/weaveworks/promrus v1.2.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	go.etcd.io/etcd v3.3.25+incompatible // indirect
	go.etcd.io/etcd/api/v3 v3.5.0 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.0 // indirect
	go.etcd.io/etcd/client/v3 v3.5.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/term v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	golang.org/x/tools v0.1.12 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/fsnotify/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.63.2 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	k8s.io/klog/v2 v2.40.1 // indirect
	k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.1.2 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

require (
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/cloudflare/cloudflare-go v0.0.0-00010101000000-000000000000
	github.com/dnaeon/go-vcr v1.2.0 // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/onsi/gomega v1.16.0 // indirect
	github.com/prometheus/prometheus v1.8.2-0.20200911110723-e83ef207b6c2
)

// Upgrade to run with gRPC 1.3.0 and above.
replace github.com/sercand/kuberesolver => github.com/sercand/kuberesolver v2.4.0+incompatible

replace github.com/hpcloud/tail => github.com/grafana/tail v0.0.0-20201004203643-7aa4e4a91f03

replace k8s.io/client-go => k8s.io/client-go v0.22.4

// We use a fork of Graylog to avoid leaking goroutine when closing the Promtail target.
replace gopkg.in/Graylog2/go-gelf.v2 => github.com/grafana/go-gelf v0.0.0-20211112153804-126646b86de8

replace github.com/cloudflare/cloudflare-go => github.com/cyriltovena/cloudflare-go v0.27.1-0.20211118103540-ff77400bcb93
