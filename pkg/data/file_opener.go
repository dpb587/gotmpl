package data

import (
	"io"
	"os"
)

type FileOpener struct {
	path string
}

func NewFileOpener(path string) *FileOpener {
	return &FileOpener{
		path: path,
	}
}

var _ Opener = &FileOpener{}

func (o FileOpener) Open() (io.ReadCloser, error) {
	return os.Open(o.path)
}
