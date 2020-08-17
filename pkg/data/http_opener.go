package data

import (
	"fmt"
	"io"
	"net/http"
)

type HTTPOpener struct {
	client *http.Client
	url    string
}

func NewHTTPOpener(client *http.Client, url string) *HTTPOpener {
	return &HTTPOpener{
		client: client,
		url:    url,
	}
}

var _ Opener = &HTTPOpener{}

func (o HTTPOpener) Open() (io.ReadCloser, error) {
	res, err := o.client.Get(o.url)
	if err != nil {
		return nil, err
	} else if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return res.Body, nil
}
