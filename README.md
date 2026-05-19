# http-from-tcp

A small Go project for learning how HTTP works on top of raw network primitives.

This project explores how HTTP is built by working directly with TCP sockets, handling raw byte streams, and parsing HTTP requests manually without relying on high-level frameworks.

---

## What this project does

- Listens for incoming TCP connections
- Reads incoming data as a stream (incremental reads, not one-shot buffers)
- Parses and validates HTTP/1.1 request lines (method, target, version)
- Parses request headers (case-insensitive keys, duplicate merging)
- Parses request bodies when `Content-Length` is present
- Includes tests for valid and invalid request formats (including split/chunked reads)

---

## Project structure

```text
.
├── cmd/
│   ├── tcplistener/
│   │   └── main.go          # TCP demo: parse one request per connection, print fields
│   └── udpsender/
│       └── main.go          # UDP stdin client (separate from HTTP path)
├── internal/
│   ├── headers/
│   │   ├── headers.go
│   │   └── headers_test.go
│   └── request/
│       ├── request.go       # Request parser + state machine
│       └── request_test.go
├── messages.txt
├── go.mod
└── go.sum
```

For a deeper architecture and gap analysis, see [REPO_AUDIT.md](./REPO_AUDIT.md).

---

## Why this project exists

Most developers interact with HTTP through frameworks or standard libraries. This project takes a lower-level approach to help you understand:

- how data is transmitted over TCP
- how streaming input is processed incrementally
- how HTTP request lines, headers, and bodies are structured
- how protocol validation works internally

This serves as a foundation for building an HTTP server from scratch.

---

## Components

### TCP Listener (`cmd/tcplistener`)

Accepts connections on `:42069`, feeds the raw `net.Conn` into the request parser, and prints the parsed request line, headers, and body. Useful for observing how a full HTTP/1.1 request arrives over a socket.

### UDP Sender (`cmd/udpsender`)

A simple CLI tool to send lines over UDP for experimentation. It shares the same port number as the TCP listener but is not part of the HTTP parsing path.

### HTTP Request Parser (`internal/request` + `internal/headers`)

Incremental HTTP/1.1 request parser driven by a small state machine:

1. Request line — method, request-target, `HTTP/1.1` only  
2. Headers — `\r\n`-delimited lines until blank line  
3. Body — `Content-Length` when present; otherwise empty body  

`RequestFromReader(io.Reader)` handles buffering and partial reads so parsing works when data arrives in arbitrary chunk sizes.

---

## Example request

```http
GET / HTTP/1.1
Host: localhost:42069
User-Agent: curl/7.81.0

```

With a body (POST):

```http
POST /submit HTTP/1.1
Host: localhost:42069
Content-Length: 13

hello world!
```

---

## Getting started

### Prerequisites

- Go 1.25+

### Clone the repository

```bash
git clone https://github.com/muafa7/http-from-tcp.git
cd http-from-tcp
```

### Run tests

```bash
go test ./...
```

### Run TCP listener

```bash
go run ./cmd/tcplistener
```

In another terminal, send a request (example with curl):

```bash
curl -v http://localhost:42069/
```

### Run UDP sender

```bash
go run ./cmd/udpsender
```

---

## Learning focus

- Go networking fundamentals
- TCP socket handling
- Stream processing and incremental parsing
- HTTP/1.1 request structure (line, headers, body)
- Protocol validation and test-driven edge cases

---

## Limitations

This is **not** a production HTTP server yet.

- No HTTP responses (status line, headers, or body back to the client)
- No routing or handlers
- Body parsing via `Content-Length` only (no chunked transfer encoding)
- No read/write timeouts, connection limits, or graceful shutdown
- TCP listener does not close connections after handling (demo only)
- Requires `\r\n` line endings (not bare `\n`)

---

## Roadmap

- [ ] HTTP response writer and minimal handlers (e.g. `GET /health`)
- [ ] Method + path router
- [ ] Per-connection timeouts and request size limits
- [ ] Graceful shutdown and proper connection lifecycle
- [ ] Integration tests against a live TCP listener
- [ ] Benchmarks and CI (`-race`)

See [REPO_AUDIT.md](./REPO_AUDIT.md) for milestone detail.

---

## Tech stack

- Go
- TCP / UDP
- HTTP/1.1
