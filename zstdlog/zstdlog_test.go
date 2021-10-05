// Copyright 2021 The Zombie Zen Log Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause

package zstdlog

import (
	"context"
	"io/ioutil"
	stdlog "log"
	"path/filepath"
	"testing"
	"time"

	"zombiezen.com/go/log"
)

func TestNew(t *testing.T) {
	testAdapter(t, New)
}

func TestSetOutput(t *testing.T) {
	testAdapter(t, func(dst log.Logger, opts *Options) *stdlog.Logger {
		src := stdlog.New(ioutil.Discard, "foo", stdlog.LstdFlags)
		SetOutput(dst, src, opts)
		return src
	})
}

func testAdapter(t *testing.T, newLogger func(log.Logger, *Options) *stdlog.Logger) {
	t.Helper()

	t.Run("DefaultOptions", func(t *testing.T) {
		const msg = "Hello, World!"
		l := new(captureLogger)

		stdlogger := newLogger(l, nil)
		stdlogger.Print(msg)

		if !l.called {
			t.Fatal("Logger.Print did not trigger call to Log")
		}
		if l.e.Msg != msg {
			t.Errorf("Msg = %q; want %q", l.e.Msg, msg)
		}
		if got, want := filepath.Base(l.e.File), "zstdlog_test.go"; got != want {
			t.Errorf("File = %q; want %q", got, msg)
		}
		if l.e.Time.IsZero() {
			t.Error("Time is not set")
		} else if now := time.Now(); l.e.Time.After(now) {
			t.Errorf("Time = %v; want before %v", l.e.Time, now)
		}
		if l.e.Level != log.Info {
			t.Errorf("Level = %v; want %v", l.e.Level, log.Info)
		}
		if l.ctx == nil {
			t.Error("nil Context passed to Log")
		}
	})

	t.Run("ErrorLevel", func(t *testing.T) {
		const msg = "Hello, World!"
		l := new(captureLogger)

		stdlogger := newLogger(l, &Options{
			Level: log.Error,
		})
		stdlogger.Print(msg)

		if !l.called {
			t.Fatal("Logger.Print did not trigger call to Log")
		}
		if l.e.Msg != msg {
			t.Errorf("Msg = %q; want %q", l.e.Msg, msg)
		}
		if l.e.Level != log.Error {
			t.Errorf("Level = %v; want %v", l.e.Level, log.Error)
		}
	})

	t.Run("NonDefaultContext", func(t *testing.T) {
		type myKey struct{}
		const msg = "Hello, World!"
		l := new(captureLogger)
		const myValue = "foo"
		ctx := context.WithValue(context.Background(), myKey{}, myValue)

		stdlogger := newLogger(l, &Options{
			Context: ctx,
		})
		stdlogger.Print(msg)

		if !l.called {
			t.Fatal("Logger.Print did not trigger call to Log")
		}
		if l.e.Msg != msg {
			t.Errorf("Msg = %q; want %q", l.e.Msg, msg)
		}
		if l.ctx == nil {
			t.Error("nil Context passed to Log")
		} else if v := l.ctx.Value(myKey{}); v != myValue {
			t.Errorf("Context.Value(myKey{}) = %#v; want %#v", v, myValue)
		}
	})
}

type captureLogger struct {
	ctx    context.Context
	e      log.Entry
	called bool
}

func (cl *captureLogger) Log(ctx context.Context, e log.Entry) {
	cl.ctx = ctx
	cl.e = e
	cl.called = true
}

func (cl *captureLogger) LogEnabled(log.Entry) bool {
	return true
}

func BenchmarkPrintf(b *testing.B) {
	stdLogger := New(new(captureLogger), nil)
	b.ResetTimer()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		stdLogger.Printf("Hello World")
	}
}

func BenchmarkDiscard(b *testing.B) {
	stdLogger := New(log.Discard, nil)
	b.ResetTimer()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		stdLogger.Printf("Hello World")
	}
}
