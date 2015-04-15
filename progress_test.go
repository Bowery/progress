// Copyright 2015 Bowery, Inc.

package progress

import (
	"bytes"
	"net/http"
	"os"
	"testing"
)

func TestCopyFileSuccess(t *testing.T) {
	file, err := os.Open("progress_test.go")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	progChan, errChan := Copy(&buf, file, stat.Size())

	for {
		select {
		case status := <-progChan:
			if status.IsFinished() {
				// Actual testing occurs here.
				if int64(buf.Len()) != stat.Size() {
					t.Error("Copied length does not match file length")
				}

				return
			}
		case err := <-errChan:
			t.Error(err)
		}
	}
}

func TestCopyGetRequestSuccess(t *testing.T) {
	res, err := http.Get("http://bowery.io")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	var buf bytes.Buffer
	progChan, errChan := Copy(&buf, res.Body, res.ContentLength)

	for {
		select {
		case status := <-progChan:
			if status.IsFinished() {
				// Actual testing occurs here.
				if int64(buf.Len()) != res.ContentLength {
					t.Error("Copied length does not match file length")
				}

				return
			}
		case err := <-errChan:
			t.Error(err)
		}
	}
}
