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

package loki

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	util "github.com/Clymene-project/Clymene/pkg/lokiutil"
	lokiflag "github.com/Clymene-project/Clymene/pkg/lokiutil/flagext"
	"github.com/Clymene-project/Clymene/pkg/version"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
)

const (
	contentType  = "application/x-protobuf"
	maxErrMsgLen = 1024
)

var UserAgent = fmt.Sprintf("promtail/%s", version.Get().Version)

type LogWriter struct {
	client         *http.Client
	logger         *zap.Logger
	url            *url.URL
	externalLabels lokiflag.LabelSet
}

func (l *LogWriter) Writelog(ctx context.Context, tenantID string, batch logstore.Batch) (int, int64, int64, error) {
	buf, entriesCount, err := batch.Encode()
	if err != nil {
		l.logger.Error("error encoding batch", zap.Error(err))
		return -1, -1, -1, err
	}
	bufBytes := int64(len(buf))
	entriesCount64 := int64(entriesCount)

	ctx, cancel := context.WithTimeout(ctx, l.client.Timeout)
	defer cancel()
	req, err := http.NewRequest("POST", l.url.String(), bytes.NewReader(buf))
	if err != nil {
		return -1, bufBytes, entriesCount64, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", UserAgent)

	// If the tenant ID is not empty promtail is running in multi-tenant mode, so
	// we should send it to Loki
	if tenantID != "" {
		req.Header.Set("X-Scope-OrgID", tenantID)
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return -1, bufBytes, entriesCount64, err
	}
	defer util.LogError("closing response body", l.logger, resp.Body.Close)

	if resp.StatusCode/100 != 2 {
		scanner := bufio.NewScanner(io.LimitReader(resp.Body, maxErrMsgLen))
		line := ""
		if scanner.Scan() {
			line = scanner.Text()
		}
		err = fmt.Errorf("server returned HTTP status %s (%d): %s", resp.Status, resp.StatusCode, line)
	}
	return resp.StatusCode, bufBytes, entriesCount64, err
}

func NewLogWriter(client *http.Client, url *url.URL, externalLabels lokiflag.LabelSet, logger *zap.Logger) *LogWriter {
	return &LogWriter{
		client:         client,
		url:            url,
		logger:         logger,
		externalLabels: externalLabels,
	}
}
