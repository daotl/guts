# GUTS - Go UTilitieS

## [sync.Mutex and sync.RWMutex](./sync/mutex.go)

Use `sync.Mutex` and `sync.RWMutex` from this package instead of `sync` so you can build your 
program with `deadlock` flag to detect deadlocks using [go-deadklock](https://github.com/sasha-s/go-deadlock).
