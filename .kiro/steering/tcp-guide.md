---
inclusion: auto
---

# TCP Package Guide

## Key Files

- `logical.go` — `Connection` struct, `NewConnection`, `WrapConnection`, `Send`, `Receive`, idle monitor, `wrappedPool`
- `pool.go` — `connectionPool` (global), `physicalPool` (client-side pooling with ref counting)
- `demux.go` — `demuxer` struct, single-reader goroutine, UUID-based message routing, chunk reassembly
- `errors.go` — sentinel errors: `errConnectionClosed`, `errConnectionDead`, `errConnectionInUse`
- `connection_test.go` — unit tests for send/receive, large messages, multiplexing, pool ref counting, concurrent use
- `connection_bench_test.go` — benchmarks for small/large messages, chunking, pooling, send-only, receive-only
- `logical_connections_test.go` — integration test for server+client multiplexed logical connections

## Connection Lifecycle

1. Client: `NewConnection(address, id)` → pool `getOrCreate` → register route in demuxer → start idle monitor
2. Server: `WrapConnection(conn, id)` → `getOrCreateWrappedPool` → register route → start idle monitor
3. Idle timeout (30s) → `close()` → unregister from demuxer → release pool ref → physical close at refCount 0

## Send Frame Format

```
[idLen:1][id:N][totalLen:4][chunkLen:4][data:M]
```
- Chunks at 64KB, each chunk gets full header
- Big-endian uint32 for lengths

## Demuxer

- One goroutine per physical connection reads frames continuously
- Routes complete messages to registered UUID channels
- `messageAssembly` tracks partial messages across chunks
- Buffered channels (cap 10) prevent demux blocking

## Important Constraints

- Only one concurrent `Send` and one concurrent `Receive` per `Connection` (atomic CAS guards)
- `Receive` reads from channel, not socket — always goes through demuxer
- Client connections auto-reconnect; server (wrapped) connections do not
- Pool uses reference counting — physical connection closes when last logical connection releases
- `close()` is unexported — connections are closed via idle timeout or context cancellation
- Maximum message size: 16MB (`maxMessageSize`) — enforced on both send and receive
- Maximum pending incomplete message streams per connection: 64 (`maxPendingStreams`)
- No TLS — plaintext TCP only. Use a TLS proxy or tunnel for encryption
