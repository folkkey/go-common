package gocommon

import "encoding/json"

func TypeConverter[T any](source any) (*T, error) {
	var result T
	b, err := json.Marshal(&source)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}
