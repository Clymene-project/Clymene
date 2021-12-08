// Copyright (c) 2021 The Jaeger Authors.
//
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

package mappings

import (
	"bytes"
	"embed"
	"github.com/Clymene-project/Clymene/pkg/es"
	"strings"
)

//go:embed *.json
// MAPPINGS contains embedded index templates.
var MAPPINGS embed.FS

// MappingBuilder holds parameters required to render an elasticsearch index template
type MappingBuilder struct {
	TemplateBuilder es.TemplateBuilder
	Shards          int64
	Replicas        int64
	EsVersion       uint
	IndexPrefix     string
	UseILM          bool
	ILMPolicyName   string
}

// GetMapping returns the rendered mapping based on elasticsearch version
func (mb *MappingBuilder) GetMapping(mapping string) (string, error) {
	if mb.EsVersion == 7 {
		return mb.fixMapping(mapping + "-7.json")
	}
	return mb.fixMapping(mapping + ".json")
}

// GetMetricsMappings returns Metrics mappings
func (mb *MappingBuilder) GetMetricsMappings() (string, error) {
	metricMapping, err := mb.GetMapping("clymene-metrics")
	if err != nil {
		return "", err
	}
	return metricMapping, nil
}

func loadMapping(name string) string {
	s, _ := MAPPINGS.ReadFile(name)
	return string(s)
}

func (mb *MappingBuilder) fixMapping(mapping string) (string, error) {

	tmpl, err := mb.TemplateBuilder.Parse(loadMapping(mapping))
	if err != nil {
		return "", err
	}
	writer := new(bytes.Buffer)

	if mb.IndexPrefix != "" && !strings.HasSuffix(mb.IndexPrefix, "-") {
		mb.IndexPrefix += "-"
	}
	if err := tmpl.Execute(writer, mb); err != nil {
		return "", err
	}

	return writer.String(), nil
}
