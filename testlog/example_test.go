// Copyright 2017 The Zombie Zen Log Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause

package testlog_test

import (
	"context"
	"os"
	"testing"

	"zombiezen.com/go/log"
	"zombiezen.com/go/log/testlog"
)

var (
	m *testing.M
	t *testing.T
)

func ExampleMain() {
	// Inside your TestMain. m is a *testing.M.
	testlog.Main(nil)
	os.Exit(m.Run())
}

func ExampleWithTB() {
	// Inside a test. t is a *testing.T.
	// You must have set the default logger using Main or log.SetDefault.
	ctx := testlog.WithTB(context.Background(), t)
	log.Infof(ctx, "Log to *testing.T")
}
