// Copyright 2021 The Zombie Zen Log Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause

//go:build !go1.16
// +build !go1.16

package zstdlog

import (
	"context"
	stdlog "log"

	"zombiezen.com/go/log"
)

func setDefaultOutput(dst log.Logger, opts *Options) {
	stdlog.SetFlags(stdlogFlags)
	stdlog.SetOutput(newWriter(dst, opts))
}
