package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	stringChannel := make(chan string)

	buff := make([]byte, 8)
	var currentLine string

	go func() {
		defer f.Close()
		defer close(stringChannel)
		for {
			n, err := f.Read(buff)
			if err == io.EOF {
				break
			}
			if err != nil {
				return
			}

			if n > 0 {
				chunk := string(buff[:n])
				parts := strings.Split(chunk, "\n")
				for _, part := range parts[:len(parts)-1] {
					stringChannel <- currentLine + part
					currentLine = ""
				}
				currentLine += parts[len(parts)-1]
			}
		}

		if currentLine != "" {
			stringChannel <- currentLine
		}
	}()

	return stringChannel
}

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatal(err)
	}

	for result := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", result)
	}
}
