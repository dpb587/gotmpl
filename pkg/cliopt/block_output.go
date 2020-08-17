package cliopt

import (
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type BlockOutput struct {
	Block string
	Path  string
}

func (d *BlockOutput) UnmarshalFlag(data string) error {
	sp := strings.SplitN(data, "=", 2)

	if len(sp) == 2 {
		d.Block = sp[1]
		d.Path = sp[0]
	} else {
		d.Path = sp[0]
	}

	return nil
}

func (d *BlockOutput) BlockLabel(defaultBlock string) string {
	if d.Block == "" {
		return defaultBlock
	}

	return d.Block
}

func (d *BlockOutput) PathLabel() string {
	if d.Path == "-" {
		return "/dev/stdout"
	}

	return d.Path
}

func (d *BlockOutput) Open() (io.Writer, func() error, error) {
	if d.Path == "/dev/stdout" || d.Path == "-" {
		return os.Stdout, nil, nil
	} else if d.Path == "/dev/stderr" {
		return os.Stderr, nil, nil
	}

	fh, err := os.OpenFile(d.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "opening %s", d.Path)
	}

	return fh, fh.Close, nil
}

type BlockOutputList []BlockOutput
