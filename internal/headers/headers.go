package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	dataString := string(data)

	idx := strings.Index(dataString, "\r\n")
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	line := dataString[:idx]

	colonIdx := strings.Index(line, ":")
	if colonIdx == -1 {
		return 0, false, errors.New("not valid format")
	}
	if colonIdx == 0 {
		return 0, false, errors.New("invalid header key")
	}
	if line[colonIdx-1] == ' ' {
		return 0, false, errors.New("not valid spacing")
	}

	key := strings.ToLower(strings.TrimSpace(line[:colonIdx]))
	value := strings.TrimSpace(line[colonIdx+1:])

	if key == "" {
		return 0, false, errors.New("invalid header key")
	}

	const validSpecials = "!#$%&'*+-.^_`|~"

	for _, ch := range key {
		isLower := ch >= 'a' && ch <= 'z'
		isDigit := ch >= '0' && ch <= '9'
		isSpecial := strings.ContainsRune(validSpecials, ch)

		if !isLower && !isDigit && !isSpecial {
			return 0, false, errors.New("invalid header key")
		}
	}

	if currentValue, ok := h[key]; ok {
		h[key] = currentValue + "," + value
	} else {
		h[key] = value
	}

	return len(line) + 2, false, nil
}

func (h Headers) Get(key string) (value string, found bool) {
	lowerKey := strings.ToLower(key)

	if keyValue, ok := h[lowerKey]; ok {
		return keyValue, true
	}

	return "", false
}
