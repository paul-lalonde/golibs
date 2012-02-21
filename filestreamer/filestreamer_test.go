// Copyright 2012 Twitter, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filestreamer

import (
	"testing"
	"os"
	"fmt"
	"path/filepath"
	"time"
)

func write(data string) (string, os.Error) {
	var (
		err os.Error
		file *os.File
	)
	fileName := fmt.Sprintf("filestreamer-test-%v.txt", time.Nanoseconds())
	tempPath := filepath.Join(os.TempDir(), fileName)
	if file, err = os.Create(tempPath); err != nil {
		return "", err
	}
	defer file.Close()
	if _, err = file.WriteString(data); err != nil {
		return "", err
	}
	return tempPath, nil
}

func clean(path string) os.Error {
	return os.Remove(path)
}

var TEST_DATA string = `1	2	3
4	5	6
7	8	9
10	11	12
`
func TestConsume(t *testing.T) {
	var (
		path string
		err os.Error
		exit bool
		tally int
	)
	if path, err = write(TEST_DATA); err != nil {
		t.Fatal(err)
	}
	defer clean(path)
	streamer := NewStreamer(path)
	lines, errors := streamer.Stream()
	exit = false
	tally = 0
	for exit == false{
		select {
		case <-lines:
			tally += 1
		case err = <-errors:
			if err != os.EOF {
				t.Fatal(err)
			} else {
				exit = true
			}
		}
	}
	if tally != 4 {
		t.Fatalf("Read %v lines, expected %v", tally, 4)
	}
}
