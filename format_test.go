// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"bytes"
	"context"
	"runtime"
	"testing"
)

func TestLogf(t *testing.T) {
	tests := []struct {
		name  string
		f     func(context.Context, Logger, string, ...interface{})
		level Level
	}{
		{"Info", Infof, Info},
		{"Debug", Debugf, Debug},
	}
	const wantMsg = "Hello, World!"
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cl := captureLogger{}
			test.f(context.Background(), &cl, "%s", wantMsg)
			if cl.e.Msg != wantMsg {
				t.Errorf("e.Msg = %q; want %q", cl.e.Msg, wantMsg)
			}
			if cl.e.Level != test.level {
				t.Errorf("e.Level = %v; want %v", cl.e.Level, test.level)
			}
			if cl.e.Time.IsZero() {
				t.Error("e.Time is zero")
			}
			if _, file, line, ok := runtime.Caller(0); ok {
				if cl.e.File != file {
					t.Errorf("e.File = %q; want %q", cl.e.File, file)
				}
				if cl.e.Line <= 0 || cl.e.Line >= line {
					t.Errorf("e.Line = %d; want in range (0, %d)", cl.e.Line, line)
				}
			}
		})
		t.Run(test.name+"_disabled", func(t *testing.T) {
			cl := captureLogger{disabled: true}
			test.f(context.Background(), &cl, "%s", wantMsg)
			if cl.called {
				t.Error("called Log when Logger disabled")
			}
		})
	}
}

type captureLogger struct {
	e        Entry
	called   bool
	disabled bool
}

func (cl *captureLogger) Log(_ context.Context, e Entry) {
	cl.e = e
	cl.called = true
}

func (cl *captureLogger) LogEnabled(Entry) bool { return !cl.disabled }

func BenchmarkDiscardLogf(b *testing.B) {
	ctx := context.Background()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Infof(ctx, Discard, "Hello, %v!", "World")
	}
}

func BenchmarkWriterLogf(b *testing.B) {
	ctx := context.Background()
	buf := new(bytes.Buffer)
	logger := New(buf, "", 0, nil)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		Infof(ctx, logger, "Hello, %v!", "World")
	}
}
