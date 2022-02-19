package config

import (
	"flag"
	"fmt"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/limit"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/positions"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/scrapeconfig"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/server"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/targets/file"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config for promtail, describing what files to watch.
type Config struct {
	ServerConfig    server.Config         `yaml:"server,omitempty"`
	PositionsConfig positions.Config      `yaml:"positions,omitempty"`
	ScrapeConfig    []scrapeconfig.Config `yaml:"scrape_configs,omitempty"`
	TargetConfig    file.Config           `yaml:"target_config,omitempty"`
	LimitConfig     limit.Config          `yaml:"limit_config,omitempty"`
}

// RegisterFlagsWithPrefix with prefix registers flags where every name is prefixed by
// prefix. If prefix is a non-empty string, prefix should end with a period.
func (c *Config) RegisterFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	c.ServerConfig.RegisterFlagsWithPrefix(prefix, f)
	c.PositionsConfig.RegisterFlagsWithPrefix(prefix, f)
	c.TargetConfig.RegisterFlagsWithPrefix(prefix, f)
	c.LimitConfig.RegisterFlagsWithPrefix(prefix, f)
}

// RegisterFlags registers flags.
func (c *Config) RegisterFlags(prefix string, f *flag.FlagSet) {
	c.RegisterFlagsWithPrefix(prefix, f)
}

func (c Config) String() string {
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

// LoadFile parses the given YAML file into a Config.
func LoadFile(filename string, logger *zap.Logger) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(string(content))
	if err != nil {
		return nil, errors.Wrapf(err, "parsing YAML file %s", filename)
	}

	return cfg, nil
}

// Load parses the YAML input s into a Config.
func Load(s string) (*Config, error) {
	cfg := &Config{}
	cfg.RegisterFlags("", flag.CommandLine)

	err := yaml.UnmarshalStrict([]byte(s), cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
