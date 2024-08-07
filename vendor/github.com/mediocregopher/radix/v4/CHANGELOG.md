Changelog from v4.0.0 and up. v3 changelog can be found in its branch.

# v4.0.0

Below is documented all breaking changes between v3 and v4. There are further
enhancements which don't qualify as breaking changes which may not be documented
here.

**Major Changes**

* Stop using `...opts` pattern for optional parameters across all types, and
  switch instead to a `(Config{}).New` kind of pattern.

* Add `MultiClient` interface which is implemented by `Sentinel` and `Cluster`,
  `Client` has been modified to be implemented only by clients which point at a
  single redis instance (`Conn` and `Pool`). Methods on all affected
  client types have been modified to fit these new interfaces.

  * `Cluster.NewScanner` has been replaced by `ScannerConfig.NewMulti`.

* `Conn` has been completely re-designed. It is now always thread-safe. When
  multiple `Action`s are performed against a single `Conn` concurrently the
  `Conn` will automatically pipeline the `Action`'s read/writes, as appropriate.

  * `Pipeline` has been re-designed as a result as well.

  * `CmdAction` has been removed.

* `Pool` has been completely rewritten to better take advantage of connection
  sharing (previously called "implicit pipelining" in v3) and the new `Conn`
  design.

  * `EvalScript` and `Pipeline` now support connection sharing.

  * Since most `Action`s can be shared on the same `Conn` the `Pool` no longer
    runs the risk of being depleted during too many concurrent `Action`s, and so
    no longer needs to dynamically create/destroy `Conn`s.

  * A Pool size of 0 is no longer supported.

* Brand new `resp/resp3` package which implements the [RESP3][resp3] protocol.
  The new package features more consistent type mappings between go and redis
  and support for streaming types.

* Usage of `context.Context` in many places.

  * Add `context.Context` parameter to `Client.Do`, `PubSub` methods,
    `Scanner.Next`, and `WithConn`.

  * Add `context.Context` parameter to all `Client` and `Conn` creation functions.

  * Add `context.Context` parameter to `Action.Perform` (previously called
    `Action.Run`).


**Minor Changes**

* Remove usage of `xerrors` package.

* Rename `resp.ErrDiscarded` to `resp.ErrConnUsable`, and change some of the
  semantics around using the error. A `resp.ErrConnUnusable` convenience
  function has been added as well.

* `resp.LenReader` now uses `int` instead of `int64` to signify length.

* `resp.Marshaler` and `resp.Unmarshaler` now take an `Opts` argument, to give
  the caller more control over things like byte pools and potentially other
  functionality in the future.

* `resp.Unmarshaler` now takes a `resp.BufferedReader`, rather than
  `*bufio.Reader`. Generally `resp.BufferedReader` will be implemented by a
  `*bufio.Reader`, but this gives move flexibility.

* `Stub` and `PubSubStub` have been renamed to `NewStubConn` and
  `NewPubSubStubConn`, respectively.

* Rename `MaybeNil` to just `Maybe`, and change its semantics a bit.

* The `trace` package has been significantly updated to reflect changes to
  `Pool` and other `Client`s.

* Refactor the `StreamReader` interface to be simpler to use.

[resp3]: https://github.com/antirez/RESP3
