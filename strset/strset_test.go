package strset

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := map[string]bool{"a": true, "b": true, "c": true}

	a := New(m)
	if !reflect.DeepEqual(map[string]bool(a), m) {
		t.Error(a, m)
	}

	// コピーしてあるか。
	a["d"] = true
	if reflect.DeepEqual(map[string]bool(a), m) {
		t.Error(a, m)
	}
}

func TestFromSlice(t *testing.T) {
	m := map[string]bool{"a": true, "b": true, "c": true}

	l := []string{}
	for elem := range m {
		l = append(l, elem)
	}

	a := FromSlice(l)
	if !reflect.DeepEqual(map[string]bool(a), m) {
		t.Error(a, m)
	}
}

func TestJson(t *testing.T) {
	a := New(map[string]bool{"a": true, "b": true, "c": true})

	buff, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	} else if buff[0] != '[' {
		// JSON 配列じゃない。
		t.Error(string(buff))
	}

	var b StringSet
	if err := json.Unmarshal(buff, &b); err != nil {
		t.Fatal(err, string(buff))
	}

	if !reflect.DeepEqual(b, a) {
		t.Error(b, a)
	}
}

// 何かの中に入ってても大丈夫か。
func TestNestedJson(t *testing.T) {
	type testType struct {
		S StringSet
	}

	var a testType
	a.S = New(map[string]bool{"a": true, "b": true, "c": true})

	buff, err := json.Marshal(&a)
	if err != nil {
		t.Fatal(err)
	}

	var b testType
	if err := json.Unmarshal(buff, &b); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(b, a) {
		t.Error(b, a)
	}
}

// テストするなら、mongodb をたてる必要あり。
const (
	testMongoDb   = "edo-test"
	testMongoColl = "strset-test"
)

var testMongoAddr = "localhost"

func init() {
	if testMongoAddr != "" {
		// 実際にサーバーが立っているかどうか調べる。
		// 立ってなかったらテストはスキップ。
		conn, err := mgo.Dial(testMongoAddr)
		if err != nil {
			testMongoAddr = ""
		} else {
			conn.Close()
		}
	}
}

func TestBson(t *testing.T) {
	if testMongoAddr == "" {
		t.SkipNow()
	}

	conn, err := mgo.Dial(testMongoAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	type testType struct {
		K string    `bson:"key"`
		S StringSet `bson:"set"`
	}

	var a testType
	a.K = strconv.FormatInt(time.Now().UnixNano(), 16)
	a.S = New(map[string]bool{"a": true, "b": true, "c": true})

	if err := conn.DB(testMongoDb).C(testMongoColl).Insert(&a); err != nil {
		t.Fatal(err)
	}

	var b testType
	if err := conn.DB(testMongoDb).C(testMongoColl).Find(bson.M{"key": a.K}).One(&b); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(b, a) {
		t.Error(b, a)
	}
}

func TestCopy(t *testing.T) {
	a := New(map[string]bool{"a": true, "b": true, "c": true})
	b := a.Copy()
	if !reflect.DeepEqual(b, a) {
		t.Error(b, a)
	}

	// コピーしてあるか。
	b["d"] = true
	if reflect.DeepEqual(b, a) {
		t.Error(b, a)
	}
}

func TestOneOf(t *testing.T) {
	a := New(map[string]bool{"a": true, "b": true, "c": true})

	if v := a.OneOf(); !a[v] {
		t.Error(v, a)
	}
}