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

## Why this project exists

Most applications interact with HTTP through frameworks and standard library helpers. This project takes a lower-level approach so you can better understand:

- how bytes arrive over a socket
- how message boundaries are handled
- how HTTP request lines are structured
- how protocol validation works before building a full server

It is a good foundation for eventually building a minimal HTTP server from scratch.

## Current components

### TCP listener

Accepts connections and reads incoming data line by line to simulate how raw HTTP requests arrive over a socket.

### UDP sender

Simple CLI tool to send messages over UDP for testing and experimentation.

### HTTP request parser

Parses and validates HTTP/1.1 request lines, including method, path, and version.

## Example request

GET / HTTP/1.1
Host: localhost:42069

## Getting started

### Prerequisites

- Go 1.25+

### Clone the repo

git clone https://github.com/muafa7/http-from-tcp.git
cd http-from-tcp

### Run tests

go test ./...

### Run the TCP listener

go run ./cmd/tcplistener

### Run the UDP sender

go run ./cmd/udpsender

## Learning focus

- Go networking basics
- TCP socket handling
- streamed input processing
- HTTP request parsing
- protocol validation

## Limitations

- no full HTTP server yet
- no header/body parsing beyond request line
- no routing or response handling

## Future improvements

- header parsing
- request body support
- HTTP response writing
- routing system
- full HTTP server implementation

## Tech stack

- Go
- TCP / UDP
- HTTP/1.1
