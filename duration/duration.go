package duration

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// JSON にしたときに 72h3m0.5s みたいな文字列になる time.Duration。
type Duration time.Duration

func (this Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(this).String())
}

func (this *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*this = Duration(d)
	return nil
}

func (this Duration) GetBSON() (interface{}, error) {
	return time.Duration(this).String(), nil
}

func (this *Duration) SetBSON(raw bson.Raw) error {
	var s string
	if err := raw.Unmarshal(&s); err != nil {
		return err
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*this = Duration(d)
	return nil
}
