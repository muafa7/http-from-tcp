package main

import (
	"fmt"
	"log"
	"net"

	"github.com/muafa7/http-from-tcp/internal/request"
)

func main() {
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

		result, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("Error %s\n", err)
			continue
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", result.RequestLine.Method)
		fmt.Printf("- Target: %s\n", result.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", result.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range result.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
		fmt.Println("Body:")
		fmt.Printf("%s\n", string(result.Body))

		fmt.Println("connection closed")

	}
}
