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
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

// HTTPHandler implements HTTP CollectorService.
type HTTPHandler struct {
	prefix    string
	logWriter logstore.Writer
	logger    *zap.Logger
}

func (h *HTTPHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(h.prefix+"/logs", h.RequestLogs).Methods(http.MethodPost)
}

func NewHTTPHandler(logger *zap.Logger, logWriter logstore.Writer) *HTTPHandler {
	return &HTTPHandler{
		logger:    logger,
		logWriter: logWriter,
		prefix:    "/api",
	}
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
