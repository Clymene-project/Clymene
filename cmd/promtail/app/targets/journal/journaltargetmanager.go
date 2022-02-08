//go:build !linux || !cgo
// +build !linux !cgo

package journal

import (
	"github.com/Clymene-project/Clymene/cmd/promtail/app/api"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/positions"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/scrapeconfig"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/targets/target"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// JournalTargetManager manages a series of JournalTargets.
// nolint:revive
type JournalTargetManager struct{}

// NewJournalTargetManager returns nil as JournalTargets are not supported
// on this platform.
func NewJournalTargetManager(
	reg prometheus.Registerer,
	logger *zap.Logger,
	positions positions.Positions,
	client api.EntryHandler,
	scrapeConfigs []scrapeconfig.Config,
) (*JournalTargetManager, error) {
	logger.Warn("WARNING!!! Journal target was configured but support for reading the systemd journal is not compiled into this build of promtail!")
	return &JournalTargetManager{}, nil
}

// Ready always returns false for JournalTargetManager on non-Linux
// platforms.
func (tm *JournalTargetManager) Ready() bool {
	return false
}

// Stop is a no-op on non-Linux platforms.
func (tm *JournalTargetManager) Stop() {}

// ActiveTargets always returns nil on non-Linux platforms.
func (tm *JournalTargetManager) ActiveTargets() map[string][]target.Target {
	return nil
}

// AllTargets always returns nil on non-Linux platforms.
func (tm *JournalTargetManager) AllTargets() map[string][]target.Target {
	return nil
}
