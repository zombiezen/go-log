// Copyright 2017 The Zombie Zen Log Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause

package log

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// Infof writes an info message to the default Logger.  Its arguments are handled in the manner of fmt.Sprintf.
func Infof(ctx context.Context, format string, args ...interface{}) {
	if false {
		// Enable printf checking in go vet.
		_ = fmt.Sprintf(format, args...)
	}
	logf(ctx, Default(), Info, format, args)
}

// Debugf writes a debug message to the default Logger.  Its arguments are handled in the manner of fmt.Sprintf.
func Debugf(ctx context.Context, format string, args ...interface{}) {
	if false {
		// Enable printf checking in go vet.
		_ = fmt.Sprintf(format, args...)
	}
	logf(ctx, Default(), Debug, format, args)
}

// Warnf writes a warning message to the default Logger.  Its arguments are handled in the manner of fmt.Sprintf.
func Warnf(ctx context.Context, format string, args ...interface{}) {
	if false {
		// Enable printf checking in go vet.
		_ = fmt.Sprintf(format, args...)
	}
	logf(ctx, Default(), Warn, format, args)
}

// Errorf writes an error message to the default Logger.  Its arguments are handled in the manner of fmt.Sprintf.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	if false {
		// Enable printf checking in go vet.
		_ = fmt.Sprintf(format, args...)
	}
	logf(ctx, Default(), Error, format, args)
}

// Logf writes a message to a Logger.  Its arguments are handled in the manner of fmt.Sprintf.
func Logf(ctx context.Context, logger Logger, level Level, format string, args ...interface{}) {
	if false {
		// Enable printf checking in go vet.
		_ = fmt.Sprintf(format, args...)
	}
	logf(ctx, logger, level, format, args)
}

func logf(ctx context.Context, logger Logger, level Level, format string, args []interface{}) {
	ent := Entry{Time: time.Now(), Level: level}
	if _, file, line, ok := runtime.Caller(2); ok {
		ent.File = file
		ent.Line = line
	}
	if !logger.LogEnabled(ent) {
		return
	}
	ent.Msg = fmt.Sprintf(format, args...)
	if n := len(ent.Msg); n > 0 && ent.Msg[n-1] == '\n' {
		ent.Msg = ent.Msg[:n-1]
	}
	logger.Log(ctx, ent)
}
