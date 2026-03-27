# http-from-tcp

A small Go project for learning how HTTP works on top of raw network primitives.

This repository explores the building blocks behind HTTP by working directly with TCP sockets, reading streaming input manually, and parsing the HTTP request line without relying on higher-level web frameworks.

## What this project does

- listens for incoming TCP connections
- reads incoming data as a stream
- splits incoming bytes into newline-delimited messages
- parses and validates an HTTP/1.1 request line
- includes tests for valid and invalid request parsing cases

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
