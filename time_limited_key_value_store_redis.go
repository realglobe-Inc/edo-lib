package driver

import (
	"github.com/garyburd/redigo/redis"
	"github.com/realglobe-Inc/go-lib-rg/erro"
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

type redisTimeLimitedKeyValueStore struct {
	pool *redis.Pool

	tag string
	Marshal
	Unmarshal
	GetStamp

	staleDur time.Duration
	expiDur  time.Duration
}

func NewRedisTimeLimitedKeyValueStore(pool *redis.Pool, tag string, marshal Marshal, unmarshal Unmarshal, getStmp GetStamp, staleDur, expiDur time.Duration) TimeLimitedKeyValueStore {
	return newRedisTimeLimitedKeyValueStore(pool, tag, marshal, unmarshal, getStmp, staleDur, expiDur)
}

func newRedisTimeLimitedKeyValueStore(pool *redis.Pool, tag string, marshal Marshal, unmarshal Unmarshal, getStmp GetStamp, staleDur, expiDur time.Duration) *redisTimeLimitedKeyValueStore {
	if getStmp == nil {
		getStmp = func(val interface{}) *Stamp {
			m, _ := val.(map[string]interface{})
			date, _ := m["date"].(time.Time)
			dig, _ := m["digest"].(string)
			return &Stamp{Date: date, Digest: dig}
		}
	}
	return &redisTimeLimitedKeyValueStore{
		pool:      pool,
		tag:       tag,
		Marshal:   marshal,
		Unmarshal: unmarshal,
		GetStamp:  getStmp,
		staleDur:  staleDur,
		expiDur:   expiDur,
	}
}

func (this *redisTimeLimitedKeyValueStore) getStamp(val interface{}) *Stamp {
	now := time.Now()
	stmp := this.GetStamp(val)
	stmp.StaleDate = now.Add(this.staleDur)
	stmp.ExpiDate = now.Add(this.expiDur)
	return stmp
}

func (this *redisTimeLimitedKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	conn := this.pool.Get()
	defer conn.Close()

	buff, err := redis.Bytes(conn.Do("GET", this.tag+key))
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil, nil
		} else {
			return nil, nil, erro.Wrap(err)
		}
	}

	val, err = this.Unmarshal(buff)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	newCaStmp = this.getStamp(val)
	if caStmp != nil && !caStmp.Older(newCaStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, newCaStmp, nil
	}

	// 要求元のキャッシュより新しそう。
	return val, newCaStmp, nil
}

func (this *redisTimeLimitedKeyValueStore) Put(key string, val interface{}, expiDate time.Time) (newCaStmp *Stamp, err error) {
	buff, err := this.Marshal(val)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	conn := this.pool.Get()
	defer conn.Close()

	newCaStmp = this.getStamp(val)

	milExpiDur := int64(expiDate.Sub(time.Now()) / time.Millisecond)
	if milExpiDur <= 0 {
		return newCaStmp, nil
	}
	if _, err := conn.Do("PSETEX", this.tag+key, milExpiDur, buff); err != nil {
		return nil, erro.Wrap(err)
	}
	return newCaStmp, nil
}

func (this *redisTimeLimitedKeyValueStore) Remove(key string) error {
	conn := this.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", this.tag+key); err != nil {
		return erro.Wrap(err)
	}
	return nil
}
