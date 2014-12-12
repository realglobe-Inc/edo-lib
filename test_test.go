package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"reflect"
)

const (
	testLabel = "edo-test"
	testKey   = "test-key"
)

var testData = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

var testVal = map[string]interface{}{"array": []interface{}{"elem-1", "elem-2"}}

// JSON を通して等しいかどうか調べる。
func jsonEqual(v1 interface{}, v2 interface{}) (equal bool) {
	b1, err := json.Marshal(v1)
	if err != nil {
		log.Err(erro.Wrap(err))
		return false
	}
	var w1 interface{}
	if err := json.Unmarshal(b1, &w1); err != nil {
		log.Err(erro.Wrap(err))
		return false
	}

	b2, err := json.Marshal(v2)
	if err != nil {
		log.Err(erro.Wrap(err))
		return false
	}
	var w2 interface{}
	if err := json.Unmarshal(b2, &w2); err != nil {
		log.Err(erro.Wrap(err))
		return false
	}

	return reflect.DeepEqual(w1, w2)
}

func jsonUnmarshal(data []byte) (interface{}, error) {
	var res interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, erro.Wrap(err)
	}
	return res, nil
}
