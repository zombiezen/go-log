// Copyright 2021 The Zombie Zen Log Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause

package zstdlog_test

import (
	"context"
	stdlog "log"
	"net/http"
	"os"

	"zombiezen.com/go/log"
	"zombiezen.com/go/log/zstdlog"
)

// The most common usage of this package is to send log messages from the
// standard library log package to the zombiezen.com/go/log package.
func Example() {
	// Configure the zombiezen.com/go/log default logger as usual:
	log.SetDefault(log.New(os.Stdout, "", log.ShowLevel, nil))

	// SetDefaultOutput sets the standard library global logger's output to the
	// zombiezen.com/go/log global logger.
	zstdlog.SetDefaultOutput(log.Default(), nil)

	// Now log.Print will be sent to zombiezen.com/go/log with info level:
	stdlog.Print("Hello, World!")
	// Output:
	// INFO: Hello, World!
}

func ExampleNew() {
	ctx := context.Background()

	// You can use zstdlog.New to create a new standard library logger that emits
	// error-level logs.
	errorLog := zstdlog.New(log.Default(), &zstdlog.Options{
		Context: ctx,
		Level:   log.Error,
	})

	// Then you can use that logger in places that use a *log.Logger instead of a
	// more general interface.
	srv := &http.Server{
		Addr:     ":8080",
		ErrorLog: errorLog,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Errorf(ctx, "%v", err)
	}
}
