package request

import (
	"errors"
	"io"
	"strings"

	"github.com/muafa7/http-from-tcp/internal/headers"
)

const bufferSize = 8

type parseState int

const (
	stateInitialized parseState = iota
	stateParsingHeaders
	stateDone
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       parseState
}

func (r *Request) parse(data []byte) (int, error) {
	if r.state == stateDone {
		return 0, nil
	}

	totalBytesParsed := 0

	for r.state != stateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return totalBytesParsed, err
		}

		if n == 0 {
			return totalBytesParsed, nil
		}

		totalBytesParsed += n
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case stateInitialized:
		requestLine, numOfBytes, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if numOfBytes == 0 {
			return 0, nil
		}

		r.RequestLine = requestLine
		r.Headers = headers.Headers{}
		r.state = stateParsingHeaders
		return numOfBytes, nil

	case stateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if n == 0 {
			return 0, nil
		}

		if done {
			r.state = stateDone
		}

		return n, nil

	case stateDone:
		return 0, nil

	default:
		return 0, errors.New("unknown state")
	}
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(data []byte) (RequestLine, int, error) {
	dataString := string(data)
	idx := strings.Index(dataString, "\r\n")
	if idx == -1 {
		return RequestLine{}, 0, nil
	}

	line := dataString[:idx]

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return RequestLine{}, 0, errors.New("Not valid request line")
	}

	method := parts[0]
	requestTarget := parts[1]
	versionPart := parts[2]

	if method == "" {
		return RequestLine{}, 0, errors.New("invalid method")
	}

	for _, ch := range method {
		if ch < 'A' || ch > 'Z' {
			return RequestLine{}, 0, errors.New("invalid method")
		}
	}

	versionParts := strings.Split(versionPart, "/")
	if len(versionParts) != 2 {
		return RequestLine{}, 0, errors.New("invalid HTTP version")
	}

	if versionParts[0] != "HTTP" || versionParts[1] != "1.1" {
		return RequestLine{}, 0, errors.New("unsupported HTTP version")
	}

	return RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, len(line) + 2, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{
		state: stateInitialized,
	}

	buf := make([]byte, bufferSize)
	readToIndex := 0

	for req.state != stateDone {
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf[:readToIndex])
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if numBytesRead > 0 {
			readToIndex += numBytesRead

			numBytesParsed, parseErr := req.parse(buf[:readToIndex])
			if parseErr != nil {
				return nil, parseErr
			}

			if numBytesParsed > 0 {
				copy(buf, buf[numBytesParsed:readToIndex])
				readToIndex -= numBytesParsed
			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
	}

	if req.state != stateDone {
		return nil, errors.New("incomplete request line")
	}

	return req, nil
}
