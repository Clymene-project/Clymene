// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"github.com/Clymene-project/Clymene/cmd/agent/app/discovery"
	"github.com/Clymene-project/Clymene/cmd/agent/app/model/labels"
	"github.com/Clymene-project/Clymene/cmd/agent/app/relabel"
	"go.uber.org/zap"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/alecthomas/units"
	"github.com/pkg/errors"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	yaml "gopkg.in/yaml.v2"
)

var (
	patRulePath = regexp.MustCompile(`^[^*]*(\*[^/]*)?$`)
)

// Load parses the YAML input s into a Config.
func Load(s string, logger *zap.Logger) (*Config, error) {
	cfg := &Config{}
	// If the entire config body is empty the UnmarshalYAML method is
	// never called. We thus have to set the DefaultConfig at the entry
	// point as well.
	*cfg = DefaultConfig

	err := yaml.UnmarshalStrict([]byte(s), cfg)
	if err != nil {
		return nil, err
	}

	for i, v := range cfg.GlobalConfig.ExternalLabels {
		newV := os.Expand(v.Value, func(s string) string {
			if s == "$" {
				return "$"
			}
			if v := os.Getenv(s); v != "" {
				return v
			}
			logger.Warn("msg Empty environment variable", zap.String("name", s))
			return ""
		})
		if newV != v.Value {
			logger.Debug("msg External label replaced", zap.String("label", v.Name), zap.String("input", v.Value), zap.String("output", newV))
			v.Value = newV
			cfg.GlobalConfig.ExternalLabels[i] = v
		}
	}
	return cfg, nil
}

// LoadFile parses the given YAML file into a Config.
func LoadFile(filename string, logger *zap.Logger) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(string(content), logger)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing YAML file %s", filename)
	}

	cfg.SetDirectory(filepath.Dir(filename))
	return cfg, nil
}

// The defaults applied before parsing the respective config sections.
var (
	// DefaultConfig is the default top-level configuration.
	DefaultConfig = Config{
		GlobalConfig: DefaultGlobalConfig,
	}

	// DefaultGlobalConfig is the default global configuration.
	DefaultGlobalConfig = GlobalConfig{
		ScrapeInterval:     model.Duration(1 * time.Minute),
		ScrapeTimeout:      model.Duration(10 * time.Second),
		EvaluationInterval: model.Duration(1 * time.Minute),
	}

	// DefaultScrapeConfig is the default scrape configuration.
	DefaultScrapeConfig = ScrapeConfig{
		// ScrapeTimeout and ScrapeInterval default to the
		// configured globals.
		MetricsPath:      "/metrics",
		Scheme:           "http",
		HonorLabels:      false,
		HonorTimestamps:  true,
		HTTPClientConfig: config.DefaultHTTPClientConfig,
	}
)

// Config is the top-level configuration for Prometheus's config files.
type Config struct {
	GlobalConfig  GlobalConfig    `yaml:"global"`
	RuleFiles     []string        `yaml:"rule_files,omitempty"`
	ScrapeConfigs []*ScrapeConfig `yaml:"scrape_configs,omitempty"`
}

// SetDirectory joins any relative file paths with dir.
func (c *Config) SetDirectory(dir string) {
	c.GlobalConfig.SetDirectory(dir)
	for i, file := range c.RuleFiles {
		c.RuleFiles[i] = config.JoinDir(dir, file)
	}
	for _, c := range c.ScrapeConfigs {
		c.SetDirectory(dir)
	}
}

func (c Config) String() string {
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultConfig
	// We want to set c to the defaults and then overwrite it with the input.
	// To make unmarshal fill the plain data struct rather than calling UnmarshalYAML
	// again, we have to hide it using a type indirection.
	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	// If a global block was open but empty the default global config is overwritten.
	// We have to restore it here.
	if c.GlobalConfig.isZero() {
		c.GlobalConfig = DefaultGlobalConfig
	}

	for _, rf := range c.RuleFiles {
		if !patRulePath.MatchString(rf) {
			return errors.Errorf("invalid rule file path %q", rf)
		}
	}
	// Do global overrides and validate unique names.
	jobNames := map[string]struct{}{}
	for _, scfg := range c.ScrapeConfigs {
		if scfg == nil {
			return errors.New("empty or null scrape config section")
		}
		// First set the correct scrape interval, then check that the timeout
		// (inferred or explicit) is not greater than that.
		if scfg.ScrapeInterval == 0 {
			scfg.ScrapeInterval = c.GlobalConfig.ScrapeInterval
		}
		if scfg.ScrapeTimeout > scfg.ScrapeInterval {
			return errors.Errorf("scrape timeout greater than scrape interval for scrape config with job name %q", scfg.JobName)
		}
		if scfg.ScrapeTimeout == 0 {
			if c.GlobalConfig.ScrapeTimeout > scfg.ScrapeInterval {
				scfg.ScrapeTimeout = scfg.ScrapeInterval
			} else {
				scfg.ScrapeTimeout = c.GlobalConfig.ScrapeTimeout
			}
		}

		if _, ok := jobNames[scfg.JobName]; ok {
			return errors.Errorf("found multiple scrape configs with job name %q", scfg.JobName)
		}
		jobNames[scfg.JobName] = struct{}{}
	}

	return nil
}

// GlobalConfig configures values that are used across other configuration
// objects.
type GlobalConfig struct {
	// How frequently to scrape targets by default.
	ScrapeInterval model.Duration `yaml:"scrape_interval,omitempty"`
	// The default timeout when scraping targets.
	ScrapeTimeout model.Duration `yaml:"scrape_timeout,omitempty"`
	// How frequently to evaluate rules by default.
	EvaluationInterval model.Duration `yaml:"evaluation_interval,omitempty"`
	// File to which PromQL queries are logged.
	QueryLogFile string `yaml:"query_log_file,omitempty"`
	// The labels to add to any timeseries that this Prometheus instance scrapes.
	ExternalLabels labels.Labels `yaml:"external_labels,omitempty"`
}

// SetDirectory joins any relative file paths with dir.
func (c *GlobalConfig) SetDirectory(dir string) {
	c.QueryLogFile = config.JoinDir(dir, c.QueryLogFile)
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *GlobalConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Create a clean global config as the previous one was already populated
	// by the default due to the YAML parser behavior for empty blocks.
	gc := &GlobalConfig{}
	type plain GlobalConfig
	if err := unmarshal((*plain)(gc)); err != nil {
		return err
	}

	for _, l := range gc.ExternalLabels {
		if !model.LabelName(l.Name).IsValid() {
			return errors.Errorf("%q is not a valid label name", l.Name)
		}
		if !model.LabelValue(l.Value).IsValid() {
			return errors.Errorf("%q is not a valid label value", l.Value)
		}
	}

	// First set the correct scrape interval, then check that the timeout
	// (inferred or explicit) is not greater than that.
	if gc.ScrapeInterval == 0 {
		gc.ScrapeInterval = DefaultGlobalConfig.ScrapeInterval
	}
	if gc.ScrapeTimeout > gc.ScrapeInterval {
		return errors.New("global scrape timeout greater than scrape interval")
	}
	if gc.ScrapeTimeout == 0 {
		if DefaultGlobalConfig.ScrapeTimeout > gc.ScrapeInterval {
			gc.ScrapeTimeout = gc.ScrapeInterval
		} else {
			gc.ScrapeTimeout = DefaultGlobalConfig.ScrapeTimeout
		}
	}
	if gc.EvaluationInterval == 0 {
		gc.EvaluationInterval = DefaultGlobalConfig.EvaluationInterval
	}
	*c = *gc
	return nil
}

// isZero returns true iff the global config is the zero value.
func (c *GlobalConfig) isZero() bool {
	return c.ExternalLabels == nil &&
		c.ScrapeInterval == 0 &&
		c.ScrapeTimeout == 0 &&
		c.EvaluationInterval == 0 &&
		c.QueryLogFile == ""
}

// ScrapeConfig configures a scraping unit for Prometheus.
type ScrapeConfig struct {
	// The job name to which the job label is set by default.
	JobName string `yaml:"job_name"`
	// Indicator whether the scraped metrics should remain unmodified.
	HonorLabels bool `yaml:"honor_labels,omitempty"`
	// Indicator whether the scraped timestamps should be respected.
	HonorTimestamps bool `yaml:"honor_timestamps"`
	// A set of query parameters with which the target is scraped.
	Params url.Values `yaml:"params,omitempty"`
	// How frequently to scrape the targets of this scrape config.
	ScrapeInterval model.Duration `yaml:"scrape_interval,omitempty"`
	// The timeout for scraping targets of this config.
	ScrapeTimeout model.Duration `yaml:"scrape_timeout,omitempty"`
	// The HTTP resource path on which to fetch metrics from targets.
	MetricsPath string `yaml:"metrics_path,omitempty"`
	// The URL scheme with which to fetch metrics from targets.
	Scheme string `yaml:"scheme,omitempty"`
	// An uncompressed response body larger than this many bytes will cause the
	// scrape to fail. 0 means no limit.
	BodySizeLimit units.Base2Bytes `yaml:"body_size_limit,omitempty"`
	// More than this many samples post metric-relabeling will cause the scrape to
	// fail.
	SampleLimit uint `yaml:"sample_limit,omitempty"`
	// More than this many targets after the target relabeling will cause the
	// scrapes to fail.
	TargetLimit uint `yaml:"target_limit,omitempty"`
	// More than this many labels post metric-relabeling will cause the scrape to
	// fail.
	LabelLimit uint `yaml:"label_limit,omitempty"`
	// More than this label name length post metric-relabeling will cause the
	// scrape to fail.
	LabelNameLengthLimit uint `yaml:"label_name_length_limit,omitempty"`
	// More than this label value length post metric-relabeling will cause the
	// scrape to fail.
	LabelValueLengthLimit uint `yaml:"label_value_length_limit,omitempty"`

	// We cannot do proper Go type embedding below as the parser will then parse
	// values arbitrarily into the overflow maps of further-down types.

	ServiceDiscoveryConfigs discovery.Configs       `yaml:"-"`
	HTTPClientConfig        config.HTTPClientConfig `yaml:",inline"`

	// List of target relabel configurations.
	RelabelConfigs []*relabel.Config `yaml:"relabel_configs,omitempty"`
	// List of metric relabel configurations.
	MetricRelabelConfigs []*relabel.Config `yaml:"metric_relabel_configs,omitempty"`
}

// SetDirectory joins any relative file paths with dir.
func (c *ScrapeConfig) SetDirectory(dir string) {
	c.ServiceDiscoveryConfigs.SetDirectory(dir)
	c.HTTPClientConfig.SetDirectory(dir)
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *ScrapeConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultScrapeConfig
	if err := discovery.UnmarshalYAMLWithInlineConfigs(c, unmarshal); err != nil {
		return err
	}
	if len(c.JobName) == 0 {
		return errors.New("job_name is empty")
	}

	// The UnmarshalYAML method of HTTPClientConfig is not being called because it's not a pointer.
	// We cannot make it a pointer as the parser panics for inlined pointer structs.
	// Thus we just do its validation here.
	if err := c.HTTPClientConfig.Validate(); err != nil {
		return err
	}

	// Check for users putting URLs in target groups.
	if len(c.RelabelConfigs) == 0 {
		if err := checkStaticTargets(c.ServiceDiscoveryConfigs); err != nil {
			return err
		}
	}

	for _, rlcfg := range c.RelabelConfigs {
		if rlcfg == nil {
			return errors.New("empty or null target relabeling rule in scrape config")
		}
	}
	for _, rlcfg := range c.MetricRelabelConfigs {
		if rlcfg == nil {
			return errors.New("empty or null metric relabeling rule in scrape config")
		}
	}

	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (c *ScrapeConfig) MarshalYAML() (interface{}, error) {
	return discovery.MarshalYAMLWithInlineConfigs(c)
}

// StorageConfig configures runtime reloadable configuration options.
type StorageConfig struct {
	ExemplarsConfig *ExemplarsConfig `yaml:"exemplars,omitempty"`
}

type TracingClientType string

const (
	TracingClientHTTP TracingClientType = "http"
	TracingClientGRPC TracingClientType = "grpc"
)

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (t *TracingClientType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*t = TracingClientType("")
	type plain TracingClientType
	if err := unmarshal((*plain)(t)); err != nil {
		return err
	}

	if *t != TracingClientHTTP && *t != TracingClientGRPC {
		return fmt.Errorf("expected tracing client type to be to be %s or %s, but got %s",
			TracingClientHTTP, TracingClientGRPC, *t,
		)
	}

	return nil
}

// TracingConfig configures the tracing options.
type TracingConfig struct {
	ClientType       TracingClientType `yaml:"client_type,omitempty"`
	Endpoint         string            `yaml:"endpoint,omitempty"`
	SamplingFraction float64           `yaml:"sampling_fraction,omitempty"`
	WithSecure       bool              `yaml:"with_secure,omitempty"`
	TLSConfig        config.TLSConfig  `yaml:"tls_config,omitempty"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (t *TracingConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*t = TracingConfig{}
	type plain TracingConfig
	if err := unmarshal((*plain)(t)); err != nil {
		return err
	}

	if t.Endpoint == "" {
		return errors.New("tracing endpoint must be set")
	}

	// Fill in gRPC client as default if none is set.
	if t.ClientType == "" {
		t.ClientType = TracingClientGRPC
	}

	return nil
}

// ExemplarsConfig configures runtime reloadable configuration options.
type ExemplarsConfig struct {
	// MaxExemplars sets the size, in # of exemplars stored, of the single circular buffer used to store exemplars in memory.
	// Use a value of 0 or less than 0 to disable the storage without having to restart Prometheus.
	MaxExemplars int64 `yaml:"max_exemplars,omitempty"`
}

func checkStaticTargets(configs discovery.Configs) error {
	for _, cfg := range configs {
		sc, ok := cfg.(discovery.StaticConfig)
		if !ok {
			continue
		}
		for _, tg := range sc {
			for _, t := range tg.Targets {
				if err := CheckTargetAddress(t[model.AddressLabel]); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// CheckTargetAddress checks if target address is valid.
func CheckTargetAddress(address model.LabelValue) error {
	// For now check for a URL, we may want to expand this later.
	if strings.Contains(string(address), "/") {
		return errors.Errorf("%q is not a valid hostname", address)
	}
	return nil
}

type Builder struct {
	ConfigFile   string
	HostPort     int
	NewSDManager bool
}

// NewConfigBuilder creates a new configfile builder.
func NewConfigBuilder() *Builder {
	return &Builder{}
}

func ReloadConfig(filename string, logger *zap.Logger, rls ...func(*Config) error) (err error) {
	logger.Info("Loading configuration file", zap.String("filename", filename))
	conf, err := LoadFile(filename, logger)
	if err != nil {
		return errors.Wrapf(err, "couldn't load configuration (--config.file=%q)", filename)
	}

	failed := false
	for _, rl := range rls {
		if err := rl(conf); err != nil {
			logger.Error("Failed to apply configuration", zap.Error(err))
			failed = true
		}
	}
	if failed {
		return errors.Errorf("one or more errors occurred while applying the new configuration (--config.file=%q)", filename)
	}

	logger.Info("Completed loading of configuration file", zap.String("filename", filename))
	return nil
}
