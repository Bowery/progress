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

	isCopied := false
	for !isCopied {
		select {
		case status := <-progChan:
			if status.IsFinished() {
				isCopied = true
				break
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

	isCopied := false
	for !isCopied {
		select {
		case status := <-progChan:
			if status.IsFinished() {
				isCopied = true
				break
			}
		case err := <-errChan:
			t.Error(err)
		}
	}
}
