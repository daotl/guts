# GUTS - Go UTilitieS

[![Go Reference](https://pkg.go.dev/badge/github.com/daotl/guts.svg)](https://pkg.go.dev/github.com/daotl/guts)

## [net](./net/net.go)

- Connect
- ProtocolAndAddress
- GetFreePort

## [rand](./rand)

Random number, string, bytes generation. 

## [service](./service)

Package [service/goprocesssrv](./service/goprocess/service.go) is a service implementation based on
[goprocess](https://github.com/jbenet/goprocess) that supports declaring and ensuring runtime
dependencies between services.

Package [service/suturesrv](./service/suture/service.go) is a service implementation based on
[Suture](https://github.com/thejerf/suture).

## [sync/mutex](./sync/mutex.go)

Use `sync.Mutex` and `sync.RWMutex` from this package instead of `sync` so you can build your
program with `deadlock` flag to detect deadlocks
using [go-deadklock](https://github.com/sasha-s/go-deadlock).

## [sync.ResultNotifier](./sync/result_notifier.go)

ResultNotifier can be used for a goroutine to notify others that a job is done
and pass the error occurred or nil if no error.
It's based on `sync.Cond`.
