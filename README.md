# zombiezen.com/go/log

[![GoDoc](https://godoc.org/zombiezen.com/go/log?status.svg)][godoc.org]
[![Build Status](https://travis-ci.org/zombiezen/go-log.svg?branch=master)][travis]

This package is the concrete implementation of the design described in
[DESIGN.md][], formerly a proposal for the Go project.  While it was not
suitable for being an "official" package, I still wanted to use it for my own
projects.  If others find it helpful, then I'm happy.

I don't consider the API stable at the moment; I still want some time for this
API to bake.  I will be following [SemVer][] tags, but in v0: pin to specific
revisions for the time being.

[godoc.org]: https://godoc.org/zombiezen.com/go/log
[travis]: https://travis-ci.org/zombiezen/go-log
[DESIGN.md]: https://github.com/zombiezen/go-log/blob/master/DESIGN.md
[SemVer]: http://semver.org/

## Install

```bash
go get -u zombiezen.com/go/log
```

## Documentation

See the docs on [godoc.org][].

## License

BSD - see [LICENSE](https://github.com/zombiezen/go-log/blob/master/LICENSE).
