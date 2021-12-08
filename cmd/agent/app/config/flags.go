package config

import (
	"flag"
	"github.com/Clymene-project/Clymene/ports"
	"github.com/spf13/viper"
)

const (
	configFile    = "config.file"
	defaultConfig = "/etc/clymene/clymene.yml"
	httpPort      = "http.port"
)

// AddFlags adds flags for Options.
func (b *Builder) AddFlags(flags *flag.FlagSet) {
	flags.String(configFile, defaultConfig, "configuration file path.")
	flags.Int(httpPort, ports.AgentReloadHTTP, "http port")
}

// InitFromViper initializes Options with properties retrieved from Viper.
func (b *Builder) InitFromViper(v *viper.Viper) *Builder {
	b.ConfigFile = v.GetString(configFile)
	b.HostPort = v.GetInt(httpPort)
	return b
}
