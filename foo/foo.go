// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Sample application to demonstrate the log package.

package main

import (
	"context"
	"os"

	"zombiezen.com/go/log"
)

var myLog = log.Default()

func init() {
	ctx := context.Background()
	log.Infof(ctx, myLog, "Hello during package init!")
}

func main() {
	ctx := context.Background()
	initLog()
	log.Infof(ctx, myLog, "Hello after init!")
}

func initLog() {
	stderrLog := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile, nil)
	log.SetDefault(stderrLog)
}
