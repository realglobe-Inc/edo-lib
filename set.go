package util

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

// JSON にしたときに要素の配列になる文字列集合型。
type StringSet struct {
	m map[string]bool
}

func NewStringSet(m map[string]bool) *StringSet {
	s := &StringSet{map[string]bool{}}
	for elem, ok := range m {
		if ok {
			s.m[elem] = true
		}
	}
	return s
}

func (this *StringSet) Put(elem string) {
	this.m[elem] = true
}

func (this *StringSet) Contains(elem string) bool {
	return this.m[elem]
}

func (this *StringSet) Remove(elem string) {
	delete(this.m, elem)
}

func (this *StringSet) Len() int {
	return len(this.m)
}

func (this *StringSet) Elements() map[string]bool {
	return this.m
}

func (this *StringSet) MarshalJSON() ([]byte, error) {
	a := []string{}
	for elem := range this.m {
		a = append(a, elem)
	}
	return json.Marshal(a)
}

func (this *StringSet) UnmarshalJSON(data []byte) error {
	var a []string
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	this.m = map[string]bool{}
	for _, elem := range a {
		this.m[elem] = true
	}
	return nil
}

func (this *StringSet) GetBSON() (interface{}, error) {
	a := []string{}
	for elem := range this.m {
		a = append(a, elem)
	}
	return a, nil
}

func (this *StringSet) SetBSON(raw bson.Raw) error {
	var a []string
	if err := raw.Unmarshal(&a); err != nil {
		return err
	}
	this.m = map[string]bool{}
	for _, elem := range a {
		this.m[elem] = true
	}
	return nil
}
