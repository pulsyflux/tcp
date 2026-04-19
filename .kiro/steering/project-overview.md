---
inclusion: auto
---

# TCP Package Overview

A high-performance TCP connection abstraction library written in Go.

## Module Info

- Go module: `github.com/pulsyflux/tcp` (go 1.25)
- Go package: `tcp`
- Dependency: `github.com/google/uuid v1.6.0`
- Published: `https://pkg.go.dev/github.com/pulsyflux/tcp`
- Source: `https://github.com/pulsyflux/tcp`

## Architecture

Provides multiplexed TCP connections with:
- Global connection pool — multiple logical connections share one physical TCP socket per address
- UUID-based message routing via a shared demuxer goroutine per physical connection
- Message framing: `[ID Length (1B)][ID (N B)][Total Length (4B)][Chunk Length (4B)][Data (M B)]`
- Chunked I/O (64KB chunks) for large messages
- Auto-reconnect on client side, no reconnect on server side (wrapped connections)
- Idle timeout (30s default) with automatic cleanup
- Thread-safe: atomic read/write guards, RWMutex on pool and routes

Key types:
- `Connection` — logical connection (client via `NewConnection`, server via `WrapConnection`)
- `demuxer` — single-reader goroutine per physical connection, routes by UUID
- `physicalPool` — client-side connection pooling with reference counting
- `wrappedPool` — server-side pooling for accepted connections

## Public API

- `NewConnection(address string, id uuid.UUID) *Connection` — client-side, auto-reconnect
- `WrapConnection(conn net.Conn, id uuid.UUID) *Connection` — server-side, no reconnect
- `(*Connection).Send(data []byte) error`
- `(*Connection).Receive() ([]byte, error)`

## Testing

- `go test -v` — run all tests
- `go test -bench=. -benchmem` — run benchmarks
- Tests use `127.0.0.1:0` for OS-assigned ports

## Performance

- ~39µs small message latency, ~8µs send/receive, supports 1MB+ messages
