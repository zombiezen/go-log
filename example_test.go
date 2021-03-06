// Copyright 2017 The Zombie Zen Log Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause

// Sample application to demonstrate the log package.

package log_test

import (
	"context"
	"os"

	"zombiezen.com/go/log"
)

func Example() {
	// Initialize the global logger.
	// This should only happen in main and before any log statements.
	stdoutLog := log.New(os.Stdout, "", 0, nil)
	log.SetDefault(stdoutLog)

	// Once the logger is set, you can log from anywhere.
	ctx := context.Background()
	log.Infof(ctx, "Hello, World!")

	// Output:
	// Hello, World!
}

// TODO(someday): ExampleLevelFilter should be an output test, but SetDefault
// can only be called once.

func ExampleLevelFilter() {
	log.SetDefault(&log.LevelFilter{
		Min:    log.Warn, // Only show warnings or above
		Output: log.New(os.Stdout, "", 0, nil),
	})
	ctx := context.Background()
	log.Infof(ctx, "This won't show up.")
	log.Warnf(ctx, "Only Warn or higher will show up.")
}
