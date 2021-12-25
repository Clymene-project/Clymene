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
)

// The opentsdb factory was developed based on opentsdb's tcollector.

const (
	configPrefix                = "opentsdb"
	suffixDryRun                = ".dry-run"
	suffixHost                  = ".host"
	suffixEvictInterval         = ".evict-interval"
	suffixDeDupInterval         = ".dedup-interval"
	suffixDeDupOnlyZero         = ".dedup-only-zero"
	suffixAllowedInactivityTime = ".allowed-inactivity-time"
	suffixMaxTags               = ".max-tags"
	suffixHttpPassword          = ".http-password"
	suffixHttpUsername          = ".http-username"
	suffixReconnectInterval     = ".reconnect-interval"
	suffixPort                  = ".port"
	suffixHttp                  = ".http"
	suffixHttpApiPath           = ".http-api-path"
	suffixSSL                   = ".ssl"
	suffixStdin                 = ".stdin"
	suffixHosts                 = ".hosts"
	suffixNamespacePrefix       = ".namespace-prefix"

	defaultDryRun                = false
	defaultEvictInterval         = 6000
	defaultDeDupInterval         = 300
	defaultHost                  = "localhost"
	defaultDeDupOnlyZero         = false
	defaultAllowedInactivityTime = 600
	defaultMaxTags               = 8
	defaultHttpPassword          = false
	defaultHttpUsername          = false
	defaultReconnectInterval     = 0
	defaultPort                  = 4242
	defaultHttp                  = false
	defaultHttpApiPath           = "api/put"
	defaultSSL                   = false
	defaultStdin                 = false
	defaultHosts                 = ""
	defaultNamespacePrefix       = ""
)

type Options struct {
	evictInterval         int
	deDupInterval         int
	deDupOnlyZero         bool
	allowedInactivityTime int
	dryRun                bool
	maxTags               int
	httpPassword          bool
	httpUsername          bool
	reconnectInterval     int
	port                  int
	http                  bool
	httpApiPath           string
	host                  string
	ssl                   bool
	stdin                 bool
	hosts                 string
	namespacePrefix       string
}

func (o *Options) AddFlags(flagSet *flag.FlagSet) {
	flagSet.Bool(
		configPrefix+suffixDryRun,
		defaultDryRun,
		"Don't actually send anything to the TSD, just print the datapoints.",
	)
	flagSet.Int(
		configPrefix+suffixEvictInterval,
		defaultEvictInterval,
		"Number of seconds after which to remove cached values of old data points to save memory",
	)
	flagSet.String(
		configPrefix+suffixHost,
		defaultHost,
		"Hostname to use to connect to the TSD.",
	)
	flagSet.Int(
		configPrefix+suffixDeDupInterval,
		defaultDeDupInterval,
		"Number of seconds in which successive duplicate datapoints are suppressed before sending to the TSD. Use zero to disable",
	)
	flagSet.Bool(
		configPrefix+suffixDeDupOnlyZero,
		defaultDeDupOnlyZero,
		"Only dedup 0 values.",
	)
	flagSet.Int(
		configPrefix+suffixAllowedInactivityTime,
		defaultAllowedInactivityTime,
		"How long to wait for datapoints before assuming a collector is dead and restart it",
	)
	flagSet.Int(
		configPrefix+suffixMaxTags,
		defaultMaxTags,
		"The maximum number of tags to send to our TSD Instances",
	)
	flagSet.Bool(
		configPrefix+suffixHttpPassword,
		defaultHttpPassword,
		"Password to use for HTTP Basic Auth when sending the data via HTTP",
	)
	flagSet.Bool(
		configPrefix+suffixHttpUsername,
		defaultHttpUsername,
		"Username to use for HTTP Basic Auth when sending the data via HTTP",
	)
	flagSet.Int(
		configPrefix+suffixReconnectInterval,
		defaultReconnectInterval,
		"Number of seconds after which the connection to the TSD hostname reconnects itself. This is useful when the hostname is a multiple A record (RRDNS)",
	)
	flagSet.Int(
		configPrefix+suffixPort,
		defaultPort,
		"Port to connect to the TSD instance on",
	)
	flagSet.Bool(
		configPrefix+suffixHttp,
		defaultHttp,
		"Send the data via the http interface",
	)
	flagSet.String(
		configPrefix+suffixHttpApiPath,
		defaultHttpApiPath,
		"URL path to use for HTTP requests to TSD.",
	)
	flagSet.Bool(
		configPrefix+suffixSSL,
		defaultSSL,
		"Enable SSL - used in conjunction with http",
	)
	flagSet.Bool(
		configPrefix+suffixStdin,
		defaultStdin,
		"Run once, read and dedup data points from stdin",
	)
	flagSet.String(
		configPrefix+suffixHosts,
		defaultHosts,
		"List of host:port to connect to tsd's (comma separated)",
	)
	flagSet.String(
		configPrefix+suffixNamespacePrefix,
		defaultNamespacePrefix,
		"Prefix to prepend to all metric names collected",
	)

}

func (o *Options) InitFromViper(v *viper.Viper) {
	o.dryRun = v.GetBool(configPrefix + suffixDryRun)
	o.evictInterval = v.GetInt(configPrefix + suffixEvictInterval)
	o.deDupInterval = v.GetInt(configPrefix + suffixDeDupInterval)
	o.deDupOnlyZero = v.GetBool(configPrefix + suffixDeDupOnlyZero)
	o.allowedInactivityTime = v.GetInt(configPrefix + suffixAllowedInactivityTime)
	o.maxTags = v.GetInt(configPrefix + suffixMaxTags)
	o.httpPassword = v.GetBool(configPrefix + suffixHttpPassword)
	o.httpUsername = v.GetBool(configPrefix + suffixHttpUsername)
	o.reconnectInterval = v.GetInt(configPrefix + suffixReconnectInterval)
	o.port = v.GetInt(configPrefix + suffixPort)
	o.http = v.GetBool(configPrefix + suffixHttp)
	o.httpApiPath = v.GetString(configPrefix + suffixHttpApiPath)
	o.host = v.GetString(configPrefix + suffixHost)
	o.ssl = v.GetBool(configPrefix + suffixSSL)
	o.stdin = v.GetBool(configPrefix + suffixStdin)
	o.hosts = v.GetString(configPrefix + suffixHosts)
	o.namespacePrefix = v.GetString(configPrefix + suffixNamespacePrefix)
}
