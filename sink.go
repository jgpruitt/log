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

package log

import (
	"bytes"
	"os"
)

// A Sink is a consumer of log messages which will likely render them to some I/O device
type Sink interface {
	// Log consumes a Msg and presumably writes it to a device.
	// Log will not be called concurrently.
	Log(m *Msg)
}

// A StdOutSink is a Sink which writes to STDOUT.
// It will ignore log messages at levels WARNING, ERROR, PANIC, and FATAL.
// It will write log messages at levels DEBUG and NORMAL.
type StdOutSink struct {
	buf bytes.Buffer
}

// Log writes a log message to STDOUT
func (s *StdOutSink) Log(m *Msg) {
	// ignore the levels listed
	if m.Level&(WARNING|ERROR|PANIC|FATAL) > 0 {
		return
	}
	buf := &s.buf
	m.PrintDate(buf)
	buf.WriteString(" ")
	m.PrintTime(buf)
	buf.WriteString(" ")
	m.PrintLevel(buf)
	buf.WriteString(" ")
	m.PrintFileLine(buf)
	buf.WriteString(" ")
	m.PrintMsg(buf)
	buf.WriteTo(os.Stdout)
	buf.Reset()
}

// A StdErrSink is a Sink which writes to STDERR.
// It will ignore log messages at levels DEBUG and NORMAL.
// It will write log messages at levels WARNING, ERROR, PANIC, and FATAL.
type StdErrSink struct {
	buf bytes.Buffer
}

// Log writes a log message to STDERR
func (s *StdErrSink) Log(m *Msg) {
	// ignore the levels listed
	if m.Level&(DEBUG|NORMAL) > 0 {
		return
	}
	buf := &s.buf
	m.PrintDate(buf)
	buf.WriteString(" ")
	m.PrintTime(buf)
	buf.WriteString(" ")
	m.PrintLevel(buf)
	buf.WriteString(" ")
	m.PrintFileLine(buf)
	buf.WriteString(" ")
	m.PrintMsg(buf)
	buf.WriteTo(os.Stderr)
	buf.Reset()
}