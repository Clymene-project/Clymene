package db

import (
	"github.com/Clymene-project/Clymene/plugin/storage/tdengine/db/async"
	"go.uber.org/zap"
	"sync"

	"github.com/taosdata/driver-go/v2/common"
	"github.com/taosdata/driver-go/v2/errors"
	"github.com/taosdata/driver-go/v2/wrapper"
)

var once = sync.Once{}

func PrepareConnection(taosConfigDir string, logger *zap.Logger) {
	if len(taosConfigDir) != 0 {
		once.Do(func() {
			code := wrapper.TaosOptions(common.TSDB_OPTION_CONFIGDIR, taosConfigDir)
			err := errors.GetError(code)
			if err != nil {
				logger.Panic("config Error", zap.String("set taos config file ", taosConfigDir))
			}
		})
	}
	async.GlobalAsync = async.NewAsync(async.NewHandlerPool(10000))
}
