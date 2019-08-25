// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
