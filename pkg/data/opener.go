package data

import "io"

type Opener interface {
	Open() (io.ReadCloser, error)
}
