package main

import (
	"fmt"
	"io"
	"log"
	"net"
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
	// file, err := os.Open("./messages.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error %s\n", err)
			continue
		}

		fmt.Println("connection accepted")

		for result := range getLinesChannel(conn) {
			fmt.Println(result)
		}

		fmt.Println("connection closed")

	}
}
