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

package dbmodel

import (
	"encoding/json"
	"github.com/Clymene-project/Clymene/prompb"
	"regexp"
)

type Converter struct {
	MaxTags int
}

type metricEntry map[string]interface{}

// Metrics and Tags
// The following rules apply to metric and tag values:
// 1. Strings are case sensitive, i.e. “Sys.Cpu.User” will be stored separately from “sys.cpu.user”
// 2. Spaces are not allowed
// 3. Only the following characters are allowed: a to z, A to Z, 0 to 9, -, _, ., / or Unicode letters (as per the specification)
// ref: http://opentsdb.net/docs/build/html/user_guide/writing/index.html

func (c *Converter) ConvertTsToOpenTSDBJSON(metrics []prompb.TimeSeries) ([]byte, error) {
	var openTSDBJson []metricEntry
	for _, metric := range metrics {
		entry := make(map[string]interface{})
		tags := make(map[string]interface{})
		for _, label := range metric.Labels {

			if label.Name == "__name__" {
				entry["metric"] = label.Value
				continue
			}
			if len(tags) < c.MaxTags {
				// Only the following characters are allowed: a to z, A to Z, 0 to 9, -, _, ., / or Unicode letters (as per the specification)
				match, _ := regexp.MatchString("[`~!@#$%^&*()|+=?;:'\", <>{}\\[\\]\\\\/ㄱ-ㅎ|ㅏ-ㅣ|가-힣]", label.Value)
				if !match {
					if label.Value != "" {
						tags[label.Name] = label.Value
					}
				}
			}
		}
		entry["tags"] = tags
		for _, sample := range metric.Samples {
			entry["timestamp"] = sample.Timestamp
			entry["value"] = sample.Value
		}
		openTSDBJson = append(openTSDBJson, entry)
	}
	return json.Marshal(openTSDBJson)
}
