// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package logutil implements some logging utility functions.
package logutil // import "zombiezen.com/go/log/logutil"

import (
	"context"
	"fmt"

	"zombiezen.com/go/log"
)

// Log writes a message to log.  Its arguments are handled in the manner of fmt.Sprintln.
func Log(ctx context.Context, logger log.Logger, args ...interface{}) {
	ent := log.NewEntry(1, "")
	if !logger.LogEnabled(ent) {
		return
	}
	ent.Msg = fmt.Sprintln(args...)
	ent.Msg = ent.Msg[:len(ent.Msg)-1]
	logger.Log(ctx, ent)
}

// Logf writes a message to log.  Its arguments are handled in the manner of fmt.Sprintf.
func Logf(ctx context.Context, logger log.Logger, format string, args ...interface{}) {
	ent := log.NewEntry(1, "")
	if !logger.LogEnabled(ent) {
		return
	}
	ent.Msg = fmt.Sprintf(format, args...)
	logger.Log(ctx, ent)
}
