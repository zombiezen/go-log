// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"context"
	"os"
	"sync"
	"sync/atomic"
)

var (
	defaultLogger    atomicLogger
	setDefaultLogger sync.Once

	fallback = New(os.Stderr, "", LstdFlags, nil)
)

// Default returns the global logger.
// Until SetDefault is called, the returned Logger will send all
// entries to stderr.
func Default() Logger {
	return &defaultLogger
}

// SetDefault sets the global logger.  It can only be called once.
func SetDefault(l Logger) {
	ok := false
	setDefaultLogger.Do(func() {
		if l == nil {
			panic("log.SetDefaultLogger(nil)")
		}
		defaultLogger.out.Store(l)
		ok = true
	})
	if !ok {
		panic("log.SetDefaultLogger called more than once")
	}
}

type atomicLogger struct {
	out atomic.Value
}

func (l *atomicLogger) Log(ctx context.Context, ent Entry) {
	l.logger().Log(ctx, ent)
}

func (l *atomicLogger) LogEnabled(ent Entry) bool {
	return l.logger().LogEnabled(ent)
}

func (l *atomicLogger) logger() Logger {
	out := l.out.Load()
	if out == nil {
		return fallback
	}
	return out.(Logger)
}
