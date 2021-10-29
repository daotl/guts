# GUTS - Go UTilitieS

[![Go Reference](https://pkg.go.dev/badge/github.com/daotl/guts.svg)](https://pkg.go.dev/github.com/daotl/guts)

## [net](./net/net.go)

- Connect
- ProtocolAndAddress
- GetFreePort

## [sync.Mutex and sync.RWMutex](./sync/mutex.go)

Use `sync.Mutex` and `sync.RWMutex` from this package instead of `sync` so you can build your 
program with `deadlock` flag to detect deadlocks using [go-deadklock](https://github.com/sasha-s/go-deadlock).
