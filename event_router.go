package driver

import (
	"bytes"
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io"
	"net/http"
)

// イベントの処理。
type EventRouter interface {
	// イベントを発生させる。
	Fire(usrUuid, event string, body interface{}) error
}

func NewWebEventRouter(prefix string) EventRouter {
	return newWebDriver(prefix)
}

func (rout *webDriver) Fire(usrUuid, event string, body interface{}) error {
	var bodyType string
	var buff io.Reader
	if body != nil {
		bodyJson, err := json.Marshal(body)
		if err != nil {
			return erro.Wrap(err)
		}
		buff = bytes.NewReader(bodyJson)
		bodyType = util.ContentTypeJson
	}
	resp, err := rout.Post(rout.prefix+"/"+usrUuid+event, bodyType, buff)
	if err != nil {
		return erro.Wrap(err)
	}
	defer resp.Body.Close()
	//util.LogResponse(resp, true)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	return nil
}
