// MIT License
//
// Copyright (c) 2018 John Pruitt
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package logger

import (
	"bytes"
	"errors"
	"io"
	"math"
	"os"
	"strings"
)

// A FileSink writes log messages to a file.
// It will roll content within the file to prevent the size exceeding a maximum.
type FileSink struct {
	max  int64
	keep int64
	size int64
	out  *os.File
	path string
	buf  bytes.Buffer
}

// NewFileSink constructs a new FileSink with a file at path.
// It will limit the file size to max bytes.
// When it rolls content it will copy the latest keep bytes of content from
// the end of the file to the beginning.
func NewFileSink(path string, max int64, keep int64) (file *FileSink, err error) {
	if keep < 0 {
		return nil, errors.New("keep must be greater than zero")
	}
	if max < 0 {
		return nil, errors.New("max must be greater than zero")
	}
	if math.MaxInt64/2 < keep {
		return nil, errors.New("keep is too large")
	}
	if keep+keep > max {
		return nil, errors.New("keep must less than or equal to half max")
	}

	file = &FileSink{
		max:  max,
		keep: keep,
		path: path,
	}

	// create or open the file for writing
	file.out, err = os.OpenFile(file.path, os.O_CREATE|os.O_WRONLY, os.FileMode(0666))
	if err != nil {
		return nil, err
	}

	// seek to the end of the file
	file.size, err = file.out.Seek(0, io.SeekEnd)
	if err != nil {
		file.out.Close()
		return nil, err
	}

	// print a separator line to delineate the beginning of a new log in the file
	var n int
	n, err = file.out.Write([]byte(strings.Repeat("=", 80) + "\n"))
	if err != nil {
		file.out.Close()
		return nil, err
	}
	file.size = file.size + int64(n)

	return
}

// Close closes the logFile, rendering it unusable for I/O.
func (file *FileSink) Close() error {
	out := file.out
	file.out = nil
	return out.Close()
}

// Log writes a log message to the file
func (file *FileSink) Log(m *Msg) {
	// make sure the file wasn't already closed
	if file.out == nil {
		return
	}

	buf := &file.buf
	m.PrintDate(buf)
	buf.WriteString(" ")
	m.PrintTime(buf)
	buf.WriteString(" ")
	m.PrintLevel(buf)
	buf.WriteString(" ")
	m.PrintFileLine(buf)
	buf.WriteString(" ")
	m.PrintMsg(buf)
	n, _ := file.out.Write(buf.Bytes())
	buf.Reset()

	file.size = file.size + int64(n)
	if file.size >= file.max {
		file.roll()
	}
}

// roll rolls the file content
func (file *FileSink) roll() {
	keep := file.keep
	out := file.out

	// make sure everything is really on disk
	if err := out.Sync(); err != nil {
		panic(err)
	}

	// open the file for reading
	in, err := os.OpenFile(file.path, os.O_RDONLY, os.FileMode(0666))
	if err != nil {
		panic(err)
	}
	defer in.Close()

	// seek the reader back from the end of the file last position
	if _, err = in.Seek(0-keep, io.SeekEnd); err != nil {
		panic(err)
	}

	// seek the writer to the beginning of the file
	if _, err := out.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}

	// copy bytes back to the beginning of the file
	if keep, err = io.CopyN(out, in, keep); err != nil {
		panic(err)
	}

	// make sure it's really on disk
	if err = out.Sync(); err != nil {
		panic(err)
	}

	// resize the file
	if err = out.Truncate(keep); err != nil {
		panic(err)
	}

	file.size = keep
}
