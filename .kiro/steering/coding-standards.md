---
inclusion: auto
---

# Coding Standards

## Go Conventions

- Package name: `tcp`
- Module path: `github.com/pulsyflux/tcp`
- Use `sync.RWMutex` for concurrent map access, `atomic` for counters and flags
- Use `sync.Mutex` for simple exclusive access (e.g. `wrappedPools` global)
- Error handling: return errors, don't panic
- Use `uuid.UUID` from `github.com/google/uuid` for all identifiers
- Channel buffers: capacity 10 for demux recv channels
- Use `context.Context` for goroutine lifecycle management
- Goroutine cleanup: context cancellation
- Test files: `*_test.go` in same package, benchmark files: `*_bench_test.go`
- No debug logging in library code (no `fmt.Print*` or `log.*`)

## Message Framing Protocol

All messages use this wire format:
```
[ID Length: 1 byte][ID: N bytes (UUID string)][Total Length: 4 bytes big-endian][Chunk Length: 4 bytes big-endian][Data: M bytes]
```
- Chunking at 64KB boundaries for large messages
- Each chunk gets the full header (idLen + id + totalLen + chunkLen)
- Demuxer reassembles chunks before routing to logical connection

## Thread Safety Patterns

- One `Send` and one `Receive` per connection at a time (atomic CAS guards via `int32`)
- Single demux reader goroutine per physical connection — never read from socket directly
- `connectionPool` (global pool): `sync.RWMutex` (read-heavy, write-rare)
- `demuxer.routes`: `sync.RWMutex` for concurrent route access
- `wrappedPools` (global): `sync.Mutex` for exclusive access
- `Connection.mu`: `sync.RWMutex` for state and timestamp access

## Testing Patterns

- Go: use `testing.T` and `testing.B` directly
- Go benchmarks: use `b.N` loop, report `-benchmem`
- Always use `127.0.0.1:0` for server address in tests (OS-assigned port)
- Server side in tests: use goroutines with `WrapConnection` for echo servers
