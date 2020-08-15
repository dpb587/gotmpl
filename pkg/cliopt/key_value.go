package cliopt

import "strings"

type KeyValue struct {
	Key   string
	Value string
}

func (d *KeyValue) UnmarshalFlag(data string) error {
	sp := strings.SplitN(data, "=", 2)

	if len(sp) == 2 {
		d.Key = sp[0]
		d.Value = sp[1]
	} else {
		d.Value = sp[0]
	}

	return nil
}

type KeyValueList []KeyValue
