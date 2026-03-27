# http-from-tcp

A small Go project for learning how HTTP works on top of raw network primitives.

This project explores how HTTP is built by working directly with TCP sockets, handling raw byte streams, and parsing HTTP requests manually without relying on high-level frameworks.

---

## What this project does

- Listens for incoming TCP connections
- Reads incoming data as a stream
- Splits messages by newline boundaries
- Parses and validates HTTP/1.1 request lines
- Includes tests for valid and invalid request formats

---

## Project structure

```text
.
├── cmd/
│   ├── tcplistener/
│   │   └── main.go
│   └── udpsender/
│       └── main.go
├── internal/
│   └── request/
│       ├── request.go
│       └── request_test.go
├── main.go
├── messages.txt
├── go.mod
└── go.sum
```

---

## Why this project exists

Most developers interact with HTTP through frameworks or standard libraries. This project takes a lower-level approach to help you understand:

- how data is transmitted over TCP
- how streaming input is processed
- how HTTP request lines are structured
- how protocol validation works internally

This serves as a foundation for building an HTTP server from scratch.

---

## Components

### TCP Listener
Accepts connections and reads incoming data line-by-line to simulate how raw HTTP requests arrive over a socket.

### UDP Sender
A simple CLI tool to send messages over UDP for testing and experimentation.

### HTTP Request Parser
Parses and validates HTTP/1.1 request lines, including method, path, and version.

---

## Example request

```http
GET / HTTP/1.1
Host: localhost:42069
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

### Run UDP sender

```bash
go run ./cmd/udpsender
```

---

## Learning focus

- Go networking fundamentals
- TCP socket handling
- Stream processing
- HTTP request parsing
- Protocol validation

---

## Limitations

- Not a full HTTP server yet
- No header or body parsing beyond request line
- No routing or response handling

---

## Future improvements

- Header parsing
- Request body support
- HTTP response generation
- Routing system
- Full HTTP server implementation

---

## Tech stack

- Go
- TCP / UDP
- HTTP/1.1
