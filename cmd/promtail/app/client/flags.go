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

package client

import (
	"flag"
	lokiflag "github.com/Clymene-project/Clymene/pkg/lokiutil/flagext"
	"github.com/grafana/dskit/backoff"
	"github.com/grafana/dskit/flagext"
	"github.com/spf13/viper"
	"time"
)

const (
	ConfigPrefix            = "clymene-promtail"
	suffixDryRun            = ".dry-run"
	suffixInspect           = ".inspect"
	suffixPrintConfigStdErr = ".print-config-stderr"

	suffixTenantID        = ".tenant-id"
	suffixBatchWait       = ".batch-wait"
	suffixBatchSize       = ".batch-size-bytes"
	suffixMaxRetries      = ".max-retries"
	suffixMinBackOff      = ".min-backoff"
	suffixMaxBackOff      = ".max-backoff"
	suffixTimeout         = ".timeout"
	suffixExternalLabels  = ".external-labels"
	suffixStreamLagLabels = ".stream-lag-labels"

	// config file flag is "--config.file="
	suffixConfigFile = "config.file"

	defaultDryRun            = false
	defaultInspect           = false
	defaultPrintConfigStdErr = false
	defaultConfigFile        = "/etc/promtail/config.yml"

	defaultBatchWait      = 1 * time.Second
	defaultBatchSize  int = 1024 * 1024
	defaultMinBackoff     = 500 * time.Millisecond
	defaultMaxBackoff     = 5 * time.Minute
	defaultMaxRetries int = 10
)

type Options struct {
	printConfig bool
	DryRun      bool
	ConfigFile  string
	inspect     bool

	BackoffConfig backoff.Config `yaml:"backoff_config"`
	// The labels to add to any time series or alerts when communicating with loki
	ExternalLabels  lokiflag.LabelSet      `yaml:"external_labels,omitempty"`
	TenantID        string                 `yaml:"tenant_id"`
	StreamLagLabels flagext.StringSliceCSV `yaml:"stream_lag_labels"`
	BatchWait       time.Duration
	BatchSize       int
}

func AddFlags(flags *flag.FlagSet) {
	flags.Bool(ConfigPrefix+suffixPrintConfigStdErr, defaultPrintConfigStdErr, "Dump the entire Loki config object to stderr")
	flags.Bool(ConfigPrefix+suffixDryRun, defaultDryRun, "Start Promtail but print entries instead of sending them to Loki.")
	flags.Bool(ConfigPrefix+suffixInspect, defaultInspect, "Allows for detailed inspection of pipeline stages")
	flags.String(suffixConfigFile, defaultConfigFile, "yaml file to load")
	flags.Duration(
		ConfigPrefix+suffixBatchWait,
		defaultBatchWait,
		"Maximum wait period before sending batch.",
	)
	flags.Int(
		ConfigPrefix+suffixBatchSize,
		defaultBatchSize,
		"Maximum batch size to accrue before sending.",
	)
	// Default backoff schedule: 0.5s, 1s, 2s, 4s, 8s, 16s, 32s, 64s, 128s, 256s(4.267m) For a total time of 511.5s(8.5m) before logs are lost
	flags.Int(
		ConfigPrefix+suffixMaxRetries,
		defaultMaxRetries,
		"Maximum number of retires when sending batches.",
	)
	flags.Duration(
		ConfigPrefix+suffixMinBackOff,
		defaultMinBackoff,
		"Initial backoff time between retries.",
	)
	flags.Duration(
		ConfigPrefix+suffixMaxBackOff,
		defaultMaxBackoff,
		"Maximum backoff time between retries.",
	)
	flags.String(
		ConfigPrefix+suffixExternalLabels,
		"",
		"list of external labels to add to each log (e.g: --loki.client.external-labels=lb1=v1,lb2=v2)",
	)
	flags.String(
		ConfigPrefix+suffixTenantID,
		"",
		"Tenant ID to use when pushing logs to Loki.",
	)
	flags.String(
		ConfigPrefix+suffixStreamLagLabels,
		"filename",
		"Comma-separated list of labels to use when calculating stream lag",
	)
}

// InitFromViper initializes Options with properties retrieved from Viper.
func (o *Options) InitFromViper(v *viper.Viper) *Options {
	o.ConfigFile = v.GetString(suffixConfigFile)
	o.DryRun = v.GetBool(ConfigPrefix + suffixDryRun)
	o.inspect = v.GetBool(ConfigPrefix + suffixInspect)
	o.printConfig = v.GetBool(ConfigPrefix + suffixPrintConfigStdErr)

	o.BatchWait = v.GetDuration(ConfigPrefix + suffixBatchWait)
	o.BatchSize = v.GetInt(ConfigPrefix + suffixBatchSize)
	o.BackoffConfig.MaxRetries = v.GetInt(ConfigPrefix + suffixMaxRetries)
	o.BackoffConfig.MinBackoff = v.GetDuration(ConfigPrefix + suffixMinBackOff)
	o.BackoffConfig.MaxBackoff = v.GetDuration(ConfigPrefix + suffixMaxBackOff)

	externalLabels := v.GetString(ConfigPrefix + suffixExternalLabels)
	if externalLabels != "" {
		err := o.ExternalLabels.Set(externalLabels)
		if err != nil {
			panic("loki external labels parse error")
		}
	}
	o.TenantID = v.GetString(ConfigPrefix + suffixTenantID)

	streamLagLabels := v.GetString(ConfigPrefix + suffixStreamLagLabels)
	if streamLagLabels != "" {
		err := o.StreamLagLabels.Set(v.GetString(ConfigPrefix + suffixStreamLagLabels))
		if err != nil {
			panic("loki stream lag labels parse error")
		}
	}

	return o
}
