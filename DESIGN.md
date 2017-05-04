# Design Rationale

This is an adaptation of the original public proposal/design document to explain
the rationale of the package's API.

## Goal

As a Go application author, I want to be able to control the diagnostic logging
backend a library uses for the benefit of avoiding missing log messages and
keeping a consistent log entry structure.

As a Go library author, I want to be able to log diagnostic information for the
benefit of making problems in my library easier to diagnose.

## Asynchrony

An implementation of `Logger` guarantees that if log entries are delivered, then
they are delivered in order of the calls given.  However, it is important that
implementations of `Logger` should not be required to guarantee delivery.  For
application performance, sending a log message generally should not block
request handling.  For example, logs in a cloud environment may be sent over the
network.  Blocking the request until the log is delivered would be unacceptable
in this case.  For this reason, the Log method does not return an error and the
`Context`'s deadline and cancelation should be ignored.

A common objection is that returning an error allows the application to specify
whether delivery is important, and not returning an error renders the
application unable to know if log delivery failed.  However, it's important to
bear in mind that the error here is being returned to the library, not the
application.  The application author (not the library author) is the one who
knows whether or not a log needs to be written successfully.  By including an
error return in the interface, the log backend is hoisting responsibility to
library authors to check every log error, or else the application will not know
whether logs were sent successfully or not.  A library can't choose to ignore
the error (like when `fmt.Fprintf(os.Stderr, ...)` is used), without being
rendered useless to applications that need to know when delivery fails.  Thus,
checking the error at library call sites would be unreliable at best, and is
best handled by the application out-of-band — the application may have an
alternate means of indicating failure of log delivery like metric counters.

It's worth noting that the parameters to `Log` are either safe to use from
multiple goroutines or are passed by value, so there are no [aliasing][]
concerns.

[aliasing]: https://en.wikipedia.org/wiki/Aliasing_(computing)

## Context

The application may need to pass request or other trace information down to its
`Logger` via the `Context` values.  For example, an application's logging
backend of choice could require knowledge of the request ID,
[like on Google App Engine][App Engine logs] (although the need is not specific
to Google: recording trace/request ID has come up as a common requirement for
multiple production Go users).  This is [a classic example][Correctly Use Context]
of a `Context` value: a value that is not part of the library's control
flow and needs to be plumbed through code that the application does not own.
`Loggers` might be interested in the values, but should ignore the deadline and
cancelation (as per the previous section).

It could be argued that passing in `Context` is not needed, as one could create
a per-request `Logger` that has the [curried][] values.  Indeed, this is what
[go-kit's `log.Context`][go-kit log.Context] does.  However, if you imagine a
library providing a long-lived resource type that can be shared by multiple
goroutines, then this would force the application to pass in a per-request
`Logger` on each operation versus creating a `Logger` at the resource's
creation, then using the values from `Context` to fill in the request-specific
fields at `Log` time.

As an example, assume that the application wants to include the current
request/trace ID in the log.  Passing the context as part of the `Log` call
allows you to write this:

```go
type myHandler struct {
  log log.Logger
}

func (h myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  log.Infof(r.Context, h.log, "serving!")
}
```

instead of:

```go
type myHandler struct {
  requestLogger func(context.Context) log.Logger
}

func (h myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  logger := h.requestLogger(r.Context)
  log.Infof(r.Context, logger, "serving!")
}
```

As always, if a library is unable to propagate a context (API compatibility
concerns, etc.), then it can pass in [context.TODO][].

[App Engine logs]: https://cloud.google.com/appengine/docs/go/logs/#writing_application_logs
[Correctly Use Context]: https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39#.yc7x7ge4d
[curried]: https://en.wikipedia.org/wiki/Partial_application
[go-kit log.Context]: https://godoc.org/github.com/go-kit/kit/log#Context
[context.TODO]: https://godoc.org/context#TODO

## Time

A common pattern with Go interfaces is to implement a method by wrapping another
instance of the same interface (the [decorator pattern][]).  This can be very
useful for a `Logger` (teeing, filtering, etc.).  By passing an explicit
timestamp, it is easier to compose `Logger`s and ensure they all have the same
timestamp for the same entry.  Another scenario would be if someone wants to
write a `Logger` that sends an entry to multiple underlying `Logger`s, the
timestamps would mismatch if each underlying `Logger` used `time.Now`.  To
permit these common use cases, time should be passed in as a parameter.

## Levels

(Major thanks to Sarah Adams for designing this part.)

Dave Cheney explores levels in his article [Let's talk about logging][].  This
introduces most of the rationale for the two levels chosen.  Any levels beyond
verbosity for to use encroach on being the concern of how to format the logs,
which is what the application ought to enforce instead.  For example, a network
diagnostic tool application might want to print dropped packet at level warning
to the underlying logs, whereas a long-running server application might want to
send this as level info in the underlying logs.  As such, you won't see warning
or error.

[decorator pattern]: https://en.wikipedia.org/wiki/Decorator_pattern
[Let's talk about logging]: https://dave.cheney.net/2015/11/05/lets-talk-about-logging

## Why not structured?

Structured logging has a number of well-recognized benefits, so it may be
surprising to see that the payload is a simple string instead of a
`map[string]interface{}` or perhaps even `interface{}`.  After all, structured
logging is easier to parse than text and easier to analyze.  This is true, but
this statement should not be construed to say that applications should not use
structured logging — quite the opposite.  An application could provide a library
with an implementation of `Logger` that attaches the unstructured message to a
structured logging entry, perhaps with some extra data about which package the
message came from.  This package's goal is to provide a common interface for
libraries, not applications.

For the purposes of the rest of this discussion, let's contrast this package's
design with this one:

```go
// Entry is a single log record.
type Entry struct {
	Payload map[string]interface{}

	// The rest of the fields are the same.
}
```

There are concerns to the above API that would need to be solved, and it's not
clear there are comfortable solutions.  Not all backends would be able to
support all types passed into as an `interface{}`: channels would be an easy one
to nix, but what about cyclic data structures?  Some may support integer keys in
sub-maps, but others (JSON) require string keys.  One could document that for
unsupported types, `fmt.Sprint` would be used to convert to string.  This
negates many benefits of structured logging, but does reveal an interesting
observation: any logging backend knows how to deal with unstructured text.

There's also some smaller design concerns that would have to be addressed before
an interface such as the one above would need to overcome.  There's lifetime and
mutability concerns: if the call is asynchronous (again, required for
application performance), then either the `Logger` must be able to hold onto the
map past the return of `Log` (which could introduce variable aliasing issues) or
create a copy.   There are performance concerns: creating a map forces an
allocation on the caller (the library).  What if the log backend writes a binary
format like protobufs?  The fields would have to be mapped to other fields in
the binary structure, which may introduce overhead.  And finally, there are
application formatting concerns: what are the acceptable keys that a library can
use in the map?  Can they be any UTF-8 clean string?  C identifiers?  How do I
avoid collisions between keys of the library and the application formatting?  If
I want to remap the keys that a library emits, the library now need to document
the keys that it uses.

As soon as we hit needing to document emitted keys and their types, we have
reached the point of The Codeless Code case, [Back to Basics][].  Go has a type
system that is usable for documenting and, more importantly, enforcing these
sorts of contracts of data between components.  If a library needs to emit this
sort of structured diagnostic information, then it ought to define its own API
for exposing this information to the application, and then the application can
make the choice of how to format this information (e.g. structured logs,
unstructured logs, metrics, ignoring it, etc.).
[`net/http/httptrace.ClientTrace`][ClientTrace] is a wonderful example of a
library exposing a diagnostic event API to an application.  The data structures
and events are well-defined.  These APIs are necessarily going to be
library-specific, and thus defining one generic, unified structured logging
interface would provide negative utility.  Such an interface would encourage
avoiding the type system and having underspecified library APIs.

But back to the interface proposed at the top of the document.  By this same
argument, wouldn't defining an interface to log unstructured events be just as
harmful?  I don't think so; I think that having an unstructured logger provides
net positive utility.  Sometimes a library wants to emit internal _debugging_
information: the sort of thing where you want to dump a value to your log in
production to see what's going wrong.  This is ephemeral information and should
be considered an implementation detail of the library, not part of the library's
API itself.

[Back to Basics]: http://thecodelesscode.com/case/167
[ClientTrace]: https://godoc.org/net/http/httptrace#ClientTrace

## Why not `io.Writer`?

Assuming you are convinced by the need to write unstructured logs, you may point
out that there is an existing interface that writes unstructured data:
`io.Writer`.  This might be good enough, and is the assumption that the standard
log library makes.  However, I would argue that `io.Writer` is not quite the right
abstraction because it does not capture the boundaries of a log entry (which is
what some logging backends require).  In this case, I'm defining a log entry as
a single library call to log some message.

There are three general approaches you could take to define the log entry boundary:

1.  1 `Write` call = 1 log entry.  While this sounds fine and works well with
    the `fmt.Fprint` cases, this is because these functions in their current
    implementation only make one call to `Write` per call.  This is not a
    guarantee, and there is discussion of changing this behavior.  You can
    imagine a class of bugs where someone just passes the log `io.Writer` to a
    function that doesn't guarantee this (or `fmt.Fprint` changes in a future Go
    release) and the backend gets spews of partial log entries.  This class of
    bugs is subtle and can occur from refactoring seemingly unrelated code.
    (There's also a practical concern that a real file like `os.Stderr` would
    not work well for log, because devs often don't place newlines in log
    statements, assuming that the underlying log library will add it for them.)

1.  1 line = 1 log entry.  This has the obvious problem of not supporting
    multi-line log entries — like stack traces or program outputs.

1.  Some agreed-upon marker is the log entry boundary.  Maybe this is the log
    timestamp.  Now all logging backends that want to differentiate log entries
    need to parse the underlying output of the `io.Writer` and buffer it, and
    that's not particularly performant.

When viewing these limitations, introducing an explicit interface that knows how
to deal with a log boundary provides utility over `io.Writer`.

## The `LogEnabled` method

Formatting a log message is often an expensive process: many small temporary
objects may be created and is generally unbounded work — `String()` methods can
do anything.  A major performance gain is not creating log entries when they are
not needed when logging is turned off in particular places.  `Logger.LogEnabled`
gives a signal as to whether it is even worth filling out the message based on
the `Entry`'s metadata.  Re-using `Entry` allows more filterable fields to be
added in the future without changing the `Logger` interface signature.  For
purity, you could define an `EntryMetadata` struct and embed it in `Entry`, but
in practice, this seemed clumsier than just filling out `Entry.Msg` after
calling `LogEnabled`.

This could be made into an optional interface, but since it is purely
an optimization and the simplest correct implementation is always `return true`,
keeping both methods on the interface is the simplest.

## Global logger motivation

For diagnostic logs (as opposed to audit logs), you must eventually assume that any function anywhere can log.  The evidence of this is the proliferation of package-level variables for logging or using `log.Printf` from the standard library.  The reason this is largely maligned is that the backend cannot be swapped out for another, as the creation of the `Logger` is tied to its package.  

If any function may want to log, then it quickly becomes burdensome to have to plumb through a `Logger` to basically everything.  What if a function does not have anything it needs to debug now, but might in the future?  Library authors would need to define interfaces and functions to explicitly take in a Logger, even if they don't need it.  Especially when you consider diagnostic logging to be an implementation detail, adding Logger to the API signatures leaks the abstraction.

However, if we allow a customizable global `Logger`, then we eliminate the compile-time dependency on a particular implementation and don't pollute every function's signature with `Logger`.

## Default to stderr

There are four ways to handle `Default().Log(...)` before `SetDefault` is called:

1.  Drop the entry
2.  Buffer the entries until `SetDefault` is called
3.  Panic
4.  Log to a fallback source

The first two options can be rejected almost immediately.  Dropping the entry
defeats the goal of not missing logs.  Buffering is tempting at face-value, but
consider the case where `SetDefault` is never called.  The buffer is unbounded
and there's no obvious indication of the leak.  Especially when this package is
new, this case will occur frequently, so this isn't a real option.

Panicing was seriously considered.  It makes it obvious that logging isn't
configured and avoids either of the first two options' downsides.  However, if a
library only logs in exceptional circumstances (quite common), then there's a
very probable chance that the user won't notice the misconfiguration until the
rare library condition occurs, but then the log is lost.  It also would
discourage most early adopters of such a package, since application authors
would be very likely to become aware of this choice through the panic.

So that leaves us with option 4.  Logging to stderr is the most common behavior,
and matches what the standard log library does right now, so it is no worse than
the current state of the world and would allow library authors to migrate usage
of standard log library to this package without API disruption, while giving
the application control over log output.
