// Copyright 2017 The Zombie Zen Log Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause

// Package testlog provides a Logger that writes to a *testing.T or *testing.B.
// See the examples for how to set this up.
package testlog

import (
	"context"

	"zombiezen.com/go/log"
)

// Logger writes to a *testing.T or *testing.B that comes from the Context.
type Logger struct {
	// Fallback is the Logger used if the Context does not have a TB.
	// If nil, then log.Discard is assumed.
	Fallback log.Logger
}

// Main sets the default logger to a testlog.Logger. fallback may be nil.
// Main is intended to be called in TestMain.
func Main(fallback log.Logger) {
	log.SetDefault(Logger{fallback})
}

// Log writes to the TB in ctx or l.Fallback.
func (l Logger) Log(ctx context.Context, e log.Entry) {
	tb, _ := ctx.Value(ctxKey{}).(TB)
	if tb == nil {
		if l.Fallback != nil {
			l.Fallback.Log(ctx, e)
		}
		return
	}
	switch e.Level {
	case log.Warn:
		tb.Log("WARN: " + e.Msg)
	case log.Error:
		tb.Log("ERROR: " + e.Msg)
	default:
		tb.Log(e.Msg)
	}
}

// LogEnabled always returns true.
func (l Logger) LogEnabled(e log.Entry) bool {
	return true
}

type ctxKey struct{}

// WithTB returns a Context derived from parent that will use tb to log
// when sending an entry to this package's Logger.
func WithTB(parent context.Context, tb TB) context.Context {
	return context.WithValue(parent, ctxKey{}, tb)
}

// TB is the interface provided by a *testing.T or a *testing.B
// that is needed for Logger.
type TB interface {
	Log(...interface{})
}
