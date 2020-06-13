// Copyright 2017 The Zombie Zen Log Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause

package log

import (
	"bytes"
	"context"
	"runtime"
	"testing"
)

func TestLogf(t *testing.T) {
	tests := []Level{
		Debug,
		Info,
		Warn,
		Error,
	}
	const wantMsg = "Hello, World!"
	for _, level := range tests {
		name := level.String()
		t.Run(name, func(t *testing.T) {
			cl := captureLogger{}
			Logf(context.Background(), &cl, level, "%s", wantMsg)
			if cl.e.Msg != wantMsg {
				t.Errorf("e.Msg = %q; want %q", cl.e.Msg, wantMsg)
			}
			if cl.e.Level != level {
				t.Errorf("e.Level = %v; want %v", cl.e.Level, level)
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
		t.Run(name+"_disabled", func(t *testing.T) {
			cl := captureLogger{disabled: true}
			Logf(context.Background(), &cl, level, "%s", wantMsg)
			if cl.called {
				t.Error("called Log when Logger disabled")
			}
		})
		t.Run(name+"_newline", func(t *testing.T) {
			cl := captureLogger{}
			Logf(context.Background(), &cl, level, "%s\n", wantMsg)
			if cl.e.Msg != wantMsg {
				t.Errorf("e.Msg = %q; want %q", cl.e.Msg, wantMsg)
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
		Logf(ctx, Discard, Info, "Hello, %v!", "World")
	}
}

func BenchmarkWriterLogf(b *testing.B) {
	ctx := context.Background()
	buf := new(bytes.Buffer)
	logger := New(buf, "", 0, nil)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		Logf(ctx, logger, Info, "Hello, %v!", "World")
	}
}
