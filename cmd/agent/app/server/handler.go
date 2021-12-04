package server

import (
	"github.com/Clymene-project/Clymene/cmd/agent/app/scrapeconfig"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

type APIHandler struct {
	l          *zap.Logger
	Reloader   []func(cfg *scrapeconfig.Config) error
	ConfigFile string
}

// NewAPIHandler returns a new APIHandler
func NewAPIHandler(p *HttpServerParams) *APIHandler {
	return &APIHandler{
		l:          p.Logger,
		Reloader:   p.Reloader,
		ConfigFile: p.ConfigFile,
	}
}

func (aH *APIHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/reload", aH.Reload).Methods(http.MethodGet)
}

func (aH *APIHandler) Reload(w http.ResponseWriter, r *http.Request) {
	if err := scrapeconfig.ReloadConfig(aH.ConfigFile, aH.l, aH.Reloader...); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.Wrapf(err, "error loading scrapeconfig from %q", aH.ConfigFile).Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("scrapeconfig reload, success"))
	return
}
