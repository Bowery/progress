// Copyright 2015 Bowery, Inc.

package progress

import (
	"io"
)

// transmitter monitors copying progress. It contains the current progress
// and channels to report the progress and any errors that occur.
type transmitter struct {
	src          io.Reader
	current      int64
	total        int64
	progressChan chan *Status
	errorChan    chan error
}

// newTransmitter creates a transmitter reading from src a total of n bytes.
func newTransmitter(src io.Reader, n int64) *transmitter {
	return &transmitter{
		src:          src,
		total:        n,
		progressChan: make(chan *Status),
		errorChan:    make(chan error),
	}
}

// Close closes all channels associated with the transmitter.
func (t *transmitter) Close() error {
	close(t.progressChan)
	close(t.errorChan)
	return nil
}

// Read implements io.Reader, reading from the source and incrementing
// the progress and reporting it to the progress channel.
func (t *transmitter) Read(p []byte) (int, error) {
	n, err := t.src.Read(p)
	if n > 0 {
		t.current += int64(n)
		t.progressChan <- &Status{
			Current: t.current,
			Total:   t.total,
		}
	}

	return n, err
}

// Copy tracks progress of copied data from src to dst, progress events are
// sent to the return status channel, once completed the IsFinished routine
// will report true.
func Copy(dst io.Writer, src io.Reader, n int64) (chan *Status, chan error) {
	t := newTransmitter(src, n)

	go func() {
		defer t.Close()

		_, err := io.Copy(dst, t)
		if err != nil {
			t.errorChan <- err
			return
		}

		t.progressChan <- &Status{
			Current:  t.total,
			Total:    t.total,
			finished: true,
		}
	}()

	return t.progressChan, t.errorChan
}
