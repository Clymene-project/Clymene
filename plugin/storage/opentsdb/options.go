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

package opentsdb

import (
	"flag"
	"github.com/spf13/viper"
	"time"
)

// The opentsdb factory was developed based on opentsdb's tcollector.

const (
	configPrefix = "opentsdb"

	suffixDryRun       = ".dry-run"
	suffixHost         = ".host"
	suffixMaxTags      = ".max-tags"
	suffixHttpPassword = ".http-password"
	suffixHttpUsername = ".http-username"

	suffixPort        = ".port"
	suffixHttp        = ".http"
	suffixHttpApiPath = ".http-api-path"
	suffixSSL         = ".ssl"
	suffixHosts       = ".hosts"
	suffixTimeout     = ".timeout"
	suffixMaxChunk    = ".max-chunk"

	defaultDryRun       = false
	defaultHost         = "localhost"
	defaultMaxTags      = 8
	defaultHttpPassword = ""
	defaultHttpUsername = ""
	defaultPort         = 4242
	defaultHttp         = false
	defaultHttpApiPath  = "api/put"
	defaultSSL          = false
	defaultHosts        = ""
	defaultTimeout      = 10 * time.Second
	defaultMaxChunk     = 512
)

type Options struct {
	dryRun bool

	maxTags int

	http bool
	port int

	httpPassword string
	httpUsername string
	httpApiPath  string
	host         string
	ssl          bool
	hosts        string
	timeout      time.Duration
	maxChunk     int
}

func (o *Options) AddFlags(flagSet *flag.FlagSet) {
	flagSet.Bool(
		configPrefix+suffixDryRun,
		defaultDryRun,
		"Don't actually send anything to the TSD, just print the datapoints.",
	)
	flagSet.String(
		configPrefix+suffixHost,
		defaultHost,
		"Hostname to use to connect to the TSD.",
	)
	flagSet.Int(
		configPrefix+suffixMaxTags,
		defaultMaxTags,
		"The maximum number of tags to send to our TSD Instances",
	)
	flagSet.String(
		configPrefix+suffixHttpPassword,
		defaultHttpPassword,
		"Password to use for HTTP Basic Auth when sending the data via HTTP",
	)
	flagSet.String(
		configPrefix+suffixHttpUsername,
		defaultHttpUsername,
		"Username to use for HTTP Basic Auth when sending the data via HTTP",
	)
	flagSet.Int(
		configPrefix+suffixPort,
		defaultPort,
		"Port to connect to the TSD instance on",
	)
	flagSet.Bool(
		configPrefix+suffixHttp,
		defaultHttp,
		"Send the data via the http interface (default 'false')",
	)
	flagSet.String(
		configPrefix+suffixHttpApiPath,
		defaultHttpApiPath,
		"URL path to use for HTTP requests to TSD.",
	)
	flagSet.Bool(
		configPrefix+suffixSSL,
		defaultSSL,
		"Enable SSL - used in conjunction with http (default 'false')",
	)
	flagSet.String(
		configPrefix+suffixHosts,
		defaultHosts,
		"List of host:port to connect to tsd's (comma separated)",
	)
	flagSet.Duration(
		configPrefix+suffixTimeout,
		defaultTimeout,
		"Time out when doing http insert(sec, default 10 sec)",
	)
	flagSet.Int(
		configPrefix+suffixMaxChunk,
		defaultMaxChunk,
		"The maximum request body size to support for incoming HTTP requests when chunking is enabled",
	)
}

func (o *Options) InitFromViper(v *viper.Viper) {
	o.dryRun = v.GetBool(configPrefix + suffixDryRun)

	o.http = v.GetBool(configPrefix + suffixHttp)

	// http
	o.httpPassword = v.GetString(configPrefix + suffixHttpPassword)
	o.httpUsername = v.GetString(configPrefix + suffixHttpUsername)
	o.ssl = v.GetBool(configPrefix + suffixSSL)
	o.httpApiPath = v.GetString(configPrefix + suffixHttpApiPath)
	o.timeout = v.GetDuration(configPrefix + suffixTimeout)
	o.maxChunk = v.GetInt(configPrefix + suffixMaxChunk)

	// common
	o.host = v.GetString(configPrefix + suffixHost)
	o.port = v.GetInt(configPrefix + suffixPort)
	o.hosts = v.GetString(configPrefix + suffixHosts)
	o.maxTags = v.GetInt(configPrefix + suffixMaxTags)
}
