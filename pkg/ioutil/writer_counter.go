package ioutil

import "io"

type WriterCounter struct {
	writer io.Writer
	length int
}

var _ io.Writer = &WriterCounter{}

func NewWriterCounter(writer io.Writer) *WriterCounter {
	return &WriterCounter{
		writer: writer,
	}
}

func (wc *WriterCounter) Write(p []byte) (int, error) {
	n, err := wc.writer.Write(p)
	wc.length += n

	return n, err
}

func (wc *WriterCounter) WriteLength() int {
	return wc.length
}
