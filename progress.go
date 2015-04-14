// Copyright 2015 Bowery, Inc.

package progress

import (
	"io"
)

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

func newTransmitter(dst io.Writer, src io.Reader) *transmitter {
	return &transmitter{
		Writer:       dst,
		Reader:       src,
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

// Copy tracks progress of copied data from src to dst.
func Copy(dst io.Writer, src io.Reader, n int64) (chan *Status, chan error) {
	t := newTransmitter(dst, src)
	t.total = n

	go func(t *transmitter) {
		_, err := io.Copy(dst, t)
		if err != nil {
			t.errorChan <- err
		}
	}(t)

	return t.progressChan, t.errorChan
}
