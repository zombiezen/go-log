// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package log provides an interface for logging diagnostic messages
// independent of output medium.
package log // import "zombiezen.com/go/log"

import (
	"context"
	"time"
)

// Logger is the interface that wraps the Log method.
//
// Log sends an Entry to the underlying log sink.  Each call will be sent
// to the sink in the order that Log was called, but Log may not wait for
// delivery to be acknowledged.  As such, Log should ignore any deadline
// set on the Context.  Log must be safe to call from multiple goroutines.
//
// LogEnabled returns false if Log will no-op for a particular Entry, which
// may or may not have Msg filled in.
type Logger interface {
	Log(context.Context, Entry)
	LogEnabled(Entry) bool
}

// Entry is a single log record.
type Entry struct {
	Msg string

	// Time is the timestamp that the entry was created.
	Time time.Time

	// Level is the verbosity level for the entry. A Logger may skip
	// writing an Entry based on its level.
	Level Level

	// File name and line number of the code that initiated the log as
	// reported by runtime.Caller.
	File string
	Line int
}

// Append appends a formatted entry to a buffer.
// flag controls the formatting of the entry.
// Even if ent.Msg ends in a newline, the last byte appended to buf will
// never be a newline.
func (ent Entry) Append(buf []byte, flag int) []byte {
	file, line := ent.File, ent.Line
	if file == "" {
		file, line = "???", 0
	}
	buf = formatHeader(buf, flag, ent.Time, ent.Level, file, line)
	s := ent.Msg
	if n := len(s); n == 0 || s[n-1] == '\n' {
		s = s[:n-1]
	}
	buf = append(buf, s...)
	return buf
}

// String returns a formatted entry string with the default options.
func (ent Entry) String() string {
	return string(ent.Append(nil, LstdFlags))
}

func formatHeader(buf []byte, flag int, t time.Time, level Level, file string, line int) []byte {
	if flag&LUTC != 0 {
		t = t.UTC()
	}
	if flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if flag&Ldate != 0 {
			year, month, day := t.Date()
			buf = itoa(buf, year, 4)
			buf = append(buf, '/')
			buf = itoa(buf, int(month), 2)
			buf = append(buf, '/')
			buf = itoa(buf, day, 2)
			buf = append(buf, ' ')
		}
		if flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			buf = itoa(buf, hour, 2)
			buf = append(buf, ':')
			buf = itoa(buf, min, 2)
			buf = append(buf, ':')
			buf = itoa(buf, sec, 2)
			if flag&Lmicroseconds != 0 {
				buf = append(buf, '.')
				buf = itoa(buf, t.Nanosecond()/1e3, 6)
			}
			buf = append(buf, ' ')
		}
	}
	if flag&(Llevel|Lshortfile|Llongfile) != 0 {
		if flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		switch {
		case flag&Llevel != 0 && flag&(Lshortfile|Llongfile) == 0:
			buf = append(buf, entryLevel(level)...)
		case flag&Llevel == 0 && flag&(Lshortfile|Llongfile) != 0:
			buf = append(buf, file...)
			buf = append(buf, ':')
			buf = itoa(buf, line, -1)
		default:
			buf = append(buf, entryLevel(level)...)
			buf = append(buf, ' ')
			buf = append(buf, file...)
			buf = append(buf, ':')
			buf = itoa(buf, line, -1)
		}
		buf = append(buf, ": "...)
	}
	return buf
}

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
func itoa(buf []byte, i int, wid int) []byte {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	return append(buf, b[bp:]...)
}

type discard struct{}

func (discard) Log(ctx context.Context, ent Entry) {}

func (discard) LogEnabled(Entry) bool { return false }

// Discard is the no-op Logger.
var Discard Logger = discard{}

// Level is a hint at the audience of the log.
// Lower values of Level mean more verbose.
type Level int

// Predefined logging levels.
const (
	Debug Level = -10 // messages for developers
	Info  Level = 0   // messages for users
	Warn  Level = 10  // warnings for users
	Error Level = 20  // failures
)

// String returns the constant name of the level.
func (l Level) String() string {
	switch l {
	case Debug:
		return "Debug"
	case Info:
		return "Info"
	case Warn:
		return "Warn"
	case Error:
		return "Error"
	}
	var buf []byte
	buf = append(buf, "Level("...)
	buf = itoa(buf, int(l), -1)
	buf = append(buf, ')')
	return string(buf)
}

func entryLevel(l Level) string {
	switch l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	default:
		return "???"
	}
}
