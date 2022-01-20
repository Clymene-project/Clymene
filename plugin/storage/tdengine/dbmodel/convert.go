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
	"crypto/md5"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/Clymene-project/Clymene/plugin/storage/tdengine/db/async"
	"github.com/Clymene-project/Clymene/plugin/storage/tdengine/db/tool"
	"github.com/Clymene-project/Clymene/prompb"
	"github.com/pkg/errors"
	"github.com/taosdata/driver-go/v2/common"
	tErrors "github.com/taosdata/driver-go/v2/errors"
	"github.com/taosdata/driver-go/v2/wrapper"
	"github.com/taosdata/taosadapter/thread"
	"github.com/taosdata/taosadapter/tools/pool"
	"math"
	"sort"
	"time"
	"unsafe"
)

type Converter struct {
}

func (c *Converter) processWrite(taosConn unsafe.Pointer, metrics []prompb.TimeSeries, db string) error {
	err := tool.SelectDB(taosConn, db)
	if err != nil {
		return err
	}

	sql, err := c.generateWriteSql(metrics)
	if err != nil {
		return err
	}
	err = async.GlobalAsync.TaosExecWithoutResult(taosConn, sql)
	if err != nil {
		if tErr, is := err.(*tErrors.TaosError); is {
			if tErr.Code == tErrors.MND_INVALID_TABLE_NAME {
				err := async.GlobalAsync.TaosExecWithoutResult(taosConn, "create stable if not exists metrics(ts timestamp,value double) tags (labels json)")
				if err != nil {
					return err
				}
				// retry
				err = async.GlobalAsync.TaosExecWithoutResult(taosConn, sql)
				return err
			} else {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (c *Converter) generateWriteSql(timeseries []prompb.TimeSeries) (string, error) {
	sql := pool.BytesPoolGet()
	defer pool.BytesPoolPut(sql)
	sql.WriteString("insert into ")
	tmp := pool.BytesPoolGet()
	defer pool.BytesPoolPut(tmp)
	for _, timeSeriesData := range timeseries {
		tagName := make([]string, len(timeSeriesData.Labels))
		tagMap := make(map[string]string, len(timeSeriesData.Labels))
		for i, label := range timeSeriesData.GetLabels() {
			tagName[i] = label.Name
			tagMap[label.Name] = label.Value
		}
		sort.Strings(tagName)
		tmp.Reset()
		for i, s := range tagName {
			v := tagMap[s]
			tmp.WriteString(s)
			tmp.WriteByte('=')
			tmp.WriteString(v)
			if i != len(tagName)-1 {
				tmp.WriteByte(',')
			}
		}
		labelsJson, err := json.Marshal(tagMap)
		if err != nil {
			return "", err
		}
		tableName := fmt.Sprintf("t_%x", md5.Sum(tmp.Bytes()))
		sql.WriteString(tableName)
		sql.WriteString(" using metrics tags('")
		sql.Write(labelsJson)
		sql.WriteString("') values")
		for _, sample := range timeSeriesData.Samples {
			sql.WriteString("('")
			sql.WriteString(time.Unix(0, sample.GetTimestamp()*1e6).UTC().Format(time.RFC3339Nano))
			sql.WriteString("',")
			if math.IsNaN(sample.GetValue()) {
				sql.WriteString("null")
			} else {
				fmt.Fprintf(sql, "%v", sample.GetValue())
			}
			sql.WriteString(") ")
		}
	}
	return sql.String(), nil
}

// TODO processRead, for the querier
func (c *Converter) processRead(taosConn unsafe.Pointer, req *prompb.ReadRequest, db string) (resp *prompb.ReadResponse, err error) {
	thread.Lock()
	wrapper.TaosSelectDB(taosConn, db)
	thread.Unlock()
	resp = &prompb.ReadResponse{}
	for i, query := range req.Queries {
		sql, err := c.generateReadSql(query)
		if err != nil {
			return nil, err
		}

		data, err := async.GlobalAsync.TaosExec(taosConn, sql, func(ts int64, precision int) driver.Value {
			switch precision {
			case common.PrecisionMilliSecond:
				return ts
			case common.PrecisionMicroSecond:
				return ts / 1e3
			case common.PrecisionNanoSecond:
				return ts / 1e6
			default:
				return 0
			}
		})
		if err != nil {
			return nil, err
		}
		//ts value labels time.Time float64 []byte
		group := map[string]*prompb.TimeSeries{}
		for _, d := range data.Data {
			if len(d) != 4 {
				continue
			}
			if d[0] == nil || d[1] == nil || d[2] == nil || d[3] == nil {
				continue
			}
			ts := d[0].(int64)
			value := d[1].(float64)
			var tags map[string]string
			err = json.Unmarshal(d[2].(json.RawMessage), &tags)
			if err != nil {
				return nil, err
			}
			tbName := d[3].(string)
			timeSeries, exist := group[tbName]
			if exist {
				timeSeries.Samples = append(timeSeries.Samples, prompb.Sample{
					Value:     value,
					Timestamp: ts,
				})
			} else {
				timeSeries = &prompb.TimeSeries{
					Samples: []prompb.Sample{
						{
							Value:     value,
							Timestamp: ts,
						},
					},
				}
				timeSeries.Labels = make([]prompb.Label, 0, len(tags))
				for name, tagValue := range tags {
					timeSeries.Labels = append(timeSeries.Labels, prompb.Label{
						Name:  name,
						Value: tagValue,
					})
				}
				group[tbName] = timeSeries
			}
		}
		if len(group) > 0 {
			resp.Results = append(resp.Results, &prompb.QueryResult{Timeseries: make([]*prompb.TimeSeries, 0, len(group))})
		}
		for _, series := range group {
			resp.Results[i].Timeseries = append(resp.Results[i].Timeseries, series)
		}
	}
	return resp, err
}

func (c *Converter) generateReadSql(query *prompb.Query) (string, error) {
	sql := pool.BytesPoolGet()
	defer pool.BytesPoolPut(sql)
	sql.WriteString("select *,tbname from metrics where ts >= '")
	sql.WriteString(c.ms2Time(query.GetStartTimestampMs()))
	sql.WriteString("' and ts <= '")
	sql.WriteString(c.ms2Time(query.GetEndTimestampMs()))
	sql.WriteByte('\'')
	for _, matcher := range query.GetMatchers() {
		sql.WriteString(" and ")
		k := matcher.GetName()
		v := matcher.GetValue()
		sql.WriteString("labels->'")
		sql.WriteString(k)
		switch matcher.Type {
		case prompb.LabelMatcher_EQ:
			sql.WriteString("' = '")
		case prompb.LabelMatcher_NEQ:
			sql.WriteString("' != '")
		case prompb.LabelMatcher_RE:
			sql.WriteString("' match '")
		case prompb.LabelMatcher_NRE:
			sql.WriteString("' nmatch '")
		default:
			return "", errors.New("not support match type")
		}
		sql.WriteString(v)
		sql.WriteByte('\'')
	}
	sql.WriteString(" order by ts desc")
	return sql.String(), nil
}

func (c *Converter) ms2Time(ts int64) string {
	return time.Unix(0, ts*1e6).UTC().Format(time.RFC3339Nano)
}
