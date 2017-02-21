// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"context"
	"runtime"
	"testing"
)

func TestNewEntry(t *testing.T) {
	const wantMsg = "Hello, World!"
	e := NewEntry(0, wantMsg)
	if e.Msg != wantMsg {
		t.Errorf("e.Msg = %q; want %q", e.Msg, wantMsg)
	}
	if e.Time.IsZero() {
		t.Error("e.Time is zero")
	}
	if _, file, line, ok := runtime.Caller(0); ok {
		if e.File != file {
			t.Errorf("e.File = %q; want %q", e.File, file)
		}
		if e.Line <= 0 || e.Line >= line {
			t.Errorf("e.Line = %d; want in range (0, %d)", e.Line, line)
		}
	}
}

func BenchmarkNewEntry(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		NewEntry(0, "hello, world")
	}
}

func BenchmarkDiscard(b *testing.B) {
	ctx := context.Background()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Discard.Log(ctx, NewEntry(0, "hello, world"))
	}
}
