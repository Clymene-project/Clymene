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
	"github.com/grafana/dskit/backoff"
	"github.com/grafana/dskit/flagext"
	"github.com/prometheus/common/config"
	"github.com/spf13/viper"
	"time"
)

const (
	configPrefix = "loki"

	suffixURL     = ".client.url"
	suffixTimeout = ".client.timeout"

	defaultTimeout = 10 * time.Second
	defaultURL     = "http://localhost:3100/loki/api/v1/push"
)

type Options struct {
	URL           flagext.URLValue
	Client        config.HTTPClientConfig `yaml:",inline"`
	BackoffConfig backoff.Config          `yaml:"backoff_config"`
	Timeout       time.Duration           `yaml:"timeout"`
}

func (o *Options) AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(
		configPrefix+suffixURL,
		defaultURL,
		"URL of log server",
	)
	flagSet.Duration(
		configPrefix+suffixTimeout,
		defaultTimeout,
		"Maximum time to wait for server to respond to a request",
	)
}

func (o *Options) InitFromViper(v *viper.Viper) {
	err := o.URL.Set(v.GetString(configPrefix + suffixURL))
	if err != nil {
		panic("loki url parse error")
	}
	o.Timeout = v.GetDuration(configPrefix + suffixTimeout)
}
