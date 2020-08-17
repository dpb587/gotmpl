package cliopt

import "strings"

type DataSource struct {
	Key string
	URI string
}

func (d *DataSource) UnmarshalFlag(data string) error {
	sp := strings.SplitN(data, "=", 2)

	if len(sp) == 2 {
		d.Key = sp[0]
		d.URI = sp[1]
	} else {
		d.URI = sp[0]
	}

	return nil
}

type DataSourceList []DataSource
