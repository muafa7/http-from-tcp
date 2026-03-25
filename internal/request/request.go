package request

import (
	"errors"
	"io"
	"log"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(requestLine string) (RequestLine, error) {
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return RequestLine{}, errors.New("Not valid request line")
	}

	method := parts[0]
	requestTarget := parts[1]
	versionPart := parts[2]

	if method == "" {
		return RequestLine{}, errors.New("invalid method")
	}

	for _, ch := range method {
		if ch < 'A' || ch > 'Z' {
			return RequestLine{}, errors.New("invalid method")
		}
	}

	versionParts := strings.Split(versionPart, "/")
	if len(versionParts) != 2 {
		return RequestLine{}, errors.New("invalid HTTP version")
	}

	if versionParts[0] != "HTTP" || versionParts[1] != "1.1" {
		return RequestLine{}, errors.New("unsupported HTTP version")
	}

	return RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("request error %v\n", err)
		return &Request{}, err
	}

	toString := string(request)

	if checkstring := strings.Contains(toString, "\r\n"); checkstring == false {
		log.Println("No request line provided")
		return &Request{}, errors.New("No request line provided")
	}
	split := strings.Split(toString, "\r\n")
	requestLine := split[0]

	parsedRequestLine, err := parseRequestLine(requestLine)
	if err != nil {
		log.Printf("request error %v\n", err)
		return nil, err
	}

	return &Request{
		RequestLine: parsedRequestLine,
	}, nil
}
