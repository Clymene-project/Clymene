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

package http

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
	Encoding     string
	//TLS          tlscfg.Options
}

const (
	hTTPPrefix         = "gateway.http"
	suffixUrl          = ".url"
	suffixUserAgent    = ".user.agent"
	suffixTimeout      = ".timeout"
	suffixMaxErrMsgLen = ".max-err-msg-len"

	defaultClymeneGatewayUrl = "http://localhost:15611/api/metrics"
	defaultTimeout           = 10 * time.Second
	defaultMaxErrMsgLen      = 256
)

// AddFlags adds flags for Options.
func AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(
		hTTPPrefix+suffixUrl,
		defaultClymeneGatewayUrl,
		"the clymene-gateway remote write HTTP receiver endpoint(/api/metrics)",
	)
	flagSet.Duration(
		hTTPPrefix+suffixTimeout,
		defaultTimeout,
		"Time out when doing remote write(sec, default 10 sec)",
	)
	flagSet.String(
		hTTPPrefix+suffixUserAgent,
		fmt.Sprintf("Clymene/%s", version.Get().Version),
		"User-Agent in request header",
	)
	flagSet.Int(
		hTTPPrefix+suffixMaxErrMsgLen,
		defaultMaxErrMsgLen,
		"Maximum length of error message",
	)
	//tlsFlagsConfig.AddFlags(flagSet)
}

// InitFromViper initializes Options with properties retrieved from Viper.
func (o *Options) InitFromViper(v *viper.Viper) {
	o.url = v.GetString(hTTPPrefix + suffixUrl)
	o.maxErrMsgLen = v.GetInt64(hTTPPrefix + suffixMaxErrMsgLen)
	o.timeout = v.GetDuration(hTTPPrefix + suffixTimeout)
	o.userAgent = v.GetString(hTTPPrefix + suffixUserAgent)
	//o.TLS = tlsFlagsConfig.InitFromViper(v)
}
