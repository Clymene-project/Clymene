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

package config

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/Clymene-project/Clymene/pkg/config/tlscfg"
	"github.com/Clymene-project/Clymene/pkg/es"
	eswrapper "github.com/Clymene-project/Clymene/pkg/es/wrapper"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/olivere/elastic"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"

	storageMetrics "github.com/Clymene-project/Clymene/storage/metricstore/metrics"
)

// Configuration describes the configuration properties needed to connect to an ElasticSearch cluster
type Configuration struct {
	Servers               []string       `mapstructure:"server_urls"`
	RemoteReadClusters    []string       `mapstructure:"remote_read_clusters"`
	Username              string         `mapstructure:"username"`
	Password              string         `mapstructure:"password" json:"-"`
	TokenFilePath         string         `mapstructure:"token_file"`
	AllowTokenFromContext bool           `mapstructure:"-"`
	Sniffer               bool           `mapstructure:"sniffer"` // https://github.com/olivere/elastic/wiki/Sniffing
	SnifferTLSEnabled     bool           `mapstructure:"sniffer_tls_enabled"`
	MaxDocCount           int            `mapstructure:"-"`                       // Defines maximum number of results to fetch from storage per query
	MaxMetricAge          time.Duration  `yaml:"max_metric_age" mapstructure:"-"` // configures the maximum lookback on metric reads
	NumShards             int64          `yaml:"shards" mapstructure:"num_shards"`
	NumReplicas           int64          `yaml:"replicas" mapstructure:"num_replicas"`
	Timeout               time.Duration  `validate:"min=500" mapstructure:"-"`
	BulkSize              int            `mapstructure:"-"`
	BulkWorkers           int            `mapstructure:"-"`
	BulkActions           int            `mapstructure:"-"`
	BulkFlushInterval     time.Duration  `mapstructure:"-"`
	IndexPrefix           string         `mapstructure:"index_prefix"`
	IndexDateLayout       string         `mapstructure:"index_date_layout"`
	Tags                  TagsAsFields   `mapstructure:"tags_as_fields"`
	Enabled               bool           `mapstructure:"-"`
	TLS                   tlscfg.Options `mapstructure:"tls"`
	UseReadWriteAliases   bool           `mapstructure:"use_aliases"`
	CreateIndexTemplates  bool           `mapstructure:"create_mappings"`
	UseILM                bool           `mapstructure:"use_ilm"`
	Version               uint           `mapstructure:"version"`
	LogLevel              string         `mapstructure:"log_level"`
}

// TagsAsFields holds configuration for tag schema.
// By default Jaeger stores tags in an array of nested objects.
// This configurations allows to store tags as object fields for better Kibana support.
type TagsAsFields struct {
	// Store all tags as object fields, instead nested objects
	AllAsFields bool `mapstructure:"all"`
	// Dot replacement for tag keys when stored as object fields
	DotReplacement string `mapstructure:"dot_replacement"`
	// File path to tag keys which should be stored as object fields
	File string `mapstructure:"config_file"`
	// Comma delimited list of tags to store as object fields
	Include string `mapstructure:"include"`
}

// ClientBuilder creates new es.Client
type ClientBuilder interface {
	NewClient(logger *zap.Logger, metricsFactory metrics.Factory) (es.Client, error)
	GetRemoteReadClusters() []string
	GetNumShards() int64
	GetNumReplicas() int64
	GetMaxMetricAge() time.Duration
	GetMaxDocCount() int
	GetIndexPrefix() string
	GetIndexDateLayout() string
	GetTagsFilePath() string
	GetAllTagsAsFields() bool
	GetTagDotReplacement() string
	GetUseReadWriteAliases() bool
	GetTokenFilePath() string
	IsStorageEnabled() bool
	IsCreateIndexTemplates() bool
	GetVersion() uint
	TagKeysAsFields() ([]string, error)
	GetUseILM() bool
	GetLogLevel() string
}

// NewClient creates a new ElasticSearch client
func (c *Configuration) NewClient(logger *zap.Logger, metricsFactory metrics.Factory) (es.Client, error) {
	if len(c.Servers) < 1 {
		return nil, errors.New("no servers specified")
	}
	options, err := c.getConfigOptions(logger)
	if err != nil {
		return nil, err
	}

	rawClient, err := elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}

	sm := storageMetrics.NewWriteMetrics(metricsFactory, "bulk_index")
	m := sync.Map{}

	service, err := rawClient.BulkProcessor().
		Before(func(id int64, requests []elastic.BulkableRequest) {
			m.Store(id, time.Now())
		}).
		After(func(id int64, requests []elastic.BulkableRequest, response *elastic.BulkResponse, err error) {
			start, ok := m.Load(id)
			if !ok {
				return
			}
			m.Delete(id)

			// log individual errors, note that err might be false and these errors still present
			if response != nil && response.Errors {
				for _, it := range response.Items {
					for key, val := range it {
						if val.Error != nil {
							logger.Error("Elasticsearch part of bulk request failed", zap.String("map-key", key),
								zap.Reflect("response", val))
						}
					}
				}
			}

			sm.Emit(err, time.Since(start.(time.Time)))
			if err != nil {
				var failed int
				if response == nil {
					failed = 0
				} else {
					failed = len(response.Failed())
				}
				total := len(requests)
				logger.Error("Elasticsearch could not process bulk request",
					zap.Int("request_count", total),
					zap.Int("failed_count", failed),
					zap.Error(err),
					zap.Any("response", response))
			}
		}).
		BulkSize(c.BulkSize).
		Workers(c.BulkWorkers).
		BulkActions(c.BulkActions).
		FlushInterval(c.BulkFlushInterval).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	if c.Version == 0 {
		// Determine ElasticSearch Version
		pingResult, _, err := rawClient.Ping(c.Servers[0]).Do(context.Background())
		if err != nil {
			return nil, err
		}
		esVersion, err := strconv.Atoi(string(pingResult.Version.Number[0]))
		if err != nil {
			return nil, err
		}
		logger.Info("Elasticsearch detected", zap.Int("version", esVersion))
		c.Version = uint(esVersion)
	}

	return eswrapper.WrapESClient(rawClient, service, c.Version), nil
}

// ApplyDefaults copies settings from source unless its own value is non-zero.
func (c *Configuration) ApplyDefaults(source *Configuration) {
	if len(c.RemoteReadClusters) == 0 {
		c.RemoteReadClusters = source.RemoteReadClusters
	}
	if c.Username == "" {
		c.Username = source.Username
	}
	if c.Password == "" {
		c.Password = source.Password
	}
	if !c.Sniffer {
		c.Sniffer = source.Sniffer
	}
	if c.MaxMetricAge == 0 {
		c.MaxMetricAge = source.MaxMetricAge
	}
	if c.NumShards == 0 {
		c.NumShards = source.NumShards
	}
	if c.NumReplicas == 0 {
		c.NumReplicas = source.NumReplicas
	}
	if c.BulkSize == 0 {
		c.BulkSize = source.BulkSize
	}
	if c.BulkWorkers == 0 {
		c.BulkWorkers = source.BulkWorkers
	}
	if c.BulkActions == 0 {
		c.BulkActions = source.BulkActions
	}
	if c.BulkFlushInterval == 0 {
		c.BulkFlushInterval = source.BulkFlushInterval
	}
	if !c.SnifferTLSEnabled {
		c.SnifferTLSEnabled = source.SnifferTLSEnabled
	}
	if !c.Tags.AllAsFields {
		c.Tags.AllAsFields = source.Tags.AllAsFields
	}
	if c.Tags.DotReplacement == "" {
		c.Tags.DotReplacement = source.Tags.DotReplacement
	}
	if c.Tags.Include == "" {
		c.Tags.Include = source.Tags.Include
	}
	if c.Tags.File == "" {
		c.Tags.File = source.Tags.File
	}
	if c.MaxDocCount == 0 {
		c.MaxDocCount = source.MaxDocCount
	}
	if c.LogLevel == "" {
		c.LogLevel = source.LogLevel
	}
}

// GetRemoteReadClusters returns list of remote read clusters
func (c *Configuration) GetRemoteReadClusters() []string {
	return c.RemoteReadClusters
}

// GetNumShards returns number of shards from Configuration
func (c *Configuration) GetNumShards() int64 {
	return c.NumShards
}

// GetNumReplicas returns number of replicas from Configuration
func (c *Configuration) GetNumReplicas() int64 {
	return c.NumReplicas
}

// GetMaxMetricAge returns max metric age from Configuration
func (c *Configuration) GetMaxMetricAge() time.Duration {
	return c.MaxMetricAge
}

// GetMaxDocCount returns the maximum number of documents that a query should return
func (c *Configuration) GetMaxDocCount() int {
	return c.MaxDocCount
}

// GetIndexPrefix returns index prefix
func (c *Configuration) GetIndexPrefix() string {
	return c.IndexPrefix
}

// GetIndexDateLayout returns index date layout
func (c *Configuration) GetIndexDateLayout() string {
	return c.IndexDateLayout
}

// GetTagsFilePath returns a path to file containing tag keys
func (c *Configuration) GetTagsFilePath() string {
	return c.Tags.File
}

// GetAllTagsAsFields returns true if all tags should be stored as object fields
func (c *Configuration) GetAllTagsAsFields() bool {
	return c.Tags.AllAsFields
}

// GetVersion returns Elasticsearch version
func (c *Configuration) GetVersion() uint {
	return c.Version
}

// GetTagDotReplacement returns character is used to replace dots in tag keys, when
// the tag is stored as object field.
func (c *Configuration) GetTagDotReplacement() string {
	return c.Tags.DotReplacement
}

// GetUseReadWriteAliases indicates whether read alias should be used
func (c *Configuration) GetUseReadWriteAliases() bool {
	return c.UseReadWriteAliases
}

// GetUseILM indicates whether ILM should be used
func (c *Configuration) GetUseILM() bool {
	return c.UseILM
}

// GetLogLevel returns the log-level the ES client should log at.
func (c *Configuration) GetLogLevel() string {
	return c.LogLevel
}

// GetTokenFilePath returns file path containing the bearer token
func (c *Configuration) GetTokenFilePath() string {
	return c.TokenFilePath
}

// IsStorageEnabled determines whether storage is enabled
func (c *Configuration) IsStorageEnabled() bool {
	return c.Enabled
}

// IsCreateIndexTemplates determines whether index templates should be created or not
func (c *Configuration) IsCreateIndexTemplates() bool {
	return c.CreateIndexTemplates
}

// TagKeysAsFields returns tags from the file and command line merged
func (c *Configuration) TagKeysAsFields() ([]string, error) {
	var tags []string

	// from file
	if c.Tags.File != "" {
		file, err := os.Open(filepath.Clean(c.Tags.File))
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if tag := strings.TrimSpace(line); tag != "" {
				tags = append(tags, tag)
			}
		}
		if err := file.Close(); err != nil {
			return nil, err
		}
	}

	// from params
	if c.Tags.Include != "" {
		tags = append(tags, strings.Split(c.Tags.Include, ",")...)
	}

	return tags, nil
}

// getConfigOptions wraps the configs to feed to the ElasticSearch client init
func (c *Configuration) getConfigOptions(logger *zap.Logger) ([]elastic.ClientOptionFunc, error) {

	options := []elastic.ClientOptionFunc{elastic.SetURL(c.Servers...), elastic.SetSniff(c.Sniffer),
		// Disable health check when token from context is allowed, this is because at this time
		// we don' have a valid token to do the check ad if we don't disable the check the service that
		// uses this won't start.
		elastic.SetHealthcheck(!c.AllowTokenFromContext)}
	if c.SnifferTLSEnabled {
		options = append(options, elastic.SetScheme("https"))
	}
	httpClient := &http.Client{
		Timeout: c.Timeout,
	}
	options = append(options, elastic.SetHttpClient(httpClient))
	options = append(options, elastic.SetBasicAuth(c.Username, c.Password))

	options, err := addLoggerOptions(options, c.LogLevel)
	if err != nil {
		return options, err
	}

	transport, err := GetHTTPRoundTripper(c, logger)
	if err != nil {
		return nil, err
	}
	httpClient.Transport = transport
	return options, nil
}

func addLoggerOptions(options []elastic.ClientOptionFunc, logLevel string) ([]elastic.ClientOptionFunc, error) {
	// Decouple ES logger from the log-level assigned to the parent application's log-level; otherwise, the least
	// permissive log-level will dominate.
	// e.g. --log-level=info and --es.log-level=debug would mute ES's debug logging and would require --log-level=debug
	// to show ES debug logs.
	prodConfig := zap.NewProductionConfig()
	prodConfig.Level.SetLevel(zap.DebugLevel)

	esLogger, err := prodConfig.Build()
	if err != nil {
		return options, err
	}

	// Elastic client requires a "Printf"-able logger.
	l := zapgrpc.NewLogger(esLogger)
	switch logLevel {
	case "debug":
		l = zapgrpc.NewLogger(esLogger, zapgrpc.WithDebug())
		options = append(options, elastic.SetTraceLog(l))
	case "info":
		options = append(options, elastic.SetInfoLog(l))
	case "error":
		options = append(options, elastic.SetErrorLog(l))
	default:
		return options, fmt.Errorf("unrecognized log-level: \"%s\"", logLevel)
	}
	return options, nil
}

// GetHTTPRoundTripper returns configured http.RoundTripper
func GetHTTPRoundTripper(c *Configuration, logger *zap.Logger) (http.RoundTripper, error) {
	if c.TLS.Enabled {
		ctlsConfig, err := c.TLS.Config(logger)
		if err != nil {
			return nil, err
		}
		return &http.Transport{
			TLSClientConfig: ctlsConfig,
		}, nil
	}
	var transport http.RoundTripper
	httpTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		// #nosec G402
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.TLS.SkipHostVerify},
	}
	if c.TLS.CAPath != "" {
		ctlsConfig, err := c.TLS.Config(logger)
		if err != nil {
			return nil, err
		}
		httpTransport.TLSClientConfig = ctlsConfig
		transport = httpTransport
	}

	token := ""
	if c.TokenFilePath != "" {
		if c.AllowTokenFromContext {
			logger.Warn("Token file and token propagation are both enabled, token from file won't be used")
		}
		tokenFromFile, err := loadToken(c.TokenFilePath)
		if err != nil {
			return nil, err
		}
		token = tokenFromFile
	}
	if token != "" || c.AllowTokenFromContext {
		transport = &tokenAuthTransport{
			token:                token,
			allowOverrideFromCtx: c.AllowTokenFromContext,
			wrapped:              httpTransport,
		}
	}
	return transport, nil
}

// TokenAuthTransport
type tokenAuthTransport struct {
	token                string
	allowOverrideFromCtx bool
	wrapped              *http.Transport
}

func (tr *tokenAuthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	token := tr.token
	if tr.allowOverrideFromCtx {
		headerToken, _ := metricstore.GetBearerToken(r.Context())
		if headerToken != "" {
			token = headerToken
		}
	}
	r.Header.Set("Authorization", "Bearer "+token)
	return tr.wrapped.RoundTrip(r)
}

func loadToken(path string) (string, error) {
	b, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(b), "\r\n"), nil
}
