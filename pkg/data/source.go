package data

import (
	"io/ioutil"

	"github.com/pkg/errors"
)

type Source struct {
	key    string
	opener Opener
	parser Parser
}

func NewSource(key string, opener Opener, parser Parser) *Source {
	return &Source{
		key:    key,
		opener: opener,
		parser: parser,
	}
}

func (s *Source) Key() string {
	return s.key
}

func (s *Source) Data() (interface{}, error) {
	r, err := s.opener.Open()
	if err != nil {
		return nil, errors.Wrap(err, "opening data source")
	}

	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "reading data source")
	}

	d, err := s.parser.Parse(b)
	if err != nil {
		return nil, errors.Wrap(err, "parsing data source")
	}

	return d, nil
}
