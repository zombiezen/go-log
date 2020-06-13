// Copyright 2020 The Zombie Zen Log Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause

package log

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var _ Logger = new(LevelFilter)

func TestLevelFilter(t *testing.T) {
	newEntry := func(l Level) Entry {
		return Entry{
			Level: l,
			Msg:   "Hello, World!",
			File:  "foo.go",
			Line:  123,
			Time:  time.Date(2020, time.June, 19, 0, 0, 0, 0, time.UTC),
		}
	}
	tests := []struct {
		name  string
		min   Level
		entry Entry
		want  bool
	}{
		{
			name:  "MinInfo/Info",
			min:   Info,
			entry: newEntry(Info),
			want:  true,
		},
		{
			name:  "MinInfo/Debug",
			min:   Info,
			entry: newEntry(Debug),
			want:  false,
		},
		{
			name:  "MinInfo/Error",
			min:   Info,
			entry: newEntry(Error),
			want:  true,
		},
		{
			name:  "MinWarn/Info",
			min:   Warn,
			entry: newEntry(Info),
			want:  false,
		},
		{
			name:  "MinWarn/Warn",
			min:   Warn,
			entry: newEntry(Warn),
			want:  true,
		},
		{
			name:  "MinWarn/Error",
			min:   Warn,
			entry: newEntry(Error),
			want:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sink := new(captureLogger)
			filter := LevelFilter{
				Min:    test.min,
				Output: sink,
			}

			enableEntry := test.entry
			enableEntry.Msg = ""
			if got := filter.LogEnabled(enableEntry); got != test.want {
				t.Errorf("filter.LogEnabled(entry) = %t; want %t", got, test.want)
			}
			if sink.called {
				t.Fatal("Output.Log called during LogEnabled")
			}

			filter.Log(context.Background(), test.entry)
			if sink.called && !test.want {
				t.Error("filter.Log(entry) logged an entry; wanted drop")
			} else if !sink.called && test.want {
				t.Error("filter.Log(entry) dropped an entry; wanted to log")
			}
			if diff := cmp.Diff(test.entry, sink.e); sink.called && diff != "" {
				t.Errorf("Output.Log entry (-want +got):\n%s", diff)
			}
		})
	}
}
