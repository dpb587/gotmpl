package data

import (
	"fmt"

	"github.com/pkg/errors"
)

type Data struct {
	sources []*Source
}

func NewData(sources ...*Source) *Data {
	return &Data{
		sources: sources,
	}
}

func (d *Data) Validate() error {
	keys := map[string]bool{}

	for _, ds := range d.sources {
		k := ds.Key()

		_, found := keys[k]
		if found {
			err := fmt.Errorf("duplicate data key: %s", k)

			if k == "" {
				err = fmt.Errorf("%s (keys must be used with multiple data sources)", err)
			}

			return err
		}

		keys[k] = true
	}

	if _, found := keys[""]; found && len(keys) > 0 {
		return errors.New("data source is missing a key")
	}

	return nil
}

func (d *Data) Data() (interface{}, error) {
	res := map[string]interface{}{}

	for _, ds := range d.sources {
		k := ds.Key()
		data, err := ds.Data()
		if err != nil {
			return nil, errors.Wrapf(err, "loading data (key: %s)", k)
		}

		res[k] = data
	}

	if _, found := res[""]; found && len(res) == 1 {
		return res[""], nil
	}

	return res, nil
}
