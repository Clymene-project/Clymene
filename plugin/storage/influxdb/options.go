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

package influxdb

import (
	"flag"
	"fmt"
	"github.com/Clymene-project/Clymene/pkg/config/tlscfg"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/spf13/viper"
	"strings"
	"time"
)

// https://docs.influxdata.com/influxdb/v2.0/api-guide/client-libraries/go/
// https://docs.influxdata.com/influxdb/v2.0/write-data/best-practices/optimize-writes/?t=Client+libraries

const (
	configPrefix = "influxdb"
	writePrefix  = ".write"
	httpPrefix   = ".http"

	suffixUrl    = ".url"
	suffixToken  = ".token"
	suffixOrg    = ".org"
	suffixBucket = ".bucket"

	// writeOptions
	suffixBatchSize        = writePrefix + ".batch-size"
	suffixFlushInterval    = writePrefix + ".flush-interval"
	suffixPrecision        = writePrefix + ".precision"
	suffixUseGZip          = writePrefix + ".use-gzip"
	suffixDefaultTags      = writePrefix + ".default-tags"
	suffixRetryInterval    = writePrefix + ".retry-interval"
	suffixMaxRetries       = writePrefix + ".max-retries"
	suffixRetryBufferLimit = writePrefix + ".retry-buffer-limit"
	suffixMaxRetryInterval = writePrefix + ".max-retry-interval"
	suffixMaxRetryTime     = writePrefix + ".max-retry-time"
	suffixExponentialBase  = writePrefix + ".exponential-base"

	defaultUrl    = "http://localhost:8086"
	defaultToken  = ""
	defaultOrg    = ""
	defaultBucket = ""

	// writeOptions
	defaultBatchSize        = 5_000
	defaultFlushInterval    = 1_000
	defaultPrecision        = time.Nanosecond
	defaultUseGZip          = false
	defaultDefaultTags      = ""
	defaultRetryInterval    = 5_000
	defaultMaxRetries       = 5
	defaultRetryBufferLimit = 50_000
	defaultMaxRetryInterval = 125_000
	defaultMaxRetryTime     = 180_000
	defaultExponentialBase  = 2

	// http Options
	suffixHTTPRequestTimeout  = httpPrefix + ".http-request-timeout"
	defaultHTTPRequestTimeout = 10 * time.Second
)

var tlsFlagsConfig = tlscfg.ClientFlagsConfig{
	Prefix:         configPrefix,
	ShowEnabled:    true,
	ShowServerName: true,
}

type Options struct {
	url    string
	token  string
	org    string
	bucket string

	influxdb2.Options
	TLS tlscfg.Options
}

func (o *Options) AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(
		configPrefix+suffixUrl,
		defaultUrl,
		"the influxdb url",
	)
	flagSet.String(
		configPrefix+suffixToken,
		defaultToken,
		"Use the Authorization header and the Token scheme",
	)
	flagSet.String(
		configPrefix+suffixOrg,
		defaultOrg,
		"influx organization, An organization is a workspace for a group of users.",
	)
	flagSet.String(
		configPrefix+suffixBucket,
		defaultBucket,
		"influx bucket, A bucket is a named location where time series data is stored",
	)
	flagSet.Uint(
		configPrefix+suffixBatchSize,
		defaultBatchSize,
		"Maximum number of points sent to server in single request",
	)
	flagSet.Uint(
		configPrefix+suffixFlushInterval,
		defaultFlushInterval,
		"Interval, in ms, in which is buffer flushed if it has not been already written (by reaching batch size)",
	)
	flagSet.Duration(
		configPrefix+suffixPrecision,
		defaultPrecision,
		"Precision to use in writes for timestamp. In unit of duration: time.Nanosecond, time.Microsecond, time.Millisecond, time.Second",
	)
	flagSet.Bool(
		configPrefix+suffixUseGZip,
		defaultUseGZip,
		"Whether to use GZip compression in requests",
	)
	flagSet.String(
		configPrefix+suffixDefaultTags,
		defaultDefaultTags,
		"Tags added to each point during writing. separated by , (TAG1=VALUE1,TAG2=VALUE2)",
	)
	flagSet.Uint(
		configPrefix+suffixRetryInterval,
		defaultRetryInterval,
		"Default retry interval in ms, if not sent by server",
	)
	flagSet.Uint(
		configPrefix+suffixMaxRetries,
		defaultMaxRetries,
		"Maximum count of retry attempts of failed writes",
	)
	flagSet.Uint(
		configPrefix+suffixRetryBufferLimit,
		defaultRetryBufferLimit,
		"Maximum number of points to keep for retry. Should be multiple of BatchSize",
	)
	flagSet.Uint(
		configPrefix+suffixMaxRetryInterval,
		defaultMaxRetryInterval,
		"The maximum delay between each retry attempt in milliseconds",
	)
	flagSet.Uint(
		configPrefix+suffixMaxRetryTime,
		defaultMaxRetryTime,
		"The maximum total retry timeout in millisecond",
	)
	flagSet.Uint(
		configPrefix+suffixExponentialBase,
		defaultExponentialBase,
		"The base for the exponential retry delay",
	)
	flagSet.Duration(
		configPrefix+suffixHTTPRequestTimeout,
		defaultHTTPRequestTimeout,
		"HTTP request timeout in sec",
	)
	tlsFlagsConfig.AddFlags(flagSet)
}

func (o *Options) InitFromViper(v *viper.Viper) {
	o.url = v.GetString(configPrefix + suffixUrl)
	o.token = v.GetString(configPrefix + suffixToken)
	o.org = v.GetString(configPrefix + suffixOrg)
	o.bucket = v.GetString(configPrefix + suffixBucket)

	// writeOptions
	o.SetBatchSize(v.GetUint(configPrefix + suffixBatchSize))
	o.SetFlushInterval(v.GetUint(configPrefix + suffixFlushInterval))
	o.SetPrecision(v.GetDuration(configPrefix + suffixPrecision))
	o.SetUseGZip(v.GetBool(configPrefix + suffixUseGZip))
	o.SetRetryInterval(v.GetUint(configPrefix + suffixRetryInterval))
	o.SetMaxRetries(v.GetUint(configPrefix + suffixMaxRetries))
	o.SetRetryBufferLimit(v.GetUint(configPrefix + suffixRetryBufferLimit))
	o.SetMaxRetryInterval(v.GetUint(configPrefix + suffixMaxRetryInterval))
	o.SetMaxRetryTime(v.GetUint(configPrefix + suffixMaxRetryTime))
	o.SetExponentialBase(v.GetUint(configPrefix + suffixExponentialBase))

	// Convert string to tags
	o.makeTags(v.GetString(configPrefix + suffixDefaultTags))

	// HTTP Options
	o.SetHTTPRequestTimeout(v.GetUint(configPrefix + suffixHTTPRequestTimeout))

	o.TLS = tlsFlagsConfig.InitFromViper(v)

}

func (o *Options) makeTags(StringTags string) {
	tags := strings.ReplaceAll(StringTags, " ", "") // for blank safety
	for _, tag := range strings.Split(tags, ",") {
		splitTag := strings.Split(tag, "=")
		o.AddDefaultTag(splitTag[0], splitTag[1])
	}
}

func (o *Options) checkNecessaryOptions() error {
	if o.bucket == "" {
		return fmt.Errorf("influxdb bucket is not set, it is necessary")
	}
	if o.token == "" {
		return fmt.Errorf("influxdb token is not set, it is necessary")
	}
	if o.org == "" {
		return fmt.Errorf("influxdb org is not set, it is necessary")
	}
	return nil
}
