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
	"github.com/Clymene-project/Clymene/prompb"
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
}

func (h *HTTPHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(h.prefix+"/metrics", h.RequestMetrics).Methods(http.MethodPost)
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

func NewHTTPHandler(
	logger *zap.Logger,
	metricWriter metricstore.Writer,
) *HTTPHandler {
	return &HTTPHandler{
		logger:       logger,
		metricWriter: metricWriter,
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
