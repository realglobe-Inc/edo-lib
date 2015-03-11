package driver

import (
	"github.com/garyburd/redigo/redis"
	"github.com/realglobe-Inc/go-lib/erro"
	"time"
)

func NewRedisPool(addr string, connNum int, idlDur time.Duration) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     connNum,
		IdleTimeout: idlDur,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

type redisConcurrentVolatileKeyValueStore struct {
	pool *redis.Pool

	tag string
	Marshal
	Unmarshal
	GetStamp

	staleDur time.Duration
	expiDur  time.Duration
}

func NewRedisConcurrentVolatileKeyValueStore(pool *redis.Pool, tag string, marshal Marshal, unmarshal Unmarshal,
	getStmp GetStamp, staleDur, expiDur time.Duration) ConcurrentVolatileKeyValueStore {
	return newRedisConcurrentVolatileKeyValueStore(pool, tag,
		marshal, unmarshal, getStmp, staleDur, expiDur)
}

func newRedisConcurrentVolatileKeyValueStore(pool *redis.Pool, tag string, marshal Marshal, unmarshal Unmarshal,
	getStmp GetStamp, staleDur, expiDur time.Duration) *redisConcurrentVolatileKeyValueStore {
	if getStmp == nil {
		getStmp = func(val interface{}) *Stamp {
			m, _ := val.(map[string]interface{})
			date, _ := m["date"].(time.Time)
			dig, _ := m["digest"].(string)
			return &Stamp{Date: date, Digest: dig}
		}
	}
	return &redisConcurrentVolatileKeyValueStore{
		pool:      pool,
		tag:       tag,
		Marshal:   marshal,
		Unmarshal: unmarshal,
		GetStamp:  getStmp,
		staleDur:  staleDur,
		expiDur:   expiDur,
	}
}

func (drv *redisConcurrentVolatileKeyValueStore) getStamp(val interface{}) *Stamp {
	now := time.Now()
	stmp := drv.GetStamp(val)
	stmp.StaleDate = now.Add(drv.staleDur)
	stmp.ExpiDate = now.Add(drv.expiDur)
	return stmp
}

func (drv *redisConcurrentVolatileKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	buff, err := redis.Bytes(func() (interface{}, error) {
		// パニックでも解放するように defer で、使ったらすぐ解放するように無名関数で。
		conn := drv.pool.Get()
		defer conn.Close()
		return conn.Do("GET", drv.tag+key)
	}())
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil, nil
		} else {
			return nil, nil, erro.Wrap(err)
		}
	}

	val, err = drv.Unmarshal(buff)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	newCaStmp = drv.getStamp(val)
	if caStmp != nil && !caStmp.Older(newCaStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, newCaStmp, nil
	}

	// 要求元のキャッシュより新しそう。
	return val, newCaStmp, nil
}

func (drv *redisConcurrentVolatileKeyValueStore) Put(key string, val interface{}, expiDate time.Time) (newCaStmp *Stamp, err error) {
	buff, err := drv.Marshal(val)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	newCaStmp = drv.getStamp(val)

	milExpiDur := int64(expiDate.Sub(time.Now()) / time.Millisecond)
	if milExpiDur <= 0 {
		return nil, nil
	}

	conn := drv.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("SET", drv.tag+key, buff, "PX", milExpiDur); err != nil {
		return nil, erro.Wrap(err)
	}
	return newCaStmp, nil
}

func (drv *redisConcurrentVolatileKeyValueStore) Remove(key string) error {
	conn := drv.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", drv.tag+key); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func (drv *redisConcurrentVolatileKeyValueStore) Close() error {
	drv.pool = nil
	return nil
}

func (drv *redisConcurrentVolatileKeyValueStore) Entry(eKey string) (eVal string, err error) {
	eVal, err = redis.String(func() (interface{}, error) {
		// パニックでも解放するように defer で、使ったらすぐ解放するように無名関数で。
		conn := drv.pool.Get()
		defer conn.Close()
		return conn.Do("GET", drv.tag+eKey)
	}())
	if err != nil {
		if err == redis.ErrNil {
			return "", nil
		} else {
			return "", erro.Wrap(err)
		}
	}
	return eVal, nil
}

func (drv *redisConcurrentVolatileKeyValueStore) SetEntry(eKey, eVal string, expiDate time.Time) error {
	milExpiDur := int64(expiDate.Sub(time.Now()) / time.Millisecond)
	if milExpiDur <= 0 {
		return nil
	}

	conn := drv.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("SET", drv.tag+eKey, eVal, "PX", milExpiDur); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func (drv *redisConcurrentVolatileKeyValueStore) GetAndSetEntry(key string, caStmp *Stamp, eKey, eVal string, eExpiDate time.Time) (val interface{}, newCaStmp *Stamp, err error) {
	const script = `redis.call("set",KEYS[1],ARGV[1],"PX",ARGV[2])
return redis.call("get",KEYS[2])`

	milEExpiDur := int64(eExpiDate.Sub(time.Now()) / time.Millisecond)

	buff, err := redis.Bytes(func() (interface{}, error) {
		// パニックでも解放するように defer で、使ったらすぐ解放するように無名関数で。
		conn := drv.pool.Get()
		defer conn.Close()
		return conn.Do("EVAL", script, 2, drv.tag+eKey, drv.tag+key, eVal, milEExpiDur)
	}())
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil, nil
		} else {
			return nil, nil, erro.Wrap(err)
		}
	}

	val, err = drv.Unmarshal(buff)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	newCaStmp = drv.getStamp(val)
	if caStmp != nil && !caStmp.Older(newCaStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, newCaStmp, nil
	}

	// 要求元のキャッシュより新しそう。
	return val, newCaStmp, nil
}

func (drv *redisConcurrentVolatileKeyValueStore) PutIfEntered(key string, val interface{}, expiDate time.Time, eKey, eVal string) (entered bool, newCaStmp *Stamp, err error) {
	const script = `if redis.call("get",KEYS[1]) ~= ARGV[1] then
return 0
end
redis.call("set",KEYS[2],ARGV[2],"PX",ARGV[3])
return 1`

	buff, err := drv.Marshal(val)
	if err != nil {
		return false, nil, erro.Wrap(err)
	}

	newCaStmp = drv.getStamp(val)

	milExpiDur := int64(expiDate.Sub(time.Now()) / time.Millisecond)
	if milExpiDur <= 0 {
		return false, nil, nil
	}

	if entered, err := redis.Int(func() (interface{}, error) {
		conn := drv.pool.Get()
		defer conn.Close()
		return conn.Do("EVAL", script, 3, drv.tag+eKey, drv.tag+key, "", eVal, buff, milExpiDur)
	}()); err != nil {
		return false, nil, erro.Wrap(err)
	} else if entered > 0 {
		return true, newCaStmp, nil
	} else {
		return false, nil, nil
	}
}
