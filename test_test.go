package driver

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
