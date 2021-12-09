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

import "github.com/Clymene-project/Clymene/prompb"

type Converter struct{}

// ConvertTsToJSON Change TimeSeries data to json
func (c *Converter) ConvertTsToJSON(metric prompb.TimeSeries) map[string]interface{} {
	jsonTs := make(map[string]interface{})
	for _, label := range metric.Labels {
		jsonTs[label.Name] = label.Value
	}
	for _, sample := range metric.Samples {
		jsonTs["timestamp"] = sample.Timestamp
		jsonTs["value"] = sample.Value
	}
	return jsonTs
}
