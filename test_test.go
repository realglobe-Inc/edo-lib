package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"reflect"
)

const (
	testLabel = "edo-test"

	testDir = "/test/directory"

	testUsrName = "test-user-no-namae"
	testUsrUuid = "test-user-no-uuid"

	testServUuid = "test-service-no-uuid"

	testIdpName = "test-id-provider-no-name"
	testIdpUuid = "test-id-provider-no-uuid"

	testAttrName = "test-attribute-no-name"

	testUri = "http://localhost:1234/test/uri"

	testKey = "test-key"

	testAccToken = "test-access-token"
)

var testAttr = map[string]interface{}{"array": []interface{}{"elem-1", "elem-2"}}
var testValue = testAttr

// data を JSON として、encoding/json の標準データ型にデコードする。
func jsonUnmarshal(data []byte) (interface{}, error) {
	var res interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, erro.Wrap(err)
	}
	return res, nil
}

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
