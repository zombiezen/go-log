// Copyright 2021 The Zombie Zen Log Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause

// Package zstdlog provides functions to support interoperation between the
// standard library logger and zombiezen.com/go/log.
package zstdlog

import (
	"context"
	"fmt"
	stdlog "log"
	"strconv"
	"strings"
	"time"

	"zombiezen.com/go/log"
)

// Options is the set of optional arguments to New, SetOutput, and SetDefaultOutput.
type Options struct {
	// Context is used if non-nil when logging entries. Otherwise,
	// context.Background() is used.
	Context context.Context
	// Level is used for all created entries. Defaults to log.Info (the zero value).
	Level log.Level
}

// New returns a new standard library logger that writes to the given
// zombiezen.com/go/log logger.
func New(dst log.Logger, opts *Options) *stdlog.Logger {
	w := newWriter(dst, opts)
	return stdlog.New(w, "", stdlogFlags)
}

const stdlogFlags = stdlog.Ldate |
	stdlog.Ltime |
	stdlog.Lmicroseconds |
	stdlog.LUTC |
	stdlog.Llongfile

// SetOutput configures the standard library logger src to write to the given
// zombiezen.com/go/log.Logger dst. opts may be nil, in which case it is treated
// the same as if new(Options) were passed.
func SetOutput(dst log.Logger, src *stdlog.Logger, opts *Options) {
	src.SetFlags(stdlogFlags)
	src.SetPrefix("")
	src.SetOutput(newWriter(dst, opts))
}

// SetDefaultOutput configures the default standard library logger to write to
// the given zombiezen.com/go/log.Logger dst. opts may be nil, in which case it
// is treated the same as if new(Options) were passed.
func SetDefaultOutput(dst log.Logger, opts *Options) {
	setDefaultOutput(dst, opts)
}

type writer struct {
	ctx   context.Context
	level log.Level
	dst   log.Logger
}

func newWriter(dst log.Logger, opts *Options) *writer {
	w := &writer{
		ctx: context.Background(),
		dst: dst,
	}
	if opts != nil {
		if opts.Context != nil {
			w.ctx = opts.Context
		}
		w.level = opts.Level
	}
	return w
}

func (w *writer) Write(p []byte) (int, error) {
	const layout = "2006/01/02 15:04:05.999999 "
	ps := string(p)
	var ent log.Entry
	var err error
	ent.Time, err = time.Parse(layout, ps[:len(layout)])
	if err != nil {
		return 0, fmt.Errorf("log entry %q: invalid format: %v", p, err)
	}
	ent.Time = ent.Time.Local()

	const msgSeparator = ": "
	fileLineEnd := strings.Index(ps[len(layout):], msgSeparator)
	if fileLineEnd == -1 {
		return 0, fmt.Errorf("log entry %q: invalid format", p)
	}
	fileLineEnd += len(layout)
	fileLine := ps[len(layout):fileLineEnd]
	if fileEnd := strings.LastIndex(fileLine, ":"); fileEnd == -1 {
		ent.File = fileLine
	} else {
		ent.File = fileLine[:fileEnd]
		ent.Line, err = strconv.Atoi(fileLine[fileEnd+len(":"):])
		if err != nil {
			return 0, fmt.Errorf("log entry %q: invalid format: %v", p, err)
		}
	}
	if ent.File == "???" {
		ent.File = ""
		ent.Line = 0
	}

	ent.Level = w.level
	ent.Msg = strings.TrimSuffix(ps[fileLineEnd+len(msgSeparator):], "\n")
	w.dst.Log(w.ctx, ent)
	return len(p), nil
}
