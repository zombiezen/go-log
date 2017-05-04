// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"context"
	"io"
	"sync"
)

// Flags define which text to prefix to each log entry in a Writer.
type Flags uint

// Bits or'ed together to control what's printed.
// There is no control over the order they appear (the order listed
// here) or the format they present (as described in the comments).
// The prefix is followed by a colon only when ShowFile, ShortFile, or ShowLevel
// is specified.
//
// For example, flags ShowDate | ShowTime produce,
//
//	2009/01/23 01:23:23 message
//
// while flags ShowDate | ShowTime | Microseconds | ShowFile produce,
//
//	2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
const (
	ShowDate     Flags = 1 << iota // the date in the local time zone: 2009/01/23
	ShowTime                       // the time in the local time zone: 01:23:23
	Microseconds                   // microsecond resolution: 01:23:23.123123.  assumes ShowTime.
	ShowFile                       // full file name and line number: /a/b/c/d.go:23
	ShortFile                      // final file name element and line number: d.go:23. overrides ShowFile
	UTC                            // if ShowDate or ShowTime is set, use UTC rather than the local time zone
	ShowLevel                      // level name, all caps. INFO

	StdFlags Flags = ShowDate | ShowTime | ShowLevel // initial values for the standard logger
)

const allFlags = ShowDate | ShowTime | Microseconds | ShowFile | ShortFile | UTC | ShowLevel

// String formats the flags as |-separated constant names,
// like "ShowDate|ShowFile|ShowLevel".
func (f Flags) String() string {
	if f == 0 {
		return "0"
	}
	var buf []byte
	if f&ShowDate != 0 {
		buf = append(buf, "ShowDate"...)
	}
	if f&ShowTime != 0 {
		if len(buf) > 0 {
			buf = append(buf, '|')
		}
		buf = append(buf, "ShowTime"...)
	}
	if f&Microseconds != 0 {
		if len(buf) > 0 {
			buf = append(buf, '|')
		}
		buf = append(buf, "ShowMicroseconds"...)
	}
	if f&ShowFile != 0 {
		if len(buf) > 0 {
			buf = append(buf, '|')
		}
		buf = append(buf, "ShowFile"...)
	}
	if f&ShortFile != 0 {
		if len(buf) > 0 {
			buf = append(buf, '|')
		}
		buf = append(buf, "ShortFile"...)
	}
	if f&UTC != 0 {
		if len(buf) > 0 {
			buf = append(buf, '|')
		}
		buf = append(buf, "UTC"...)
	}
	if f&ShowLevel != 0 {
		if len(buf) > 0 {
			buf = append(buf, '|')
		}
		buf = append(buf, "ShowLevel"...)
	}
	if others := f & allFlags; others != 0 {
		if len(buf) > 0 {
			buf = append(buf, '|')
		}
		buf = itoa(buf, int(others), -1)
	}
	return string(buf)
}

// A Writer implements Logger by writing lines of output to an io.Writer.
// Each logging operation makes a single call to the Writer's Write method.
// A Logger can be used simultaneously from multiple goroutines; it
// guarantees to serialize access to the Writer.
type Writer struct {
	errFunc func(context.Context, error) // called if out returns error
	prefix  string                       // prefix to write at beginning of each line
	flag    Flags                        // properties

	mu  sync.Mutex // ensures atomic writes; protects the following fields
	out io.Writer  // destination for output
	buf []byte     // for accumulating text to write
}

// New creates a new Writer. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// The flag argument defines the logging properties.
// If not nil, the errFunc argument is called when w.Write returns an error.
// errFunc must be safe to call from multiple goroutines and should be fast, as it blocks Log returning.
func New(w io.Writer, prefix string, flag Flags, errFunc func(context.Context, error)) *Writer {
	return &Writer{out: w, prefix: prefix, flag: flag, errFunc: errFunc}
}

// Log writes the output for a logging entry. A newline is appended if
// the last character of s is not already a newline.
func (w *Writer) Log(ctx context.Context, ent Entry) {
	defer w.mu.Unlock()
	w.mu.Lock()
	w.buf = append(w.buf[:0], w.prefix...)
	w.buf = ent.Append(w.buf, w.flag)
	w.buf = append(w.buf, '\n')
	_, err := w.out.Write(w.buf)
	if err != nil && w.errFunc != nil {
		w.errFunc(ctx, err)
	}
}

// LogEnabled always returns true.
func (w *Writer) LogEnabled(Entry) bool { return true }
