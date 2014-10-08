package driver

import (
	"bytes"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
	"time"
)

// GET http://example.com/{key}
// でボディに JSON で value が返り、
// PUT http://example.com/{key}
// のボディに JSON で value を入れると書き換えられ、
// DELETE http://example.com/{key}
// で消せると想定する。

type webTimeLimitedKeyValueStore struct {
	base RawDataStore
	Marshal
	Unmarshal
}

// スレッドセーフ。
func NewWebTimeLimitedKeyValueStore(prefix string, marshal Marshal, unmarshal Unmarshal) TimeLimitedKeyValueStore {
	return newSynchronizedTimeLimitedKeyValueStore(newCachingTimeLimitedKeyValueStore(newWebTimeLimitedKeyValueStore(prefix, marshal, unmarshal)))
}

// スレッドセーフ。
func newWebTimeLimitedKeyValueStore(prefix string, marshal Marshal, unmarshal Unmarshal) *webTimeLimitedKeyValueStore {
	return &webTimeLimitedKeyValueStore{newWebRawDataStore(prefix), marshal, unmarshal}
}

func (reg *webTimeLimitedKeyValueStore) Get(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
	data, newCaStmp, err := reg.base.Get(key, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if data == nil {
		return nil, newCaStmp, nil
	}

	value, err = reg.Unmarshal(data)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}
	return value, newCaStmp, nil
}

func (reg *webTimeLimitedKeyValueStore) Put(key string, value interface{}, expiDate time.Time) (*Stamp, error) {
	data, err := reg.Marshal(value)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	req, err := http.NewRequest("PUT", ((*webDriver)(reg.base.(*webRawDataStore))).prefix+"/"+key, bytes.NewReader(data))
	if err != nil {
		return nil, erro.Wrap(err)
	}
	req.Header.Set("Expires", expiDate.Format(http.TimeFormat))

	util.LogRequest(req, true)
	resp, err := ((*webDriver)(reg.base.(*webRawDataStore))).Client.Do(req)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	util.LogResponse(resp, true)

	// 200 OK, 201 Created 以外なら失敗。
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		newCaStmp, err := ParseStampFromResponseHeader(resp.Header)
		if err != nil {
			return nil, erro.Wrap(err)
		}
		return newCaStmp, nil
	default:
		return nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode))
	}
}

func (reg *webTimeLimitedKeyValueStore) Remove(key string) error {
	return reg.base.Remove(key)
}
