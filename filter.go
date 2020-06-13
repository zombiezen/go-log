// Copyright 2020 The Zombie Zen Log Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause

package log

import "context"

// A LevelFilter filters out entries below a minimum log level before sending
// them to another logger.
type LevelFilter struct {
	Min    Level
	Output Logger
}

// Log sends the entry to the filter's output if the entry's level is at least
// the filter's minimum.
func (f *LevelFilter) Log(ctx context.Context, e Entry) {
	if e.Level < f.Min {
		return
	}
	f.Output.Log(ctx, e)
}

// LogEnabled returns false if the entry's level is below the filter's minimum,
// otherwise it returns the result of f.Output.LogEnabled(e).
func (f *LevelFilter) LogEnabled(e Entry) bool {
	if e.Level < f.Min {
		return false
	}
	return f.Output.LogEnabled(e)
}
