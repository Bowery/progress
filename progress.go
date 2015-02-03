// Copyright 2015 Bowery, Inc.

package progress

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

var (
	errContentLengthNotSet = errors.New("Content-Length not set")
	errStatusNotOK         = errors.New("Status Code non 200")
)

// Client is a progress client to Get and Put files. It uses
// the default http.DefaultClient to make requests, however
// one can be set if non standard options are required.
type Client struct {
	HTTPClient *http.Client
}

// New creates a new Client.
func New() *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
	}
}

// transmitter is the monitoring object for a download
// or upload. It is comprised of the raw progress, as well
// as channels for reporting progress and errors.
type transmitter struct {
	io.Reader
	io.Writer
	current      int64
	total        int64
	progressChan chan *Status
	errorChan    chan error
}

func newTransmitter() *transmitter {
	return &transmitter{
		progressChan: make(chan *Status),
		errorChan:    make(chan error),
	}
}

// closeChannels closes all channels associated with
// the transmitter.
func (t *transmitter) closeChannels() {
	close(t.progressChan)
	close(t.errorChan)
}

// Read is the implementation of the transmitter's io.Reader
// read method. It increments current progress and reports
// progress to the associated progressChannel.
func (t *transmitter) Read(p []byte) (int, error) {
	n, err := t.Reader.Read(p)
	t.current += int64(n)
	t.progressChan <- &Status{
		Current: t.current,
		Total:   t.total,
	}
	return n, err
}

// Get downloads the contents from the provided url.
func (c *Client) Get(url *url.URL, header http.Header, dst io.Writer) (chan *Status, chan error) {
	t := newTransmitter()

	go func(t *transmitter) {
		req, err := http.NewRequest("GET", url.String(), nil)
		if err != nil {
			t.errorChan <- err
			return
		}

		if header != nil {
			req.Header = header
		}

		res, err := c.HTTPClient.Do(req)
		if err != nil {
			t.errorChan <- err
			return
		}

		if res.StatusCode != http.StatusOK {
			t.errorChan <- errStatusNotOK
			return
		}

		if res.ContentLength == -1 {
			t.errorChan <- errContentLengthNotSet
			return
		}

		t.Reader = res.Body
		t.total = res.ContentLength

		defer res.Body.Close()
		defer t.closeChannels()
		_, err = io.Copy(dst, t)
		if err != nil {
			t.errorChan <- err
		}
	}(t)

	return t.progressChan, t.errorChan
}
