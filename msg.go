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
	"fmt"
	"io"
	"time"
)

// Level describes the type of log message
type Level uint8

const (
	// DEBUG is meant for the finest grained detailed logging
	DEBUG Level = 1 << iota
	// NORMAL is meant for normal run-of-the-mill logging
	NORMAL
	// WARNING is meant for logging that may need extra attention but isn't a problem
	WARNING
	// ERROR is meant for logging that needs extra attention due to problems
	ERROR
	// PANIC is meant for logging problems just prior to panicking
	PANIC
	// FATAL is meant for logging extreme problems just prior to exiting the process
	FATAL

	// ALL is a bit mask enabling all the log Levels
	ALL = DEBUG | NORMAL | WARNING | ERROR | PANIC | FATAL
)

// String converts the receiving Level to a string.
// String returns "[?????]" if the Level is unknown.
func (lvl Level) String() string {
	switch lvl {
	case DEBUG:
		return "[DEBUG]"
	case NORMAL:
		return "[_____]"
	case WARNING:
		return "[WARN!]"
	case ERROR:
		return "[ERROR]"
	case PANIC:
		return "[PANIC]"
	case FATAL:
		return "[FATAL]"
	default:
		return "[?????]"
	}
}

// Msg is a log message
type Msg struct {
	Level Level
	Time  time.Time
	File  string
	Line  int
	Body  string
}

// PrintDate prints the date of the receiving Msg to w
// returning the number of bytes written and any write error encountered.
func (m *Msg) PrintDate(w io.Writer) (n int, err error) {
	var mth time.Month
	year, mth, day := m.Time.Date()
	return fmt.Fprintf(w, "%04d-%02d-%02d", year, int(mth), day)
}

// PrintTime prints the time of the receiving Msg to w
// returning number of bytes written and any write error encountered.
func (m *Msg) PrintTime(w io.Writer) (n int, err error) {
	hour, min, sec := m.Time.Clock()
	micro := m.Time.Nanosecond() / 1000
	return fmt.Fprintf(w, "%02d:%02d:%02d.%06d", hour, min, sec, micro)
}

// PrintLevel prints the Level of the receiving Msg to w
// returning the number of bytes written and any write error encountered.
func (m *Msg) PrintLevel(w io.Writer) (n int, err error) {
	return fmt.Fprintf(w, "%s", m.Level.String())
}

// PrintFileLine prints the File and Line number of the receiving Msg to w
// returning the number of bytes written and any write error encountered.
func (m *Msg) PrintFileLine(w io.Writer) (n int, err error) {
	return fmt.Fprintf(w, "%s:%d", m.File, m.Line)
}

// PrintMsg prints the string Msg of the receiving Msg type to w
// returning the number of bytes written an any write errors encountered.
func (m *Msg) PrintMsg(w io.Writer) (n int, err error) {
	if len(m.Body) == 0 {
		return 0, nil
	}

	if m.Body[len(m.Body)-1] != '\n' {
		return fmt.Fprintln(w, m.Body)
	}
	return fmt.Fprint(w, m.Body)
}
