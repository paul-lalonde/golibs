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
	"os"
	"bytes"
	"bufio"
)

type Streamer struct {
	path           string
	error          chan os.Error
	queue          chan string
	linebuffersize int
	charbuffersize int
}

func NewStreamer(path string) *Streamer {
	return &Streamer{
		path:           path,
		linebuffersize: 2,
		charbuffersize: 256,
	}
}

func (r *Streamer) Stream() (chan string, chan os.Error) {
	r.error = make(chan os.Error)
	r.queue = make(chan string, r.linebuffersize)
	go r.process()
	return r.queue, r.error
}

func (r *Streamer) process() {
	var (
		file *os.File
		err  os.Error
	)
	if file, err = os.Open(r.path); err != nil {
		r.error <-err
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, r.charbuffersize))
	var (
		part   []byte
		prefix bool
	)
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			r.error <-err
			break
		}
		buffer.Write(part)
		if !prefix { // Complete line has been read
			r.queue <-buffer.String()
		}
		buffer.Reset()
	}
}
