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

// Package log handles logging
package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// A Logger is a frontend to a logging system that writes log messages to one or more Sinks.
// A Logger is safe to use concurrently and serializes calls to attached Sinks
type Logger struct {
	mu      sync.Mutex
	enabled Level
	exit    func()
	sinks   []Sink
}

// NewLogger creates a new Logger which will write log messages to the Sinks passed at the Levels enabled.
// The exit func passed will be used for the FATAL Level and should ultimately kill the process.
func NewLogger(enabled Level, exit func(), sinks ...Sink) *Logger {
	return &Logger{
		enabled: enabled,
		exit:    exit,
		sinks:   sinks,
	}
}

// log creates a Msg and writes it to the attached Sinks
func (l *Logger) log(lvl Level, body string) {
	var msg = &Msg{
		Level: lvl,
		Time:  time.Now(),
		Body:  body,
	}

	var ok bool
	_, msg.File, msg.Line, ok = runtime.Caller(2)
	if !ok {
		msg.File = "???"
		msg.Line = 0
	} else {
		// try to trunc the path to $GOPATH/src/
		msg.File = filepath.ToSlash(msg.File)
		if n := strings.LastIndex(msg.File, "/src/"); n > 0 {
			msg.File = msg.File[n+5:]
		} else {
			// give up and just show the file name only
			msg.File = filepath.Base(msg.File)
		}
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	for _, sink := range l.sinks {
		sink.Log(msg)
	}
}

// DebugFunc will execute fn if the DEBUG Level is enabled on the standard Logger
// and will print the returned string to the log
func (l *Logger) DebugFunc(fn func() string) {
	if l.enabled&DEBUG != DEBUG {
		return
	}
	l.log(DEBUG, fn())
}

// Debug will print in the manner of fmt.Print to the standard Logger if the DEBUG Level is enabled
func (l *Logger) Debug(v ...interface{}) {
	if l.enabled&DEBUG != DEBUG {
		return
	}
	l.log(DEBUG, fmt.Sprint(v...))
}

// Debugln will print in the manner of fmt.Println to the standard Logger if the DEBUG Level is enabled
func (l *Logger) Debugln(v ...interface{}) {
	if l.enabled&DEBUG != DEBUG {
		return
	}
	l.log(DEBUG, fmt.Sprintln(v...))
}

// Debugf will print in the manner of fmt.Printf to the standard Logger if the DEBUG Level is enabled
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.enabled&DEBUG != DEBUG {
		return
	}
	l.log(DEBUG, fmt.Sprintf(format, v...))
}

// Print will print in the manner of fmt.Print to the standard Logger if the NORMAL Level is enabled
func (l *Logger) Print(v ...interface{}) {
	if l.enabled&NORMAL != NORMAL {
		return
	}
	l.log(NORMAL, fmt.Sprint(v...))
}

// Println will print in the manner of fmt.Println to the standard Logger if the NORMAL Level is enabled
func (l *Logger) Println(v ...interface{}) {
	if l.enabled&NORMAL != NORMAL {
		return
	}
	l.log(NORMAL, fmt.Sprintln(v...))
}

// Printf will print in the manner of fmt.Printf to the standard Logger if the NORMAL Level is enabled
func (l *Logger) Printf(format string, v ...interface{}) {
	if l.enabled&NORMAL != NORMAL {
		return
	}
	l.log(NORMAL, fmt.Sprintf(format, v...))
}

// Warning will print in the manner of fmt.Print to the standard Logger if the WARNING Level is enabled
func (l *Logger) Warning(v ...interface{}) {
	if l.enabled&WARNING != WARNING {
		return
	}
	l.log(WARNING, fmt.Sprint(v...))
}

// Warningln will print in the manner of fmt.Println to the standard Logger if the WARNING Level is enabled
func (l *Logger) Warningln(v ...interface{}) {
	if l.enabled&WARNING != WARNING {
		return
	}
	l.log(WARNING, fmt.Sprintln(v...))
}

// Warningf will print in the manner of fmt.Printf to the standard Logger if the WARNING Level is enabled
func (l *Logger) Warningf(format string, v ...interface{}) {
	if l.enabled&WARNING != WARNING {
		return
	}
	l.log(WARNING, fmt.Sprintf(format, v...))
}

// Error will print in the manner of fmt.Print to the standard Logger if the ERROR Level is enabled
func (l *Logger) Error(v ...interface{}) {
	if l.enabled&ERROR != ERROR {
		return
	}
	l.log(ERROR, fmt.Sprint(v...))
}

// Errorln will print in the manner of fmt.Println to the standard Logger if the ERROR Level is enabled
func (l *Logger) Errorln(v ...interface{}) {
	if l.enabled&ERROR != ERROR {
		return
	}
	l.log(ERROR, fmt.Sprintln(v...))
}

// Errorf will print in the manner of fmt.Printf to the standard Logger if the ERROR Level is enabled
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.enabled&ERROR != ERROR {
		return
	}
	l.log(ERROR, fmt.Sprintf(format, v...))
}

// Panic will print in the manner of fmt.Print to the standard Logger if the PANIC Level is enabled
// After logging the message, Panic will call panic()
func (l *Logger) Panic(v ...interface{}) {
	if l.enabled&PANIC != PANIC {
		return
	}
	msg := fmt.Sprint(v...)
	l.log(PANIC, msg)
	panic(msg)
}

// Panicln will print in the manner of fmt.Println to the standard Logger if the PANIC Level is enabled
// After logging the message, Panicln will call panic()
func (l *Logger) Panicln(v ...interface{}) {
	if l.enabled&PANIC != PANIC {
		return
	}
	msg := fmt.Sprintln(v...)
	l.log(PANIC, msg)
	panic(msg)
}

// Panicf will print in the manner of fmt.Printf to the standard Logger if the PANIC Level is enabled
// After logging the message, Panicf will call panic()
func (l *Logger) Panicf(format string, v ...interface{}) {
	if l.enabled&PANIC != PANIC {
		return
	}
	msg := fmt.Sprintf(format, v...)
	l.log(PANIC, msg)
	panic(msg)
}

// Fatal will print in the manner of fmt.Print to the standard Logger if the FATAL Level is enabled
// After logging the message, Fatal will call os.Exit(1)
func (l *Logger) Fatal(v ...interface{}) {
	if l.enabled&FATAL != FATAL {
		return
	}
	l.log(FATAL, fmt.Sprint(v...))
	l.exit()
}

// Fatalln will print in the manner of fmt.Println to the standard Logger if the FATAL Level is enabled
// After logging the message, Fatalln will call os.Exit(1)
func (l *Logger) Fatalln(v ...interface{}) {
	if l.enabled&FATAL != FATAL {
		return
	}
	l.log(FATAL, fmt.Sprintln(v...))
	l.exit()
}

// Fatalf will print in the manner of fmt.Printf to the standard Logger if the FATAL Level is enabled
// After logging the message, Fatalf will call os.Exit(1)
func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.enabled&FATAL != FATAL {
		return
	}
	l.log(FATAL, fmt.Sprintf(format, v...))
	l.exit()
}
