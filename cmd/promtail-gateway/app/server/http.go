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

package server

import (
	"github.com/Clymene-project/Clymene/cmd/promtail-gateway/app/handler"
	"github.com/Clymene-project/Clymene/pkg/config/tlscfg"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net"
	"net/http"
)

// HTTPServerParams to construct a new Clymene Gateway HTTP Server
type HTTPServerParams struct {
	TLSConfig tlscfg.Options
	HostPort  string
	Handler   *handler.HTTPHandler
	Logger    *zap.Logger
}

// StartHTTPServer based on the given parameters
func StartHTTPServer(params *HTTPServerParams) (*http.Server, error) {
	var server = &http.Server{Addr: params.HostPort}
	if params.TLSConfig.Enabled {
		tlsCfg, err := params.TLSConfig.Config(params.Logger) // This checks if the certificates are correctly provided
		if err != nil {
			return nil, err
		}
		server.TLSConfig = tlsCfg
	}

	listener, err := net.Listen("tcp", params.HostPort)
	if err != nil {
		return nil, err
	}

	serveHTTP(server, listener, params)
	return server, nil
}

func serveHTTP(server *http.Server, listener net.Listener, params *HTTPServerParams) {
	r := mux.NewRouter()
	apiHandler := params.Handler
	apiHandler.RegisterRoutes(r)

	server.Handler = r
	go func() {
		var err error
		if params.TLSConfig.Enabled {
			err = server.ServeTLS(listener, "", "")
		} else {
			err = server.Serve(listener)
		}
		if err != nil {
			if err != http.ErrServerClosed {
				params.Logger.Error("Could not start HTTP gateway", zap.Error(err))
			}
		}
	}()
}
