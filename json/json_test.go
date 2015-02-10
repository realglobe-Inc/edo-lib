package json

import (
	"encoding/json"
	"testing"
)

func TestStringEscape(t *testing.T) {
	s := ` !"#$%'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_\` +
		"`" + `abcdefghijklmnopqrstuvwxyz{|}~` +
		`いろは` + "\n\r\t\b\f"

	data1, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	data2 := `"` + StringEscape(s) + `"`

	var s1, s2 string
	if err := json.Unmarshal(data1, &s1); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(data2), &s2); err != nil {
		t.Fatal(err)
	}

	if s2 != s1 {
		t.Error(s)
		t.Error(s2)
		t.Error(s1)
	}
}
