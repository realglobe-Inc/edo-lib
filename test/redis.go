// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"github.com/garyburd/redigo/redis"
	"github.com/realglobe-Inc/go-lib/erro"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var (
	// 実行可能な redis-server。
	RedisPath = "redis-server"
)

type RedisServer struct {
	pool *redis.Pool
	conf string
	addr string
}

func (this *RedisServer) Pool() *redis.Pool {
	return this.pool
}

func (this *RedisServer) Address() string {
	return this.addr
}

func (this *RedisServer) Close() {
	conn := this.pool.Get()
	conn.Do("SHUTDOWN")
	conn.Close()
	this.pool.Close()
	os.Remove(this.conf)
}

func NewRedisServer() (*RedisServer, error) {
	path, err := exec.LookPath(RedisPath)
	if err != nil {
		// 実行可能な redis-server が無い。
		return nil, nil
	}

	port, err := FreePort()
	if err != nil {
		return nil, erro.Wrap(err)
	}
	file, err := ioutil.TempFile("", "edo-lib.test")
	if err != nil {
		return nil, erro.Wrap(err)
	} else if _, err := file.Write([]byte("port " + strconv.Itoa(port))); err != nil {
		return nil, erro.Wrap(err)
	} else if err := file.Close(); err != nil {
		return nil, erro.Wrap(err)
	}
	// 失敗したら設定ファイルは消す。
	delete := true
	defer func() {
		if delete {
			os.Remove(file.Name())
		}
	}()

	cmd := exec.Command(path, file.Name())
	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.Run()
	}()

	// 起動待ち。
	for i := time.Nanosecond; ; i *= 2 {
		_, err := redis.Dial("tcp", ":"+strconv.Itoa(port))
		if err == nil {
			break
		}
		select {
		case err := <-errCh:
			return nil, erro.Wrap(err)
		case <-time.After(i):
		}
	}

	delete = false
	return &RedisServer{&redis.Pool{
		MaxIdle:     5,
		IdleTimeout: time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":"+strconv.Itoa(port))
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}, file.Name(), ":" + strconv.Itoa(port)}, nil
}
