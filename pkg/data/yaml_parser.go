package data

import "gopkg.in/yaml.v2"

type yamlParser struct{}

var YAML Parser = &yamlParser{}

func (yamlParser) Parse(buf []byte) (interface{}, error) {
	var out interface{}

	err := yaml.Unmarshal(buf, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
