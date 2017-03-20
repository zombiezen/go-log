// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func TestWriter(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		flag   int
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
			flag: Ldate | Ltime | Lmicroseconds | Lshortfile,
			entry: Entry{
				Msg:  "Hello, World!\n",
				Time: time.Date(2017, time.February, 17, 1, 2, 3, 456789000, time.UTC),
				File: "foo/bar.go",
				Line: 278,
			},
			want: "2017/02/17 01:02:03.456789 bar.go:278: Hello, World!\n",
		},
		{
			name: "all flags",
			flag: Ldate | Ltime | Lmicroseconds | Llongfile,
			entry: Entry{
				Msg:  "Hello, World!\n",
				Time: time.Date(2017, time.February, 17, 1, 2, 3, 456789000, time.UTC),
				File: "foo/bar.go",
				Line: 278,
			},
			want: "2017/02/17 01:02:03.456789 foo/bar.go:278: Hello, World!\n",
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

func BenchmarkWriter(b *testing.B) {
	ctx := context.Background()
	buf := new(bytes.Buffer)
	logger := New(buf, "", LstdFlags, nil)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		logger.Log(ctx, Entry{
			Msg:  "Hello, World!",
			Time: time.Now(),
		})
	}
}
