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
	"fmt"
	"os"
)

var std = &Logger{
	enabled: ALL,
	exit: func() {
		os.Exit(1)
	},
}

// StandardLogger returns a pointer to the standard Logger
func StandardLogger() *Logger {
	return std
}

// SetEnabled enables on the standard Logger all log Levels with a bit set on the argument
func SetEnabled(enabled Level) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.enabled = enabled
}

// SetExit sets the func to be used for the FATAL Level and should ultimately kill the process.
func SetExit(exit func()) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.exit = exit
}

// SetSinks sets the Sinks to which the standard Logger will write log messages
func SetSinks(sinks ...Sink) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.sinks = sinks
}

// DebugFunc will execute fn if the DEBUG Level is enabled on the standard Logger
// and will print the returned string to the log
func DebugFunc(fn func() string) {
	if std.enabled&DEBUG != DEBUG {
		return
	}
	std.log(DEBUG, fn())
}

// Debug will print in the manner of fmt.Print to the standard Logger if the DEBUG Level is enabled
func Debug(v ...interface{}) {
	if std.enabled&DEBUG != DEBUG {
		return
	}
	std.log(DEBUG, fmt.Sprint(v...))
}

// Debugln will print in the manner of fmt.Println to the standard Logger if the DEBUG Level is enabled
func Debugln(v ...interface{}) {
	if std.enabled&DEBUG != DEBUG {
		return
	}
	std.log(DEBUG, fmt.Sprintln(v...))
}

// Debugf will print in the manner of fmt.Printf to the standard Logger if the DEBUG Level is enabled
func Debugf(format string, v ...interface{}) {
	if std.enabled&DEBUG != DEBUG {
		return
	}
	std.log(DEBUG, fmt.Sprintf(format, v...))
}

// Print will print in the manner of fmt.Print to the standard Logger if the NORMAL Level is enabled
func Print(v ...interface{}) {
	if std.enabled&NORMAL != NORMAL {
		return
	}
	std.log(NORMAL, fmt.Sprint(v...))
}

// Println will print in the manner of fmt.Println to the standard Logger if the NORMAL Level is enabled
func Println(v ...interface{}) {
	if std.enabled&NORMAL != NORMAL {
		return
	}
	std.log(NORMAL, fmt.Sprintln(v...))
}

// Printf will print in the manner of fmt.Printf to the standard Logger if the NORMAL Level is enabled
func Printf(format string, v ...interface{}) {
	if std.enabled&NORMAL != NORMAL {
		return
	}
	std.log(NORMAL, fmt.Sprintf(format, v...))
}

// Warning will print in the manner of fmt.Print to the standard Logger if the WARNING Level is enabled
func Warning(v ...interface{}) {
	if std.enabled&WARNING != WARNING {
		return
	}
	std.log(WARNING, fmt.Sprint(v...))
}

// Warningln will print in the manner of fmt.Println to the standard Logger if the WARNING Level is enabled
func Warningln(v ...interface{}) {
	if std.enabled&WARNING != WARNING {
		return
	}
	std.log(WARNING, fmt.Sprintln(v...))
}

// Warningf will print in the manner of fmt.Printf to the standard Logger if the WARNING Level is enabled
func Warningf(format string, v ...interface{}) {
	if std.enabled&WARNING != WARNING {
		return
	}
	std.log(WARNING, fmt.Sprintf(format, v...))
}

// Error will print in the manner of fmt.Print to the standard Logger if the ERROR Level is enabled
func Error(v ...interface{}) {
	if std.enabled&ERROR != ERROR {
		return
	}
	std.log(ERROR, fmt.Sprint(v...))
}

// Errorln will print in the manner of fmt.Println to the standard Logger if the ERROR Level is enabled
func Errorln(v ...interface{}) {
	if std.enabled&ERROR != ERROR {
		return
	}
	std.log(ERROR, fmt.Sprintln(v...))
}

// Errorf will print in the manner of fmt.Printf to the standard Logger if the ERROR Level is enabled
func Errorf(format string, v ...interface{}) {
	if std.enabled&ERROR != ERROR {
		return
	}
	std.log(ERROR, fmt.Sprintf(format, v...))
}

// Panic will print in the manner of fmt.Print to the standard Logger if the PANIC Level is enabled
// After logging the message, Panic will call panic()
func Panic(v ...interface{}) {
	if std.enabled&PANIC != PANIC {
		return
	}
	msg := fmt.Sprint(v...)
	std.log(PANIC, msg)
	panic(msg)
}

// Panicln will print in the manner of fmt.Println to the standard Logger if the PANIC Level is enabled
// After logging the message, Panicln will call panic()
func Panicln(v ...interface{}) {
	if std.enabled&PANIC != PANIC {
		return
	}
	msg := fmt.Sprintln(v...)
	std.log(PANIC, msg)
	panic(msg)
}

// Panicf will print in the manner of fmt.Printf to the standard Logger if the PANIC Level is enabled
// After logging the message, Panicf will call panic()
func Panicf(format string, v ...interface{}) {
	if std.enabled&PANIC != PANIC {
		return
	}
	msg := fmt.Sprintf(format, v...)
	std.log(PANIC, msg)
	panic(msg)
}

// Fatal will print in the manner of fmt.Print to the standard Logger if the FATAL Level is enabled
// After logging the message, Fatal will call os.Exit(1)
func Fatal(v ...interface{}) {
	if std.enabled&FATAL != FATAL {
		return
	}
	std.log(FATAL, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalln will print in the manner of fmt.Println to the standard Logger if the FATAL Level is enabled
// After logging the message, Fatalln will call os.Exit(1)
func Fatalln(v ...interface{}) {
	if std.enabled&FATAL != FATAL {
		return
	}
	std.log(FATAL, fmt.Sprintln(v...))
	os.Exit(1)
}

// Fatalf will print in the manner of fmt.Printf to the standard Logger if the FATAL Level is enabled
// After logging the message, Fatalf will call os.Exit(1)
func Fatalf(format string, v ...interface{}) {
	if std.enabled&FATAL != FATAL {
		return
	}
	std.log(FATAL, fmt.Sprintf(format, v...))
	os.Exit(1)
}
