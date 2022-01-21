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
	"encoding/json"
	"fmt"
	"github.com/Clymene-project/Clymene/prompb"
	"github.com/pkg/errors"
	"github.com/taosdata/taosadapter/tools/pool"
	"math"
	"regexp"
	"sort"
	"time"
)

type Converter struct {
}

func (c *Converter) GenerateWriteSql(timeseries []prompb.TimeSeries) (string, string, error) {
	sql := pool.BytesPoolGet()
	defer pool.BytesPoolPut(sql)
	sql.WriteString("insert into ")
	tmp := pool.BytesPoolGet()
	defer pool.BytesPoolPut(tmp)
	tableName := ""
	for _, timeSeriesData := range timeseries {
		tagName := make([]string, len(timeSeriesData.Labels))
		tagMap := make(map[string]string, len(timeSeriesData.Labels))
		for i, label := range timeSeriesData.GetLabels() {
			// The "/" character causes a syntax error
			match, _ := regexp.MatchString("[`\"\\[\\]\\\\]", label.Value)
			if !match {
				if label.Value != "" {
					tagName[i] = label.Name
					tagMap[label.Name] = label.Value
				}
			}
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
			return "", "", err
		}
		tableName = fmt.Sprintf("t_%x", md5.Sum(tmp.Bytes()))
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
	return tableName, sql.String(), nil
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
