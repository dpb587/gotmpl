package data

import "encoding/json"

type jsonParser struct{}

var JSON Parser = &jsonParser{}

func (jsonParser) Parse(buf []byte) (interface{}, error) {
	var out interface{}

	err := json.Unmarshal(buf, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
