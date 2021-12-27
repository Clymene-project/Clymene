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

package dryrun

import (
	"github.com/Clymene-project/Clymene/plugin/storage/opentsdb/metricstore/dbmodel"
	"github.com/Clymene-project/Clymene/prompb"
	"go.uber.org/zap"
)

type Client struct {
}
type Options struct {
}

func (c Client) SendData(series []prompb.TimeSeries) error {
	//TODO implement me
	panic("implement me")
}

func NewClient(converter *dbmodel.Converter, l *zap.Logger) *Client {

	return nil
}
