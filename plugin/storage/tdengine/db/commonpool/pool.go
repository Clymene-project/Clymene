package commonpool

import (
	"sync"
	"unsafe"

	"github.com/silenceper/pool"
	"github.com/taosdata/driver-go/v2/wrapper"
	"github.com/taosdata/taosadapter/config"
	"github.com/taosdata/taosadapter/thread"
)

type ConnectorPool struct {
	user     string
	password string
	pool     pool.Pool
}

func NewConnectorPool(user, password string) (*ConnectorPool, error) {
	a := &ConnectorPool{user: user, password: password}
	poolConfig := &pool.Config{
		InitialCap:  1,
		MaxCap:      config.Conf.Pool.MaxConnect,
		MaxIdle:     config.Conf.Pool.MaxIdle,
		Factory:     a.factory,
		Close:       a.close,
		IdleTimeout: config.Conf.Pool.IdleTimeout,
	}
	p, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		return nil, err
	}
	a.pool = p
	return a, nil
}

func (a *ConnectorPool) factory() (interface{}, error) {
	thread.Lock()
	defer thread.Unlock()
	return wrapper.TaosConnect("", a.user, a.password, "", 0)
}

func (a *ConnectorPool) close(v interface{}) error {
	if v != nil {
		thread.Lock()
		defer thread.Unlock()
		wrapper.TaosClose(v.(unsafe.Pointer))
	}
	return nil
}

func (a *ConnectorPool) Get() (unsafe.Pointer, error) {
	v, err := a.pool.Get()
	if err != nil {
		return nil, err
	}
	return v.(unsafe.Pointer), nil
}

func (a *ConnectorPool) Put(c unsafe.Pointer) error {
	wrapper.TaosResetCurrentDB(c)
	return a.pool.Put(c)
}

func (a *ConnectorPool) Close(c unsafe.Pointer) error {
	return a.pool.Close(c)
}

func (a *ConnectorPool) Release() {
	a.pool.Release()
}

func (a *ConnectorPool) verifyPassword(password string) bool {
	return password == a.password
}

var connectionMap = sync.Map{}

type Conn struct {
	TaosConnection unsafe.Pointer
	pool           *ConnectorPool
}

func (c *Conn) Put() error {
	return c.pool.Put(c.TaosConnection)
}

func GetConnection(user, password string) (*Conn, error) {
	p, exist := connectionMap.Load(user)
	if exist {
		connectionPool := p.(*ConnectorPool)
		if !connectionPool.verifyPassword(password) {
			newPool, err := NewConnectorPool(user, password)
			if err != nil {
				return nil, err
			}
			connectionPool.Release()
			connectionMap.Store(user, newPool)
			c, err := newPool.Get()
			if err != nil {
				return nil, err
			}
			return &Conn{
				TaosConnection: c,
				pool:           newPool,
			}, nil
		} else {
			c, err := connectionPool.Get()
			if err != nil {
				return nil, err
			}
			return &Conn{
				TaosConnection: c,
				pool:           connectionPool,
			}, nil
		}
	} else {
		newPool, err := NewConnectorPool(user, password)
		if err != nil {
			return nil, err
		}
		connectionMap.Store(user, newPool)
		c, err := newPool.Get()
		if err != nil {
			return nil, err
		}
		return &Conn{
			TaosConnection: c,
			pool:           newPool,
		}, nil
	}
}
