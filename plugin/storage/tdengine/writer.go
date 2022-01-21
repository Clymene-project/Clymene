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

package tdengine

import (
	"database/sql"
	"github.com/Clymene-project/Clymene/plugin/storage/tdengine/dbmodel"
	"github.com/Clymene-project/Clymene/prompb"
	"go.uber.org/zap"
)

/*
/rest/sql,		timestamp format: "2018-10-03 14:38:05.000"
/rest/sqlt, 	timestamp format: 1538548685000
/rest/sqlutc, 	timestamp format: "2018-10-03T14:38:05.000+0800"
https://github.com/taosdata/TDengine/blob/develop/tests/examples/go/taosdemo.go
*/

type MetricWriter struct {
	converter    dbmodel.Converter
	tdEngine     *sql.DB
	maxSQLLength int
	logger       *zap.Logger
}

func (m *MetricWriter) WriteMetric(metrics []prompb.TimeSeries) error {
	if len(metrics) > m.maxSQLLength {
		q := len(metrics) / m.maxSQLLength
		r := len(metrics) % m.maxSQLLength
		if r != 0 {
			q += 1
		}
		for i := 1; i <= q; i++ {
			var timeSeriesDiv []prompb.TimeSeries
			if i == 1 {
				timeSeriesDiv = metrics[:i*m.maxSQLLength]
			} else if i != q {
				timeSeriesDiv = metrics[(i-1)*m.maxSQLLength : i*m.maxSQLLength]
			} else {
				timeSeriesDiv = metrics[(i-1)*m.maxSQLLength:]
			}
			_ = m.writeMetric(timeSeriesDiv)
		}
	} else {
		return m.writeMetric(metrics)
	}
	return nil
}

func (m *MetricWriter) writeMetric(metrics []prompb.TimeSeries) error {
	_, query, err := m.converter.GenerateWriteSql(metrics)
	if err != nil {
		return err
	}
	_, err = m.tdEngine.Exec(query)
	if err != nil {
		if err.Error() == "[0x362] Table does not exist" {
			_, err = m.tdEngine.Exec("create stable if not exists metrics (ts TIMESTAMP, value DOUBLE) tags (labels json)")
			if err != nil {
				return err
			}
			// retry
			_, err = m.tdEngine.Exec(query)
			return err
		} else {
			return err
		}
	}
	return nil
}

func NewMetricWriter(tdEngine *sql.DB, maxSQLLength int, l *zap.Logger) *MetricWriter {
	return &MetricWriter{
		converter:    dbmodel.Converter{},
		tdEngine:     tdEngine,
		logger:       l,
		maxSQLLength: maxSQLLength,
	}
}
