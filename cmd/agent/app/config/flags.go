package config

import (
	"flag"
	"github.com/Clymene-project/Clymene/ports"
	"github.com/spf13/viper"
)

const (
	configFile          = "config.file"
	defaultConfig       = "/etc/clymene/clymene.yml"
	httpPort            = "http.port"
	newSDManager        = "enable.new-service-discovery-manager"
	splitLength         = "metric.split.length"
	defaultSplitLength  = 1024
	defaultNewSDManager = true
)

// AddFlags adds flags for Options.
func (b *Builder) AddFlags(flags *flag.FlagSet) {
	flags.String(configFile, defaultConfig, "configuration file path.")
	flags.Int(httpPort, ports.AgentReloadHTTP, "http port")
	flags.Bool(newSDManager, defaultNewSDManager, "use new service discovery manager")
	flags.Int(splitLength, defaultSplitLength, "split length for metric transmission")
}

// InitFromViper initializes Options with properties retrieved from Viper.
func (b *Builder) InitFromViper(v *viper.Viper) *Builder {
	b.ConfigFile = v.GetString(configFile)
	b.HostPort = v.GetInt(httpPort)
	b.NewSDManager = v.GetBool(newSDManager)
	b.SplitLength = v.GetInt(splitLength)
	return b
}
