---
inclusion: auto
---

# Coding Standards

## Go Conventions

- Package name: `tcp`
- Module path: `github.com/pulsyflux/tcp`
- Use `sync.RWMutex` for concurrent map access, `atomic` for counters and flags
- Error handling: return errors, don't panic. Log and continue for non-fatal connection errors.
- Use `uuid.UUID` from `github.com/google/uuid` for all identifiers
- Channel buffers: use appropriate capacity (10 for demux recv channels)
- Use `context.Context` for goroutine lifecycle management
- Goroutine cleanup: use `done` channels or context cancellation
- Test files: `*_test.go` in same package, benchmark files: `*_bench_test.go`

## Message Framing Protocol

All messages use this wire format:
```
[ID Length: 1 byte][ID: N bytes (UUID string)][Total Length: 4 bytes big-endian][Chunk Length: 4 bytes big-endian][Data: M bytes]
```
- Chunking at 64KB boundaries for large messages
- Demuxer reassembles chunks before routing to logical connection

## Thread Safety Patterns

- One `Send` and one `Receive` per connection at a time (atomic CAS guards)
- Single demux reader goroutine per physical connection — never read from socket directly
- Pool access: RWMutex (read-heavy, write-rare)

## Testing Patterns

- Go: table-driven tests, use `time.After` for timeout assertions
- Go benchmarks: use `b.N` loop, report `-benchmem`
- Always use `:0` for server address in tests (OS-assigned port)
