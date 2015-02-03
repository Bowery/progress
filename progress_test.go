// Copyright 2015 Bowery, Inc.

package progress

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

var (
	testClient = New()
)

func TestGetSuccessful(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(getHandlerSuccessful))
	defer server.Close()

	output, err := os.Create("tmp")
	if err != nil {
		t.Error(err)
	}
	defer output.Close()
	defer os.Remove("tmp")

	url, _ := url.Parse(server.URL)
	progChan, errChan := testClient.Get(url, nil, output)

	isDownloaded := false
	for !isDownloaded {
		select {
		case status := <-progChan:
			if status.IsFinished() {
				isDownloaded = true
				break
			}
		case err := <-errChan:
			t.Error(err)
		}
	}
}

func getHandlerSuccessful(rw http.ResponseWriter, req *http.Request) {
	fakeData := strings.Repeat("drake", 1000)
	rw.Header().Set("Content-Length", "5000")
	fmt.Fprintf(rw, fakeData)
}

func TestGetNoContentLength(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(getHandlerNoContentLength))
	defer server.Close()

	output, err := os.Create("tmp")
	if err != nil {
		t.Error(err)
	}
	defer output.Close()
	defer os.Remove("tmp")

	url, _ := url.Parse(server.URL)
	progChan, errChan := testClient.Get(url, nil, output)

	receivedExpectedError := false
	isDownloaded := false
	for !isDownloaded && !receivedExpectedError {
		select {
		case status := <-progChan:
			if status.IsFinished() {
				isDownloaded = true
				break
			}
		case err := <-errChan:
			log.Println(err)
			if err == errContentLengthNotSet {
				receivedExpectedError = true
				break
			} else {
				t.Error(err)
			}
		}
	}

	if !receivedExpectedError {
		t.Error("failed to end with expected error")
	}
}

func getHandlerNoContentLength(rw http.ResponseWriter, req *http.Request) {
	fakeData := strings.Repeat("j cole", 1000)
	fmt.Fprintf(rw, fakeData)
}
