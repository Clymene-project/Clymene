package stages

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"reflect"
	"regexp"
	"text/template"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
)

// Config Errors
const (
	ErrEmptyReplaceStageConfig = "empty replace stage configuration"
	ErrEmptyReplaceStageSource = "empty source in replace stage"
)

// ReplaceConfig contains a regexStage configuration
type ReplaceConfig struct {
	Expression string  `mapstructure:"expression"`
	Source     *string `mapstructure:"source"`
	Replace    string  `mapstructure:"replace"`
}

// validateReplaceConfig validates the config and return a regex
func validateReplaceConfig(c *ReplaceConfig) (*regexp.Regexp, error) {
	if c == nil {
		return nil, errors.New(ErrEmptyReplaceStageConfig)
	}

	if c.Expression == "" {
		return nil, errors.New(ErrExpressionRequired)
	}

	if c.Source != nil && *c.Source == "" {
		return nil, errors.New(ErrEmptyReplaceStageSource)
	}

	expr, err := regexp.Compile(c.Expression)
	if err != nil {
		return nil, errors.Wrap(err, ErrCouldNotCompileRegex)
	}
	return expr, nil
}

// replaceStage sets extracted data using regular expressions
type replaceStage struct {
	cfg        *ReplaceConfig
	expression *regexp.Regexp
	logger     *zap.Logger
}

// newReplaceStage creates a newReplaceStage
func newReplaceStage(logger *zap.Logger, config interface{}) (Stage, error) {
	cfg, err := parseReplaceConfig(config)
	if err != nil {
		return nil, err
	}
	expression, err := validateReplaceConfig(cfg)
	if err != nil {
		return nil, err
	}

	return toStage(&replaceStage{
		cfg:        cfg,
		expression: expression,
		logger:     logger.With(zap.String("component", "stage"), zap.String("type", "replace")),
	}), nil
}

// parseReplaceConfig processes an incoming configuration into a ReplaceConfig
func parseReplaceConfig(config interface{}) (*ReplaceConfig, error) {
	cfg := &ReplaceConfig{}
	err := mapstructure.Decode(config, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Process implements Stage
func (r *replaceStage) Process(labels model.LabelSet, extracted map[string]interface{}, t *time.Time, entry *string) {
	// If a source key is provided, the replace stage should process it
	// from the extracted map, otherwise should fallback to the entry
	input := entry

	if r.cfg.Source != nil {
		if _, ok := extracted[*r.cfg.Source]; !ok {
			if Debug {
				r.logger.Debug("source does not exist in the set of extracted values", zap.String("source", *r.cfg.Source))
			}
			return
		}

		value, err := getString(extracted[*r.cfg.Source])
		if err != nil {
			if Debug {
				r.logger.Debug("failed to convert source value to string", zap.String("source", *r.cfg.Source), zap.Error(err),
					zap.String("type", reflect.TypeOf(extracted[*r.cfg.Source]).String()))
			}
			return
		}

		input = &value
	}

	if input == nil {
		if Debug {
			r.logger.Debug("cannot parse a nil entry")
		}
		return
	}

	// Get string of matched captured groups. We will use this to extract all named captured groups
	match := r.expression.FindStringSubmatch(*input)
	matchAllIndex := r.expression.FindAllStringSubmatchIndex(*input, -1)

	if matchAllIndex == nil {
		if Debug {
			r.logger.Debug("regex did not match", zap.String("input", *input), zap.String("regex", r.expression.String()))
		}
		return
	}

	// All extracted values will be available for templating
	td := r.getTemplateData(extracted)

	// Initialize the template with the "replace" string defined by user
	templ, err := template.New("pipeline_template").Funcs(functionMap).Parse(r.cfg.Replace)
	if err != nil {
		if Debug {
			r.logger.Debug("template initialization error", zap.Error(err))
		}
		return
	}

	result, capturedMap, err := r.getReplacedEntry(matchAllIndex, *input, td, templ)
	if err != nil {
		if Debug {
			r.logger.Debug("failed to execute template on extracted value", zap.Error(err))
		}
		return
	}

	if r.cfg.Source != nil {
		extracted[*r.cfg.Source] = result
	} else {
		*entry = result
	}

	// All the named captured group will be extracted
	for i, name := range r.expression.SubexpNames() {
		if i != 0 && name != "" {
			if v, ok := capturedMap[match[i]]; ok {
				extracted[name] = v
			}
		}
	}
	if Debug {
		r.logger.Debug("extracted data debug in replace stage", zap.String("extracted data", fmt.Sprintf("%v", extracted)))
	}
}

func (r *replaceStage) getReplacedEntry(matchAllIndex [][]int, input string, td map[string]string, templ *template.Template) (string, map[string]string, error) {
	var result string
	previousInputEndIndex := 0
	capturedMap := make(map[string]string)
	// For a simple string like `11.11.11.11 - frank 12.12.12.12 - frank`
	// if the regex is "(\\d{2}.\\d{2}.\\d{2}.\\d{2}) - (\\S+)"
	// FindAllStringSubmatchIndex would return [[0 19 0 11 14 19] [20 37 20 31 34 37]].
	// Each inner array's first two values will be the start and end index of the entire
	// matched string and the next values will be start and end index of the matched
	// captured group. Here 0-19 is "11.11.11.11 - frank",  0-11 is "11.11.11.11" and
	// 14-19 is "frank". So, we advance by 2 index to get the next match
	for _, matchIndex := range matchAllIndex {
		for i := 2; i < len(matchIndex); i += 2 {
			if matchIndex[i] == -1 {
				continue
			}
			capturedString := input[matchIndex[i]:matchIndex[i+1]]
			buf := &bytes.Buffer{}
			td["Value"] = capturedString
			err := templ.Execute(buf, td)
			if err != nil {
				return "", nil, err
			}
			st := buf.String()
			if previousInputEndIndex == 0 || previousInputEndIndex <= matchIndex[i] {
				result += input[previousInputEndIndex:matchIndex[i]] + st
				previousInputEndIndex = matchIndex[i+1]
			}
			capturedMap[capturedString] = st
		}
	}
	return result + input[previousInputEndIndex:], capturedMap, nil
}

func (r *replaceStage) getTemplateData(extracted map[string]interface{}) map[string]string {
	td := make(map[string]string)
	for k, v := range extracted {
		s, err := getString(v)
		if err != nil {
			if Debug {
				r.logger.Debug("extracted template could not be converted to a string", zap.Error(err), zap.String("type", reflect.TypeOf(v).String()))
			}
			continue
		}
		td[k] = s
	}
	return td
}

// Name implements Stage
func (r *replaceStage) Name() string {
	return StageTypeReplace
}
