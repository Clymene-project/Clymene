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

package handler

import (
	"encoding/json"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/client"
	"github.com/Clymene-project/Clymene/prompb"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/Clymene-project/Clymene/storage/metricstore"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
)

// HTTPHandler implements HTTP CollectorService.
type HTTPHandler struct {
	logger       *zap.Logger
	metricWriter metricstore.Writer
	prefix       string
	logWriter    logstore.Writer
}

func (h *HTTPHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(h.prefix+"/metrics", h.RequestMetrics).Methods(http.MethodPost)
	router.HandleFunc(h.prefix+"/logs", h.RequestLogs).Methods(http.MethodPost)
}

func (h *HTTPHandler) RequestMetrics(w http.ResponseWriter, r *http.Request) {
	req, err := h.DecodeWriteRequest(r.Body)
	if err != nil {
		h.logger.Error("Error decoding remote write request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = h.metricWriter.WriteMetric(req.GetTimeseries())
	w.WriteHeader(http.StatusNoContent)
}

func NewHTTPHandler(logger *zap.Logger, metricWriter metricstore.Writer, logWriter logstore.Writer) *HTTPHandler {
	return &HTTPHandler{
		logger:       logger,
		metricWriter: metricWriter,
		logWriter:    logWriter,
		prefix:       "/api",
	}
}

// DecodeWriteRequest from an io.Reader into a prompb.WriteRequest, handling
// snappy decompression.
func (h *HTTPHandler) DecodeWriteRequest(r io.Reader) (*prompb.WriteRequest, error) {
	compressed, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		return nil, err
	}

	var req prompb.WriteRequest
	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (h *HTTPHandler) RequestLogs(w http.ResponseWriter, r *http.Request) {
	compressed, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Error decoding logs write request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req := &client.ProducerBatch{}
	err = json.Unmarshal(compressed, req)
	if err != nil {
		h.logger.Error("Error Unmarshal logs write request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, _, _, _ = h.logWriter.Writelog(r.Context(), req.TenantID, &req.Batch)
	w.WriteHeader(http.StatusNoContent)
}
