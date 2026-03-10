package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer f.Close()
		defer close(ch)

		curLine := ""
		for {
			data := make([]byte, 8)
			_, err := f.Read(data)
			if err != nil {
				break
			}
			data_s := string(data)
			strings := strings.Split(data_s, "\n")
			curLine += strings[0]
			if len(strings) > 1 {
				ch <- curLine
				curLine = strings[1]
			}
		}
		if len(curLine) > 0 {
			ch <- string(curLine)
		}
	}()
	return ch

}
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
		fmt.Printf("message accepted\n")
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Printf("read: %s\n", line)
		}
	}

}
