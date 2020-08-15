package data

import (
	"fmt"
	"io"
	"net/http"
)

type URLOpener struct {
	client *http.Client
	url    string
}

func NewURLOpener(client *http.Client, url string) *URLOpener {
	return &URLOpener{
		client: client,
		url:    url,
	}
}

var _ Opener = &URLOpener{}

func (o URLOpener) Open() (io.ReadCloser, error) {
	res, err := o.client.Get(o.url)
	if err != nil {
		return nil, err
	} else if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return res.Body, nil
}
