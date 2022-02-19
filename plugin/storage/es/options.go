// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package es

import (
	"flag"
	"github.com/Clymene-project/Clymene/pkg/config/tlscfg"
	"github.com/Clymene-project/Clymene/pkg/es/config"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	suffixUsername           = ".username"
	suffixPassword           = ".password"
	suffixSniffer            = ".sniffer"
	suffixSnifferTLSEnabled  = ".sniffer-tls-enabled"
	suffixTokenPath          = ".token-file"
	suffixServerURLs         = ".server-urls"
	suffixRemoteReadClusters = ".remote-read-clusters"
	suffixBulkSize           = ".bulk.size"
	suffixBulkWorkers        = ".bulk.workers"
	suffixBulkActions        = ".bulk.actions"
	suffixBulkFlushInterval  = ".bulk.flush-interval"
	suffixTimeout            = ".timeout"
	suffixIndexPrefix        = ".index-prefix"
	suffixEnabled            = ".enabled"
	suffixVersion            = ".version"
	suffixMaxDocCount        = ".max-doc-count"
	suffixLogLevel           = ".log-level"
	suffixMetricIndex        = "es.clymene-agent.index.name"
	suffixLogIndex           = "es.clymene-promtail.index.name"
	// default number of documents to return from a query (es allowed limit)
	// see search.max_buckets and index.max_result_window
	defaultMaxDocCount        = 10_000
	defaultServerURL          = "http://127.0.0.1:9200"
	defaultRemoteReadClusters = ""
	defaultMetricIndex        = "clymene-metrics"
	defaultLogIndex           = "clymene-logs"
	// default separator for Elasticsearch index date layout.
)

// TODO this should be moved next to config.Configuration struct (maybe ./flags package)

// Options contains various type of Elasticsearch configs and provides the ability
// to bind them to command line flag and apply overlays, so that some configurations
// (e.g. archive) may be underspecified and infer the rest of its parameters from primary.
type Options struct {
	Primary namespaceConfig `mapstructure:",squash"`

	others map[string]*namespaceConfig

	metricsIndex string
	logsIndex    string
}

type namespaceConfig struct {
	config.Configuration `mapstructure:",squash"`
	namespace            string
}

// NewOptions creates a new Options struct.
func NewOptions(primaryNamespace string, otherNamespaces ...string) *Options {
	defaultConfig := config.Configuration{
		Username: "",
		Password: "",
		Sniffer:  false,

		BulkSize:          5 * 1000 * 1000,
		BulkWorkers:       1,
		BulkActions:       1000,
		BulkFlushInterval: time.Millisecond * 200,

		Enabled:            true,
		Version:            0,
		Servers:            []string{defaultServerURL},
		RemoteReadClusters: []string{},
		MaxDocCount:        defaultMaxDocCount,
		LogLevel:           "error",
	}
	options := &Options{
		Primary: namespaceConfig{
			Configuration: defaultConfig,
			namespace:     primaryNamespace,
		},
		others: make(map[string]*namespaceConfig, len(otherNamespaces)),
	}

	// Other namespaces need to be explicitly enabled.
	defaultConfig.Enabled = false
	for _, namespace := range otherNamespaces {
		options.others[namespace] = &namespaceConfig{
			Configuration: defaultConfig,
			namespace:     namespace,
		}
	}

	return options
}

func (config *namespaceConfig) getTLSFlagsConfig() tlscfg.ClientFlagsConfig {
	return tlscfg.ClientFlagsConfig{
		Prefix:         config.namespace,
		ShowEnabled:    true,
		ShowServerName: true,
	}
}

// AddFlags adds flags for Options
func (opt *Options) AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(suffixMetricIndex, defaultMetricIndex, "Set the index name to save the metrics collected by clymene-agent.")
	flagSet.String(suffixLogIndex, defaultLogIndex, "Set the index name to save the log collected by clymene-promtail.")
	addFlags(flagSet, &opt.Primary)
	for _, cfg := range opt.others {
		addFlags(flagSet, cfg)
	}
}

func addFlags(flagSet *flag.FlagSet, nsConfig *namespaceConfig) {
	flagSet.String(
		nsConfig.namespace+suffixUsername,
		nsConfig.Username,
		"The username required by Elasticsearch. The basic authentication also loads CA if it is specified.")
	flagSet.String(
		nsConfig.namespace+suffixPassword,
		nsConfig.Password,
		"The password required by Elasticsearch")
	flagSet.String(
		nsConfig.namespace+suffixTokenPath,
		nsConfig.TokenFilePath,
		"Path to a file containing bearer token. This flag also loads CA if it is specified.")
	flagSet.Bool(
		nsConfig.namespace+suffixSniffer,
		nsConfig.Sniffer,
		"The sniffer config for Elasticsearch; client uses sniffing process to find all nodes automatically, disable if not required")
	flagSet.String(
		nsConfig.namespace+suffixServerURLs,
		defaultServerURL,
		"The comma-separated list of Elasticsearch servers, must be full url i.e. http://localhost:9200")
	flagSet.String(
		nsConfig.namespace+suffixRemoteReadClusters,
		defaultRemoteReadClusters,
		"Comma-separated list of Elasticsearch remote cluster names for cross-cluster querying."+
			"See Elasticsearch remote clusters and cross-cluster query api.")
	flagSet.Duration(
		nsConfig.namespace+suffixTimeout,
		nsConfig.Timeout,
		"Timeout used for queries. A Timeout of zero means no timeout")
	flagSet.Int(
		nsConfig.namespace+suffixBulkSize,
		nsConfig.BulkSize,
		"The number of bytes that the bulk requests can take up before the bulk processor decides to commit")
	flagSet.Int(
		nsConfig.namespace+suffixBulkWorkers,
		nsConfig.BulkWorkers,
		"The number of workers that are able to receive bulk requests and eventually commit them to Elasticsearch")
	flagSet.Int(
		nsConfig.namespace+suffixBulkActions,
		nsConfig.BulkActions,
		"The number of requests that can be enqueued before the bulk processor decides to commit")
	flagSet.Duration(
		nsConfig.namespace+suffixBulkFlushInterval,
		nsConfig.BulkFlushInterval,
		"A time.Duration after which bulk requests are committed, regardless of other thresholds. Set to zero to disable. By default, this is disabled.")
	flagSet.Uint(
		nsConfig.namespace+suffixVersion,
		0,
		"The major Elasticsearch version. If not specified, the value will be auto-detected from Elasticsearch.")
	flagSet.Bool(
		nsConfig.namespace+suffixSnifferTLSEnabled,
		nsConfig.SnifferTLSEnabled,
		"Option to enable TLS when sniffing an Elasticsearch Cluster ; client uses sniffing process to find all nodes automatically, disabled by default")
	flagSet.Int(
		nsConfig.namespace+suffixMaxDocCount,
		nsConfig.MaxDocCount,
		"The maximum document count to return from an Elasticsearch query. This will also apply to aggregations.")
	flagSet.String(
		nsConfig.namespace+suffixLogLevel,
		nsConfig.LogLevel,
		"The Elasticsearch client log-level. Valid levels: [debug, info, error]")

	if nsConfig.namespace == archiveNamespace {
		flagSet.Bool(
			nsConfig.namespace+suffixEnabled,
			nsConfig.Enabled,
			"Enable extra storage")
	}

	nsConfig.getTLSFlagsConfig().AddFlags(flagSet)
}

// InitFromViper initializes Options with properties from viper
func (opt *Options) InitFromViper(v *viper.Viper) {
	opt.metricsIndex = v.GetString(suffixMetricIndex)
	opt.logsIndex = v.GetString(suffixLogIndex)
	initFromViper(&opt.Primary, v)
	for _, cfg := range opt.others {
		initFromViper(cfg, v)
	}
}

func initFromViper(cfg *namespaceConfig, v *viper.Viper) {
	cfg.Username = v.GetString(cfg.namespace + suffixUsername)
	cfg.Password = v.GetString(cfg.namespace + suffixPassword)
	cfg.TokenFilePath = v.GetString(cfg.namespace + suffixTokenPath)
	cfg.Sniffer = v.GetBool(cfg.namespace + suffixSniffer)
	cfg.SnifferTLSEnabled = v.GetBool(cfg.namespace + suffixSnifferTLSEnabled)
	cfg.Servers = strings.Split(stripWhiteSpace(v.GetString(cfg.namespace+suffixServerURLs)), ",")
	cfg.BulkSize = v.GetInt(cfg.namespace + suffixBulkSize)
	cfg.BulkWorkers = v.GetInt(cfg.namespace + suffixBulkWorkers)
	cfg.BulkActions = v.GetInt(cfg.namespace + suffixBulkActions)
	cfg.BulkFlushInterval = v.GetDuration(cfg.namespace + suffixBulkFlushInterval)
	cfg.Timeout = v.GetDuration(cfg.namespace + suffixTimeout)
	cfg.IndexPrefix = v.GetString(cfg.namespace + suffixIndexPrefix)
	cfg.Enabled = v.GetBool(cfg.namespace + suffixEnabled)
	cfg.Version = uint(v.GetInt(cfg.namespace + suffixVersion))
	cfg.LogLevel = v.GetString(cfg.namespace + suffixLogLevel)
	cfg.MaxDocCount = v.GetInt(cfg.namespace + suffixMaxDocCount)

	// TODO: Need to figure out a better way for do this.
	cfg.AllowTokenFromContext = v.GetBool(metricstore.StoragePropagationKey)
	cfg.TLS = cfg.getTLSFlagsConfig().InitFromViper(v)

	remoteReadClusters := stripWhiteSpace(v.GetString(cfg.namespace + suffixRemoteReadClusters))
	if len(remoteReadClusters) > 0 {
		cfg.RemoteReadClusters = strings.Split(remoteReadClusters, ",")
	}
}

// GetPrimary returns primary configuration.
func (opt *Options) GetPrimary() *config.Configuration {
	return &opt.Primary.Configuration
}

// Get returns auxiliary named configuration.
func (opt *Options) Get(namespace string) *config.Configuration {
	nsCfg, ok := opt.others[namespace]
	if !ok {
		nsCfg = &namespaceConfig{}
		opt.others[namespace] = nsCfg
	}
	nsCfg.Configuration.ApplyDefaults(&opt.Primary.Configuration)
	if len(nsCfg.Configuration.Servers) == 0 {
		nsCfg.Servers = opt.Primary.Servers
	}
	return &nsCfg.Configuration
}

// stripWhiteSpace removes all whitespace characters from a string
func stripWhiteSpace(str string) string {
	return strings.Replace(str, " ", "", -1)
}
