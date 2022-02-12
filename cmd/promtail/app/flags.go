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

package app

import (
	"flag"
	"github.com/spf13/viper"
)

const (
	ConfigPrefix            = "clymene-promtail"
	suffixDryRun            = ".dry-run"
	suffixInspect           = ".inspect"
	suffixPrintConfigStdErr = ".print-config-stderr"

	// config file flag is "--config.file="
	suffixConfigFile = "config.file"

	defaultDryRun            = false
	defaultInspect           = false
	defaultPrintConfigStdErr = false
	defaultConfigFile        = "/etc/promtail/config.yml"
)

type Options struct {
	printConfig bool
	dryRun      bool
	configFile  string
	inspect     bool
}

func AddFlags(flags *flag.FlagSet) {
	flags.Bool(ConfigPrefix+suffixPrintConfigStdErr, defaultPrintConfigStdErr, "Dump the entire Loki config object to stderr")
	flags.Bool(ConfigPrefix+suffixDryRun, defaultDryRun, "Start Promtail but print entries instead of sending them to Loki.")
	flags.Bool(ConfigPrefix+suffixInspect, defaultInspect, "Allows for detailed inspection of pipeline stages")
	flags.String(suffixConfigFile, defaultConfigFile, "yaml file to load")
}

// InitFromViper initializes Options with properties retrieved from Viper.
func (o *Options) InitFromViper(v *viper.Viper) *Options {
	o.configFile = v.GetString(suffixConfigFile)
	o.dryRun = v.GetBool(ConfigPrefix + suffixDryRun)
	o.inspect = v.GetBool(ConfigPrefix + suffixInspect)
	o.printConfig = v.GetBool(ConfigPrefix + suffixPrintConfigStdErr)
	return o
}
