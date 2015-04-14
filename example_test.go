// Copyright 2015 Bowery, Inc.

package progress_test

import (
	"bytes"
	"os"

	"github.com/Bowery/progress"
)

func ExampleCopy() {
	file, err := os.Open("/path/to/file")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	progChan, errChan := progress.Copy(&buf, file, stat.Size())

	isCopied := false
	for !isCopied {
		select {
		case status := <-progChan:
			if status.IsFinished() {
				isCopied = true
				break
			}
		case err := <-errChan:
			panic(err)
		}
	}
}
