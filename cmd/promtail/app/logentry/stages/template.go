package stages

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"go.uber.org/zap"
	"reflect"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/prometheus/common/model"

	"golang.org/x/crypto/sha3"
)

// Config Errors
const (
	ErrEmptyTemplateStageConfig = "template stage config cannot be empty"
	ErrTemplateSourceRequired   = "template source value is required"
)

var extraFunctionMap = template.FuncMap{
	"ToLower":    strings.ToLower,
	"ToUpper":    strings.ToUpper,
	"Replace":    strings.Replace,
	"Trim":       strings.Trim,
	"TrimLeft":   strings.TrimLeft,
	"TrimRight":  strings.TrimRight,
	"TrimPrefix": strings.TrimPrefix,
	"TrimSuffix": strings.TrimSuffix,
	"TrimSpace":  strings.TrimSpace,
	"Hash": func(salt string, input string) string {
		hash := sha3.Sum256([]byte(salt + input))
		return hex.EncodeToString(hash[:])
	},
	"Sha2Hash": func(salt string, input string) string {
		hash := sha256.Sum256([]byte(salt + input))
		return hex.EncodeToString(hash[:])
	},
	"regexReplaceAll": func(regex string, s string, repl string) string {
		r := regexp.MustCompile(regex)
		return r.ReplaceAllString(s, repl)
	},
	"regexReplaceAllLiteral": func(regex string, s string, repl string) string {
		r := regexp.MustCompile(regex)
		return r.ReplaceAllLiteralString(s, repl)
	},
}

var functionMap = sprig.TxtFuncMap()

func init() {
	for k, v := range extraFunctionMap {
		functionMap[k] = v
	}
}

// TemplateConfig configures template value extraction
type TemplateConfig struct {
	Source   string `mapstructure:"source"`
	Template string `mapstructure:"template"`
}

// validateTemplateConfig validates the templateStage config
func validateTemplateConfig(cfg *TemplateConfig) (*template.Template, error) {
	if cfg == nil {
		return nil, errors.New(ErrEmptyTemplateStageConfig)
	}
	if cfg.Source == "" {
		return nil, errors.New(ErrTemplateSourceRequired)
	}

	return template.New("pipeline_template").Funcs(functionMap).Parse(cfg.Template)
}

// newTemplateStage creates a new templateStage
func newTemplateStage(logger *zap.Logger, config interface{}) (Stage, error) {
	cfg := &TemplateConfig{}
	err := mapstructure.Decode(config, cfg)
	if err != nil {
		return nil, err
	}
	t, err := validateTemplateConfig(cfg)
	if err != nil {
		return nil, err
	}

	return toStage(&templateStage{
		cfgs:     cfg,
		logger:   logger,
		template: t,
	}), nil
}

// templateStage will mutate the incoming entry and set it from extracted data
type templateStage struct {
	cfgs     *TemplateConfig
	logger   *zap.Logger
	template *template.Template
}

// Process implements Stage
func (o *templateStage) Process(labels model.LabelSet, extracted map[string]interface{}, t *time.Time, entry *string) {
	if o.cfgs == nil {
		return
	}
	td := make(map[string]interface{})
	for k, v := range extracted {
		s, err := getString(v)
		if err != nil {
			if Debug {
				o.logger.Debug("extracted template could not be converted to a string", zap.Error(err), zap.String("type", reflect.TypeOf(v).String()))
			}
			continue
		}
		td[k] = s
		if k == o.cfgs.Source {
			td["Value"] = s
		}
	}
	td["Entry"] = *entry

	buf := &bytes.Buffer{}
	err := o.template.Execute(buf, td)
	if err != nil {
		if Debug {
			o.logger.Debug("failed to execute template on extracted value", zap.Error(err))
		}
		return
	}
	st := buf.String()
	// If the template evaluates to an empty string, remove the key from the map
	if st == "" {
		delete(extracted, o.cfgs.Source)
	} else {
		extracted[o.cfgs.Source] = st
	}
}

// Name implements Stage
func (o *templateStage) Name() string {
	return StageTypeTemplate
}
