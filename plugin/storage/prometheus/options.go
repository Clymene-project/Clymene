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

package prometheus

import (
	"flag"
	"fmt"
	"github.com/Clymene-project/Clymene/pkg/version"
	"github.com/spf13/viper"
	"time"
)

type Options struct {
	url          string
	userAgent    string
	timeout      time.Duration
	maxErrMsgLen int64
}

const (
	configPrefix       = "prometheus.remote"
	suffixUrl          = ".url"
	suffixUserAgent    = ".user.agent"
	suffixTimeout      = ".timeout"
	suffixmaxErrMsgLen = ".max.err.msg.len"

	defaultPrometheusUrl = "http://localhost:9090/api/v1/write"
	defaultTimeout       = 10 * time.Second
	defaultMaxErrMsgLen  = 256
)

func (o *Options) AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(
		configPrefix+suffixUrl,
		defaultPrometheusUrl,
		"the prometheus remote write receiver endpoint(/api/v1/write)",
	)
	flagSet.Duration(
		configPrefix+suffixTimeout,
		defaultTimeout,
		"Time out when doing remote write(sec, default 10 sec)",
	)
	flagSet.String(
		configPrefix+suffixUserAgent,
		fmt.Sprintf("Clymene/%s", version.Get().Version),
		"User-Agent in request header",
	)
	flagSet.Int(
		configPrefix+suffixmaxErrMsgLen,
		defaultMaxErrMsgLen,
		"Maximum length of error message",
	)
}

func (o *Options) InitFromViper(v *viper.Viper) {
	o.url = v.GetString(configPrefix + suffixUrl)
	o.maxErrMsgLen = v.GetInt64(configPrefix + suffixmaxErrMsgLen)
	o.timeout = v.GetDuration(configPrefix + suffixTimeout)
	o.userAgent = v.GetString(configPrefix + suffixUserAgent)
}
