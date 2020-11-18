// Copyright 2017 The Zombie Zen Log Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause

package log

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestWriter(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		flag   Flags
		entry  Entry
		want   string
	}{
		{
			name: "basic no-header",
			entry: Entry{
				Msg:  "Hello, World!",
				Time: time.Date(2017, time.February, 17, 1, 2, 3, 456789000, time.UTC),
				File: "foo/bar.go",
				Line: 278,
			},
			want: "Hello, World!\n",
		},
		{
			name: "basic no-header empty",
			entry: Entry{
				Msg:  "",
				Time: time.Date(2017, time.February, 17, 1, 2, 3, 456789000, time.UTC),
				File: "foo/bar.go",
				Line: 278,
			},
			want: "\n",
		},
		{
			name: "debug level",
			flag: ShowLevel,
			entry: Entry{
				Msg:   "Hello, World!",
				Level: Debug,
				Time:  time.Date(2017, time.February, 17, 1, 2, 3, 456789000, time.UTC),
				File:  "foo/bar.go",
				Line:  278,
			},
			want: "DEBUG: Hello, World!\n",
		},
		{
			name: "info level",
			flag: ShowLevel,
			entry: Entry{
				Msg:   "Hello, World!",
				Level: Info,
				Time:  time.Date(2017, time.February, 17, 1, 2, 3, 456789000, time.UTC),
				File:  "foo/bar.go",
				Line:  278,
			},
			want: "INFO: Hello, World!\n",
		},
		{
			name: "warn level",
			flag: ShowLevel,
			entry: Entry{
				Msg:   "Hello, World!",
				Level: Warn,
				Time:  time.Date(2017, time.February, 17, 1, 2, 3, 456789000, time.UTC),
				File:  "foo/bar.go",
				Line:  278,
			},
			want: "WARN: Hello, World!\n",
		},
		{
			name: "error level",
			flag: ShowLevel,
			entry: Entry{
				Msg:   "Hello, World!",
				Level: Error,
				Time:  time.Date(2017, time.February, 17, 1, 2, 3, 456789000, time.UTC),
				File:  "foo/bar.go",
				Line:  278,
			},
			want: "ERROR: Hello, World!\n",
		},
		{
			name:   "prefix no-header",
			prefix: "///prefix|||",
			entry: Entry{
				Msg:  "Hello, World!",
				Time: time.Date(2017, time.February, 17, 1, 2, 3, 456789000, time.UTC),
				File: "foo/bar.go",
				Line: 278,
			},
			want: "///prefix|||Hello, World!\n",
		},
		{
			name: "entry message has trailing newline",
			entry: Entry{
				Msg:  "Hello, World!\n",
				Time: time.Date(2017, time.February, 17, 1, 2, 3, 456789000, time.UTC),
				File: "foo/bar.go",
				Line: 278,
			},
			want: "Hello, World!\n",
		},
		{
			name: "short file",
			flag: ShowDate | ShowTime | Microseconds | ShortFile,
			entry: Entry{
				Msg:  "Hello, World!\n",
				Time: time.Date(2017, time.February, 17, 1, 2, 3, 456789000, time.UTC),
				File: "foo/bar.go",
				Line: 278,
			},
			want: "2017/02/17 01:02:03.456789 bar.go:278: Hello, World!\n",
		},
		{
			name: "time flags",
			flag: ShowDate | ShowTime | Microseconds,
			entry: Entry{
				Msg:   "Hello, World!\n",
				Level: Warn,
				Time:  time.Date(2017, time.February, 17, 1, 2, 3, 456789000, time.UTC),
				File:  "foo/bar.go",
				Line:  278,
			},
			want: "2017/02/17 01:02:03.456789 Hello, World!\n",
		},
		{
			name: "all flags",
			flag: ShowDate | ShowTime | Microseconds | ShowFile | ShowLevel,
			entry: Entry{
				Msg:   "Hello, World!\n",
				Level: Warn,
				Time:  time.Date(2017, time.February, 17, 1, 2, 3, 456789000, time.UTC),
				File:  "foo/bar.go",
				Line:  278,
			},
			want: "2017/02/17 01:02:03.456789 WARN foo/bar.go:278: Hello, World!\n",
		},
		{
			name: "all flags debug",
			flag: ShowDate | ShowTime | Microseconds | ShowFile | ShowLevel,
			entry: Entry{
				Msg:   "Hello, World!\n",
				Level: Debug,
				Time:  time.Date(2017, time.February, 17, 1, 2, 3, 456789000, time.UTC),
				File:  "foo/bar.go",
				Line:  278,
			},
			want: "2017/02/17 01:02:03.456789 DEBUG foo/bar.go:278: Hello, World!\n",
		},
	}
	ctx := context.Background()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := New(buf, test.prefix, test.flag, nil)
			logger.Log(ctx, test.entry)
			if s := buf.String(); s != test.want {
				t.Errorf("log output = %q; want %q", s, test.want)
			}
		})
	}
}

func TestFlagString(t *testing.T) {
	tests := []struct {
		f    Flags
		want []string
	}{
		{
			f:    0,
			want: []string{"0"},
		},
		{
			f:    ShowDate,
			want: []string{"ShowDate"},
		},
		{
			f:    ShowDate | ShowTime,
			want: []string{"ShowDate", "ShowTime"},
		},
		{
			f:    ShowDate | ShowTime | Microseconds | ShowFile | ShortFile | UTC | ShowLevel,
			want: []string{"ShowDate", "ShowTime", "Microseconds", "ShowFile", "ShortFile", "UTC", "ShowLevel"},
		},
		{
			f:    ShowDate | 1<<31,
			want: []string{"ShowDate", "2147483648"},
		},
	}
	stringLess := func(s1, s2 string) bool {
		return s1 < s2
	}
	for _, test := range tests {
		got := test.f.String()
		if !cmp.Equal(strings.Split(got, "|"), test.want, cmpopts.SortSlices(stringLess)) {
			t.Errorf("Flags(%#x).String() = %q; want %q", uint(test.f), got, strings.Join(test.want, "|"))
		}
	}
}

func BenchmarkWriter(b *testing.B) {
	ctx := context.Background()
	buf := new(bytes.Buffer)
	logger := New(buf, "", StdFlags, nil)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		logger.Log(ctx, Entry{
			Msg:  "Hello, World!",
			Time: time.Now(),
		})
	}
}
