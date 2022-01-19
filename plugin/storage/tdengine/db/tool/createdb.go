package tool

import (
	"unsafe"

	"github.com/taosdata/driver-go/v2/errors"
	"github.com/taosdata/driver-go/v2/wrapper"
	"github.com/taosdata/taosadapter/db/async"
	"github.com/taosdata/taosadapter/httperror"
	"github.com/taosdata/taosadapter/thread"
	"github.com/taosdata/taosadapter/tools/pool"
)

func CreateDBWithConnection(connection unsafe.Pointer, db string) error {
	b := pool.BytesPoolGet()
	defer pool.BytesPoolPut(b)
	b.WriteString("create database if not exists ")
	b.WriteString(db)
	b.WriteString(" precision 'ns' update 2")
	err := async.GlobalAsync.TaosExecWithoutResult(connection, b.String())
	if err != nil {
		return err
	}
	return nil
}

func SelectDB(taosConnect unsafe.Pointer, db string) error {
	thread.Lock()
	code := wrapper.TaosSelectDB(taosConnect, db)
	thread.Unlock()
	if code != httperror.SUCCESS {
		if int32(code)&0xffff == errors.TSC_DB_NOT_SELECTED || int32(code)&0xffff == errors.MND_INVALID_DB {
			err := CreateDBWithConnection(taosConnect, db)
			if err != nil {
				return err
			}
			thread.Lock()
			wrapper.TaosSelectDB(taosConnect, db)
			thread.Unlock()
		} else {
			return errors.GetError(code)
		}
	}
	return nil
}
