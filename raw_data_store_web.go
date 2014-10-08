package driver

import (
	"bytes"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io/ioutil"
	"net/http"
)

// GET http://example.com/{key}
// でボディに JSON で value が返り、
// PUT http://example.com/{key}
// のボディに JSON で value を入れると書き換えられ、
// DELETE http://example.com/{key}
// で消せると想定する。

type webRawDataStore webDriver

// スレッドセーフ。
func NewWebRawDataStore(prefix string) RawDataStore {
	// TODO キャッシュの並列化。
	return newSynchronizedRawDataStore(newCachingRawDataStore(newWebRawDataStore(prefix)))
}

// スレッドセーフ。
func newWebRawDataStore(prefix string) *webRawDataStore {
	return (*webRawDataStore)(newWebDriver(prefix))
}

func (reg *webRawDataStore) Get(key string, caStmp *Stamp) (data []byte, newCaStmp *Stamp, err error) {
	req, err := http.NewRequest("GET", reg.prefix+"/"+key, nil)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	if caStmp != nil {
		WriteStampToRequestHeader(caStmp, req.Header)
	}

	util.LogRequest(req, true)
	resp, err := reg.Client.Do(req)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	util.LogResponse(resp, true)

	// 404 Not Found なら無し。
	// 304 Not Modified なら変更無し。
	// 200 OK 以外なら失敗。
	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, nil, nil
	case http.StatusNotModified:
		newCaStmp, err := ParseStampFromResponseHeader(resp.Header)
		if err != nil {
			return nil, nil, erro.Wrap(err)
		}
		return nil, newCaStmp, nil
	case http.StatusOK:
		newCaStmp, err := ParseStampFromResponseHeader(resp.Header)
		if err != nil {
			return nil, nil, erro.Wrap(err)
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, erro.Wrap(err)
		}
		return data, newCaStmp, nil
	default:
		return nil, nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode))
	}
}

func (reg *webRawDataStore) Put(key string, data []byte) (*Stamp, error) {
	req, err := http.NewRequest("PUT", reg.prefix+"/"+key, bytes.NewReader(data))
	if err != nil {
		return nil, erro.Wrap(err)
	}

	util.LogRequest(req, true)
	resp, err := reg.Client.Do(req)
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

func (reg *webRawDataStore) Remove(key string) error {
	req, err := http.NewRequest("DELETE", reg.prefix+"/"+key, nil)
	if err != nil {
		return erro.Wrap(err)
	}

	util.LogRequest(req, true)
	resp, err := reg.Client.Do(req)
	if err != nil {
		return erro.Wrap(err)
	}
	defer resp.Body.Close()
	util.LogResponse(resp, true)

	// 200 OK 以外なら失敗。
	if resp.StatusCode != http.StatusOK {
		return erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode))
	}

	return nil
}
