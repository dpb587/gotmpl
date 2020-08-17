package data

import (
	"io"
	"io/ioutil"
)

type ReaderOpener struct {
	reader io.Reader
}

func NewReaderOpener(reader io.Reader) *ReaderOpener {
	return &ReaderOpener{
		reader: reader,
	}
}

var _ Opener = &ReaderOpener{}

func (o ReaderOpener) Open() (io.ReadCloser, error) {
	return ioutil.NopCloser(o.reader), nil
}
