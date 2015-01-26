package util

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

// JSON にしたときに要素の配列になる文字列集合型。
type StringSet map[string]bool

// コピーするだけ。
func NewStringSet(m map[string]bool) StringSet {
	s := map[string]bool{}
	for elem, ok := range m {
		if ok {
			s[elem] = true
		}
	}
	return StringSet(s)
}

func (this StringSet) MarshalJSON() ([]byte, error) {
	a := []string{}
	for elem, ok := range this {
		if ok {
			a = append(a, elem)
		}
	}
	return json.Marshal(a)
}

func (this *StringSet) UnmarshalJSON(data []byte) error {
	var a []string
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	s := map[string]bool{}
	for _, elem := range a {
		s[elem] = true
	}
	*this = StringSet(s)
	return nil
}

func (this StringSet) GetBSON() (interface{}, error) {
	a := []string{}
	for elem, ok := range this {
		if ok {
			a = append(a, elem)
		}
	}
	return a, nil
}

func (this *StringSet) SetBSON(raw bson.Raw) error {
	var a []string
	if err := raw.Unmarshal(&a); err != nil {
		return err
	}
	s := map[string]bool{}
	for _, elem := range a {
		s[elem] = true
	}
	*this = StringSet(s)
	return nil
}

func (this StringSet) Copy() StringSet {
	c := StringSet{}
	for elem, ok := range this {
		if ok {
			c[elem] = true
		}
	}
	return c
}

func OneOfStringSet(s StringSet) string {
	for elem := range s {
		return elem
	}
	return ""
}
