/*
 * Copyright (c) 2021 The Clymene Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package loki

import (
	"flag"
	lokiflag "github.com/Clymene-project/Clymene/pkg/lokiutil/flagext"
	"github.com/grafana/dskit/backoff"
	"github.com/grafana/dskit/flagext"
	"github.com/prometheus/common/config"
	"github.com/spf13/viper"
	"time"
)

const (
	configPrefix          = "loki"
	suffixURL             = ".client.url"
	suffixBatchWait       = ".client.batch-wait"
	suffixBatchSize       = ".client.batch-size-bytes"
	suffixMaxRetries      = ".client.max-retries"
	suffixMinBackOff      = ".client.min-backoff"
	suffixMaxBackOff      = ".client.max-backoff"
	suffixTimeout         = ".client.timeout"
	suffixExternalLabels  = ".client.external-labels"
	suffixTenantID        = ".client.tenant-id"
	suffixStreamLagLabels = ".client.stream-lag-labels"

	defaultURL            = "http://localhost:3100/loki/api/v1/push"
	defaultBatchWait      = 1 * time.Second
	defaultBatchSize  int = 1024 * 1024
	defaultMinBackoff     = 500 * time.Millisecond
	defaultMaxBackoff     = 5 * time.Minute
	defaultMaxRetries int = 10
	defaultTimeout        = 10 * time.Second
)

type Options struct {
	URL       flagext.URLValue
	BatchWait time.Duration
	BatchSize int

	Client config.HTTPClientConfig `yaml:",inline"`

	BackoffConfig backoff.Config `yaml:"backoff_config"`
	// The labels to add to any time series or alerts when communicating with loki
	ExternalLabels lokiflag.LabelSet `yaml:"external_labels,omitempty"`
	Timeout        time.Duration     `yaml:"timeout"`

	// The tenant ID to use when pushing logs to Loki (empty string means
	// single tenant mode)
	TenantID string `yaml:"tenant_id"`

	StreamLagLabels flagext.StringSliceCSV `yaml:"stream_lag_labels"`
}

func (o *Options) AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(
		configPrefix+suffixURL,
		defaultURL,
		"URL of log server",
	)
	flagSet.Duration(
		configPrefix+suffixBatchWait,
		defaultBatchWait,
		"Maximum wait period before sending batch.",
	)
	flagSet.Int(
		configPrefix+suffixBatchSize,
		defaultBatchSize,
		"Maximum batch size to accrue before sending.",
	)
	// Default backoff schedule: 0.5s, 1s, 2s, 4s, 8s, 16s, 32s, 64s, 128s, 256s(4.267m) For a total time of 511.5s(8.5m) before logs are lost
	flagSet.Int(
		configPrefix+suffixMaxRetries,
		defaultMaxRetries,
		"Maximum number of retires when sending batches.",
	)
	flagSet.Duration(
		configPrefix+suffixMinBackOff,
		defaultMinBackoff,
		"Initial backoff time between retries.",
	)
	flagSet.Duration(
		configPrefix+suffixMaxBackOff,
		defaultMaxBackoff,
		"Maximum backoff time between retries.",
	)
	flagSet.Duration(
		configPrefix+suffixTimeout,
		defaultTimeout,
		"Maximum time to wait for server to respond to a request",
	)
	flagSet.String(
		configPrefix+suffixExternalLabels,
		"",
		"list of external labels to add to each log (e.g: --loki.client.external-labels=lb1=v1,lb2=v2)",
	)
	flagSet.String(
		configPrefix+suffixTenantID,
		"",
		"Tenant ID to use when pushing logs to Loki.",
	)
	o.StreamLagLabels = []string{"filename"}
	flagSet.String(
		configPrefix+suffixStreamLagLabels,
		"",
		"Comma-separated list of labels to use when calculating stream lag",
	)
}

func (o *Options) InitFromViper(v *viper.Viper) {
	err := o.URL.Set(v.GetString(configPrefix + suffixURL))
	if err != nil {
		panic("loki url parse error")
	}
	o.BatchWait = v.GetDuration(configPrefix + suffixBatchWait)
	o.BatchSize = v.GetInt(configPrefix + suffixBatchSize)
	o.BackoffConfig.MaxRetries = v.GetInt(configPrefix + suffixMaxRetries)
	o.BackoffConfig.MinBackoff = v.GetDuration(configPrefix + suffixMinBackOff)
	o.BackoffConfig.MaxBackoff = v.GetDuration(configPrefix + suffixMaxBackOff)
	o.Timeout = v.GetDuration(configPrefix + suffixTimeout)
	err = o.ExternalLabels.Set(v.GetString(configPrefix + suffixExternalLabels))
	if err != nil {
		panic("loki external labels parse error")
	}
	o.TenantID = v.GetString(configPrefix + suffixTenantID)
	err = o.StreamLagLabels.Set(v.GetString(configPrefix + suffixStreamLagLabels))
	if err != nil {
		panic("loki stream lag labels parse error")
	}
	o.StreamLagLabels = []string{"filename"}
}
