# TCP

A minimal, high-performance TCP connection abstraction with connection pooling, multiplexing, and auto-reconnect.

```
go get github.com/pulsyflux/tcp
```

## Features

- Global connection pool with automatic sharing of physical connections
- UUID-based multiplexing — multiple logical connections over one socket
- Auto-reconnect on client side
- Idle timeout with automatic cleanup
- Thread-safe, zero configuration

## API

```go
// Client-side: creates connection with auto-reconnect
func NewConnection(address string, id uuid.UUID) *Connection

// Server-side: wraps accepted connection (no reconnect)
func WrapConnection(conn net.Conn, id uuid.UUID) *Connection

func (c *Connection) Send(data []byte) error
func (c *Connection) Receive() ([]byte, error)
```

## Quick Start

### Client

```go
import "github.com/pulsyflux/tcp"

c := tcp.NewConnection("localhost:8080", uuid.New())
c.Send([]byte("hello"))
data, _ := c.Receive()
```

### Server

```go
listener, _ := net.Listen("tcp", ":8080")
conn, _ := listener.Accept()
wrapped := tcp.WrapConnection(conn, uuid.New())

data, _ := wrapped.Receive()
wrapped.Send(data)
```

### Multiplexing

```go
// Both share the same physical TCP connection
conn1 := tcp.NewConnection("localhost:8080", uuid.MustParse("00000000-0000-0000-0000-000000000001"))
conn2 := tcp.NewConnection("localhost:8080", uuid.MustParse("00000000-0000-0000-0000-000000000002"))

conn1.Send([]byte("user data"))
conn2.Send([]byte("admin data"))

// Each receives only its own messages
userData, _ := conn1.Receive()
adminData, _ := conn2.Receive()
```

## Documentation

Detailed docs live in [`.kiro/steering/`](.kiro/steering/):

- [Project Overview](.kiro/steering/project-overview.md) — architecture, module info, performance summary
- [TCP Package Guide](.kiro/steering/tcp-guide.md) — internals, connection lifecycle, frame format, demuxer, constraints
- [Coding Standards](.kiro/steering/coding-standards.md) — conventions, wire protocol, thread safety, testing patterns
- [CI Workflow Guide](.kiro/steering/ci-workflow-guide.md) — GitHub Actions and local CI with Act

## License

See LICENSE file for details.
