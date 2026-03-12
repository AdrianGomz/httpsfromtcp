package main

import (
	"fmt"
	"http/internal/request"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("error", err)
		}
		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error getting request")
		}

		fmt.Println("Request line")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
	}

}
