// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// Infof writes an info message to log.  Its arguments are handled in the manner of fmt.Sprintf.
func Infof(ctx context.Context, logger Logger, format string, args ...interface{}) {
	logf(ctx, logger, Info, format, args)
}

// Debugf writes a debug message to log.  Its arguments are handled in the manner of fmt.Sprintf.
func Debugf(ctx context.Context, logger Logger, format string, args ...interface{}) {
	logf(ctx, logger, Debug, format, args)
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
	// TODO: remove trailing newline, if any
	logger.Log(ctx, ent)
}
