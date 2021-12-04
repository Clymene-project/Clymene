package server

import (
	"github.com/Clymene-project/Clymene/cmd/agent/app/config"
	"github.com/Clymene-project/Clymene/pkg/recoveryhandler"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// HttpServerParams to construct a new Jaeger Collector HTTP Server
type HttpServerParams struct {
	HostPort   string
	Logger     *zap.Logger
	Reloader   []func(cfg *config.Config) error
	ConfigFile string
}

// StartHTTPServer based on the given parameters
func StartHTTPServer(params *HttpServerParams) (*http.Server, error) {
	params.Logger.Info("Starting agent HTTP server", zap.String("http host-port", params.HostPort))

	listener, err := net.Listen("tcp", params.HostPort)
	if err != nil {
		return nil, err
	}

	server := &http.Server{Addr: params.HostPort}
	serveHTTP(server, listener, params)

	return server, nil
}

func serveHTTP(server *http.Server, listener net.Listener, params *HttpServerParams) {
	r := mux.NewRouter()

	ReloadApiHandler := NewAPIHandler(params)
	ReloadApiHandler.RegisterRoutes(r)

	recoveryHandler := recoveryhandler.NewRecoveryHandler(params.Logger, true)
	server.Handler = recoveryHandler(r)
	go func() {
		if err := server.Serve(listener); err != nil {
			if err != http.ErrServerClosed {
				params.Logger.Fatal("Could not start HTTP collector", zap.Error(err))
			}
		}
	}()
}
