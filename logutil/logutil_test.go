// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logutil

import (
	"bytes"
	"context"
	"testing"

	"zombiezen.com/go/log"
)

func BenchmarkDiscardLogf(b *testing.B) {
	ctx := context.Background()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Logf(ctx, log.Discard, "Hello, %v!", "World")
	}
}

func BenchmarkWriterLogf(b *testing.B) {
	ctx := context.Background()
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0, nil)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		Logf(ctx, logger, "Hello, %v!", "World")
	}
}
